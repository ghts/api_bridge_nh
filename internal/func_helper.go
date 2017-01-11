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
// #include <windows.h>
// #include "./c_func.h"
import "C"

import (
	"github.com/ghts/lib"
	ini "gopkg.in/ini.v1"

	"bytes"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"reflect"
	"strings"
	"testing"
)

var 초기화_완료 = lib.New안전한_bool(false)

func F초기화() (에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M에러: &에러})

	if 초기화_완료.G값() {
		return nil
	} else if 에러 = 초기화_완료.S값(true); 에러 != nil {
		return nil
	}

	lib.F조건부_패닉(!fDLL존재함(), "DLL파일(%v)을 찾을 수 없습니다.", wmca_dll)
	lib.F조건부_패닉(lib.F파일_없음(f설정화일_경로()), "설정화일(%v)을 찾을 수 없습니다.", f설정화일_경로())

	Go루틴_모음 := [](func(chan lib.T신호) error){Go루틴_TR처리, Go루틴_소켓TR_중계, Go루틴_API_실시간_정보_중계}
	ch초기화 := make(chan lib.T신호, len(Go루틴_모음))

	for _, Go루틴 := range Go루틴_모음 {
		go Go루틴(ch초기화)
	}

	for _, _ = range Go루틴_모음 {
		<-ch초기화
	}

	lib.New소켓_질의(lib.P주소_NH_TR, lib.CBOR, lib.P10초).S질의(&(lib.S질의값_단순TR{TR구분: lib.TR접속})).G응답()

	return nil
}

func f계좌_비밀번호() [44]byte {
	계좌_비밀번호 := [44]byte{}
	if lib.F테스트_모드_실행_중() {
		lib.F바이트_복사_문자열(계좌_비밀번호[:], lib.P긴_공백문자열)
		return 계좌_비밀번호
	}

	lib.F패닉("NH OpenAPI 계좌 비밀번호 제공 함수 작성할 것.")
	return [44]byte{}
}

func f거래_비밀번호() [44]byte {
	거래_비밀번호 := [44]byte{}
	if lib.F테스트_모드_실행_중() {
		lib.F바이트_복사_문자열(거래_비밀번호[:], lib.P긴_공백문자열)
		return 거래_비밀번호
	}

	lib.F패닉("NH OpenAPI 거래 비밀번호 제공 함수 작성할 것.")
	return [44]byte{}
}

func f2등락부호(바이트_모음 [1]byte) uint8 {
	값 := uint8(바이트_모음[0])

	switch 값 {
	case NH상한, NH상승, NH보합, NH하한, NH하락:
		return 값
	// 이하 내용은 게시판에 답변에 따름.
	case 1, 6:
		return NH상한
	case 2, 7:
		return NH상승
	case 3, 0:
		return NH보합
	case 4, 8:
		return NH하한
	case 5, 9:
		return NH하락
	}

	go문자열 := lib.F2문자열(바이트_모음)
	//go문자열 := lib.F2문자열(값)

	switch go문자열 {
	case "1", "6":
		return NH상한
	case "2", "7":
		return NH상승
	case "3", "0":
		return NH보합
	case "4", "8":
		return NH하한
	case "5", "9":
		return NH하락
	}

	lib.F변수값_확인(바이트_모음)
	lib.F변수값_확인(값)
	lib.F변수값_확인(go문자열)
	에러 := lib.New에러("예상치 못한 등락부호 값.")
	panic(에러)

	return 0xFF
}

func f등락부호2정수(등락부호 uint8) int64 {
	switch 등락부호 {
	case NH상한, NH상승:
		return int64(1)
	case NH보합:
		return int64(0)
	case NH하락, NH하한:
		return int64(-1)
	default:
		에러 := lib.New에러("등락부호가 예상된 값이 아님. %v", 등락부호)
		panic(에러)
	}
}

func f올바른_등락부호(값 uint8) bool {
	switch 값 {
	case NH상한, NH상승, NH보합, NH하락, NH하한:
		return true
	default:
		return false
	}
}

func f테스트_등락부호(t *testing.T, 등락부호 uint8, 값, 비교대상, 상한, 하한 int64) {
	switch 등락부호 {
	case NH상한:
		lib.F테스트_같음(t, 값, 상한)
	case NH상승:
		lib.F테스트_참임(t, 값 > 비교대상, 값, 비교대상)
	case NH보합:
		if 값 != 0 && 비교대상 != 0 {
			lib.F테스트_같음(t, 값, 비교대상)
		}
	case NH하락:
		lib.F테스트_참임(t, 값 < 비교대상, 값, 비교대상)
	case NH하한:
		lib.F테스트_같음(t, 값, 하한)
	default:
		lib.F문자열_출력("등락부호가 예상된 값이 아님. %v", 등락부호)
		t.FailNow()
	}
}

func f테스트_등락율(t *testing.T, 부호 uint8, 등락율 float64) {
	switch 부호 {
	case NH상한, NH상승:
		lib.F테스트_참임(t, 등락율 > 0)
	case NH보합:
		// 이게 구체적으로 무슨 의미?? 일단은 임의로 변동폭 10% 이내라고 가정함.
		lib.F테스트_참임(t, 등락율 < 10 && 등락율 > -10)
	case NH하락, NH하한:
		lib.F테스트_참임(t, 등락율 < 0)
	default:
		lib.F문자열_출력("등락부호가 예상된 값이 아님. %v", 부호)
		t.FailNow()
	}
}

func f반복되면_패닉(블록_이름 string, 전체_길이 int, 구조체_길이 uintptr) {
	if 전체_길이 == 0 && 구조체_길이 > 0 {
		lib.F문자열_출력("데이터 길이가 0임. 데이터 구조체 형식이 잘못됨.")
		// '전체_길이' 값이 제대로 수신되지 않음.
		return
	}

	수량 := 전체_길이 / int(구조체_길이)

	if 수량 != 1 {
		에러 := lib.New에러("반복되는 구조체임. %v", 블록_이름)
		panic(에러)
	}
}

func f호출(함수명 string, 인수 ...uintptr) bool {
	if !fDLL존재함() {
		return false
	}

	// Call()의 2번째 반환값은 '윈도우 + C언어'조합에서는 필요없는 듯함.
	// 인터넷에서 찾은 예제 코드들은 모두 2번째 반환값을 '_' 처리함.
	반환값, _, 에러 := windows.NewLazyDLL(wmca_dll).NewProc(함수명).Call(인수...)

	// C언어에서 BOOL의 정의는 0이면 false,그 이외의 값은 true임.
	// 일반적인 프로그래밍 언어는 true부터 먼저 확인해도 되지만
	// C언어의 BOOL은 0인지 (즉, false인지)부터 확인해야 함. (순서에 유의)
	switch {
	case !strings.Contains(에러.Error(), 실행_성공):
		lib.F에러_출력(에러)
		return false
	case 반환값 == 0:
		return false
	default:
		return true
	}
}

func fDLL존재함() bool {
	if 에러 := windows.NewLazyDLL(wmca_dll).Load(); 에러 == nil {
		return true
	}

	// 현재 PATH에 wmca.dll 이 포함되어 있지 않음.
	// github.com/ghts/dependency/NH_OpenAPI 를 포함시킬 것.
	DLL파일_예상경로 := os.Getenv("GOPATH") + "/src/github.com/ghts/ghts_dependency/NH_OpenAPI/" + wmca_dll
	if lib.F파일_없음(DLL파일_예상경로) {
		return false
	}

	lib.F에러2패닉(lib.F실행경로_추가(DLL파일_예상경로))

	return lib.F2참거짓(windows.NewLazyDLL(wmca_dll).Load(), nil, true)
}

func f접속_정보() (아이디, 암호, 공인인증_암호 string) {
	if lib.F파일_없음(f설정화일_경로()) {
		버퍼 := new(bytes.Buffer)
		버퍼.WriteString("NH 설정화일 없음\n")
		버퍼.WriteString("%v가 존재하지 않습니다.\n")
		버퍼.WriteString("sample_config.ini를 참조하여 새로 생성하십시오.")
		lib.F패닉(fmt.Errorf(버퍼.String(), f설정화일_경로()))
	}

	cfg파일, 에러 := ini.Load(f설정화일_경로())
	lib.F에러2패닉(에러)

	섹션, 에러 := cfg파일.GetSection("NH_OpenApi_LogIn_Info")
	lib.F에러2패닉(에러)

	키_ID, 에러 := 섹션.GetKey("ID")
	lib.F에러2패닉(에러)
	아이디 = 키_ID.String()

	키_PWD, 에러 := 섹션.GetKey("PWD")
	lib.F에러2패닉(에러)
	암호 = 키_PWD.String()

	키_CertPWD, 에러 := 섹션.GetKey("CertPWD")
	lib.F에러2패닉(에러)
	공인인증_암호 = 키_CertPWD.String()

	return
}

// ZeroMQ에서 1개의 메시지 안에 여러 개의 값을 저장할 수 있으며,
// 그럴 경우 각각의 값은 프레임이라고 불리우고, 공백문자("")에 의해서 구분된다.
//  메시지의 첫번째 프레임을 추출해서 1번째 반환값에 지정하고, 나머지는 2번째 반환값에 지정.
//  1번째 프레임의 다음 프레임에 내용이 없으면 중간 구분 프레임을 사용함.
func f첫번째_프레임_추출(메시지 [][]byte) (첫부분 []byte, 나머지 [][]byte) {
	첫부분 = 메시지[0]

	if len(메시지) > 1 && len(메시지[1]) == 0 {
		나머지 = 메시지[2:] // 구분자는 건너뜀.
	} else {
		나머지 = 메시지[1:]
	}

	return
}

func f올바른_RT코드(RT코드 string) bool {
	switch RT코드 {
	case lib.NH_RT코스피_호가_잔량,
		lib.NH_RT코스닥_호가_잔량,
		lib.NH_RT코스피_시간외_호가_잔량,
		lib.NH_RT코스닥_시간외_호가_잔량,
		lib.NH_RT코스피_예상_호가_잔량,
		lib.NH_RT코스닥_예상_호가_잔량,
		lib.NH_RT코스피_체결, lib.NH_RT코스닥_체결,
		lib.NH_RT코스피_선물_호가,
		lib.NH_RT코스피_선물_이론가,
		lib.NH_RT코스피_선물_미결제_약정,
		lib.NH_RT코스피_선물_체결,
		lib.NH_RT코스피_선물_스프레드_호가,
		lib.NH_RT코스피_선물_스프레드_체결,
		lib.NH_RT코스피_옵션_호가,
		lib.NH_RT코스피_옵션_체결,
		lib.NH_RT코스피_옵션_이론가,
		lib.NH_RT코스피_옵션_미결제_약정,
		lib.NH_RT주식선물_호가,
		lib.NH_RT주식선물_체결,
		lib.NH_RT주식선물_이론가,
		lib.NH_RT주식선물_미결제약정,
		lib.NH_RT주식선물_스프레드_호가,
		lib.NH_RT주식선물_스프레드_체결,
		lib.NH_RT_ELW_체결, lib.NH_RT_ELW_호가,
		lib.NH_RT_ELW_이론가, lib.NH_RT_ELW_투자지표,
		lib.NH_RT_ELW_실시간_거래원,
		lib.NH_RT코스피_선물_예상_체결,
		lib.NH_RT코스피_옵션_예상_체결,
		lib.NH_RT코스피_주식선물_예상_체결,
		lib.NH_RT코스피_ETF_NAV,
		lib.NH_RT코스닥_ETF_NAV,
		lib.NH_RT선물_단계별_상하한가,
		lib.NH_RT옵션_단계별_상하한가,
		lib.NH_RT주식선물_단계별_상하한가,
		lib.NH_RT코스피_업종지수, lib.NH_RT코스닥_업종지수,
		lib.NH_RT주문_접수, lib.NH_RT주문_체결:
		return true
	default:
		panic(lib.New에러("잘못된 RT 코드. %v", RT코드))
		return false
	}
}

func f2NH주문유형(호가유형 lib.T호가유형, 주문조건 lib.T주문조건) NH주문유형 {
	switch 주문조건 {
	case lib.P주문조건_없음:
		switch 호가유형 {
		case lib.P호가유형_지정가:
			return NH주문유형_지정가
		case lib.P호가유형_시장가:
			return NH주문유형_시장가
		case lib.P호가유형_최유리_지정가:
			return NH주문유형_최유리_지정가
		case lib.P호가유형_최우선_지정가:
			return NH주문유형_최우선_지정가
		case lib.P호가유형_시간외종가_장개시전:
			return NH주문유형_장전_시간외_전일종가
		case lib.P호가유형_시간외종가:
			return NH주문유형_장후_시간외_금일종가
		case lib.P호가유형_시간외단일가:
			return NH주문유형_시간외_단일가
		}
	case lib.P주문조건_FOK:
		switch 호가유형 {
		case lib.P호가유형_지정가:
			return NH주문유형_FOK_지정가
		case lib.P호가유형_시장가:
			return NH주문유형_FOK_시장가
		case lib.P호가유형_최유리_지정가:
			return NH주문유형_FOK_최유리_지정가
		}
	case lib.P주문조건_IOC:
		switch 호가유형 {
		case lib.P호가유형_지정가:
			return NH주문유형_IOC_지정가
		case lib.P호가유형_시장가:
			return NH주문유형_IOC_시장가
		case lib.P호가유형_최유리_지정가:
			return NH주문유형_IOC_최유리_지정가
		}
	}

	lib.F패닉("매매 유형을 결정할 수 없습니다. 주문조건 %v, 호가유형 %v", 주문조건, 호가유형)

	return NH주문유형("")
}

func f2NH정정구분(잔량_일부 lib.T잔량_일부) string {
	switch 잔량_일부 {
	case lib.P잔량:
		return NH정정구분_잔량
	case lib.P일부:
		return NH정정구분_일부
	}

	lib.F패닉("정정 구분을 결정할 수 없습니다. 잔량_일부 %v", 잔량_일부)

	return ""
}

func f2NH취소구분(잔량_일부 lib.T잔량_일부) string {
	switch 잔량_일부 {
	case lib.P잔량:
		return NH취소구분_잔량
	case lib.P일부:
		return NH취소구분_일부
	}

	lib.F패닉("취소 구분을 결정할 수 없습니다. 잔량_일부 %v", 잔량_일부)

	return ""
}

func f2신용거래_구분(값 int) lib.T신용거래_구분 {
	switch 값 {
	case 10:
		return lib.P신용거래_해당없음
	case 21:
		return lib.P신용거래_자기융자신규
	case 22:
		return lib.P신용거래_자기융자상환
	case 23:
		return lib.P신용거래_자기대주신규
	case 24:
		return lib.P신용거래_자기대주상환
	case 31:
		return lib.P신용거래_유통융자신규
	case 32:
		return lib.P신용거래_유통융자상환
	case 33:
		return lib.P신용거래_유통대주신규
	case 34:
		return lib.P신용거래_유통대주상환
	case 61:
		return lib.P신용거래_청약대출상환
	case 62:
		return lib.P신용거래_보통대출상환
	case 63:
		return lib.P신용거래_매입대출신규
	case 64:
		return lib.P신용거래_매입대출상환
	}

	lib.F패닉("예상하지 못한 신용거래 구분코드. %v", 값)

	return lib.P신용거래_해당없음
}

func f2매수_매도(값 interface{}) lib.T매수_매도 {
	정수값, 에러 := lib.F2정수(값)
	lib.F에러2패닉(에러)

	switch 정수값 {
	case NH매수:
		return lib.P매수
	case NH매도:
		return lib.P매도
	default:
		lib.F패닉("예상하지 못한 값. %v", 정수값)
	}

	return lib.P매수
}

func f2계좌_인덱스(계좌번호 string) int {
	계좌_인덱스_맵_잠금.RLock()
	계좌_인덱스, ok := 계좌_인덱스_맵[계좌번호]
	계좌_인덱스_맵_잠금.RUnlock()

	switch {
	case !ok:
		return -1
	case 계좌_인덱스 < 1:
		lib.New에러("예상하지 못한 계좌 인덱스. %v", 계좌_인덱스)
	}

	return 계좌_인덱스
}

func f계좌번호by인덱스(인덱스 int) string {
	계좌_인덱스_맵_잠금.RLock()
	defer 계좌_인덱스_맵_잠금.RUnlock()

	for 계좌번호, 계좌_인덱스 := range 계좌_인덱스_맵 {
		if 계좌_인덱스 == 인덱스 {
			return 계좌번호
		}
	}

	return ""
}

func f주문_응답_실시간_정보_구독() {
	// 이전에 실행된 다른 테스트의 잔여물 제거.
	// 이전 테스트에서 남겨진 실시간 정보 비워냄.
	for i := 0; i < len(ch주문_응답); i++ {
		<-ch주문_응답
	}

	질의값 := new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR실시간_정보_구독
	질의값.TR코드 = lib.NH_RT주문_접수

	lib.F대기(lib.P300밀리초)
	lib.New채널_질의(ch실시간_정보_구독_및_해지, lib.P5초, 1).S질의(질의값).G응답()

	질의값 = new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR실시간_정보_구독
	질의값.TR코드 = lib.NH_RT주문_체결

	lib.F대기(lib.P300밀리초)
	lib.New채널_질의(ch실시간_정보_구독_및_해지, lib.P5초, 1).S질의(질의값).G응답()
}

func f주문_응답_실시간_정보_해지() {
	질의값 := new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR실시간_정보_구독
	질의값.TR코드 = lib.NH_RT주문_접수

	// 거래소에서 보내주는 주문 응답 실시간 데이터 구독 해지.
	lib.F대기(lib.P300밀리초)
	lib.New채널_질의(ch실시간_정보_구독_및_해지, lib.P5초, 1).S질의(lib.NH_RT주문_접수)

	lib.F대기(lib.P300밀리초)
	lib.New채널_질의(ch실시간_정보_구독_및_해지, lib.P5초, 1).S질의(lib.NH_RT주문_체결)

	// 현재 테스트의 잔여물 실시간 정보 제거. (이후 다른 테스트에 끼치는 영향 최소화)
	for i := 0; i < len(ch주문_응답); i++ {
		<-ch주문_응답
	}
}

func f설정화일_경로() string {
	현재_패키지명 := reflect.TypeOf(NH수신_데이터_블록{}).PkgPath()

	return os.Getenv("GOPATH") + "\\src\\" + 현재_패키지명 + "\\config.ini"
}
