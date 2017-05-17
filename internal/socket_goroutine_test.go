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
	"bytes"
	"github.com/ghts/lib"
	"github.com/go-mangos/mangos"
	"strings"
	"testing"
	"time"
)

func TestF자료형_상수(t *testing.T) {
	lib.F테스트_같음(t, nh주식_현재가_조회_기본, lib.F자료형_문자열(lib.NH주식_현재가_조회_기본_정보{}))
	lib.F테스트_같음(t, nh주식_현재가_조회_변동, lib.F자료형_문자열(make([]*lib.NH주식_현재가_조회_변동_거래량_정보, 0)))
	lib.F테스트_같음(t, nh주식_현재가_조회_동시호가, lib.F자료형_문자열(lib.NH주식_현재가_조회_동시호가_정보{}))

	lib.F테스트_같음(t, nhETF_현재가_조회_기본, lib.F자료형_문자열(lib.NH_ETF_현재가_조회_기본_정보{}))
	lib.F테스트_같음(t, nhETF_현재가_조회_변동, lib.F자료형_문자열(make([]*lib.NH_ETF_현재가_조회_변동_거래량_정보, 0)))
	lib.F테스트_같음(t, nhETF_현재가_조회_동시호가, lib.F자료형_문자열(lib.NH_ETF_현재가_조회_동시호가_정보{}))
	lib.F테스트_같음(t, nhETF_현재가_조회_ETF, lib.F자료형_문자열(lib.NH_ETF_현재가_조회_ETF정보{}))
	lib.F테스트_같음(t, nhETF_현재가_조회_지수, lib.F자료형_문자열(lib.NH_ETF_현재가_조회_지수_정보{}))

	lib.F테스트_같음(t, nh호가_잔량, lib.F자료형_문자열(lib.NH호가_잔량{}))
	lib.F테스트_같음(t, nh시간외_호가_잔량, lib.F자료형_문자열(lib.NH시간외_호가잔량{}))
	lib.F테스트_같음(t, nh예상_호가_잔량, lib.F자료형_문자열(lib.NH예상_호가잔량{}))
	lib.F테스트_같음(t, nh체결, lib.F자료형_문자열(lib.NH체결{}))
	lib.F테스트_같음(t, nhETF_NAV, lib.F자료형_문자열(lib.NH_ETF_NAV{}))
	lib.F테스트_같음(t, nh업종_지수, lib.F자료형_문자열(lib.NH업종지수{}))
}

func TestTR소켓_주식_현재가_조회(t *testing.T) {
	if !lib.F한국증시_정규시장_거래시간임() {
		t.SkipNow()
	}

	lib.F대기(lib.P3초)
	lib.F테스트_에러없음(t, f접속_확인())

	변환_형식 := lib.F임의_변환형식()
	소켓_질의 := lib.New소켓_질의(lib.P주소_NH_TR, 변환_형식, lib.P30초)

	질의값 := new(lib.S질의값_단일종목)
	질의값.TR구분 = lib.TR조회
	질의값.TR코드 = lib.NH_TR주식_현재가_조회
	질의값.M종목코드 = lib.F임의_종목_코스피_주식().G코드()

	응답_메시지 := 소켓_질의.S질의(질의값).G응답()
	lib.F테스트_에러없음(t, 응답_메시지.G에러())
	lib.F테스트_같음(t, 응답_메시지.G길이(), 1)

	조회_응답값 := new(lib.NH주식_현재가_조회_응답)
	lib.F테스트_에러없음(t, 응답_메시지.G값(0, 조회_응답값))

	f주식_현재가_조회_기본_정보_테스트(t, 조회_응답값.M기본_정보, 질의값.M종목코드)
	f주식_현재가_조회_변동_거래량_정보_테스트(t, 조회_응답값.M기본_정보, 조회_응답값.M변동_거래량_정보)
	f주식_현재가_조회_동시호가_정보_테스트(t, 조회_응답값.M기본_정보, 조회_응답값.M동시호가_정보)
}

func TestTR소켓_ETF_현재가_조회(t *testing.T) {
	if !lib.F한국증시_정규시장_거래시간임() {
		t.SkipNow()
	}

	lib.F대기(lib.P3초)
	lib.F테스트_에러없음(t, f접속_확인())

	변환_형식 := lib.F임의_변환형식()
	소켓_질의 := lib.New소켓_질의(lib.P주소_NH_TR, 변환_형식, lib.P30초)

	질의값 := new(lib.S질의값_단일종목)
	질의값.TR구분 = lib.TR조회
	질의값.TR코드 = lib.NH_TR_ETF_현재가_조회
	질의값.M종목코드 = lib.F임의_종목_ETF().G코드()

	응답_메시지 := 소켓_질의.S질의(질의값).G응답()
	lib.F테스트_에러없음(t, 응답_메시지.G에러())
	lib.F테스트_같음(t, 응답_메시지.G길이(), 1)

	조회_응답값 := new(lib.NH_ETF_현재가_조회_응답)
	lib.F테스트_에러없음(t, 응답_메시지.G값(0, 조회_응답값))

	f_ETF_현재가_조회_기본_정보_테스트(t, 조회_응답값.M기본_정보, 질의값.M종목코드)
	f_ETF_현재가_조회_변동_거래_정보_테스트(t, 조회_응답값.M기본_정보, 조회_응답값.M변동_거래량_정보)
	f_ETF_현재가_조회_동시호가_정보_테스트(t, 조회_응답값.M기본_정보, 조회_응답값.M동시호가_정보)
	f_ETF_현재가_조회_ETF자료_테스트(t, 조회_응답값.M기본_정보, 조회_응답값.ETF_정보)
	f_ETF_현재가_조회_지수_정보_테스트(t, 조회_응답값.M기본_정보, 조회_응답값.ETF_정보, 조회_응답값.M지수_정보)
}

func TestTR소켓_실시간_서비스_등록_및_해지(t *testing.T) {
	if !lib.F한국증시_정규시장_거래시간임() {
		return // 실시간 정보 테스트는 거래시간에만 하도록 함.
	}

	lib.F대기(lib.P3초)
	lib.F테스트_에러없음(t, f접속_확인())
	변환_형식 := lib.F임의_변환형식()

	TR소켓, 에러 := lib.New소켓REQ(lib.P주소_NH_TR)
	lib.F테스트_에러없음(t, 에러)
	defer TR소켓.Close()

	소켓SUB_CBOR, 에러 := lib.New소켓SUB(lib.P주소_NH_실시간_CBOR)
	lib.F테스트_에러없음(t, 에러)
	defer 소켓SUB_CBOR.Close()

	소켓SUB_JSON, 에러 := lib.New소켓SUB(lib.P주소_NH_실시간_JSON)
	lib.F테스트_에러없음(t, 에러)
	defer 소켓SUB_JSON.Close()

	소켓SUB_MsgPack, 에러 := lib.New소켓SUB(lib.P주소_NH_실시간_MsgPack)
	lib.F테스트_에러없음(t, 에러)
	defer 소켓SUB_MsgPack.Close()

	// 실시간 서비스 등록
	종목모음_코스피, 에러 := lib.F종목모음_코스피()
	lib.F테스트_에러없음(t, 에러)

	종목코드_모음 := lib.F종목코드_추출(종목모음_코스피, 20)

	질의값 := new(lib.S질의값_복수종목)
	질의값.TR구분 = lib.TR실시간_정보_구독
	질의값.TR코드 = lib.NH_RT코스피_체결
	질의값.M종목코드_모음 = 종목코드_모음

	// 여기에서 에러 발생.
	소켓_질의 := lib.New소켓_질의(lib.P주소_NH_TR, 변환_형식, lib.P30초)
	응답_메시지 := 소켓_질의.S질의(질의값).G응답()
	lib.F테스트_에러없음(t, 응답_메시지.G에러())
	lib.F테스트_같음(t, 응답_메시지.G길이(), 1)

	// 실시간 정보 수신 확인
	버퍼 := new(bytes.Buffer)
	for _, 종목코드 := range 질의값.M종목코드_모음 {
		버퍼.WriteString(종목코드)
	}

	전체_종목코드 := 버퍼.String()

	for i := 0; i < 10; i++ {
		체결_정보 := new(lib.NH체결)
		응답 := lib.New소켓_메시지_내용없음().S소켓_수신(소켓SUB_CBOR, lib.P무기한)
		lib.F테스트_에러없음(t, 응답.G에러())
		lib.F테스트_에러없음(t, 응답.G값(0, 체결_정보))
		lib.F테스트_다름(t, strings.TrimSpace(체결_정보.M종목코드), "")
		lib.F테스트_참임(t, strings.Contains(전체_종목코드, 체결_정보.M종목코드))
	}

	// 실시간 서비스 해지
	lib.F대기(lib.P300밀리초)

	질의값 = new(lib.S질의값_복수종목)
	질의값.TR구분 = lib.TR실시간_정보_해지
	질의값.TR코드 = lib.NH_RT코스피_체결
	질의값.M종목코드_모음 = 종목코드_모음

	응답 := 소켓_질의.S질의(질의값).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())
	lib.F테스트_같음(t, 응답.G길이(), 1)

	var 응답_구분 lib.TR응답_구분
	lib.F테스트_에러없음(t, 응답.G값(0, &응답_구분))
	lib.F테스트_같음(t, 응답_구분, lib.TR응답_완료)
}

func f업종지수_모음_코스피() []*lib.S종목 {
	업종지수_모음 := make([]*lib.S종목, 0)

	업종지수_코드_모음 := []string{"00", "01", "02", "03", "04", "05",
		"06", "07", "08", "09", "10", "11", "12", "13", "14", "15",
		"16", "17", "18", "19", "20", "21", "22", "23", "24", "25",
		"26", "27", "28", "29", "30", "32", "39", "40", "41", "42",
		"43", "44", "45", "46", "47", "48", "49"}

	for _, 코드 := range 업종지수_코드_모음 {
		업종지수_모음 = append(업종지수_모음, lib.New업종지수(코드, "", lib.P시장구분_코스피))
	}

	return 업종지수_모음
}

func f업종지수_모음_코스닥() []*lib.S종목 {
	업종지수_모음 := make([]*lib.S종목, 0)

	업종지수_코드_모음 := []string{"01", "03", "04", "06", "07",
		"08", "10", "11", "12", "13", "14", "15", "16", "17", "18",
		"19", "20", "21", "22", "23", "24", "25", "26", "27", "28",
		"29", "30", "31", "32", "33", "34", "35", "36", "37", "38",
		"39", "40", "43", "44", "45", "46", "47", "48", "49"}

	for _, 코드 := range 업종지수_코드_모음 {
		업종지수_모음 = append(업종지수_모음, lib.New업종지수(코드, "", lib.P시장구분_코스닥))
	}

	return 업종지수_모음
}

type sRT테스트_항목 struct {
	검사_진행 bool
	질의_인수 *sRT질의_인수
	실행_함수 func(t *testing.T, 질의_인수 *sRT질의_인수)
}

type sRT질의_인수 struct {
	질의값           *lib.S질의값_복수종목
	TR소켓          mangos.Socket
	소켓SUB_CBOR    mangos.Socket
	소켓SUB_MsgPack mangos.Socket
	소켓SUB_JSON    mangos.Socket
	t             *testing.T
}

func TestRT소켓_실시간_정보_수신(t *testing.T) {
	lib.F대기(lib.P3초)

	TR소켓, 에러 := lib.New소켓REQ(lib.P주소_NH_TR)
	lib.F테스트_에러없음(t, 에러)
	defer TR소켓.Close()

	소켓SUB_CBOR, 에러 := lib.New소켓SUB(lib.P주소_NH_실시간_CBOR)
	lib.F테스트_에러없음(t, 에러)
	defer 소켓SUB_CBOR.Close()

	소켓SUB_MsgPack, 에러 := lib.New소켓SUB(lib.P주소_NH_실시간_MsgPack)
	lib.F테스트_에러없음(t, 에러)
	defer 소켓SUB_MsgPack.Close()

	소켓SUB_JSON, 에러 := lib.New소켓SUB(lib.P주소_NH_실시간_JSON)
	lib.F테스트_에러없음(t, 에러)
	defer 소켓SUB_JSON.Close()

	종목모음_코스피, 에러 := lib.F종목모음_코스피()
	lib.F테스트_에러없음(t, 에러)

	종목모음_코스닥, 에러 := lib.F종목모음_코스닥()
	lib.F테스트_에러없음(t, 에러)

	종목모음_ETF, 에러 := lib.F종목모음_ETF()
	lib.F테스트_에러없음(t, 에러)

	업종지수_모음_코스피 := f업종지수_모음_코스피()
	업종지수_모음_코스닥 := f업종지수_모음_코스닥()

	테스트_항목_모음 := []sRT테스트_항목{
		sRT테스트_항목{
			검사_진행: lib.F한국증시_정규시장_거래시간임(),
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "h1",
				M종목코드_모음: lib.F종목코드_추출(종목모음_코스피, 20)})},
			실행_함수: h1_k3_수신_테스트},
		sRT테스트_항목{
			검사_진행: lib.F한국증시_정규시장_거래시간임(),
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "k3",
				M종목코드_모음: lib.F종목코드_추출(종목모음_코스닥, 20)})},
			실행_함수: h1_k3_수신_테스트},
		sRT테스트_항목{
			검사_진행: lib.F한국증시_시간외_종가매매_시간임(),
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "h2",
				M종목코드_모음: lib.F종목코드_추출(종목모음_코스피, 20)})},
			실행_함수: h2_k4_수신_테스트},
		sRT테스트_항목{
			검사_진행: lib.F한국증시_시간외_종가매매_시간임(),
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "k4",
				M종목코드_모음: lib.F종목코드_추출(종목모음_코스닥, 20)})},
			실행_함수: h2_k4_수신_테스트},
		sRT테스트_항목{
			검사_진행: lib.F한국증시_동시호가_시간임(),
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "h3",
				M종목코드_모음: lib.F종목코드_추출(종목모음_코스피, 20)})},
			실행_함수: h3_k5_수신_테스트},
		sRT테스트_항목{
			검사_진행: lib.F한국증시_동시호가_시간임(),
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "k5",
				M종목코드_모음: lib.F종목코드_추출(종목모음_코스닥, 20)})},
			실행_함수: h3_k5_수신_테스트},
		sRT테스트_항목{
			검사_진행: lib.F한국증시_정규시장_거래시간임(),
			// || lib.F한국증시_시간외_종가매매_시간임(),
			// || lib.F한국증시_시간외_단일가매매_시간임,
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "j8",
				M종목코드_모음: lib.F종목코드_추출(종목모음_코스피, 20)})},
			실행_함수: j8_k8_수신_테스트},
		sRT테스트_항목{
			검사_진행: lib.F한국증시_정규시장_거래시간임(),
			// || lib.F한국증시_시간외_종가매매_시간임(),
			// || lib.F한국증시_시간외_단일가매매_시간임,
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "k8",
				M종목코드_모음: lib.F종목코드_추출(종목모음_코스닥, 20)})},
			실행_함수: j8_k8_수신_테스트},
		sRT테스트_항목{
			검사_진행: lib.F한국증시_정규시장_거래시간임(),
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "j1",
				M종목코드_모음: lib.F종목코드_추출(종목모음_ETF, 20)})},
			실행_함수: j1_j0_수신_테스트},
		// 'j0'에 대해서 수신되는 값이 없어서 잠정 보류.
		//sRT테스트_항목{
		//	검사_진행: lib.F한국증시_거래시간임(),
		//	질의_인수: &sRT질의_인수{질의값: NewNH실시간_정보_질의("j0", 종목모음_ETF)},
		//	실행_함수: j1_j0_수신_테스트},
		// 업종지수 모음 함수를 못 찾겠음. 테스트 건너뜀.
		sRT테스트_항목{
			검사_진행: lib.F한국증시_정규시장_거래시간임(),
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "u1",
				M종목코드_모음: lib.F종목코드_추출(업종지수_모음_코스피, 20)})},
			실행_함수: u1_k1_수신_테스트},
		sRT테스트_항목{
			검사_진행: lib.F한국증시_정규시장_거래시간임(),
			질의_인수: &sRT질의_인수{질의값: &(lib.S질의값_복수종목{
				TR구분:     lib.TR실시간_정보_구독,
				TR코드:     "k1",
				M종목코드_모음: lib.F종목코드_추출(업종지수_모음_코스닥, 20)})},
			실행_함수: u1_k1_수신_테스트}}

	for _, 테스트_항목 := range 테스트_항목_모음 {
		if !테스트_항목.검사_진행 {
			//lib.F문자열_출력("체크포인트 skip %v ", 테스트_항목.질의_인수.질의값.RT코드)
			continue
		}

		테스트_항목.질의_인수.TR소켓 = TR소켓
		테스트_항목.질의_인수.소켓SUB_CBOR = 소켓SUB_CBOR
		테스트_항목.질의_인수.소켓SUB_MsgPack = 소켓SUB_MsgPack
		테스트_항목.질의_인수.소켓SUB_JSON = 소켓SUB_JSON
		테스트_항목.질의_인수.t = t

		테스트_항목.실행_함수(t, 테스트_항목.질의_인수)

		// 이전 데이터가 계속 오는 경우가 발견됨.
		lib.F대기(lib.P3초)
	}
}

func f실시간_서비스_수신(질의_인수 *sRT질의_인수) []interface{} {
	var 에러 error
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M에러: &에러})

	t := 질의_인수.t
	lib.F테스트_에러없음(t, f접속_확인())
	변환_형식 := lib.F임의_변환형식()
	소켓_질의 := lib.New소켓_질의(lib.P주소_NH_TR, 변환_형식, lib.P30초)

	// 실시간 서비스 등록
	응답 := 소켓_질의.S질의(질의_인수.질의값).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())

	// 데이터 수신
	수신값_모음 := make([]interface{}, 0)
	소켓SUB_모음 := []mangos.Socket{질의_인수.소켓SUB_CBOR, 질의_인수.소켓SUB_MsgPack, 질의_인수.소켓SUB_JSON}

	for _, 소켓SUB := range 소켓SUB_모음 {
		수신_메시지 := lib.New소켓_메시지_내용없음().S소켓_수신(소켓SUB, lib.P무기한)

		switch 수신_메시지.G자료형_문자열(0) {
		case nh호가_잔량:
			값 := new(lib.NH호가_잔량)
			에러 = 수신_메시지.G값(0, &값)
			lib.F테스트_에러없음(t, 에러)
			수신값_모음 = append(수신값_모음, 값)
		case nh시간외_호가_잔량:
			값 := new(lib.NH시간외_호가잔량)
			에러 = 수신_메시지.G값(0, &값)
			lib.F테스트_에러없음(t, 에러)
			수신값_모음 = append(수신값_모음, 값)
		case nh예상_호가_잔량:
			값 := new(lib.NH예상_호가잔량)
			에러 = 수신_메시지.G값(0, &값)
			lib.F테스트_에러없음(t, 에러)
			수신값_모음 = append(수신값_모음, 값)
		case nh체결:
			값 := new(lib.NH체결)
			에러 = 수신_메시지.G값(0, &값)
			lib.F테스트_에러없음(t, 에러)
			수신값_모음 = append(수신값_모음, 값)
		case nhETF_NAV:
			값 := new(lib.NH_ETF_NAV)
			에러 = 수신_메시지.G값(0, &값)
			lib.F테스트_에러없음(t, 에러)
			수신값_모음 = append(수신값_모음, 값)
		case nh업종_지수:
			값 := new(lib.NH업종지수)
			에러 = 수신_메시지.G값(0, &값)
			lib.F테스트_에러없음(t, 에러)
			수신값_모음 = append(수신값_모음, 값)
		default:
			lib.New에러("예상하지 못한 자료형. %v", 수신_메시지.G자료형_문자열(0))
			t.FailNow()
		}
	}

	질의_인수.질의값.TR구분 = lib.TR실시간_정보_해지
	응답 = 소켓_질의.S질의(질의_인수.질의값).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())

	lib.F대기(lib.P5초) // 실시간 정보 해지가 실제 효력을 발생하도록 시간을 줌.

	// 소켓에 누적된 메시지 제거
	for _, 소켓SUB := range 소켓SUB_모음 {
		소켓SUB.SetOption(mangos.OptionRecvDeadline, lib.P100밀리초)

		for {
			if _, 에러 := 소켓SUB.Recv(); 에러 != nil {
				소켓SUB.SetOption(mangos.OptionRecvDeadline, lib.P무기한)
				break
			}
		}
	}

	return 수신값_모음
}

func h1_k3_수신_테스트(t *testing.T, 질의_인수 *sRT질의_인수) {
	데이터_모음 := f실시간_서비스_수신(질의_인수)
	lib.F테스트_다름(t, 데이터_모음, nil)
	lib.F테스트_같음(t, len(데이터_모음), 3)

	for _, 데이터 := range 데이터_모음 {
		s, ok := 데이터.(*lib.NH호가_잔량)
		lib.F테스트_참임(t, ok, 데이터)
		lib.F테스트_같음(t, s.M종목코드, lib.F2인터페이스_모음(질의_인수.질의값.M종목코드_모음)...)
		lib.F테스트_참임(t, s.M시각.After(time.Now().Add(-30*lib.P1초)))
		lib.F테스트_참임(t, s.M시각.Before(time.Now().Add(20*lib.P1초)))
		lib.F테스트_같음(t, len(s.M매도_잔량_모음), len(s.M매도_호가_모음))
		lib.F테스트_같음(t, len(s.M매수_잔량_모음), len(s.M매수_호가_모음))

		for i, 매도_잔량 := range s.M매도_잔량_모음 {
			lib.F테스트_참임(t, 매도_잔량 > 0, i, 매도_잔량)
		}

		for i, 매도_호가 := range s.M매도_호가_모음 {
			if i == 0 {
				lib.F테스트_참임(t, 매도_호가 > 0, i, 매도_호가)
				continue
			}

			lib.F테스트_참임(t,
				매도_호가 > s.M매도_호가_모음[i-1],
				s.M종목코드, s.M시각,
				s.M매도_호가_모음, s.M매도_잔량_모음,
				s.M매수_호가_모음, s.M매수_잔량_모음)
		}

		for i, 매수_잔량 := range s.M매수_잔량_모음 {
			lib.F테스트_참임(t, 매수_잔량 > 0, i, 매수_잔량)
		}

		for i, 매수_호가 := range s.M매수_호가_모음 {
			if i == 0 {
				lib.F테스트_참임(t, 매수_호가 > 0, i, 매수_호가)
				continue
			}

			lib.F테스트_참임(t,
				s.M매수_호가_모음[i-1] > 매수_호가,
				s.M종목코드, s.M시각,
				s.M매도_호가_모음, s.M매도_잔량_모음,
				s.M매수_호가_모음, s.M매수_잔량_모음)
		}

		lib.F테스트_참임(t, s.M누적_거래량 >= 0, s.M누적_거래량)
	}
}

func h2_k4_수신_테스트(t *testing.T, 질의_인수 *sRT질의_인수) {
	데이터_모음 := f실시간_서비스_수신(질의_인수)
	lib.F테스트_참임(t, 데이터_모음 != nil)
	lib.F테스트_참임(t, len(데이터_모음) == 3, 데이터_모음)

	종목코드_모음 := make([]string, 0)

	for _, 종목코드 := range 질의_인수.질의값.M종목코드_모음 {
		종목코드_모음 = append(종목코드_모음, 종목코드)
	}

	for _, 데이터 := range 데이터_모음 {
		s, ok := 데이터.(*lib.NH시간외_호가잔량)
		lib.F테스트_참임(t, ok, s)
		lib.F테스트_같음(t, s.M종목코드, lib.F2인터페이스_모음(종목코드_모음)...)
		lib.F테스트_참임(t, s.M시각.After(time.Now().Add(-30*time.Second)))
		lib.F테스트_참임(t, s.M시각.Before(time.Now().Add(20*time.Second)))
		lib.F테스트_참임(t, s.M총_매도호가_잔량 >= 0, s.M총_매도호가_잔량)
		lib.F테스트_참임(t, s.M총_매수호가_잔량 >= 0, s.M총_매수호가_잔량)
	}
}

func h3_k5_수신_테스트(t *testing.T, 질의_인수 *sRT질의_인수) {
	데이터_모음 := f실시간_서비스_수신(질의_인수)
	lib.F테스트_참임(t, 데이터_모음 != nil)
	lib.F테스트_참임(t, len(데이터_모음) == 3, 데이터_모음)

	종목코드_모음 := make([]string, 0)

	for _, 종목코드 := range 질의_인수.질의값.M종목코드_모음 {
		종목코드_모음 = append(종목코드_모음, 종목코드)
	}

	for _, 데이터 := range 데이터_모음 {
		lib.F테스트_참임(t, 데이터 != nil)

		s, ok := 데이터.(*lib.NH예상_호가잔량)
		lib.F테스트_참임(t, ok, s)
		lib.F테스트_같음(t, s.M종목코드, lib.F2인터페이스_모음(종목코드_모음)...)
		lib.F테스트_참임(t, s.M시각.After(time.Now().Add(-30*time.Second)))
		lib.F테스트_참임(t, s.M시각.Before(time.Now().Add(20*time.Second)))
		lib.F테스트_같음(t, int(s.M동시호가_구분), 0, 1, 2, 3, 4, 5, 6)
		lib.F테스트_참임(t, s.M예상_체결가 >= 0, s.M예상_체결가)
		lib.F테스트_참임(t, s.M예상_체결가 <= s.M매도_호가 || s.M매도_호가 == 0)
		lib.F테스트_참임(t, s.M예상_체결가 >= s.M매수_호가 || s.M예상_체결가 == 0)
		lib.F테스트_참임(t, f올바른_등락부호(s.M예상_등락부호), s.M예상_등락부호)
		lib.F테스트_참임(t, s.M예상_등락폭 >= 0, s.M예상_등락폭) // 절대값 ?
		lib.F테스트_참임(t, s.M예상_등락율 >= 0, s.M예상_등락율) // 절대값 ?
		lib.F테스트_참임(t, s.M예상_등락율 <= 30, s.M예상_등락율)

		if s.M예상_등락폭 != 0 &&
			s.M예상_체결가 != 0 &&
			s.M예상_등락율 != 0 {
			예상_등락율_근사값 := lib.F2절대값(
				float64(s.M예상_등락폭) / float64(s.M예상_체결가) * 100)

			lib.F테스트_참임(t,
				lib.F오차(lib.F2절대값(s.M예상_등락율), 예상_등락율_근사값) < 5 ||
					lib.F오차율(lib.F2절대값(s.M예상_등락율), 예상_등락율_근사값) < 10)
		}

		lib.F테스트_참임(t, s.M예상_체결수량 >= 0, s.M예상_체결수량)

		if s.M매도_호가잔량 > 0 {
			lib.F테스트_참임(t, s.M매도_호가 > 0)
		} else {
			lib.F테스트_참임(t, s.M매도_호가 == 0)
		}

		if s.M매수_호가잔량 > 0 {
			lib.F테스트_참임(t, s.M매수_호가 > 0)
		} else {
			lib.F테스트_참임(t, s.M매수_호가 == 0)
		}

		if s.M예상_체결가 > 0 && s.M매도_호가 > 0 {
			lib.F테스트_참임(t, lib.F오차율(s.M예상_체결가, s.M매도_호가) < 60)
		}

		if s.M예상_체결가 > 0 && s.M매수_호가 > 0 {
			lib.F테스트_참임(t, lib.F오차율(s.M예상_체결가, s.M매수_호가) < 60)
		}
	}
}

func j8_k8_수신_테스트(t *testing.T, 질의_인수 *sRT질의_인수) {
	데이터_모음 := f실시간_서비스_수신(질의_인수)
	lib.F테스트_참임(t, 데이터_모음 != nil)
	lib.F테스트_참임(t, len(데이터_모음) == 3, 데이터_모음)

	종목코드_모음 := make([]string, 0)

	for _, 종목코드 := range 질의_인수.질의값.M종목코드_모음 {
		종목코드_모음 = append(종목코드_모음, 종목코드)
	}

	for _, 데이터 := range 데이터_모음 {
		s, ok := 데이터.(*lib.NH체결)

		lib.F테스트_참임(t, ok, s)
		lib.F테스트_같음(t, s.M종목코드, lib.F2인터페이스_모음(종목코드_모음)...)
		lib.F테스트_참임(t, s.M시각.After(time.Now().Add(-30*time.Second)))
		lib.F테스트_참임(t, s.M시각.Before(time.Now().Add(20*time.Second)))
		lib.F테스트_참임(t, f올바른_등락부호(s.M등락부호), s.M등락부호)
		lib.F테스트_참임(t, s.M등락폭 >= 0) // 절대값?
		lib.F테스트_참임(t, s.M현재가 >= 0, s.M현재가)

		등락율_근사값 := float64(s.M등락폭) / float64(s.M현재가) * 100
		오차율 := lib.F오차율(s.M등락율, 등락율_근사값)
		오차 := lib.F오차(s.M등락율, 등락율_근사값)
		lib.F테스트_참임(t, 등락율_근사값 <= 30, 등락율_근사값, s.M등락폭, s.M현재가)
		lib.F테스트_참임(t, lib.F2절대값(s.M등락율) <= 30, s.M등락율)
		lib.F테스트_참임(t, 오차율 < 10 || 오차 < 5,
			s.M종목코드, s.M시각, s.M현재가, s.M등락폭, s.M등락율, 등락율_근사값)
		lib.F테스트_참임(t, s.M고가 >= s.M현재가 || s.M고가 == 0, s.M현재가, s.M고가)
		lib.F테스트_참임(t, s.M저가 <= s.M현재가 || s.M현재가 == 0, s.M현재가, s.M저가)

		if s.M고가 != 0 {
			lib.F테스트_참임(t, s.M고가 >= s.M현재가)
			lib.F테스트_참임(t, s.M고가 >= s.M저가)
			//lib.F테스트_참임(t, s.M고가 >= s.M매도_호가, s.M고가, s.M매도_호가)
			lib.F테스트_참임(t, s.M고가 >= s.M매수_호가)
			lib.F테스트_참임(t, s.M고가 >= s.M시가)
			lib.F테스트_참임(t, s.M고가 >= s.M가중_평균_가격)
		}

		if s.M매도_호가 != 0 {
			lib.F테스트_참임(t, s.M매도_호가 >= s.M현재가)
			lib.F테스트_참임(t, s.M매도_호가 >= s.M저가)
			lib.F테스트_참임(t, s.M매도_호가 >= s.M매수_호가)
		}

		if s.M현재가 != 0 {
			lib.F테스트_참임(t, s.M현재가 >= s.M저가)
			lib.F테스트_참임(t, s.M현재가 >= s.M매수_호가)
		}

		if s.M매수_호가 != 0 {
			lib.F테스트_참임(t, s.M매수_호가 >= s.M저가)
		}

		if s.M가중_평균_가격 != 0 {
			lib.F테스트_참임(t, s.M가중_평균_가격 >= s.M저가)
		}

		lib.F테스트_참임(t, s.M누적_거래량 >= 0)
		lib.F테스트_참임(t, s.M전일대비_거래량_비율 >= 0)
		lib.F테스트_참임(t, s.M변동_거래량 > 0, s.M변동_거래량)

		//거래_대금_근사값 := s.M현재가 * s.M변동_거래량
		거래_대금_근사값 := s.M현재가 * s.M누적_거래량 / 1000000
		오차율 = lib.F오차율(s.M거래_대금_100만, 거래_대금_근사값)
		오차 = lib.F오차(s.M거래_대금_100만, 거래_대금_근사값)
		lib.F테스트_참임(t, 오차율 < 10 || 오차 < 3,
			오차율, 오차, s.M거래_대금_100만, 거래_대금_근사값,
			s.M현재가, s.M누적_거래량, s.M변동_거래량)
		lib.F테스트_같음(t, s.M시장구분, lib.P시장구분_코스피, lib.P시장구분_코스닥)
	}
}

func j1_j0_수신_테스트(t *testing.T, 질의_인수 *sRT질의_인수) {
	데이터_모음 := f실시간_서비스_수신(질의_인수)
	lib.F테스트_참임(t, 데이터_모음 != nil)
	lib.F테스트_참임(t, len(데이터_모음) == 3, 데이터_모음)

	종목코드_모음 := make([]string, 0)

	for _, 종목코드 := range 질의_인수.질의값.M종목코드_모음 {
		종목코드_모음 = append(종목코드_모음, 종목코드)
	}

	for _, 데이터 := range 데이터_모음 {
		lib.F테스트_참임(t, 데이터 != nil)

		s, ok := 데이터.(*lib.NH_ETF_NAV)
		lib.F테스트_참임(t, ok, s)
		lib.F테스트_같음(t, s.M종목코드, lib.F2인터페이스_모음(종목코드_모음)...)
		lib.F테스트_참임(t, s.M시각.After(time.Now().Add(-30*time.Second)))
		lib.F테스트_참임(t, s.M시각.Before(time.Now().Add(20*time.Second)))
		lib.F테스트_참임(t, f올바른_등락부호(s.M등락부호), s.M등락부호)
		lib.F테스트_참임(t, s.M등락폭 >= 0) // 절대값?

		if s.M고가_NAV != 0 {
			lib.F테스트_참임(t, s.M고가_NAV >= s.M현재가_NAV)
			lib.F테스트_참임(t, s.M고가_NAV >= s.M시가_NAV)
			lib.F테스트_참임(t, s.M고가_NAV >= s.M저가_NAV)
		}

		if s.M현재가_NAV != 0 {
			lib.F테스트_참임(t, s.M현재가_NAV >= s.M저가_NAV)
		}

		if s.M시가_NAV != 0 {
			lib.F테스트_참임(t, s.M시가_NAV >= s.M저가_NAV)
		}

		lib.F테스트_참임(t, f올바른_등락부호(s.M추적오차_부호), s.M추적오차_부호)
		lib.F테스트_참임(t, f올바른_등락부호(s.M괴리율_부호), s.M괴리율_부호)
		lib.F테스트_참임(t, s.M추적오차 >= 0)
		lib.F테스트_참임(t, s.M괴리율 >= 0)

		lib.F메모("괴리율 테스트를 정확하게 하려면 ETF_NAV, 현재가를 동시에 얻어야 한다.")
	}
}

func u1_k1_수신_테스트(t *testing.T, 질의_인수 *sRT질의_인수) {
	데이터_모음 := f실시간_서비스_수신(질의_인수)
	lib.F테스트_참임(t, 데이터_모음 != nil)
	lib.F테스트_참임(t, len(데이터_모음) == 3, 데이터_모음)

	종목코드_모음 := make([]string, 0)

	for _, 종목코드 := range 질의_인수.질의값.M종목코드_모음 {
		종목코드_모음 = append(종목코드_모음, 종목코드)
	}

	for _, 데이터 := range 데이터_모음 {
		s, ok := 데이터.(*lib.NH업종지수)
		lib.F테스트_참임(t, ok, s)
		lib.F테스트_같음(t, s.M업종코드, lib.F2인터페이스_모음(종목코드_모음)...)
		lib.F테스트_참임(t, s.M시각.After(time.Now().Add(-30*time.Second)))
		lib.F테스트_참임(t, s.M시각.Before(time.Now().Add(20*time.Second)))
		lib.F테스트_참임(t, f올바른_등락부호(s.M등락부호), s.M등락부호)
		lib.F테스트_참임(t, s.M등락폭 >= 0) // 절대값?
		lib.F테스트_참임(t, s.M거래량 >= 0)
		lib.F테스트_참임(t, s.M거래_대금 >= 0)

		if s.M최고값 != 0 {
			lib.F테스트_참임(t, s.M최고값 >= s.M현재값)
			lib.F테스트_참임(t, s.M최고값 >= s.M개장값)
			lib.F테스트_참임(t, s.M최고값 >= s.M최저값)
		}

		최근_개장일, 에러 := lib.F한국증시_최근_개장일()
		lib.F테스트_에러없음(t, 에러)
		lib.F테스트_참임(t, s.M최고값_시각.After(최근_개장일))
		lib.F테스트_참임(t, s.M최고값_시각.Before(time.Now().Add(20*time.Second)))

		if s.M현재값 != 0 {
			lib.F테스트_참임(t, s.M현재값 >= s.M최저값)
		}

		if s.M개장값 != 0 {
			lib.F테스트_참임(t, s.M개장값 >= s.M최저값)
		}

		lib.F테스트_참임(t, s.M최저값_시각.After(최근_개장일))
		lib.F테스트_참임(t, s.M최저값_시각.Before(time.Now().Add(20*time.Second)))
		lib.F테스트_참임(t, s.M거래_비중 >= 0, s.M거래_비중)
		lib.F테스트_참임(t, s.M거래_비중 <= 100, s.M거래_비중)
		lib.F테스트_참임(t, s.M지수_등락율 >= 0, s.M지수_등락율) // 절대값??
	}
}

func TestTR소켓_접속됨(t *testing.T) {
	lib.F메모("단일 테스트에서는 잘 되지만, 일괄 테스트에서는 에러가 발생함. 원인 불명")

	t.SkipNow()

	lib.F대기(lib.P3초)

	질의값 := new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR접속됨

	응답 := lib.New소켓_질의(lib.P주소_NH_TR, lib.F임의_변환형식(), lib.P30초).S질의(질의값).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())
	lib.F테스트_같음(t, 응답.G길이(), 1)

	var 접속_여부 bool
	lib.F테스트_에러없음(t, 응답.G값(0, &접속_여부))
	lib.F테스트_같음(t, 접속_여부, f접속됨())
}

func TestTR소켓_접속(t *testing.T) {
	lib.F대기(lib.P3초)
	변환_형식 := lib.F임의_변환형식()
	소켓_질의 := lib.New소켓_질의(lib.P주소_NH_TR, 변환_형식, lib.P30초)

	lib.F체크포인트()

	for f접속됨() {
		lib.F체크포인트("우선 접속 해제 해야 함.")

		질의값 := new(lib.S질의값_단순TR)
		질의값.TR구분 = lib.TR접속_해제

		응답 := 소켓_질의.S질의(질의값).G응답()

		lib.F테스트_에러없음(t, 응답.G에러())
		lib.F테스트_거짓임(t, f접속됨())
		lib.F대기(lib.P1초)

		lib.F체크포인트("접속 해제.")
	}

	lib.F체크포인트("접속 시작")

	질의값 := new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR접속

	응답 := 소켓_질의.S질의(질의값).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())
	lib.F테스트_같음(t, 응답.G길이(), 1)

	로그인_정보 := new(lib.NH로그인_정보)

	lib.F체크포인트(응답.G자료형_문자열(0))

	lib.F테스트_에러없음(t, 응답.G값(0, 로그인_정보))
	lib.F테스트_다름(t, 로그인_정보.M접속_ID, "")
	lib.F테스트_다름(t, 로그인_정보.M접속_서버, "")

	접속_시각 := 로그인_정보.M접속_시각
	p지금 := time.Now()
	p10초전 := p지금.Add(-1 * lib.P10초)
	p10초후 := p지금.Add(lib.P10초)

	lib.F테스트_참임(t, 접속_시각.After(p10초전), p지금, 접속_시각)
	lib.F테스트_참임(t, 접속_시각.Before(p10초후), p지금, 접속_시각)
	lib.F테스트_참임(t, f접속됨())
}

func TestTR소켓_접속_해지(t *testing.T) {
	lib.F메모("접속_해지 테스트 중 종종 에러가 발생함.")
	t.SkipNow()

	lib.F대기(lib.P3초)
	lib.F테스트_에러없음(t, f접속_확인())
	lib.F테스트_참임(t, f접속됨())

	변환_형식 := lib.F임의_변환형식()
	소켓_질의 := lib.New소켓_질의(lib.P주소_NH_TR, 변환_형식, lib.P30초)

	질의값 := new(lib.S질의값_단순TR)
	질의값.TR구분 = lib.TR접속_해제

	응답 := 소켓_질의.S질의(질의값).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())
	lib.F테스트_거짓임(t, f접속됨())
}

func TestTR소켓_실시간_서비스_모두_해지(t *testing.T) {
	lib.F메모("실시간_서비스_모두_해지 테스트 중 종종 에러가 발생하지만, 자주 사용하는 기능이 아니어서 테스트 건너뜀.")
	t.SkipNow()

	lib.F대기(lib.P3초)
	변환_형식 := lib.F임의_변환형식()
	소켓_질의 := lib.New소켓_질의(lib.P주소_NH_TR, 변환_형식, lib.P30초)

	if f접속됨() {
		질의값_접속해제 := new(lib.S질의값_단순TR)
		질의값_접속해제.TR구분 = lib.TR접속_해제
		응답 := 소켓_질의.S질의(질의값_접속해제).G응답()
		lib.F테스트_에러없음(t, 응답.G에러())
	}

	질의값_접속 := new(lib.S질의값_단순TR)
	질의값_접속.TR구분 = lib.TR접속
	응답 := 소켓_질의.S질의(질의값_접속).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())

	// 실시간 서비스 등록
	종목모음_코스피, 에러 := lib.F종목모음_코스피()
	lib.F테스트_에러없음(t, 에러)

	질의값_구독 := new(lib.S질의값_복수종목)
	질의값_구독.TR구분 = lib.TR실시간_정보_구독
	질의값_구독.TR코드 = lib.NH_RT코스피_체결
	질의값_구독.M종목코드_모음 = lib.F종목코드_추출(종목모음_코스피, 20)
	응답 = 소켓_질의.S질의(질의값_구독).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())
	lib.F대기(lib.P300밀리초)

	// 실시간 서비스 모두 해지
	질의값_일괄해지 := new(lib.S질의값_단순TR)
	질의값_일괄해지.TR구분 = lib.TR실시간_정보_일괄_해지
	응답 = 소켓_질의.S질의(질의값_일괄해지).G응답()
	lib.F테스트_에러없음(t, 응답.G에러())

	var 응답_구분 lib.TR응답_구분
	lib.F테스트_에러없음(t, 응답.G값(0, &응답_구분))
	lib.F테스트_같음(t, 응답_구분, lib.TR응답_완료)
}

