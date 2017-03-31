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
// #include "./c_func.h"
import "C"

import (
	"bytes"
	"github.com/ghts/lib"
	"runtime"
	"time"
	"unsafe"
)

var TR식별번호 = lib.New안전한_일련번호()
var TR처리_Go루틴_실행_중 = lib.New안전한_bool(false)
var ch조회 = make(chan lib.I채널_질의, 1000)
var ch주문 = make(chan lib.I채널_질의, 1000)
var ch실시간_정보_구독_및_해지 = make(chan lib.I채널_질의, 100)
var ch실시간_정보_일괄_해지 = make(chan lib.I채널_질의, 100)
var ch접속 = make(chan lib.I채널_질의, 10)
var ch접속_해제 = make(chan lib.I채널_질의, 10)
var ch접속됨 = make(chan lib.I채널_질의, 10)
var ch실시간_정보 = make(chan interface{}, 10000)
var ch주문_응답 = make(chan *lib.NH주문_응답, 1000)
var 대기항목_맵 = make(map[int64]*lib.S콜백_대기항목)

// NH OpenAPI은 thread-safe하다고 명시되어 있지 않으므로 thread-unsafe하다고 봐야함.
// API호출 및 콜백 처리가 1번에 1개씩만 이루어지도록 하기 위하여 Go루틴을 사용함.
func Go루틴_TR처리(ch초기화 chan lib.T신호) (에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M에러: &에러})

	if TR처리_Go루틴_실행_중.G값() {
		ch초기화 <- lib.P신호_초기화
		return nil
	} else if 에러 = TR처리_Go루틴_실행_중.S값(true); 에러 != nil {
		ch초기화 <- lib.P신호_초기화
		return 에러
	}

	defer TR처리_Go루틴_실행_중.S값(false)

	// Win32 API는 싱글 스레드를 기반으로 했으므로 현재 스레드를 고정.
	// F단일_스레드_모드()로 대체하면 XingAPI가 동작하지 않음.
	runtime.LockOSThread()

	ch정기_점검 := time.NewTicker(time.Second)
	ch종료 := lib.F공통_종료_채널()
	ch초기화 <- lib.P신호_초기화

	for {
		select {
		case 질의 := <-ch조회:
			f조회TR_처리(질의)
		case 질의 := <-ch주문:
			f주문TR_처리(질의)
		case 질의 := <-ch실시간_정보_구독_및_해지:
			f실시간_정보_TR_처리(질의)
		case 질의 := <-ch접속됨:
			질의.S응답(lib.New채널_메시지(f접속됨()))
		case 질의 := <-ch접속:
			f접속_TR처리(질의)
		case 질의 := <-ch접속_해제:
			f접속해제_TR처리(질의)
		case 질의 := <-ch실시간_정보_일괄_해지:
			f실시간_정보_모두_해지_TR처리(질의)
		case <-ch정기_점검.C:
			f정기_점검()
		case <-ch종료:
			f종료TR_처리()
			return nil
		default: // 처리해야 할 TR이 없을 경우 대기하지 않도록 하기 위함.
			lib.F실행권한_양보()
		}

		C.ProcessWindowsMessage(1) // 윈도우 메시지 처리. 인수는 컴파일 에러방지용.
	}
}

func f조회TR_처리(질의 lib.I채널_질의) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		질의.S응답(lib.New채널_메시지_에러(r))
	}})

	lib.F조건부_패닉(!f접속됨(), "NH API에 접속되어 있지 않습니다.")
	lib.F에러2패닉(질의.G검사(1))

	// 'I소켓_메시지'는 'I채널_질의'를 통해서 전달될 때, '[]*S바이트_변환_매개체' 형태로 바뀐다.
	질의값, 에러 := lib.F바이트_변환값_해석(질의.G값(0).([]*lib.S바이트_변환_매개체)[0])
	lib.F에러2패닉(에러)

	// 계좌번호 인덱스란 로그인시 수신한 계좌번호의 순서.
	// 계좌 인덱스는 숫자형 ‘1’부터 시작되며 ‘0’은 계좌번호를 사용하지 않는다는 의미.
	// 현재 c1101, c1151 TR만 지원하므로, 계좌 인덱스가 굳이 필요 없음.
	var 계좌_인덱스 = -1
	var 길이 int
	var 변환_형식 = 질의.G값(0).([]*lib.S바이트_변환_매개체)[0].G변환_형식()
	var c데이터 unsafe.Pointer
	defer func() {
		if c데이터 != nil {
			C.free(c데이터)
		}
	}()

	TR코드 := 질의값.(lib.I질의값).G_TR코드()
	switch TR코드 {
	case lib.NH_TR주식_현재가_조회:
		종목코드 := 질의값.(*lib.S질의값_단일종목).M종목코드
		c데이터 = NewTc1101InBlock(종목코드)
		길이 = int(unsafe.Sizeof(C.Tc1101InBlock{}))
		계좌_인덱스 = 0 // 투자 정보성 조회에서는 계좌번호를 지정하지 않아야 함.
	case lib.NH_TR_ETF_현재가_조회:
		종목코드 := 질의값.(*lib.S질의값_단일종목).M종목코드
		c데이터 = NewTc1151InBlock(종목코드)
		길이 = int(unsafe.Sizeof(C.Tc1151InBlock{}))
		계좌_인덱스 = 0
	case lib.NH_TR주식_잔고_조회, lib.NH_TR주식_주문_체결_조회,
		lib.NH_TR매도_가능_수량, lib.NH_TR매수_가능_수량,
		lib.NH_TR개별_주식_매도_수량__현금주식_only,
		lib.NH_TR선물_현재가_조회, lib.NH_TR선물_스프레드_현재가_조회,
		lib.NH_TR옵션_현재가_조회, lib.NH_TR주식선물_현재가_조회,
		lib.NH_TR주식선물_스프레드_현재가_조회, lib.NH_TR선물_옵션_코드_조회:
		lib.F패닉(" %v : 구현되지 않음.", TR코드)
	case lib.NH_TR_ELW_현재가_조회,
		lib.NH_TR선물_옵션_잔고_조회,
		lib.NH_TR선물_옵션_주문_체결_조회,
		lib.NH_TR선물_옵션_주문_가능_수량_조회,
		lib.NH_TR선물_옵션_청산_가능_수량_조회:
		lib.F패닉("%v : 선물, 옵션, ELW 매매 관련 조회 기능 구현할 계획 없음.", TR코드)
	default:
		lib.F패닉("%v : 존재하지 않는 TR코드", TR코드)
	}

	lib.F조건부_패닉(길이 == 0, "입력 데이터 길이가 0임.")
	lib.F조건부_패닉(계좌_인덱스 < 0, "계좌 인덱스가 음수임. %v", 계좌_인덱스)

	대기_항목 := lib.New콜백_대기항목(TR식별번호.G값(), lib.TR조회, TR코드, 질의, 변환_형식)
	대기항목_맵[대기_항목.G식별번호()] = 대기_항목

	실행결과_참거짓 := f일반TR_실행(대기_항목.G식별번호(), TR코드, c데이터, 길이, 계좌_인덱스)
	lib.F조건부_패닉(!실행결과_참거짓, "일반TR 실행 에러발생.")
}

func f주문TR_처리(질의 lib.I채널_질의) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		질의.S응답(lib.New채널_메시지_에러(r))
	}})

	lib.F조건부_패닉(!f접속됨(), "NH API에 접속되어 있지 않습니다.")
	lib.F에러2패닉(질의.G검사(1))

	// 'I소켓_메시지'는 'I채널_질의'를 통해서 전달될 때, '[]*S바이트_변환_매개체' 형태로 바뀐다.
	질의값, 에러 := lib.F바이트_변환값_해석(질의.G값(0).([]*lib.S바이트_변환_매개체)[0])
	lib.F에러2패닉(에러)

	// 계좌번호 인덱스란 로그인시 수신한 계좌번호의 순서.
	// 계좌 인덱스는 숫자형 ‘1’부터 시작되며 ‘0’은 계좌번호를 사용하지 않는다는 의미.
	var 계좌_인덱스 = -1
	var 길이 int
	var 변환_형식 = 질의.G값(0).([]*lib.S바이트_변환_매개체)[0].G변환_형식()
	var c데이터 unsafe.Pointer
	defer func() {
		if c데이터 != nil {
			C.free(c데이터)
		}
	}()

	TR코드 := 질의값.(lib.I질의값).G_TR코드()
	switch TR코드 {
	case lib.NH_TR주식_매수:
		주문 := 질의값.(*lib.S질의값_정상주문)
		c데이터 = NewTc8102InBlock(주문)
		길이 = int(unsafe.Sizeof(C.Tc8102InBlock{}))
		계좌_인덱스 = f2계좌_인덱스(주문.M계좌번호)
	case lib.NH_TR주식_매도:
		주문 := 질의값.(*lib.S질의값_정상주문)
		c데이터 = NewTc8101InBlock(주문)
		길이 = int(unsafe.Sizeof(C.Tc8101InBlock{}))
		계좌_인덱스 = f2계좌_인덱스(주문.M계좌번호)
	case lib.NH_TR주식_정정:
		주문 := 질의값.(*lib.S질의값_정정주문_NH)
		c데이터 = NewTc8103InBlock(주문)
		길이 = int(unsafe.Sizeof(C.Tc8103InBlock{}))
		계좌_인덱스 = f2계좌_인덱스(주문.M계좌번호)
	case lib.NH_TR주식_취소:
		주문 := 질의값.(*lib.S질의값_취소주문_NH)
		c데이터 = NewTc8104InBlock(주문)
		길이 = int(unsafe.Sizeof(C.Tc8104InBlock{}))
		계좌_인덱스 = f2계좌_인덱스(주문.M계좌번호)
	case lib.NH_TR_ELW_매도, lib.NH_TR_ELW_매수, lib.NH_TR_ELW_정정_취소,
		lib.NH_TR선물_옵션_매도_매수_주문,
		lib.NH_TR선물_옵션_정정_취소_주문:
		lib.F패닉("선물, 옵션, ELW 매매기능 구현할 계획 없음. '%v'", TR코드)
	default:
		lib.F패닉("예상하지 못한 TR코드. '%v'", TR코드)
	}

	lib.F조건부_패닉(길이 == 0, "입력 데이터 길이가 0임.")
	lib.F조건부_패닉(계좌_인덱스 < 0, "계좌 인덱스가 음수임. %v", 계좌_인덱스)

	대기_항목 := lib.New콜백_대기항목(TR식별번호.G값(), lib.TR주문, TR코드, 질의, 변환_형식)
	대기항목_맵[대기_항목.G식별번호()] = 대기_항목

	실행결과_참거짓 := f일반TR_실행(대기_항목.G식별번호(), TR코드, c데이터, 길이, 계좌_인덱스)
	lib.F조건부_패닉(!실행결과_참거짓, "일반TR 실행 에러발생.")
}

func f실시간_정보_TR_처리(질의 lib.I채널_질의) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{
		M함수with패닉내역: func(r interface{}) {
			질의.S응답(lib.New채널_메시지_에러(r))
		}})

	lib.F에러2패닉(질의.G검사(1))
	바이트_변환값 := 질의.G값(0).([]*lib.S바이트_변환_매개체)[0]
	질의값_인터페이스, 에러 := lib.F바이트_변환값_해석(바이트_변환값)
	lib.F에러2패닉(에러)

	질의값, ok := 질의값_인터페이스.(lib.I질의값)
	lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 질의값_인터페이스)

	TR구분 := 질의값.G_TR구분()
	TR코드 := 질의값.G_TR코드()
	전체_종목코드 := ""
	단위_길이 := 0

	lib.F조건부_패닉(!f접속됨(), "서버에 접속되지 않았습니다.")
	lib.F조건부_패닉(!f올바른_RT코드(TR코드), "잘못된 RT코드. %v", TR코드)

	switch 질의값_인터페이스.(type) {
	case *lib.S질의값_단순TR:
		단위_길이 = 0
	case *lib.S질의값_복수종목:
		질의값 := 질의값_인터페이스.(*lib.S질의값_복수종목)

		종목코드_모음 := 질의값.M종목코드_모음
		lib.F조건부_패닉(len(종목코드_모음) == 0, "종목코드 내용 없음.")

		단위_길이 = len(종목코드_모음[0])
		lib.F조건부_패닉(단위_길이 == 0, "종목코드 길이가 0. %v", 종목코드_모음)

		종목코드_집합 := lib.New문자열_집합()
		버퍼 := new(bytes.Buffer)
		for _, 종목코드 := range 종목코드_모음 {
			switch {
			case 단위_길이 != len(종목코드):
				lib.F패닉("불규칙한 종목코드 길이. 예상값 %v, 실제값 %v, 종목코드 '%v'",
					단위_길이, len(종목코드), 종목코드)
			case 종목코드_집합.G포함(종목코드):
				continue // 중복은 건너뜀.
			}

			버퍼.WriteString(종목코드)
			종목코드_집합.S추가(종목코드) // 중복 검사 용도.
		}
		전체_종목코드 = 버퍼.String()

		lib.F조건부_패닉(len(전체_종목코드) == 0, "종목 내용 없음.")
	default:
		lib.F패닉("예상하지 못한 자료형.")
	}

	var 실행_결과 bool
	switch TR구분 {
	case lib.TR실시간_정보_구독:
		실행_결과 = f실시간_정보_구독(TR코드, 전체_종목코드, 단위_길이)
	case lib.TR실시간_정보_해지:
		실행_결과 = f실시간_정보_해지(TR코드, 전체_종목코드, 단위_길이)
	default:
		lib.F패닉("예상하지 못한 질의_구분. %v", TR구분)
	}

	lib.F조건부_패닉(!실행_결과, "에러 발생. %v %v %v", TR코드, 단위_길이, 전체_종목코드)

	질의.S응답(lib.New채널_메시지(lib.TR응답_완료))
}

func f접속_TR처리(질의 lib.I채널_질의) {
	if f접속됨() {
		lib.New에러("이미 접속되어 있음.")
		질의.S응답(lib.New채널_메시지(lib.TR응답_완료))
		return
	}

	ID, 암호, 공인인증_암호 := f접속_정보()

	if lib.F테스트_모드_실행_중() {
		공인인증_암호 = ""
	}

	대기_항목 := lib.New콜백_대기항목(TR식별번호.G값(), lib.TR접속, "", 질의, lib.P변환형식_기본값)
	대기항목_맵[대기_항목.G식별번호()] = 대기_항목

	ok := f접속(ID, 암호, 공인인증_암호)
	if !ok {
		에러 := lib.New에러("접속 실패.")
		질의.S응답(lib.New채널_메시지_에러(에러))
		return
	}

	// 이후 f접속_콜백_처리()에서 계속 진행됨.
}

func f접속해제_TR처리(질의 lib.I채널_질의) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		질의.S응답(lib.New채널_메시지_에러(r))
	}})

	if !f접속됨() {
		// 접속되지 않았으니 굳이 접속 해제할 필요 없음.
		질의.S응답(lib.New채널_메시지(lib.TR응답_완료))
		return
	}

	//lib.F에러2패닉(질의.G검사(1))

	대기_항목 := lib.New콜백_대기항목(
		TR식별번호.G값(),
		lib.TR접속_해제, "", 질의, lib.P변환형식_기본값)

	대기항목_맵[대기_항목.G식별번호()] = 대기_항목

	lib.F조건부_패닉(!f접속_해제(), "접속 해제 실패.")
}

func f실시간_정보_모두_해지_TR처리(질의 lib.I채널_질의) {
	switch {
	case !f접속됨():
		질의.S응답(lib.New채널_메시지_에러(lib.New에러("NH API에 접속되어 있지 않습니다.")))
		return
	case len(질의.G값_모음()) != 1:
		질의.S응답(lib.New채널_메시지_에러("질의값 길이가 1이 아닙니다."))
		return
	}

	//대기항목 := New콜백_대기(NH실시간_정보_모두_해지, "", 질의)
	//대기항목_맵[대기항목.G식별번호()] = 대기항목

	if f실시간_서비스_모두_해지() {
		질의.S응답(lib.New채널_메시지(lib.TR응답_완료))
	} else {
		질의.S응답(lib.New채널_메시지_에러("실시간 서비스 모두 해지 실패."))
		//delete(대기항목_맵, 대기항목.G식별번호())
	}

	return
}

func f정기_점검() {
	f유효기간_지난_항목_정리()
	f접속_확인()
}

func f접속_확인() error {
	if f접속됨() {
		return nil
	}

	질의 := lib.New채널_질의(ch접속, lib.P10초, 1)

	f접속_TR처리(질의)

	for i := 0; i < 2; i++ {
		응답 := 질의.G응답() // 메시지와 로그인 정보 2가지 응답이 가능함.

		if f접속됨() { // 접속되기만 하면 응답 종류에 상관없이 완료 처리.
			lib.F문자열_출력("접속 확인 : 접속 성공")
			return nil
		}

		lib.F에러_출력(응답.G에러())
		lib.F변수값_확인(응답.G값_모음())
	}

	if !f접속됨() {
		lib.F문자열_출력("접속 확인 : 접속 실패")
	}

	return nil
}

func f유효기간_지난_항목_정리() {
	// 대기 항목 중에서 유효기간이 지난 항목은 정리
	지금 := time.Now()

	for 키, 대기_항목 := range 대기항목_맵 {
		if 지금.After(대기_항목.G유효기간()) {
			delete(대기항목_맵, 키)
		}
	}
}

func f종료TR_처리() {
	lib.New채널_질의(ch실시간_정보_일괄_해지, lib.P30초, 1).S질의().G응답()
	lib.New채널_질의(ch접속_해제, lib.P30초, 1).S질의().G응답()
}
