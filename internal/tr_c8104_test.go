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

	"testing"
)

func TestC8104_주식_취소주문(t *testing.T) {
	if !lib.F한국증시_정규시장_거래시간임() {
		t.SkipNow()
	}

	lib.F테스트_에러없음(t, f접속_확인())

	const 수량_정상주문 = int64(25)

	종목 := lib.New종목("069500", "KODEX 200", lib.P시장구분_ETF)
	일일_종목정보, 에러 := lib.F일일_종목정보(종목)
	lib.F테스트_에러없음(t, 에러)

	lib.F테스트_에러없음(t, 에러)

	계좌번호 := f계좌번호by인덱스(1)
	lib.F테스트_다름(t, 계좌번호, "")

	f주문_응답_실시간_정보_구독()

	주문번호_정상주문 := f매수주문_테스트_도우미(t, 계좌번호, 종목,
		lib.P호가유형_지정가, 일일_종목정보.M하한가, 수량_정상주문)

	// 일부 취소
	반복_횟수 := lib.F임의값_생성기().Intn(int(수량_정상주문) - 1)
	for i := 0; i < 반복_횟수; i++ {
		lib.F체크포인트(i)
		f취소주문_테스트_도우미(t, 계좌번호, 종목, lib.P일부, 주문번호_정상주문, 1)
	}

	lib.F체크포인트()

	// 전량 취소
	lib.F문자열_출력("\n전량 취소")
	f취소주문_테스트_도우미(t, 계좌번호, 종목, lib.P잔량, 주문번호_정상주문, 0)

	f주문_응답_실시간_정보_해지()
}
