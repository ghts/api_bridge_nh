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

import (
	"github.com/ghts/lib"

	"strings"
	"testing"
	"time"
)

func TestC8101_C8102_주식_매수_매도(t *testing.T) {
	if !lib.F한국증시_정규시장_거래시간임() {
		t.SkipNow()
	}

	lib.F테스트_에러없음(t, f접속_확인())

	f주문_응답_실시간_정보_구독()

	const 매수주문량 = 1
	const 매도주문량 = 1

	종목 := lib.New종목("069500", "KODEX 200", lib.P시장구분_ETF)
	일일종목정보, 에러 := lib.F일일_종목정보(종목)
	lib.F테스트_에러없음(t, 에러)

	계좌번호 := f계좌번호by인덱스(1)
	lib.F테스트_다름(t, 계좌번호, "")

	매수주문_응답 := new(lib.NH주식_정상주문_응답)
	매도주문_응답 := new(lib.NH주식_정상주문_응답)

	// 매수
	매수주문_질의값 := new(lib.S질의값_정상주문)
	매수주문_질의값.TR구분 = lib.TR주문
	매수주문_질의값.TR코드 = lib.NH_TR주식_매수
	매수주문_질의값.M증권사 = lib.P증권사_NH
	매수주문_질의값.M계좌번호 = 계좌번호
	매수주문_질의값.M매수_매도 = lib.P매수
	매수주문_질의값.M종목코드 = 종목.G코드()
	매수주문_질의값.M주문수량 = 매수주문량
	매수주문_질의값.M주문단가 = 일일종목정보.M상한가
	매수주문_질의값.M호가유형 = lib.P호가유형_지정가 // 모의투자에서는 시장가 매매가 지원되지 않음.
	매수주문_질의값.M주문조건 = lib.P주문조건_없음

	소켓_메시지, 에러 := lib.New소켓_메시지(lib.F임의_변환형식(), 매수주문_질의값)
	lib.F테스트_에러없음(t, 에러)

	질의 := lib.New채널_질의(ch주문, lib.P10초, 1).S질의(소켓_메시지)

	TR데이터_수신 := false
	TR메시지_수신 := false
	TR완료_수신 := false

	for {
		응답 := 질의.G응답()
		if 응답.G에러() != nil {
			if strings.Contains(응답.G에러().Error(), "타임 아웃") {
				break
			}

			lib.F에러_출력(응답.G에러())
			t.FailNow()
		}

		구분, ok := 응답.G값(0).(lib.TR응답_구분)
		lib.F테스트_참임(t, ok, 응답)

		lib.F테스트_같음(t, 구분, lib.TR응답_데이터, lib.TR응답_메시지, lib.TR응답_완료)

		lib.F문자열_출력("매수 %v", 구분.String())

		switch 구분 {
		case lib.TR응답_데이터:
			lib.F테스트_같음(t, 응답.G길이(), 2)

			변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
			lib.F테스트_참임(t, ok, 응답.G값(1))

			switch 변환값.G자료형_문자열() {
			case lib.F자료형_문자열(lib.NH주식_정상주문_응답{}):
				lib.F테스트_거짓임(t, 변환값.IsNil())
				lib.F테스트_에러없음(t, 변환값.G값(매수주문_응답))

				lib.F테스트_참임(t, 매수주문_응답.M주문번호 >= 0, 매수주문_응답.M주문번호)
				lib.F테스트_같음(t, 매수주문_응답.M매수_매도, lib.P매수)
				lib.F테스트_참임(t, 매수주문_응답.M주문_단가 > 0, 매수주문_응답.M주문_단가)
				lib.F테스트_같음(t, 매수주문_응답.M주문_수량, 매수주문량)
			default:
				lib.F문자열_출력("예상치 못한 자료형. %v", 변환값.G자료형_문자열())
				t.FailNow()
			}

			TR데이터_수신 = true
		case lib.TR응답_메시지:
			lib.F테스트_같음(t, 응답.G길이(), 3)

			코드 := 응답.G값(1).(string)  // 코드
			메시지 := 응답.G값(2).(string) // 내용

			lib.F문자열_출력("%v : %v", 코드, 메시지)

			lib.F테스트_참임(t, strings.Contains(메시지, "조회완료")) // 확인 후 수정할 것.

			TR메시지_수신 = true
		case lib.TR응답_완료:
			lib.F테스트_같음(t, 응답.G길이(), 2)

			변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
			lib.F테스트_참임(t, ok, 응답.G값(1))

			if !변환값.IsNil() {
				lib.New에러("예상하지 못한 경우. %v", 변환값.G자료형_문자열())
			}

			TR완료_수신 = true
		}

		if TR데이터_수신 && TR메시지_수신 && TR완료_수신 {
			break
		}
	}

	// 매도
	매도주문_질의값 := new(lib.S질의값_정상주문)
	매도주문_질의값.TR구분 = lib.TR주문
	매도주문_질의값.TR코드 = lib.NH_TR주식_매도
	매도주문_질의값.M증권사 = lib.P증권사_NH
	매도주문_질의값.M계좌번호 = 계좌번호
	매도주문_질의값.M매수_매도 = lib.P매도
	매도주문_질의값.M종목코드 = 종목.G코드()
	매도주문_질의값.M주문수량 = 매도주문량
	매도주문_질의값.M주문단가 = 일일종목정보.M하한가
	매도주문_질의값.M호가유형 = lib.P호가유형_지정가 // 모의투자에서는 시장가 매매가 지원되지 않음.
	매도주문_질의값.M주문조건 = lib.P주문조건_없음

	소켓_메시지, 에러 = lib.New소켓_메시지(lib.F임의_변환형식(), 매도주문_질의값)
	lib.F테스트_에러없음(t, 에러)

	질의 = lib.New채널_질의(ch주문, lib.P30초, 1).S질의(소켓_메시지)

	for {
		응답 := 질의.G응답()
		lib.F테스트_에러없음(t, 응답.G에러())

		구분, ok := 응답.G값(0).(lib.TR응답_구분)
		lib.F테스트_참임(t, ok, 응답)

		lib.F테스트_같음(t, 구분, lib.TR응답_데이터, lib.TR응답_메시지, lib.TR응답_완료)

		lib.F문자열_출력("매도 %v", 구분.String())

		switch 구분 {
		case lib.TR응답_데이터:
			lib.F테스트_같음(t, 응답.G길이(), 2)

			변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
			lib.F테스트_참임(t, ok, 응답.G값(1))

			switch 변환값.G자료형_문자열() {
			case lib.F자료형_문자열(lib.NH주식_정상주문_응답{}):
				lib.F테스트_에러없음(t, 변환값.G값(매도주문_응답))
				lib.F테스트_참임(t, 매도주문_응답.M주문번호 > 0, 매도주문_응답.M주문번호)
				lib.F테스트_같음(t, 매도주문_응답.M매수_매도, lib.P매도)
				lib.F테스트_참임(t, 매도주문_응답.M주문_단가 > 0, 매도주문_응답.M주문_단가)
				lib.F테스트_같음(t, 매도주문_응답.M주문_수량, 매도주문량)
			default:
				lib.F문자열_출력("예상치 못한 자료형. %v", 변환값.G자료형_문자열())
				t.FailNow()
			}

			continue
		case lib.TR응답_메시지:
			lib.F테스트_같음(t, 응답.G길이(), 3)

			코드, ok := 응답.G값(1).(string) // 코드
			lib.F테스트_참임(t, ok)

			메시지, ok := 응답.G값(2).(string) // 내용
			lib.F테스트_참임(t, ok)

			lib.F문자열_출력("%v : %v", 코드, 메시지)

			lib.F테스트_참임(t, strings.Contains(메시지, "조회완료")) // 확인 후 수정할 것.

			continue
		case lib.TR응답_완료:
			lib.F테스트_같음(t, 응답.G길이(), 2)

			변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
			lib.F테스트_참임(t, ok, 응답.G값(1))

			if 변환값.G자료형_문자열() != "<nil>" {
				lib.New에러("예상하지 못한 자료형. %v", 변환값.G자료형_문자열())
			}

			break
		}

		break
	}

	// 실시간 정보 수신 (접수, 체결)
	누적_매수량 := int64(0)
	누적_매도량 := int64(0)
	매수주문_접수 := false
	매도주문_접수 := false

	for {
		주문_응답 := <-ch주문_응답

		lib.F테스트_다름(t, 주문_응답, nil)
		//lib.F문자열_출력(주문_응답.RT코드)

		lib.F테스트_같음(t, 주문_응답.RT코드, lib.NH_RT주문_접수, lib.NH_RT주문_체결)
		p1분전 := time.Now().Add(-1 * lib.P1분)
		p1분후 := time.Now().Add(lib.P1분)

		switch 주문_응답.RT코드 {
		case lib.NH_RT주문_접수:
			lib.F테스트_같음(t, 주문_응답.M주문응답_구분, lib.P주문응답_정상)
			lib.F테스트_참임(t, 주문_응답.M주문번호 > 0, 주문_응답.M주문번호)
			lib.F테스트_같음(t, 주문_응답.M종목코드, 종목.G코드())
			lib.F테스트_같음(t, 주문_응답.M매수_매도, lib.P매수, lib.P매도)
			lib.F테스트_같음(t, 주문_응답.M수량, 1)
			lib.F테스트_참임(t, 주문_응답.M가격 > 0)
			lib.F테스트_참임(t, 주문_응답.M시각.After(p1분전), 주문_응답.M시각, time.Now())
			lib.F테스트_참임(t, 주문_응답.M시각.Before(p1분후), 주문_응답.M시각, time.Now())
			lib.F테스트_같음(t, 주문_응답.M신용거래_구분, lib.P신용거래_해당없음)
			lib.F테스트_같음(t, 주문_응답.M대출일, time.Time{})

			switch 주문_응답.M매수_매도 {
			case lib.P매수:
				lib.F문자열_출력("매수 주문 접수.")
				매수주문_접수 = true
			case lib.P매도:
				lib.F문자열_출력("매도 주문 접수.")
				매도주문_접수 = true
			}
		case lib.NH_RT주문_체결:
			lib.F테스트_같음(t, 주문_응답.M주문응답_구분, lib.P주문응답_체결)
			lib.F테스트_참임(t, 주문_응답.M주문번호 > 0, 주문_응답.M주문번호)
			lib.F테스트_같음(t, 주문_응답.M종목코드, 종목.G코드())
			lib.F테스트_같음(t, 주문_응답.M매수_매도, lib.P매수, lib.P매도)
			lib.F테스트_같음(t, 주문_응답.M수량, 1)
			lib.F테스트_참임(t, 주문_응답.M가격 > 0)
			lib.F테스트_참임(t, 주문_응답.M시각.After(p1분전), 주문_응답.M시각, time.Now())
			lib.F테스트_참임(t, 주문_응답.M시각.Before(p1분후), 주문_응답.M시각, time.Now())
			lib.F테스트_같음(t, 주문_응답.M신용거래_구분, lib.P신용거래_해당없음)
			lib.F테스트_같음(t, 주문_응답.M대출일, time.Time{})

			switch 주문_응답.M매수_매도 {
			case lib.P매수:
				lib.F문자열_출력("매수 주문 체결. %v", 주문_응답.M수량)
				누적_매수량 = +주문_응답.M수량
			case lib.P매도:
				lib.F문자열_출력("매도 주문 체결. %v", 주문_응답.M수량)
				누적_매도량 = +주문_응답.M수량
			}
		}

		if 매수주문_접수 && 매도주문_접수 &&
			누적_매수량 == 매수주문량 &&
			누적_매도량 == 매도주문량 {
			break
		}
	}

	f주문_응답_실시간_정보_해지()
}
