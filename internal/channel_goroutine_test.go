/* Copyright (C) 2015-2018 김운하(UnHa Kim)  unha.kim@kuh.pe.kr

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

Copyright (C) 2015~2017년 UnHa Kim (unha.kim@kuh.pe.kr)

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
	//ini "gopkg.in/ini.v1"
	//"strings"
	//"testing"
	"testing"
	"gopkg.in/ini.v1"
	"strings"
)

func TestCh접속됨(t *testing.T) {
	lib.F대기(lib.P3초)

	lib.F문자열_출력_일시정지_시작()
	defer lib.F문자열_출력_일시정지_해제()

	질의값 := new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR접속됨

	// 소켓으로 질의를 수신한다는 것을 감안한 변환.
	소켓_메시지, 에러 := lib.New소켓_메시지(lib.CBOR, 질의값)
	lib.F테스트_에러없음(t, 에러)

	응답 := lib.New채널_질의(ch접속됨, lib.P30초, 1).S질의(소켓_메시지).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())
	lib.F테스트_같음(t, 응답.G길이(), 1)

	참거짓, ok := 응답.G값(0).(bool)
	lib.F테스트_참임(t, ok)
	lib.F테스트_같음(t, 참거짓, f접속됨())
}

func TestCh접속(t *testing.T) {
	if f접속됨() {
		// 다른 테스트를 진행하려면 접속이 선행되어야 하므로,
		// 굳이 접속을 별도로 테스트 하는 것은 무의미한 듯 함.
		return
	}

	lib.F대기(lib.P3초)

	질의 := lib.New채널_질의(ch접속, lib.P30초, 3).S질의()

	for {
		응답 := 질의.G응답()
		lib.F테스트_에러없음(t, 응답.G에러())

		TR구분, ok := 응답.G값(0).(lib.TR응답_구분)
		lib.F테스트_참임(t, ok, 응답)

		lib.F테스트_같음(t, TR구분, lib.TR응답_완료)

		switch TR구분 {
		case lib.TR응답_완료:
			lib.F테스트_같음(t, 응답.G길이(), 2)

			cfg파일, 에러 := ini.Load("config.ini")
			lib.F테스트_에러없음(t, 에러)

			NH섹션, 에러 := cfg파일.GetSection("NH_OpenApi_LogIn_Info")
			lib.F테스트_에러없음(t, 에러)

			키_ID, 에러 := NH섹션.GetKey("ID")
			lib.F테스트_에러없음(t, 에러)
			ID := 키_ID.String()

			// 2달마다 테스트용 계좌번호를 새로 신청해야 하며, 그 때마다 값이 변함.
			//테스트용_계좌번호, 에러 := 섹션.GetKey("TestAccountNo")
			//lib.F에러_패닉(에러)
			//계좌번호 := 테스트용_계좌번호.String()

			변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
			lib.F테스트_참임(t, ok)

			로그인_정보 := new(lib.NH로그인_정보)
			에러 = 변환값.G값(로그인_정보)
			lib.F테스트_에러없음(t, 에러)

			lib.F테스트_같음(t, 로그인_정보.M접속_ID, ID)
			lib.F테스트_참임(t, strings.HasPrefix(로그인_정보.M접속_서버, "mt")) //??
			lib.F테스트_참임(t, len(로그인_정보.M계좌_목록) == 1)
			lib.F테스트_참임(t, 로그인_정보.M계좌_목록[0].M계좌_인덱스 == 1)
			lib.F테스트_참임(t, f접속됨())

			return
		case lib.TR응답_메시지:
			lib.F테스트_같음(t, 응답.G길이(), 3)

			_, ok := 응답.G값(1).(string) // 코드
			lib.F테스트_참임(t, ok)

			메시지, ok := 응답.G값(2).(string) // 메시지
			lib.F테스트_참임(t, ok)

			lib.F테스트_참임(t, strings.Contains(메시지, "로그인"))
			lib.F테스트_참임(t, strings.Contains(메시지, "성공"))
		}
	}
}

/* 접속 테스트 과정에서 이미 테스트 진행되었음.
func TestCh접속_해제(t *testing.T) {
	lib.F대기(lib.P3초)
	lib.F테스트_에러없음(t, f접속_확인())

	질의값 := new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR접속_해제

	소켓_메시지, 에러 := lib.New소켓_메시지(lib.F임의_변환형식(), 질의값)
	lib.F테스트_에러없음(t, 에러)

	// 접속 해제
	응답 := lib.New채널_질의(ch접속_해제, lib.P30초, 1).S질의(소켓_메시지).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())
	lib.F테스트_같음(t, 응답.G길이(), 1)

	구분, ok := 응답.G값(0).(lib.TR응답_구분)
	lib.F테스트_참임(t, ok)
	lib.F테스트_같음(t, 구분, lib.TR응답_완료)
	lib.F테스트_거짓임(t, f접속됨())
}
*/
