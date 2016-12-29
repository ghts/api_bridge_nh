/* Copyright (C) 2015-2016 김운하(UnHa Kim)  unha.kim@kuh.pe.kr

이 파일은 GHTS의 일부입니다.

이 프로그램은 자유 소프트웨어입니다.
소프트웨어의 피양도자는 자유 소프트웨어 재단이 공표한 GNU LGPL 2.1판
규정에 따라 프로그램을 개작하거나 재배포할 수 있습니다.

이 프로그램은 유용하게 사용될 수 있으리라는 희망에서 배포되고 있지만,
특정한 목적에 적합하다거나, 이익을 안겨줄 수 있다는 묵시적인 보증을 포함한
어떠한 형태의 보증도 제공하지 않습니다.
보다 자세한 사항에 대해서는 GNU LGPL 2.1판을 참고하시기 바랍니다.
GNU LGPL 2.1판은 이 프로그램과 함께 제공됩니다.
만약, 이 문서가 누락되어 있다면 자유 소프트웨어 재단으로 문의하시기 바랍니다.
(자유 소프트웨어 재단 : Free Software Foundation, Inc.,
59 Temple Place - Suite 330, Boston, MA 02111-1307, USA)

Copyright (C) 2015년 UnHa Kim (unha.kim@kuh.pe.kr)

This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, version 2.1 of the License.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. */

package internal

// #cgo CFLAGS: -m32 -Wall
// #include <stdlib.h>
// #include "./c_type.h"
import "C"

import (
	"github.com/ghts/lib"

	"strings"
	"time"
	"unsafe"
)

//export OnConnected_Go
func OnConnected_Go(c *C.LOGINBLOCK) { f콜백_접속(c) }

//export OnDisconnected_Go
func OnDisconnected_Go() { f콜백_접속_해제() }

//export OnMessage_Go
func OnMessage_Go(c *C.OUTDATABLOCK) { f콜백_메시지(c) }

//export OnTrData_Go
func OnTrData_Go(c *C.OUTDATABLOCK) { f콜백_TR데이터(c) }

//export OnComplete_Go
func OnComplete_Go(c *C.OUTDATABLOCK) { f콜백_TR완료(c) }

//export OnRealTimeData_Go
func OnRealTimeData_Go(c *C.OUTDATABLOCK) { f콜백_실시간_데이터(c) }

//export OnError_Go
func OnError_Go(c *C.OUTDATABLOCK) {
	lib.F문자열_출력("OnError_Go")
	f에러_콜백_처리(c)
}

//export OnSocketError_Go
func OnSocketError_Go(에러_코드 C.int) { f소켓_에러_콜백_처리(int(에러_코드)) }

// 콜백(역호출)으로 수신한 데이터를 실제로 처리하는 함수(핸들러?)들

func f콜백_접속(c *C.LOGINBLOCK) {
	로그인_정보 := New로그인_정보(c)

	// 계좌 인덱스 정보 저장
	계좌_인덱스_맵_잠금.Lock()
	for 계좌번호, _ := range 계좌_인덱스_맵 {
		delete(계좌_인덱스_맵, 계좌번호)
	}

	for _, 계좌_정보 := range 로그인_정보.M계좌_목록 {
		계좌_인덱스_맵[계좌_정보.M계좌_번호] = 계좌_정보.M계좌_인덱스
	}
	계좌_인덱스_맵_잠금.Unlock()

	// 콜백 대기 응답
	for _, 대기_항목 := range 대기항목_맵 {
		if 대기_항목.TR구분() == lib.TR접속 {
			대기_항목.G질의().S응답(lib.New채널_메시지(lib.TR응답_완료, 로그인_정보))
		}
	}
}

func f콜백_접속_해제() {
	defer f자원_정리()

	접속_해제_에러 := lib.New에러("접속 해제됨. %v", time.Now())

	for _, 대기_항목 := range 대기항목_맵 {
		switch 대기_항목.TR구분() {
		case lib.TR접속_해제: // 접속 해제 요청이 성공했으므로, 에러가 아님
			대기_항목.G질의().S응답(lib.New채널_메시지(lib.TR응답_완료))
		case lib.TR종료, lib.TR접속: // 종료 작업의 일부이므로 에러가 아님.
			continue
		default: // 나머지 모든 경우에 대해서 에러 회신.
			lib.F문자열_출력("접속 해제 콜백. 예상하지 못한 경우. '%v'", 대기_항목.TR구분())
			대기_항목.G질의().S응답(lib.New채널_메시지_에러(접속_해제_에러))
		}
	}

	// 질의.S회신()은 비동기 방식으로 동작하므로 지우면 안 될 듯 함.
	// 어차피 유효기간이 지나면 자동으로 삭제됨.
	//대기항목_맵 = make(map[int]s콜백_대기)	// 맵을 재생성해서 모든 항목 삭제.
}

func f콜백_메시지(c *C.OUTDATABLOCK) {
	// free()가 2번 이상 중복되면 메모리 에러 발생함.
	// c포인터는 생성자에서 free() 하도록 하자.
	//defer C.free(unsafe.Pointer(c))
	데이터 := New수신_메시지_블록(c)

	lib.F문자열_출력("메시지 : %v\t%v",
		strings.TrimSpace(데이터.M메시지_코드),
		strings.TrimSpace(데이터.M메시지_내용))

	// 해당되는 조회 질의가 존재하면 처리.
	대기_항목, 존재함 := 대기항목_맵[데이터.M식별번호]

	if 존재함 {
		대기_항목.G질의().S응답(lib.New채널_메시지(lib.TR응답_메시지, 데이터.M메시지_코드, 데이터.M메시지_내용))
		return
	}

	접속_대기_항목_찾음 := false
	for _, 대기_항목 := range 대기항목_맵 {
		if 대기_항목.TR구분() == lib.TR접속 {
			접속_대기_항목_찾음 = true
			대기_항목.G질의().S응답(lib.New채널_메시지(lib.TR응답_메시지, 데이터.M메시지_코드, 데이터.M메시지_내용))
		}
	}

	if !접속_대기_항목_찾음 {
		panic(lib.New에러("콜백 메시지 : 대기 질의 존재하지 않으며, 접속 질의 대기 항목도 없음."))
	}
}

func f콜백_TR데이터(c *C.OUTDATABLOCK) {
	식별번호 := int64(c.TrIdNo)

	대기_항목, 존재함 := 대기항목_맵[식별번호]
	if !존재함 {
		panic(lib.New에러("콜백 조회 : 대기 질의 존재하지 않음."))
	}

	응답값 := f_TR데이터_변환(c.DataStruct)
	if 응답값 == nil {
		return
	}

	변환값, 에러 := lib.New바이트_변환_매개체(대기_항목.G변환_형식(), 응답값)
	if 에러 != nil {
		lib.F에러_출력(에러)
	}

	대기_항목.G질의().S응답(lib.New채널_메시지(lib.TR응답_데이터, 변환값))
}

func f콜백_실시간_데이터(c *C.OUTDATABLOCK) {
	// BlockName은 *byte 형태이라서 F2문자열()에서 처리할 수 없음.
	c수신값 := c.DataStruct
	RT코드 := C.GoString(c수신값.BlockName)[:2]
	전체_길이 := int(c수신값.Length)
	데이터 := c수신값.DataString

	if 전체_길이 == 0 || 데이터 == nil {
		return
	}

	switch RT코드 {
	case lib.NH_RT코스피_호가_잔량, lib.NH_RT코스닥_호가_잔량:
		ch실시간_정보 <- NewNH호가_잔량(데이터)
	case lib.NH_RT코스피_시간외_호가_잔량, lib.NH_RT코스닥_시간외_호가_잔량:
		ch실시간_정보 <- NewNH시간외_호가잔량(데이터)
	case lib.NH_RT코스피_예상_호가_잔량, lib.NH_RT코스닥_예상_호가_잔량:
		ch실시간_정보 <- NewNH예상_호가잔량(데이터)
	case lib.NH_RT코스피_체결:
		ch실시간_정보 <- NewNH체결_코스피(데이터)
	case lib.NH_RT코스닥_체결:
		ch실시간_정보 <- NewNH체결_코스닥(데이터)
	case lib.NH_RT코스피_ETF_NAV, lib.NH_RT코스닥_ETF_NAV:
		ch실시간_정보 <- NewNH_ETF_NAV(데이터)
	case lib.NH_RT코스피_업종지수, lib.NH_RT코스닥_업종지수:
		ch실시간_정보 <- NewNH업종_지수(데이터)
	case lib.NH_RT주문_접수:
		ch주문_응답 <- NewNH주문_접수(데이터)
	case lib.NH_RT주문_체결:
		ch주문_응답 <- NewNH주문_체결(데이터)
	default:
		에러 := lib.New에러("예상치 못한 RT코드. %v", RT코드)
		panic(에러)
	}
}

func f콜백_TR완료(c *C.OUTDATABLOCK) {
	식별번호 := int64(c.TrIdNo)

	대기_항목, 존재함 := 대기항목_맵[식별번호]
	if !존재함 {
		panic(lib.New에러("콜백 조회 : 대기 질의 존재하지 않음."))
	}

	변환값, 에러 := lib.New바이트_변환_매개체(lib.P변환형식_기본값, f_TR데이터_변환(c.DataStruct))
	if 에러 != nil {
		lib.F에러_출력(에러)
	}

	대기_항목.G질의().S응답(lib.New채널_메시지(lib.TR응답_완료, 변환값))

	// TR완료 콜백이 와도 이후에 데이터를 수신하는 경우가 많으므로, 대기항목 맵을 삭제하면 안 된다.
	//delete(대기항목_맵, 식별번호)
}

func f에러_콜백_처리(c *C.OUTDATABLOCK) {
	식별번호 := int64(c.TrIdNo)
	대기_항목, 존재함 := 대기항목_맵[식별번호]
	if !존재함 {
		panic(lib.New에러("콜백 에러 : 해당 질의 존재하지 않음."))
	}

	c수신값 := c.DataStruct
	블록_이름 := C.GoString(c수신값.BlockName) // *byte
	바이트_모음 := C.GoBytes(unsafe.Pointer(c수신값.DataString), c수신값.Length)
	에러_메시지 := lib.F2문자열_CP949(바이트_모음)
	수신_에러 := lib.New에러("에러 발생\n%v\n%v", 블록_이름, 에러_메시지)

	대기_항목.G질의().S응답(lib.New채널_메시지_에러(수신_에러))
}

func f소켓_에러_콜백_처리(에러_코드 int) {
	회신 := lib.New채널_메시지_에러(lib.New에러("소켓 에러 발생. 에러코드 : %v", 에러_코드))

	// 모든 대기 중 질의에 대해서 에러 회신?
	for _, 대기_항목 := range 대기항목_맵 {
		대기_항목.G질의().S응답(회신)
	}
}

func f_TR데이터_변환(c *C.RECEIVED) interface{} {
	블록_이름 := C.GoString(c.BlockName)
	전체_길이 := int(c.Length)
	데이터 := c.DataString

	if 전체_길이 == 0 {
		return nil
	}

	switch 블록_이름 {
	case "c1101OutBlock":
		return New주식_현재가_조회_기본_정보(데이터)
	case "c1101OutBlock2":
		수량 := 전체_길이 / int(unsafe.Sizeof(C.Tc1101OutBlock2{}))

		// 큰 배열로 캐스팅 한 다음에 슬라이스를 취함.
		// 충분히 큰 숫자이면 아무 것이나 상관없으며, 반드시 반드시 10000이어야 하는 것은 아님.
		// Go위키에서는 '1 << 30'을 사용하지만, 너무 큰 수를 사용하니까 메모리 범위를 벗어난다고 에러 발생.
		슬라이스 := (*[10000]C.Tc1101OutBlock2)(unsafe.Pointer(데이터))[:수량:수량]
		go슬라이스 := make([]*lib.NH주식_현재가_조회_변동_거래량_정보, 수량)

		for i := 0; i < 수량; i++ {
			c := 슬라이스[i]
			go슬라이스[i] = New주식_현재가_조회_변동_거래량_정보(&c)
		}

		return go슬라이스
	case "c1101OutBlock3":
		return New주식_현재가_조회_동시호가(데이터)
	case "c1151OutBlock":
		return New_ETF_현재가_조회_기본_정보(데이터)
	case "c1151OutBlock2":
		수량 := 전체_길이 / int(unsafe.Sizeof(C.Tc1151OutBlock2{}))

		// 큰 배열로 캐스팅4 한 다음에 슬라이스를 취함.
		// 충분히 큰 숫자이면 아무 것이나 상관없으며, 반드시 반드시 10000이어야 하는 것은 아님.
		// Go위키에서는 '1 << 30'을 사용하지만, 너무 큰 수를 사용하니까 메모리 범위를 벗어난다고 에러 발생.
		슬라이스 := (*[10000]C.Tc1151OutBlock2)(unsafe.Pointer(데이터))[:수량:수량]
		go슬라이스 := make([]*lib.NH_ETF_현재가_조회_변동_거래량_정보, 수량)

		for i := 0; i < 수량; i++ {
			c := 슬라이스[i]
			go슬라이스[i] = New_ETF_현재가_조회_변동_거래_정보(&c)
			//C.free(unsafe.Pointer(&c))
		}

		return go슬라이스
	case "c1151OutBlock3":
		return New_ETF_현재가_조회_동시호가(데이터)
	case "c1151OutBlock4":
		return New_ETF_현재가_조회_ETF자료(데이터)
	case "c1151OutBlock5":
		return New_ETF_현재가_조회_지수_정보(데이터)
	case "c8101OutBlock":
		return New매도주문_응답(데이터)
	case "c8102OutBlock":
		return New매수주문_응답(데이터)
	case "c8103OutBlock":
		return New정정주문_응답(데이터)
	case "c8104OutBlock":
		return New취소주문_응답(데이터)
	case "c8101InBlock", "c8102InBlock", "c8103InBlock", "c8104InBlock":
		// 필요없는 입력값 블록이 수신되면 처리하지 않고 건너뜀.
		return nil
	default:
		lib.F문자열_출력("예상하지 못한 블록 이름. %v", 블록_이름)
		return nil
	}
}
