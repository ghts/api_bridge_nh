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

	"math"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
	//"encoding/gob"
	//"os"
)

func TestC1151_ETF_현재가(t *testing.T) {
	lib.F대기(lib.P3초)
	f_ETF_현재가_조회_테스트(t, lib.F임의_종목_ETF())
}

//func TestCh_TR_ETF_현재가_전종목(t *testing.T) {
//	f접속_확인()
//
//	파일_닫아야_함 := false
//	파일명 := "test_completed.txt"
//	파일 := new(os.File)
//	맵 := make(map[string]lib.S비어있음)
//	에러 := lib.F_nil에러()
//
//	defer func(){
//		if 파일_닫아야_함 {
//			파일.Close()
//		}
//	}()
//
//	_, 에러 = os.Stat(파일명)
//
//	if 에러 != nil && os.IsNotExist(에러) {
//		파일, 에러 = os.Create(파일명)
//		lib.F테스트_에러없음(t, 에러)
//		파일_닫아야_함 = true
//
//		인코더 := gob.NewEncoder(파일)
//		에러 = 인코더.Encode(맵)
//		lib.F테스트_에러없음(t, 에러)
//
//		for 파일.Sync() != nil {
//			time.Sleep(500 * time.Millisecond)
//		}
//
//		파일.Close()
//		파일_닫아야_함 = false
//	} else {
//		lib.F테스트_에러없음(t, 에러)
//
//		파일, 에러 = os.Open(파일명)
//		lib.F테스트_에러없음(t, 에러)
//		파일_닫아야_함 = true
//
//		디코더 := gob.NewDecoder(파일)
//		에러 = 디코더.Decode(&맵)
//
//		if 에러 != nil && 에러.Error() == "EOF" {
//			맵 = make(map[string]lib.S비어있음)
//		} else {
//			lib.F테스트_에러없음(t, 에러)
//		}
//
//		파일.Close()
//		파일_닫아야_함 = false
//	}
//
//	종목_모음 := lib.F샘플_종목_모음_ETF()
//
//	for _, 종목 := range 종목_모음 {
//		_, 이미_테스트_됨 := 맵[종목.G코드()]
//
//		if 이미_테스트_됨 {
//			continue
//		}
//
//		f_ETF_현재가_조회_테스트(t, 종목)
//
//		맵[종목.G코드()] = lib.S비어있음{}
//
//
//		파일, 에러 := os.Create(파일명)
//		lib.F테스트_에러없음(t, 에러)
//		파일_닫아야_함 = true
//
//		인코더 := gob.NewEncoder(파일)
//		에러 = 인코더.Encode(맵)
//		lib.F테스트_에러없음(t, 에러)
//
//		for 파일.Sync() != nil {
//			time.Sleep(500 * time.Millisecond)
//		}
//
//		파일.Close()
//		파일_닫아야_함 = false
//
//		// 서버 TR 제한을 피하기 위함.
//		time.Sleep(time.Second)
//	}
//}

func f_ETF_현재가_조회_테스트(t *testing.T, 종목 *lib.S종목) {
	lib.F테스트_에러없음(t, f접속_확인())

	기본_정보 := new(lib.NH_ETF_현재가_조회_기본_정보)
	변동_정보_모음 := make([]*lib.NH_ETF_현재가_조회_변동_거래량_정보, 0)
	동시호가_정보 := new(lib.NH_ETF_현재가_조회_동시호가_정보)
	ETF_정보 := new(lib.NH_ETF_현재가_조회_ETF정보)
	지수_정보 := new(lib.NH_ETF_현재가_조회_지수_정보)

	질의값 := new(lib.S질의값_단일종목)
	질의값.TR구분 = lib.TR조회
	질의값.TR코드 = lib.NH_TR_ETF_현재가_조회
	질의값.M종목코드 = 종목.G코드()

	// 소켓으로 질의를 수신한다는 것을 감안한 변환.
	소켓_메시지, 에러 := lib.New소켓_메시지(lib.CBOR, 질의값)
	lib.F테스트_에러없음(t, 에러)

	질의 := lib.New채널_질의(ch조회, lib.P30초, 1).S질의(소켓_메시지)

	수신완료_기본정보, 수신완료_변동정보, 수신완료_동시호가정보 := false, false, false
	수신완료_ETF정보, 수신완료_지수정보 := false, false

	for {
		if 수신완료_기본정보 && 수신완료_변동정보 && 수신완료_동시호가정보 &&
			수신완료_ETF정보 && 수신완료_지수정보 {
			break
		}

		응답 := 질의.G응답()
		lib.F테스트_에러없음(t, 응답.G에러())

		TR응답_구분, ok := 응답.G값(0).(lib.TR응답_구분)
		lib.F테스트_참임(t, ok, 응답)

		switch TR응답_구분 {
		case lib.TR응답_데이터:
			lib.F테스트_같음(t, 응답.G길이(), 2)

			변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
			lib.F테스트_참임(t, ok)
			자료형 := 변환값.G자료형_문자열()

			switch {
			case strings.Contains(자료형, "ETF_현재가_조회_기본"):
				수신완료_기본정보 = true
				lib.F테스트_에러없음(t, 변환값.G값(기본_정보))
			case strings.Contains(자료형, "ETF_현재가_조회_변동"):
				수신완료_변동정보 = true
				lib.F테스트_에러없음(t, 변환값.G값(&변동_정보_모음))
			case strings.Contains(자료형, "ETF_현재가_조회_동시호가"):
				수신완료_동시호가정보 = true
				lib.F테스트_에러없음(t, 변환값.G값(동시호가_정보))
			case strings.Contains(자료형, "ETF_현재가_조회_ETF"):
				수신완료_ETF정보 = true
				lib.F테스트_에러없음(t, 변환값.G값(ETF_정보))
			case strings.Contains(자료형, "ETF_현재가_조회_지수"):
				수신완료_지수정보 = true
				lib.F테스트_에러없음(t, 변환값.G값(지수_정보))
			default:
				lib.F문자열_출력("예상치 못한 자료형. %v", 자료형)
				t.FailNow()
			}
		case lib.TR응답_메시지:
			lib.F테스트_같음(t, 응답.G길이(), 3)

			코드, ok := 응답.G값(1).(string) // 코드
			lib.F테스트_참임(t, ok)

			메시지, ok := 응답.G값(2).(string)
			lib.F테스트_참임(t, ok)

			lib.F문자열_출력("메시지 : %v, %v", 코드, 메시지)

			//lib.F테스트_참임(t, strings.Contains(메시지, "TRU") ||
			//		strings.Contains(메시지, "조회완료") ||
			//		strings.Contains(메시지, "정보가 없습니다"), 메시지)
		case lib.TR응답_완료:
			lib.F테스트_같음(t, 응답.G길이(), 2)

			변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
			lib.F테스트_참임(t, ok, 응답.G값(1))

			if 변환값.G자료형_문자열() != "<nil>" {
				lib.New에러("예상하지 못한 자료형. %v", 변환값.G자료형_문자열())
			}
		default:
			lib.F문자열_출력("\n*** 예상치 못한 TR 응답 구분 : %v ***", TR응답_구분)
			lib.F변수값_확인(TR응답_구분)
			lib.F변수값_확인(응답)
			t.FailNow()
		}
	}

	//lib.F문자열_출력("*** 종목코드 %v ***", 종목.G코드())
	//lib.F문자열_출력("*** 시각 %v ***", 기본_정보.M시각)

	f_ETF_현재가_조회_기본_정보_테스트(t, 기본_정보, 종목.G코드())
	f_ETF_현재가_조회_변동_거래_정보_테스트(t, 기본_정보, 변동_정보_모음)
	f_ETF_현재가_조회_동시호가_정보_테스트(t, 기본_정보, 동시호가_정보)
	f_ETF_현재가_조회_ETF자료_테스트(t, 기본_정보, ETF_정보)

	lib.F문자열_출력("ETF 기반 지수 테스트 건너뜀.")
	//f_ETF_현재가_조회_지수_정보_테스트(t, 기본_정보, ETF_정보, 지수_정보)
}

func f_ETF_현재가_조회_기본_정보_테스트(t *testing.T,
	s *lib.NH_ETF_현재가_조회_기본_정보, 종목코드 string) {
	종목, 에러 := lib.F종목by코드(종목코드)
	lib.F에러2패닉(에러)

	지금 := time.Now()
	금일_0시 := time.Date(지금.Year(), 지금.Month(), 지금.Day(), 0, 0, 0, 0, 지금.Location())
	삼분후 := 지금.Add(3 * time.Minute)
	개장일_0시, 에러 := lib.F한국증시_최근_개장일()
	lib.F에러2패닉(에러)

	//일주전 := 지금.Add(-7 * 24 * time.Hour)
	//일년전 := time.Date(지금.Year()-1, 지금.Month(), 지금.Day(), 0, 0, 0, 0, 지금.Location())
	이백년전 := 지금.Add(-200 * 365 * 24 * time.Hour)

	lib.F테스트_참임(t, s != nil, "기본 자료를 수신하지 못함.")
	lib.F테스트_같음(t, s.M종목코드, 종목.G코드())
	lib.F테스트_참임(t, utf8.ValidString(s.M종목명))

	종목명_수신값 := strings.Replace(s.M종목명, " ", "", -1)
	종목명_비교값 := strings.Replace(종목.G이름(), " ", "", -1)
	lib.F테스트_참임(t,
		strings.Contains(종목명_수신값, 종목명_비교값) ||
			strings.Contains(종목명_비교값, 종목명_수신값),
		"'"+종목명_수신값+"'", "'"+종목명_비교값+"'")
	lib.F테스트_참임(t, s.M등락율 >= 0) // 절대값임.

	f테스트_등락부호(t, s.M등락부호, s.M현재가, s.M전일_종가, s.M상한가, s.M하한가)
	lib.F테스트_같음(t, s.M현재가,
		s.M전일_종가+f등락부호2정수(s.M등락부호)*s.M등락폭)

	if s.M현재가 != 0 && s.M등락폭 != 0 && s.M등락율 != 0 {
		등락율_근사값 := math.Abs(float64(s.M등락폭)) / float64(s.M현재가) * 100
		lib.F테스트_참임(t, lib.F오차율(s.M등락율, 등락율_근사값) < 10 ||
			lib.F오차(s.M등락율, 등락율_근사값) < 0.1,
			s.M등락율, 등락율_근사값, s.M등락폭, s.M현재가)
	}

	lib.F테스트_참임(t, s.M거래량 >= 0)
	lib.F테스트_참임(t, s.M전일대비_거래량_비율 >= 0)

	if s.M유동_주식수_1000주 != 0 {
		유동주_회전율_근사값 := float64(s.M거래량) /
			float64(s.M유동_주식수_1000주*1000) * 100
		유동주_회전율_근사값 = math.Trunc(유동주_회전율_근사값*100) / 100
		lib.F테스트_참임(t,
			lib.F오차(s.M유동주_회전율, 유동주_회전율_근사값) < 1 ||
				lib.F오차율(s.M유동주_회전율, 유동주_회전율_근사값) < 10,
			s.M유동주_회전율, 유동주_회전율_근사값)
	}

	if s.M거래대금_100만 != 0 && s.M거래량 != 0 && s.M현재가 != 0 {
		거래대금_근사값 := s.M거래량 * s.M현재가 / 1000000
		lib.F테스트_참임(t, lib.F오차율(s.M거래대금_100만, 거래대금_근사값) < 10 ||
			lib.F오차(s.M거래대금_100만, 거래대금_근사값) <= 1,
			s.M거래대금_100만, 거래대금_근사값)
	}

	if s.M거래량 > 0 {
		// 거래량이 0이면 저가, 고가 모두 0임.
		lib.F테스트_참임(t, s.M상한가 >= s.M고가, s.M상한가)
		lib.F테스트_참임(t, s.M고가 >= s.M시가, s.M시가, s.M고가)
		lib.F테스트_참임(t, s.M고가 >= s.M저가, s.M저가, s.M고가)
		lib.F테스트_참임(t, s.M저가 <= s.M시가, s.M저가, s.M시가)
		lib.F테스트_참임(t, s.M저가 >= s.M하한가, s.M하한가, s.M저가)
		lib.F테스트_참임(t, s.M현재가 <= s.M고가, s.M현재가, s.M고가)
		lib.F테스트_참임(t, s.M현재가 >= s.M저가, s.M저가, s.M현재가)
	}

	lib.F테스트_참임(t, s.M하한가 >= 0, s.M하한가)
	lib.F테스트_참임(t, s.M52주_고가 >= s.M20일_고가)
	lib.F테스트_참임(t, s.M20일_고가 >= s.M5일_고가)
	lib.F테스트_참임(t, s.M5일_고가 >= s.M5일_저가)
	lib.F테스트_참임(t, s.M20일_저가 >= s.M52주_저가)
	lib.F테스트_참임(t, s.M5일_저가 >= s.M20일_저가)
	lib.F테스트_참임(t, s.M시각.After(개장일_0시), s.M시각)
	lib.F테스트_참임(t, s.M시각.Before(개장일_0시.Add(18*time.Hour)), s.M시각)

	if lib.F한국증시_정규시장_거래시간임() {
		lib.F테스트_참임(t, s.M시각.Before(삼분후))
	} else {
		lib.F테스트_참임(t, s.M시각.Hour() == 15 ||
			s.M시각.Hour() == 16, s.M시각)
	}

	f테스트_등락부호(t, s.M시가대비_등락부호, s.M현재가, s.M시가, s.M상한가, s.M하한가)

	if s.M시가 != 0 {
		lib.F테스트_참임(t, s.M시가+s.M시가대비_등락폭 == s.M현재가,
			s.M시가, s.M현재가, s.M시가대비_등락폭) // 시가대비_등락폭 자체에 부호가 반영되어 있음.
	}

	for i, 매도_잔량 := range s.M매도_잔량_모음 {
		lib.F테스트_참임(t, 매도_잔량 >= 0)

		if 매도_잔량 == 0 {
			continue
		}

		// 매도_잔량 > 0
		매도_호가 := s.M매도_호가_모음[i]
		lib.F테스트_참임(t, 매도_호가 <= s.M상한가)
		lib.F테스트_참임(t, 매도_호가 >= s.M하한가)

		if i == 0 {
			lib.F테스트_참임(t, 매도_호가 >= s.M현재가)
		} else {
			lib.F테스트_참임(t, 매도_호가 > s.M매도_호가_모음[i-1])
		}
	}

	for i, 매수_잔량 := range s.M매수_잔량_모음 {
		lib.F테스트_참임(t, 매수_잔량 >= 0)

		if 매수_잔량 == 0 {
			continue
		}

		// 매수_잔량 > 0
		매수_호가 := s.M매수_호가_모음[i]
		lib.F테스트_참임(t, 매수_호가 <= s.M상한가)
		lib.F테스트_참임(t, 매수_호가 >= s.M하한가)

		if i == 0 {
			lib.F테스트_참임(t, 매수_호가 >= s.M현재가)
		} else {
			lib.F테스트_참임(t, 매수_호가 < s.M매수_호가_모음[i-1],
				매수_호가, s.M매수_호가_모음[i-1], i)
		}
	}

	lib.F테스트_참임(t, s.M시간외_매도_잔량 >= 0)
	lib.F테스트_참임(t, s.M시간외_매수_잔량 >= 0)
	lib.F테스트_참임(t, s.M피봇_2차_저항 >= s.M피봇_1차_저항,
		s.M피봇_2차_저항, s.M피봇_1차_저항)
	lib.F테스트_참임(t, s.M피봇_1차_저항 >= s.M피봇_가격)
	lib.F테스트_참임(t, s.M피봇_가격 >= s.M피봇_1차_지지)
	lib.F테스트_참임(t, s.M피봇_1차_지지 >= s.M피봇_2차_지지)
	lib.F테스트_참임(t, utf8.ValidString(s.M시장_구분))
	lib.F테스트_같음(t, s.M시장_구분, "코스피", "코스닥")
	lib.F테스트_참임(t, utf8.ValidString(s.M업종명))
	lib.F테스트_참임(t, utf8.ValidString(s.M자본금_규모))
	lib.F테스트_참임(t, s.M자본금_규모 == "", s.M자본금_규모)
	lib.F테스트_참임(t, utf8.ValidString(s.M결산월))
	lib.F테스트_참임(t, strings.Contains(s.M결산월, "월 결산"))

	for _, 추가_정보 := range s.M추가_정보_모음 {
		lib.F테스트_참임(t, utf8.ValidString(추가_정보))
	}

	lib.F테스트_같음(t, s.M서킷_브레이커_구분, "", "CB발동", "CB해제", "장종료")
	lib.F테스트_참임(t, s.M액면가 == 0, s.M액면가) // ETF는 액면가가 없는 건가?

	switch {
	case strings.Contains(종목.G이름(), "레버리지"):
		lib.F테스트_참임(t, lib.F오차율(s.M상한가, float64(s.M전일_종가)*1.6) < 5,
			s.M상한가, int64(float64(s.M전일_종가)*1.6), 종목.G이름())
		lib.F테스트_참임(t, lib.F오차율(s.M하한가, float64(s.M전일_종가)*0.4) < 5,
			s.M하한가, int64(float64(s.M전일_종가)*0.4), 종목.G이름())
	default:
		lib.F테스트_참임(t, lib.F오차율(s.M상한가, float64(s.M전일_종가)*1.3) < 5,
			s.M상한가, int64(float64(s.M전일_종가)*1.3), 종목.G이름())
		lib.F테스트_참임(t, lib.F오차율(s.M하한가, float64(s.M전일_종가)*0.7) < 5,
			s.M하한가, int64(float64(s.M전일_종가)*0.7), 종목.G이름())
	}

	lib.F테스트_참임(t, s.M대용가 < s.M전일_종가)
	lib.F테스트_참임(t, s.M대용가 > int64(float64(s.M전일_종가)*0.5))
	lib.F테스트_참임(t, s.M공모가 >= 0)
	lib.F테스트_참임(t, s.M유동_주식수_1000주 >= 0)

	시가총액_근사값 := s.M현재가 * s.M상장_주식수 / 100000000
	lib.F테스트_참임(t, lib.F오차율(s.M시가_총액_억, 시가총액_근사값) < 10 ||
		lib.F오차(s.M시가_총액_억, 시가총액_근사값) <= 1,
		s.M시가_총액_억, 시가총액_근사값, s.M현재가, s.M상장_주식수)
	lib.F테스트_참임(t, s.M거래원_정보_수신_시각.Before(삼분후))

	// 	ETF 거래원 정보 수신 시각에서 종종 에러가 발생함.
	lib.F테스트_참임(t, s.M거래원_정보_수신_시각.IsZero() ||
		s.M거래원_정보_수신_시각.After(개장일_0시.Add(8*time.Hour+59*time.Minute)),
		종목.G코드(), 종목.G이름(), s.M시각, s.M거래원_정보_수신_시각)

	lib.F테스트_참임(t, s.M거래원_정보_수신_시각.Before(개장일_0시.Add(18*time.Hour)),
		s.M거래원_정보_수신_시각)

	매도_거래량_합계 := int64(0)
	for i, 매도_거래량 := range s.M매도_거래량_모음 {
		lib.F테스트_참임(t, 매도_거래량 >= 0, i, 매도_거래량)

		if 매도_거래량 == 0 {
			continue
		}

		매도_거래량_합계 += 매도_거래량
		매도_거래원 := s.M매도_거래원_모음[i]
		lib.F테스트_참임(t, len(매도_거래원) > 0)
		lib.F테스트_참임(t, utf8.ValidString(매도_거래원), 매도_거래원)
	}

	매수_거래량_합계 := int64(0)
	for i, 매수_거래량 := range s.M매수_거래량_모음 {
		lib.F테스트_참임(t, 매수_거래량 >= 0, i, 매수_거래량)

		if 매수_거래량 == 0 {
			continue
		}

		매수_거래량_합계 += 매수_거래량
		매수_거래원 := s.M매수_거래원_모음[i]
		lib.F테스트_참임(t, len(매수_거래원) > 0)
		lib.F테스트_참임(t, utf8.ValidString(매수_거래원), 매수_거래원)
	}

	lib.F테스트_참임(t, s.M외국인_매도_거래량 >= 0)
	lib.F테스트_참임(t, s.M외국인_매수_거래량 >= 0)
	lib.F테스트_참임(t, s.M외국인_시간.After(개장일_0시), s.M외국인_시간)
	lib.F테스트_참임(t, s.M외국인_시간.Before(개장일_0시.Add(23*time.Hour)), s.M외국인_시간)
	lib.F테스트_참임(t, s.M외국인_지분율 >= 0 && s.M외국인_지분율 <= 100)
	lib.F테스트_참임(t, s.M신용잔고율 >= 0)
	lib.F테스트_참임(t, s.M신용잔고율 <= 100)
	//lib.F테스트_참임(t, s.M유상_기준일.After(이백년전) || s.M유상_기준일.IsZero(), s.M유상_기준일)
	//lib.F테스트_참임(t, s.M무상_기준일.After(이백년전) || s.M무상_기준일.IsZero(), s.M무상_기준일)
	lib.F테스트_참임(t, s.M유상_배정_비율 >= 0)
	lib.F테스트_참임(t, s.M유상_배정_비율 <= 100)
	lib.F테스트_참임(t, s.M무상_배정_비율 >= 0)
	lib.F테스트_참임(t, s.M무상_배정_비율 <= 100)
	lib.F테스트_참임(t, s.M상장일.After(이백년전))
	lib.F테스트_참임(t, s.M상장일.Before(금일_0시))
	lib.F테스트_참임(t, s.M상장_주식수 > 0)
	lib.F테스트_참임(t, s.M전체_거래원_매도_합계 >= 매도_거래량_합계)
	lib.F테스트_참임(t, s.M전체_거래원_매수_합계 >= 매수_거래량_합계)
}

func f_ETF_현재가_조회_변동_거래_정보_테스트(t *testing.T,
	기본_정보 *lib.NH_ETF_현재가_조회_기본_정보,
	변동_정보_모음 []*lib.NH_ETF_현재가_조회_변동_거래량_정보) {
	if len(변동_정보_모음) == 0 {
		lib.F문자열_출력("변동 자료를 수신하지 못함. %v", 기본_정보.M종목코드)
		return
	}

	거래량_잔량 := 기본_정보.M거래량
	지금 := time.Now()
	삼분후 := 지금.Add(3 * time.Minute)
	개장일_0시, 에러 := lib.F한국증시_최근_개장일()
	lib.F에러2패닉(에러)

	for i, s := range 변동_정보_모음 {
		lib.F테스트_참임(t, s.M시각.Before(삼분후))
		lib.F테스트_참임(t, s.M시각.After(개장일_0시.Add(9*time.Hour)))
		lib.F테스트_참임(t, s.M시각.Before(개장일_0시.Add(18*time.Hour)), s.M시각)

		if lib.F한국증시_정규시장_거래시간임() {
			lib.F테스트_참임(t, s.M시각.Before(삼분후))
		} else {
			lib.F테스트_참임(t, s.M시각.Hour() >= 9, s.M시각)
			lib.F테스트_참임(t, s.M시각.Hour() <= 16, s.M시각)
		}

		lib.F테스트_참임(t, s.M현재가 <= 기본_정보.M상한가)
		lib.F테스트_참임(t, s.M현재가 >= 기본_정보.M하한가)
		lib.F테스트_참임(t, f올바른_등락부호(s.M등락부호))
		lib.F테스트_같음(t, f등락부호2정수(s.M등락부호)*s.M등락폭,
			s.M현재가-기본_정보.M전일_종가)

		if lib.F한국증시_정규시장_거래시간임() {
			lib.F테스트_참임(t, s.M매도_호가 >= s.M현재가 ||
				s.M매도_호가 == 0, s.M매도_호가, s.M현재가)
			lib.F테스트_참임(t, s.M매수_호가 <= s.M현재가 ||
				s.M매수_호가 == 0, s.M매수_호가, s.M현재가)
		} else {
			if s.M매도_호가 > 0 {
				lib.F테스트_참임(t, s.M매도_호가 <= 기본_정보.M상한가)
				lib.F테스트_참임(t, s.M매도_호가 >= 기본_정보.M하한가)
			}

			if s.M매수_호가 > 0 {
				lib.F테스트_참임(t, s.M매수_호가 <= 기본_정보.M상한가)
				lib.F테스트_참임(t, s.M매수_호가 >= 기본_정보.M하한가)
			}
		}

		// 걸러낸 자료로 인한 오차 수정.
		if i == 0 && s.M거래량 != 거래량_잔량 {
			거래량_잔량 = s.M거래량
		}

		lib.F테스트_같음(t, s.M거래량, 거래량_잔량)
		거래량_잔량 -= s.M변동_거래량
	}
}

func f_ETF_현재가_조회_동시호가_정보_테스트(t *testing.T,
	기본_정보 *lib.NH_ETF_현재가_조회_기본_정보,
	s *lib.NH_ETF_현재가_조회_동시호가_정보) {
	lib.F테스트_참임(t, s != nil)
	lib.F테스트_같음(t, int(s.M동시호가_구분), 0, 1, 2, 3, 4, 5, 6)

	if s.M동시호가_구분 == 0 {
		return
	}

	if s.M예상_체결가 == 0 {
		lib.F변수값_확인(기본_정보.M종목코드, 기본_정보.M시각, s.M동시호가_구분)

		return
	}

	lib.F테스트_참임(t, f올바른_등락부호(s.M예상_체결_부호), s.M예상_체결_부호)
	lib.F테스트_참임(t, s.M예상_체결가 <= 기본_정보.M상한가)
	lib.F테스트_참임(t, s.M예상_체결가 >= 기본_정보.M하한가)
	lib.F테스트_참임(t, lib.F오차율(s.M예상_체결가, 기본_정보.M현재가) < 10)
	lib.F테스트_같음(t, f등락부호2정수(s.M예상_체결_부호)*s.M예상_등락폭,
		s.M예상_체결가-기본_정보.M전일_종가)

	if s.M예상_등락폭 != 0 && s.M예상_등락율 != 0 {
		예상_등락율_계산값 := lib.F2절대값(s.M예상_등락폭) /
			float64(s.M예상_체결가) * 100
		lib.F테스트_참임(t, lib.F오차율(s.M예상_등락율, 예상_등락율_계산값) < 10)
	}

	lib.F테스트_참임(t, s.M예상_체결_수량 >= 0)
}

func f_ETF_현재가_조회_ETF자료_테스트(t *testing.T,
	기본_정보 *lib.NH_ETF_현재가_조회_기본_정보,
	s *lib.NH_ETF_현재가_조회_ETF정보) {
	lib.F테스트_참임(t, s != nil)
	lib.F테스트_같음(t, s.ETF구분, P코스피, P코스닥)
	lib.F테스트_참임(t, s.NAV > 0, s.NAV)
	lib.F테스트_참임(t, f올바른_등락부호(s.NAV등락부호))

	NAV_근사값 := s.M전일NAV + float64(f등락부호2정수(s.NAV등락부호))*s.NAV등락폭
	lib.F테스트_참임(t, lib.F오차율(s.NAV, NAV_근사값) < 3)
	lib.F테스트_참임(t, math.Abs(float64(s.M괴리율)) < 10, s.M괴리율)

	괴리율_계산값 := lib.F2절대값(float64(기본_정보.M현재가)-s.NAV) / float64(s.NAV) * 100
	오차율 := lib.F오차율(s.M괴리율, 괴리율_계산값)
	오차 := math.Abs(s.M괴리율 - 괴리율_계산값)

	lib.F메모("NAV 산출시간을 알게 되면 괴리율 검사 오차를 줄일 것.")
	lib.F테스트_참임(t, 오차율 < 10 || 오차 < 3,
		s.NAV, 기본_정보.M현재가, 괴리율_계산값, s.M괴리율)

	lib.F테스트_참임(t, f올바른_등락부호(s.M괴리율_부호), s.M괴리율_부호)
	//lib.F변수값_확인(s.M괴리율_부호, s.M괴리율)
	//f테스트_등락율(t, s.M괴리율_부호, s.M괴리율)
	lib.F테스트_참임(t, s.M괴리율 >= 0, s.M괴리율)

	lib.F테스트_참임(t, s.M설정단위_당_현금_배당액 >= 0)
	lib.F테스트_참임(t, s.M구성_종목수 > 0)
	lib.F테스트_참임(t, s.M순자산_총액_억 > 0)

	lib.F메모("NAV 산출시간을 알게 되면 괴리율 검사 오차를 줄일 것.")
	//	순자산_총액_억 := float64(s.M순자산_총액_억)
	//	시가_총액_억 := float64(기본_정보.M시가_총액_억)
	//	시가_총액_억_계산값_1 := 순자산_총액_억 * (1 +
	//		float64(f등락부호2정수(s.M괴리율_부호)) * math.Abs(s.M괴리율/100))
	//	시가_총액_억_계산값_2 := 순자산_총액_억 * (1 -
	//		float64(f등락부호2정수(s.M괴리율_부호)) * math.Abs(s.M괴리율/100))
	//	순자산_총액_억_계산값_1 := 시가_총액_억 * (1 +
	//		float64(f등락부호2정수(s.M괴리율_부호)) * math.Abs(s.M괴리율/100))
	//	순자산_총액_억_계산값_2 := 시가_총액_억 * (1 -
	//		float64(f등락부호2정수(s.M괴리율_부호)) * math.Abs(s.M괴리율/100))
	//	lib.F테스트_참임(t,
	//		lib.F오차율(순자산_총액_억, 순자산_총액_억_계산값_1) < 3 ||
	//		lib.F오차(순자산_총액_억, 순자산_총액_억_계산값_1) < 2 ||
	//		lib.F오차율(순자산_총액_억, 순자산_총액_억_계산값_2) < 3 ||
	//		lib.F오차(순자산_총액_억, 순자산_총액_억_계산값_2) < 2 ||
	//		lib.F오차율(시가_총액_억, 시가_총액_억_계산값_1) < 3 ||
	//		lib.F오차(시가_총액_억, 시가_총액_억_계산값_1) < 2 ||
	//		lib.F오차율(시가_총액_억, 시가_총액_억_계산값_2) < 3 ||
	//		lib.F오차(시가_총액_억, 시가_총액_억_계산값_2) < 2,
	//		순자산_총액_억, 순자산_총액_억_계산값_1, 순자산_총액_억_계산값_2,
	//		시가_총액_억, 시가_총액_억_계산값_1, 시가_총액_억_계산값_2, s.M괴리율)

	// 추적 오차율 : 지수와 NAV의 수익률 차이.
	lib.F테스트_참임(t, math.Abs(float64(s.M추적_오차율)) >= 0 &&
		math.Abs(float64(s.M추적_오차율)) < 20, s.M추적_오차율)

	//lib.F문자열_출력("ETF는 지정가 주문이라서 수량이 0이 될 수도 있다고 함. 수량이 0인 주문이 가능한 건가?")
	// 수량이 0인 주문은 없는 것이나 마찬가지이니 자체적으로 걸러내자.
	for i, 매도_잔량 := range s.LP_매도_잔량_모음 {
		if i == 0 {
			continue
		}

		if 매도_잔량 > 0 {
			lib.F테스트_참임(t, s.LP_매도_잔량_모음[i-1] > 0,
				i, s.LP_매도_잔량_모음[i-1], 매도_잔량)
		}
	}

	for i, 매수_잔량 := range s.LP_매수_잔량_모음 {
		if i == 0 {
			continue
		}

		if 매수_잔량 > 0 {
			lib.F테스트_참임(t, s.LP_매수_잔량_모음[i-1] > 0,
				i, s.LP_매수_잔량_모음[i-1], 매수_잔량)
		}
	}

	lib.F테스트_같음(t, s.ETF_복제_방법_구분_코드, P실물복제, P합성복제)
	lib.F테스트_같음(t, s.ETF_상품_유형_코드, P일반형, P파생형, "")
}

func f_ETF_현재가_조회_지수_정보_테스트(t *testing.T,
	기본_정보 *lib.NH_ETF_현재가_조회_기본_정보,
	ETF_정보 *lib.NH_ETF_현재가_조회_ETF정보,
	s *lib.NH_ETF_현재가_조회_지수_정보) {

	// MUG 2101 화면 참조
	코스피_업종_맵 := make(map[string]string)
	코스피_업종_맵["00"] = "KRX 100"
	코스피_업종_맵["01"] = "코스피지수"
	코스피_업종_맵["02"] = "대형주"
	코스피_업종_맵["03"] = "중형주"
	코스피_업종_맵["04"] = "소형주"
	코스피_업종_맵["05"] = "음식료품"
	코스피_업종_맵["06"] = "섬유,의복"
	코스피_업종_맵["07"] = "종이,목재"
	코스피_업종_맵["08"] = "화학"
	코스피_업종_맵["09"] = "의약품"
	코스피_업종_맵["10"] = "비금속광물"
	코스피_업종_맵["11"] = "철광,금속"
	코스피_업종_맵["12"] = "기계"
	코스피_업종_맵["13"] = "전기,전자"
	코스피_업종_맵["14"] = "의료정밀"
	코스피_업종_맵["15"] = "운수장비"
	코스피_업종_맵["16"] = "유통업"
	코스피_업종_맵["17"] = "전기가스업"
	코스피_업종_맵["18"] = "건설업"
	코스피_업종_맵["19"] = "운수창고"
	코스피_업종_맵["20"] = "통신업"
	코스피_업종_맵["21"] = "금융업"
	코스피_업종_맵["22"] = "은행"
	코스피_업종_맵["24"] = "증권"
	코스피_업종_맵["25"] = "보험"
	코스피_업종_맵["26"] = "서비스업"
	코스피_업종_맵["27"] = "제조업"
	코스피_업종_맵["28"] = "코스피 200"
	코스피_업종_맵["29"] = "코스피 100"
	코스피_업종_맵["30"] = "코스피 50"
	코스피_업종_맵["32"] = "코스피 배당"
	코스피_업종_맵["37"] = "KP200 산업재"
	코스피_업종_맵["38"] = "KP200 건강관리"
	코스피_업종_맵["39"] = "KP200 건설기계"
	코스피_업종_맵["40"] = "KP200 조선운송"
	코스피_업종_맵["41"] = "KP200 철강소재"
	코스피_업종_맵["42"] = "KP200 에너지화학"
	코스피_업종_맵["43"] = "KP200 정보기술"
	코스피_업종_맵["44"] = "KP200 금융"
	코스피_업종_맵["45"] = "KP200 생활소비재"
	코스피_업종_맵["46"] = "KP200 경기소비재"
	코스피_업종_맵["47"] = "동일가중 KP200"
	코스피_업종_맵["48"] = "동일가중 KP100"
	코스피_업종_맵["49"] = "동일가중 KP50"

	코스닥_업종_맵 := make(map[string]string)
	코스닥_업종_맵["01"] = "코스닥지수"
	코스닥_업종_맵["03"] = "기타서비스"
	코스닥_업종_맵["04"] = "코스닥 IT"
	코스닥_업종_맵["06"] = "제조"
	코스닥_업종_맵["07"] = "건설"
	코스닥_업종_맵["08"] = "유통"
	코스닥_업종_맵["10"] = "운송"
	코스닥_업종_맵["11"] = "금융"
	코스닥_업종_맵["12"] = "통신방송서비스"
	코스닥_업종_맵["13"] = "IT S/W & SVC"
	코스닥_업종_맵["14"] = "IT H/W"
	코스닥_업종_맵["15"] = "음식료,담배"
	코스닥_업종_맵["16"] = "섬유,의류"
	코스닥_업종_맵["17"] = "종이,목재"
	코스닥_업종_맵["18"] = "출판,매체복제"
	코스닥_업종_맵["19"] = "화학"
	코스닥_업종_맵["20"] = "제약"
	코스닥_업종_맵["21"] = "비금속"
	코스닥_업종_맵["22"] = "금속"
	코스닥_업종_맵["23"] = "기계,장비"
	코스닥_업종_맵["24"] = "일반전기전자"
	코스닥_업종_맵["25"] = "의료,정밀기기"
	코스닥_업종_맵["26"] = "운송장비,부품"
	코스닥_업종_맵["27"] = "기타 제조"
	코스닥_업종_맵["28"] = "통신서비스"
	코스닥_업종_맵["29"] = "방송서비스"
	코스닥_업종_맵["30"] = "인터넷"
	코스닥_업종_맵["31"] = "디지털컨텐츠"
	코스닥_업종_맵["32"] = "소프트웨어"
	코스닥_업종_맵["33"] = "컴퓨터서비스"
	코스닥_업종_맵["34"] = "통신장비"
	코스닥_업종_맵["35"] = "정보기기"
	코스닥_업종_맵["36"] = "반도체"
	코스닥_업종_맵["37"] = "IT부품"
	코스닥_업종_맵["38"] = "KOSDAQ 100"
	코스닥_업종_맵["39"] = "KOSDAQ MID 300"
	코스닥_업종_맵["40"] = "KOSDAQ SMALL"
	코스닥_업종_맵["43"] = "코스닥 스타"
	코스닥_업종_맵["44"] = "오락,문화"
	코스닥_업종_맵["45"] = "프리미어"
	코스닥_업종_맵["46"] = "우량기업부"
	코스닥_업종_맵["47"] = "벤처기업부"
	코스닥_업종_맵["48"] = "중견기업부"
	코스닥_업종_맵["49"] = "기술성장기업부"

	// 게시판 답변에 따름.
	KRX지수_코드_맵 := make(map[string]string)
	KRX지수_코드_맵["0001"] = "KRX 100"
	KRX지수_코드_맵["0002"] = "동일가중 KRX"
	KRX지수_코드_맵["0101"] = "KRX Autos, KRX 자동차"
	KRX지수_코드_맵["0102"] = "KRX Semicon, KRX 반도체"
	KRX지수_코드_맵["0103"] = "KRX Health Care, KRX 건강"
	KRX지수_코드_맵["0104"] = "KRX Banks, KRX 은행"
	KRX지수_코드_맵["0105"] = "KRX IT"
	KRX지수_코드_맵["0106"] = "KRX Energy&Chem, KRX 에너지 화학"
	KRX지수_코드_맵["0107"] = "KRX Steels, KRX 철강"
	KRX지수_코드_맵["0108"] = "KRX Consumer, KRX 소비재"
	KRX지수_코드_맵["0109"] = "KRX Media&Tele, KRX 미디어통신"
	KRX지수_코드_맵["0110"] = "KRX Construct, KRX 건설"
	KRX지수_코드_맵["0111"] = "KRX Financials, KRX 금융"
	KRX지수_코드_맵["0112"] = "KRX Securities, KRX 증권"
	KRX지수_코드_맵["0113"] = "KRX Shipbuild, KRX 조선"
	KRX지수_코드_맵["0114"] = "KRX Insurance, KRX 보험"
	KRX지수_코드_맵["0115"] = "KRX Transport, KRX 운송"
	KRX지수_코드_맵["0116"] = "KRX Retail, KRX 소매"
	KRX지수_코드_맵["0117"] = "KRX Leisure, KRX 레저"
	KRX지수_코드_맵["0201"] = "KRX SRI, 사회책임투자지수"
	KRX지수_코드_맵["0202"] = "KRX SRI ECO, 환경책임투자지수"
	KRX지수_코드_맵["0203"] = "KRX Green, 녹색산업지수"
	KRX지수_코드_맵["0204"] = "KRX SRI Governance, 지배구조책임투자지수"

	// 게시판 답변 내용에 따름. 채권지수 세부코드는 뭐지?
	채권지수_코드_맵 := make(map[string]lib.NH_ETF_현재가_조회_지수_정보)
	채권지수_코드_맵["KISR01-1"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "KISR01",
		M채권지수_세부_코드: "1",
		M지수_이름:      "KTB Index(총수익)"}
	채권지수_코드_맵["KISR01-2"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "KISR01",
		M채권지수_세부_코드: "2",
		M지수_이름:      "KTB Index(순가격)"}
	채권지수_코드_맵["KISR01-3"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "KISR01",
		M채권지수_세부_코드: "3",
		M지수_이름:      "KTB Index(시장가격)"}
	채권지수_코드_맵["MKFR01-1"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "MKFR01",
		M채권지수_세부_코드: "1",
		M지수_이름:      "MKF 국고채지수(총수익), MKF TB Index(총수익)"}
	채권지수_코드_맵["MKFR01-2"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "MKFR01",
		M채권지수_세부_코드: "2",
		M지수_이름:      "MKF 국고채지수(순가격), MKF TB Index(순가격)"}
	채권지수_코드_맵["MKFR01-3"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "MKFR01",
		M채권지수_세부_코드: "3",
		M지수_이름:      "MKF 국고채지수(시장가격), MKF TB Index(시장가격)"}
	채권지수_코드_맵["MSBI01-1"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "MSBI01",
		M채권지수_세부_코드: "1",
		M지수_이름:      "MK 통안채지수(총수익), MK MSB Index(총수익)"}
	채권지수_코드_맵["MSBI01-2"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "MSBI01",
		M채권지수_세부_코드: "2",
		M지수_이름:      "MK 통안채지수(순가격), MK MSB Index(순가격)"}
	채권지수_코드_맵["MSBI01-3"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "MSBI01",
		M채권지수_세부_코드: "3",
		M지수_이름:      "MK 통안채지수(시장가격), MK MSB Index(시장가격)"}
	채권지수_코드_맵["MSBI03-1"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "MSBI03",
		M채권지수_세부_코드: "1",
		M지수_이름:      "MK 머니마켓 지수(총수익), MK 머니마켓(총수익)"}
	채권지수_코드_맵["MSBI03-2"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "MSBI03",
		M채권지수_세부_코드: "2",
		M지수_이름:      "MK 머니마켓 지수(순가격), MK 머니마켓(순가격)"}
	채권지수_코드_맵["MSBI03-3"] = lib.NH_ETF_현재가_조회_지수_정보{
		M채권지수_코드:    "MSBI03",
		M채권지수_세부_코드: "3",
		M지수_이름:      "MK 머니마켓 지수(시장가격), MK 머니마켓(시장가격)"}

	// MUG 3900 화면 참조
	해외지수_맵 := make(map[string][2]string)
	해외지수_맵["EPEWT"] = [2]string{"MSCI 대만", "MSCI Taiwan"}
	해외지수_맵["EPSPX"] = [2]string{"S&P 500", "S&P 500"}
	해외지수_맵["EPNDX"] = [2]string{"나스닥 100", "NASDAQ 100"}
	해외지수_맵["EWBVSP"] = [2]string{"IBOVESPA", "IBOVESPA"}
	해외지수_맵["EWBVSP"] = [2]string{"STI", "STI"}
	해외지수_맵["EWFTSE"] = [2]string{"FTSE 100", "FTSE 100"}
	해외지수_맵["EPEWJ"] = [2]string{"MSCI 일본", "MSCI JAPAN"}
	해외지수_맵["EWHSCE"] = [2]string{"HSCEI", "HSCEI (Hang Seng China Enterprise Index)"}
	해외지수_맵["000300"] = [2]string{"CSI 300", "CSI300 (China Securities Index 300)"}
	해외지수_맵["AEXX"] = [2]string{"네덜란드 지수", "AEX (Amsterdam Exchange Index)"}
	해외지수_맵["AOI"] = [2]string{"호주 지수", "AOI (Australian All Ordinaries Index)"}
	해외지수_맵["ATX"] = [2]string{"오스트리아 지수", "ATX (Austrian Traded Index)"}
	해외지수_맵["BFX"] = [2]string{"벨기에 지수", "BEL-20"}
	해외지수_맵["BSESN"] = [2]string{"인도 지수", "BSE 30, SENSEX (Bombay Stock Exchange Sensitive Index)"}
	해외지수_맵["BVSP"] = [2]string{"브라질 지수", "Bovespa Index"}
	해외지수_맵["CASE30"] = [2]string{"이집트 지수", "Cairo & Alexandria Stock Exchange Index"}
	해외지수_맵["DAX"] = [2]string{"독일 지수", "Frankfurt DAX (Deutscher Aktienindex : German stock index)"}
	해외지수_맵["ESFG"] = [2]string{"E-mini S&P 지수선물", "E-mini S&P 500 Future Globex Index"}
	해외지수_맵["GSPTSE"] = [2]string{"캐나다 지수", "TSX Composite Index, Toronto Stock Exchange, CAD"}
	해외지수_맵["HSCE"] = [2]string{"홍콩H 지수", "HangSeng China Enterprises Index"}
	해외지수_맵["HSI"] = [2]string{"항셍 지수", "HangSeng Index"}
	해외지수_맵["IPSA"] = [2]string{"칠레 지수", "IPSA (Índice de Precio Selectivo de Acciones)"}
	해외지수_맵["JALSH"] = [2]string{"남아프리카공화국 지수", "JALSH(JSE Africa All Share Index)"}
	해외지수_맵["JCI"] = [2]string{"인도네시아 지수", "Jakarta Composite Index"}
	해외지수_맵["KFX"] = [2]string{"덴마크 지수", "KFX"}
	해외지수_맵["KLSE"] = [2]string{"말레이시아 지수", "Bursa Malaysia KLCI(Kuala Lumpur Composite Index)"}
	해외지수_맵["KOSPI"] = [2]string{"한국 코스피 지수", "KOSPI(Korea Composite Stock Price Index)"}
	해외지수_맵["MERV"] = [2]string{"아르헨티나 지수", "MERVAL(MERcado de VALores : Stock Market) Index"}
	해외지수_맵["MTMS"] = [2]string{"러시아 지수", "Moscow Times (Russia)"}
	해외지수_맵["MXX"] = [2]string{"멕시코 지수", "Mexico IPC(Indice de Precios y Cotizaciones) Index"}
	해외지수_맵["N225"] = [2]string{"일본 니케이225 지수", "Nikkei 225 Index"}
	해외지수_맵["NQIG"] = [2]string{"E-mini NASDAQ 100 선물지수", "E-mini NASDAQ-100 Futures"}
	해외지수_맵["NZ50"] = [2]string{"뉴질랜드 지수", "NZX 50 Index"}
	해외지수_맵["PARI"] = [2]string{"프랑스 지수", "Paris CAC 40"}
	해외지수_맵["PSI"] = [2]string{"필리핀 지수", "PSE(Phillipine Stock Exchange) Composite Index"}
	해외지수_맵["RDXUSD"] = [2]string{"러시아 블루칩 지수", "RDX(Russian Depositary Index) in USD"}
	해외지수_맵["SET"] = [2]string{"태국 지수", "Thailand SET(Stock Exchang of Thailand) Index"}
	해외지수_맵["SHANG"] = [2]string{"중국상해종합지수", "Shanghai Composite Index"}
	해외지수_맵["SPFG"] = [2]string{"S&P 500 지수선물", "S&P 500 Index Futures"}
	해외지수_맵["SPFR"] = [2]string{"S&P 500 지수선물", "S&P 500 Index Futures"}
	해외지수_맵["SSEA"] = [2]string{"상해A 지수", "Shanghai A Share"}
	해외지수_맵["SSEB"] = [2]string{"상해B 지수", "Shanghai B Share"}
	해외지수_맵["SSMI"] = [2]string{"스위스 지수", "Swiss Market Index"}
	해외지수_맵["STI"] = [2]string{"싱가포르 지수", "Straits Times Index"}
	해외지수_맵["SX5E"] = [2]string{"유로 STOXX50 지수", "Euro Stoxx 50"}
	해외지수_맵["SXAXP"] = [2]string{"스웨덴 지수", "Sweden Stockholm General Index"}
	해외지수_맵["SZSA"] = [2]string{"심천A 지수", "Shenzhen A Share"}
	해외지수_맵["SZSB"] = [2]string{"심천B 지수", "Shenzhen B Share"}
	해외지수_맵["TWI"] = [2]string{"대만 지수", "TSE(Taiwan Stock Exchange) Weighted Index"}
	해외지수_맵["VEB"] = [2]string{"베네주엘라 지수", "Venezuelan Bolivar VE"}
	해외지수_맵["VNI"] = [2]string{"베트남 지수", "Vietnam Ho Chi Minh Stock Index"}

	// Mug 2108 화면. 게시판 답변 내용에 따름.
	기타_업종_코드_맵 := make(map[string]string)
	기타_업종_코드_맵["001"] = "코스피200 선물지수, F-KOSPI200"
	기타_업종_코드_맵["002"] = "코스피200 선물인버스지수"
	기타_업종_코드_맵["003"] = "미국달러 선물지수"
	기타_업종_코드_맵["004"] = "미국달러 선물인버스지수"
	기타_업종_코드_맵["005"] = "코스피200 레버리지지수"
	기타_업종_코드_맵["006"] = "코스피200 커버드콜지수, KOSPI 200 Cov Call"
	기타_업종_코드_맵["007"] = "코스피200 프로텍티브풋지수"
	기타_업종_코드_맵["008"] = "국채 선물지수(3년), 3년 국채선물 지수"
	기타_업종_코드_맵["009"] = "국채 선물인버스지수(3년), 3년 국채선물 인버스 지수"
	기타_업종_코드_맵["010"] = "국채 선물지수(10년), 10년 국채선물 지수"
	기타_업종_코드_맵["011"] = "국채 선물인버스지수(10년), 10년 국채선물 인버스 지수"
	기타_업종_코드_맵["012"] = "MSCI Korea Index"
	기타_업종_코드_맵["013"] = "주식골드지수"
	기타_업종_코드_맵["014"] = "코스피200 리스크컨트롤 6%지수"
	기타_업종_코드_맵["015"] = "코스피200 리스크컨트롤 8%지수"
	기타_업종_코드_맵["016"] = "코스피200 리스크컨트롤 10%지수"
	기타_업종_코드_맵["017"] = "코스피200 리스크컨트롤 12%지수"
	기타_업종_코드_맵["018"] = "주식미국채DAE지수"
	기타_업종_코드_맵["019"] = "V-KOSPI200지수 (코스피200 변동성 지수)"
	기타_업종_코드_맵["020"] = "주식국채혼합형(주식형)지수"
	기타_업종_코드_맵["021"] = "주식국채혼합형(채권형)지수"
	기타_업종_코드_맵["022"] = "코스피200 DAE 지수"
	기타_업종_코드_맵["023"] = "코스피200 고배당 지수, KP200 고배당 지수"
	기타_업종_코드_맵["024"] = "코스피200 저변동성 지수, KP200 저변동성 지수"
	기타_업종_코드_맵["025"] = "미국달러 선물 레버리지 지수"
	기타_업종_코드_맵["026"] = "코스피 고배당50 지수"
	기타_업종_코드_맵["027"] = "코스피 배당성장50 지수"
	기타_업종_코드_맵["028"] = "코스피 우선주 지수"
	기타_업종_코드_맵["029"] = "KRX 고배당 50"
	기타_업종_코드_맵["030"] = "코스피 200 선물 플러스지수"
	기타_업종_코드_맵["031"] = "K200 USD 선물 바이셀지수"
	기타_업종_코드_맵["032"] = "USD K200 선물 바이셀지수"
	기타_업종_코드_맵["033"] = "코스피 선물매수 콜매도지수"
	기타_업종_코드_맵["034"] = "코스피 선물매도 풋매도지수"
	기타_업종_코드_맵["035"] = "WISE 삼성그룹 인덱스"
	기타_업종_코드_맵["036"] = "WISE 로우볼 지수"
	기타_업종_코드_맵["037"] = "WISE 셀렉트 배당 지수"
	기타_업종_코드_맵["038"] = "WISE K150 Quant 인덱스"
	기타_업종_코드_맵["039"] = "Big Vol 지수"
	기타_업종_코드_맵["040"] = "BNP High Dividend Yield Europe Equity Long TR EUR"
	기타_업종_코드_맵["041"] = "코스피 200 내재가치 지수"
	기타_업종_코드_맵["042"] = "스마트 리밸런싱 250/3 A 지수"
	기타_업종_코드_맵["043"] = "WISE 롱숏 K150 로우볼지수"
	기타_업종_코드_맵["044"] = "WISE BIG5 동일가중 TR 지수"
	기타_업종_코드_맵["045"] = "WISE 스마트베타 Quality 지수"
	기타_업종_코드_맵["046"] = "WISE 스마트베타 Momentum 지수"
	기타_업종_코드_맵["047"] = "WISE 스마트베타 Value 지수"
	기타_업종_코드_맵["048"] = "WISE Monthly Best 11 지수"
	기타_업종_코드_맵["049"] = "코스피 200 선물 인버스-2X지수"
	기타_업종_코드_맵["050"] = "코스피 200 선물 인버스-3X지수"
	기타_업종_코드_맵["051"] = "미국달러선물 인버스-2X지수"
	기타_업종_코드_맵["052"] = "미국달러선물 인버스-3X지수"
	기타_업종_코드_맵["053"] = "KTOP30"
	기타_업종_코드_맵["054"] = "코스닥150"
	기타_업종_코드_맵["055"] = "코스닥200 중소형주 지수"
	기타_업종_코드_맵["056"] = "코스닥200 필수소비재 지수"
	기타_업종_코드_맵["057"] = "FnGuide 에너지 TOP5 지수"
	기타_업종_코드_맵["058"] = "FnGuide 필수소비재 TOP5 지수"
	기타_업종_코드_맵["059"] = "FnGuide 조선 TOP5 지수"
	기타_업종_코드_맵["060"] = "FnGuide 소프트웨어 TOP5 IT 지수"
	기타_업종_코드_맵["061"] = "FnGuide IT 하드웨어 TOP5 지수"
	기타_업종_코드_맵["062"] = "FnGuide 운송 TOP5 지수"
	기타_업종_코드_맵["063"] = "FnGuide 자동차 TOP5 지수"
	기타_업종_코드_맵["064"] = "FnGuide 의료 TOP5 지수"
	기타_업종_코드_맵["065"] = "FnGuide 화학 TOP5 지수"
	기타_업종_코드_맵["066"] = "FnGuide 바이오 TOP5 지수"
	기타_업종_코드_맵["067"] = "FnGuide 제약 TOP5 지수"
	기타_업종_코드_맵["068"] = "FnGuide 건설 TOP5 지수"
	기타_업종_코드_맵["069"] = "코스닥150 동일가중지수"
	기타_업종_코드_맵["070"] = "코스닥150 레버리지지수"
	기타_업종_코드_맵["071"] = "엔선물 지수"
	기타_업종_코드_맵["072"] = "엔선물 레버리지지수"
	기타_업종_코드_맵["073"] = "엔선물 인버스지수"
	기타_업종_코드_맵["074"] = "엔선물 인버스-2X지수"
	기타_업종_코드_맵["075"] = "엔선물 인버스-3X지수"
	기타_업종_코드_맵["076"] = "유로선물지수"
	기타_업종_코드_맵["077"] = "유로선물지수 레버리지지수"
	기타_업종_코드_맵["078"] = "유로선물지수 인버스지수"
	기타_업종_코드_맵["079"] = "유로선물지수 인버스-2X지수"
	기타_업종_코드_맵["080"] = "유로선물지수 인버스-3X지수"

	//	lib.F문자열_출력("\n" +
	//		"ETF 구분 : '%v'\n" +
	//		"업종 코드 : '%v'\n" +
	//		"KRX지수 코드 : '%v'\n" +
	//		"지수 이름 : '%v'\n" +
	//		"해외 지수 코드 : '%v'\n" +
	//		"기타 업종 코드 : '%v'\n" +
	//		"채권 지수 코드 : '%v'\n" +
	//		"채권 지수 세부 코드 : '%v'\n" +
	//		"지수 : '%v', 부호 : '%v', 등락폭 : '%v'\n" +
	//	 	"채권 지수 : '%v', 부호 : '%v', 등락폭 : '%v'\n",
	//	 	ETF_정보.ETF구분,
	//		s.M업종코드,
	//		s.KRX지수_코드,
	//		s.M지수_이름,
	//		s.M해외_지수_코드,
	//		s.M기타_업종_코드,
	//		s.M채권_지수_코드,
	//		s.M채권_지수_세부_코드,
	//		s.M지수, s.M지수_등락부호, s.M지수_등락폭,
	//		s.M채권_지수, s.M채권_지수_등락부호, s.M채권_지수_등락폭)

	이름_비교값 := ""
	존재함 := false

	검사_제외_종목코드_모음 := []string{
		"091210", // KRX 100 : 12, 0000 -> 00, 0001
		"100910", // KRX 100 : 12, 0000 -> 00, 0001
		"108440", // 코스닥 스타지수 : 코스피, 03, 0000 -> 코스닥, 43, ''
		"108480", // 0.0 -> MKF 코스닥엘리트30
		"108630", // 0.0 -> MKF 스타우량
		"117690", // 테스트 로직 한계임. 건너뛸 것
		"122390", // 코스닥 프리미어 : 코스피, 05, 0000 -> 코스닥, 45, ''
		"133690", // NASDAQ 100 : 해외지수코드 '2' -> 'EPNDX'
		"139280", // 0.0 -> 코스피 200 필수소비재, 기타 지수 코드 '' -> 056
		"168580", // 0.0 -> CSI 300, 해외 지수 코드 '' -> '000300'
		"169950", // 0.0 -> FTSE China A50 Index
		"170350", // 0.0 -> FnGuide 베타플러스 지수
		"174350", // 0.0 -> FnGuide 로우볼 지수
		"174360", // 0.0 -> CSI 100
		"181450", // 0.0 -> Markit iBoxx USD Liquid High Yield Index
		"181480", // 0.0 -> Dow Jones US Real Estate Index
		"182480", // 0.0 -> MSCI US REIT Index
		"182490", // 0.0 -> Markit iBoxx USD Liquid High Yield Index
		"183700", // 0.0 -> 주식국채혼합(채권형)지수, 기타 업종 코드 -> '021'
		"183710", // 0.0 -> 주식국채혼합(주식형)지수, 기타 업종 코드 -> '020'
		"185680", // 0.0 -> S&P Biotechnology Select Industry Index
		"189400", // 0.0 -> MSCI AC World Daily TR Net USD
		"190150", // 0.0 -> KAP Barbell Index
		"190160", // 0.0 -> KAP Money Market Index
		"190620", // 0.0 -> KIS MSB 단기 Index
		"192090", // 0.0 -> CSI 300, 해외 지수 코드 '' -> '000300'
		"195930", // 0.0 -> EURO STOXX50, 해외 지수 코드 '' -> 'SX5E'
		"195970", // 0.0 -> MSCI EAFE Index
		"195980", // 0.0 -> MSCI EM Index
		"196230", // 0.0 -> KIS MSB 5M Index(총수익)
		"200020", // 0.0 -> S&P Select Sector Technology Index
		"200030", // 0.0 -> S&P Select Sector Industrial Index
		"200040", // 0.0 -> S&P Select Sector Financials Index
		"200050", // 0.0 -> MSCI Germany Index
		"200250", // 0.0 -> CNX NIFTY INDEX(PR)
		"203780", // 0.0 -> NASDAQ Biotechnology
		"204480", // 0.0 -> CSI 300, 해외 지수 코드 '' -> '000300'
		"213610", // 0.0 -> WISE 삼성그룹밸류 인덱스
		"214980", // 0.0 -> KRW Cash PLUS 지수 (총수익)
		"217780", // 0.0 -> CSI 300, 해외 지수 코드 '' -> '000300'
		"217790", // 0.0 -> FnGuide Contrarian Index
		"219900", // 0.0 -> CSI 300, 해외 지수 코드 '' -> '000300'
		"225050", // 0.0 -> EURO STOXX50, 해외 지수 코드 '' -> 'SX5E'
		"225060", // 0.0 -> MSCI EM Index
		"225130", // 0.0 -> S&P WCI GOLD Excess Return Index
		"226810", // 0.0 -> KAP 단기채권지수(총수익)
		"226980", // 0.0 -> 코스피 200 중소형주 지수
		"227540", // 0.0 -> 코스피 200 건강관리 : '', '' -> '38', '0000'
		"227550", // 0.0 -> 코스피 200 산업재 : '', '' -> '37', '0000'
		"227570", // 0.0 -> FnGuide 퀄리티 밸류 지수
		"227930", // 0.0 -> 코스닥 150, 기타 업종 코드 '' -> '54'
		"228790", // 0.0 -> WISE 화장품 지수
		"228800", // 0.0 -> WISE 여행레저 지수
		"228810", // 0.0 -> WISE 미디어컨텐츠 지수
		"228820", // 0.0 -> KTOP 30, 기타 업종 코드 '' -> '53'
		"229200", // 0.0 -> 코스닥 150, 기타 업종 코드 '' -> '54'
		"229720", // 0.0 -> KTOP 30, 기타 업종 코드 '' -> '53'
		"099140", // HSCEI, 해외 지수 코드 '' -> 'EWHSCE'
		"156080", // MSCI Korea Index, 기타 업종 코드 '' -> '012'
		"157490", // FnGuide 소프트웨어 지수, 기타 업종 코드 '' -> '060' ??
		"157510", // FnGuide 자동차 지수, 기타 업종 코드 '' -> '063' ??
		"157520", // FnGuide 화학 지수, 기타 업종 코드 '' -> '065' ??
		"225040", // S&P 500, 해외 지수 코드 '' -> 'EPSPX'
		"226490", // 코스피 지수, 업종지수코드 '' -> '01'
		"227830", // 코스피 지수, 업종지수코드 '' -> '01'
		"139230", // 잘못된 지수 이름 '코스피200조선운송' -> '코스피 200 중공업'
		"143850", // 해외지수코드 '' -> 'SPFG' 혹은 'SPFR'
		"204420", // HSCEI, 해외 지수 코드 '' -> 'EWHSCE'
		"204450", // HSCEI, 해외 지수 코드 '' -> 'EWHSCE'
		"219480", // S&P500 선물지수(TR), 해외 지수 코드 '' -> 'SPFR' 혹은 'SPFG'
		"222180", // WISE 스마트 베타 Value, 기타 업종 코드 '' -> '047'
		"222190", // WISE 스마트 베타 Momentum, 기타 업종 코드 '' -> '046'
		"225030", // S&P500 선물지수(TR), 해외 지수 코드 '' -> 'SPFR' 혹은 'SPFG'
		"226380", // 잘못된 지수 이름 'FnGuide 바이오 TOP5' -> 'FnGuide 한류스타 지수'
		""}

	// KRX, 네이버 오류
	//105010 지수 이름 철자 틀림 : BNY Latin Amereica 35
	//160580 지수 이름 철자 틀림 : S&P CSCI Cash Copper Index

	// NH 지수 이름 오류
	//139230 : '코스피200조선운송' -> '코스피 200 중공업'
	//226380 : 'FnGuide 바이오 TOP5' -> 'FnGuide 한류스타 지수'

	for _, 종목_코드 := range 검사_제외_종목코드_모음 {
		if 종목_코드 == 기본_정보.M종목코드 {
			return
		}
	}

	비교값 := f지수_이름_검색(기본_정보.M종목코드)

	if s.M지수_이름 == "0.0" {
		lib.F문자열_출력("\n**** 지수 이름 없음 ****\n"+
			"종목코드 : '%v', 검색값 : '%v'\n"+
			"**** 지수 이름 없음 ****\n\n",
			기본_정보.M종목코드, 비교값)
		t.FailNow()
	}

	if !f지수_이름_같음(s.M지수_이름, 비교값) &&
		s.M지수_이름 != "0.0" {
		lib.F문자열_출력("\n* 잘못된 지수 이름 * 종목코드 '%v', 수신값 '%v', 비교값: '%v'",
			기본_정보.M종목코드, f지수_이름_정리(s.M지수_이름), f지수_이름_정리(비교값))
		//t.FailNow()
	}

	switch {
	case s.M업종코드 == "" && s.KRX지수_코드 == "" &&
		s.M해외지수_코드 == "" && s.M기타_업종_코드 == "" &&
		s.M채권지수_코드 == "" && s.M채권지수_세부_코드 == "":
		lib.F테스트_참임(t, utf8.ValidString(s.M지수_이름), s.M지수_이름)

		원본 := f지수_이름_정리(s.M지수_이름)

		for 코드, 비교값 := range 코스피_업종_맵 {
			if f지수_이름_같음(원본, 비교값) {
				if strings.HasPrefix(원본, "FNGUIDE") {
					continue
				}

				lib.F문자열_출력(
					"\n**** 지수 코드 누락 ****\n"+
						"종목코드 : '%v', 지수_이름 : '%v', 코스피_업종_코드 : '%v'\n"+
						"**** 지수 코드 누락 ****\n\n",
					기본_정보.M종목코드, s.M지수_이름, 코드)
				t.FailNow()
			}
		}

		for 코드, 비교값 := range 코스닥_업종_맵 {
			if f지수_이름_같음(원본, 비교값) {
				if strings.HasPrefix(원본, "FNGUIDE") {
					continue
				}

				lib.F문자열_출력(
					"\n**** 지수 코드 누락 ****\n"+
						"종목코드 : '%v', 지수_이름 : '%v', 코스닥_업종_코드 : '%v'\n"+
						"**** 지수 코드 누락 ****\n\n",
					기본_정보.M종목코드, s.M지수_이름, 코드)
				t.FailNow()
			}
		}

		for 코드, 비교값 := range KRX지수_코드_맵 {
			if f지수_이름_같음(원본, 비교값) {
				lib.F문자열_출력(
					"\n**** 지수 코드 누락 ****\n"+
						"종목코드 : '%v', 지수_이름 : '%v', KRX지수_코드 : '%v'\n"+
						"**** 지수 코드 누락 ****\n\n",
					기본_정보.M종목코드, s.M지수_이름, 코드)
				t.FailNow()
			}
		}

		for 코드, 채권_지수_정보 := range 채권지수_코드_맵 {
			if f지수_이름_같음(원본, 채권_지수_정보.M지수_이름) {
				lib.F문자열_출력(
					"\n**** 지수 코드 누락 ****\n"+
						"종목코드 : '%v', 지수_이름 : '%v', 채권_지수_코드 : '%v'\n"+
						"**** 지수 코드 누락 ****\n\n",
					기본_정보.M종목코드, s.M지수_이름, 코드)
				t.FailNow()
			}
		}

		for 코드, 비교값 := range 해외지수_맵 {
			if f지수_이름_같음(원본, 비교값[0]) {
				lib.F문자열_출력(
					"\n**** 지수 코드 누락 ****\n"+
						"종목코드 : '%v', 지수_이름 : '%v', 해외 지수_코드 : '%v'\n"+
						"**** 지수 코드 누락 ****\n\n",
					기본_정보.M종목코드, s.M지수_이름, 코드)
				t.FailNow()
			}

			if f지수_이름_같음(원본, 비교값[1]) {
				lib.F문자열_출력(
					"\n**** 지수 코드 누락 ****\n"+
						"종목코드 : '%v', 지수_이름 : '%v', 해외 지수_코드 : '%v'\n"+
						"**** 지수 코드 누락 ****\n\n",
					기본_정보.M종목코드, s.M지수_이름, 코드)
				t.FailNow()
			}
		}

		for 코드, 비교값 := range 기타_업종_코드_맵 {
			if f지수_이름_같음(원본, 비교값) {
				lib.F문자열_출력(
					"\n**** 지수 코드 누락 ****\n"+
						"종목코드 : '%v', 지수_이름 : '%v', 기타_업종_코드 : '%v'\n"+
						"**** 지수 코드 누락 ****\n\n",
					기본_정보.M종목코드, s.M지수_이름, 코드)
				t.FailNow()
			}
		}

		return
	case s.M업종코드 != "" && s.M업종코드 != "00":
		switch ETF_정보.ETF구분 {
		case P코스피:
			이름_비교값, 존재함 = 코스피_업종_맵[s.M업종코드]
		case P코스닥:
			이름_비교값, 존재함 = 코스닥_업종_맵[s.M업종코드]
		}

		lib.F테스트_같음(t, ETF_정보.ETF구분, P코스피, P코스닥)
		lib.F테스트_참임(t, 존재함, s.M업종코드, ETF_정보.ETF구분)
	case s.KRX지수_코드 != "":
		이름_비교값, 존재함 = KRX지수_코드_맵[s.KRX지수_코드]
		lib.F테스트_참임(t, 존재함, s.KRX지수_코드)
	case s.M채권지수_코드 != "":
		비교값, 존재함 := 채권지수_코드_맵[s.M채권지수_코드+"-"+s.M채권지수_세부_코드]

		lib.F테스트_참임(t, 존재함, s.M채권지수_코드)
		lib.F테스트_참임(t, s.M채권지수_코드 == 비교값.M채권지수_코드,
			s.M채권지수_코드, 비교값.M채권지수_코드)
		lib.F테스트_참임(t, s.M채권지수_세부_코드 == 비교값.M채권지수_세부_코드,
			s.M채권지수_코드, s.M채권지수_세부_코드, s.M지수_이름,
			비교값.M채권지수_코드, 비교값.M채권지수_세부_코드, 비교값.M지수_이름)

		이름_비교값 = 비교값.M지수_이름
	case s.M해외지수_코드 != "":
		해외_지수_정보, 존재함 := 해외지수_맵[s.M해외지수_코드]
		lib.F테스트_참임(t, 존재함, s.M해외지수_코드)

		이름_비교값 = 해외_지수_정보[1]
	case s.M기타_업종_코드 != "":
		이름_비교값, 존재함 = 기타_업종_코드_맵[s.M기타_업종_코드]
		lib.F테스트_참임(t, 존재함, s.M기타_업종_코드)
	}

	lib.F테스트_참임(t, utf8.ValidString(s.M지수_이름))

	이름_수신값 := f지수_이름_정리(s.M지수_이름)
	이름_비교값 = f지수_이름_정리(이름_비교값)

	lib.F테스트_참임(t,
		strings.Contains(이름_수신값, 이름_비교값) ||
			strings.Contains(이름_비교값, 이름_수신값),
		이름_수신값, 이름_비교값)

	if s.M채권지수_코드 == "" {
		lib.F테스트_참임(t, s.M지수 > 0)
		lib.F테스트_참임(t, f올바른_등락부호(s.M지수_등락부호), s.M지수_등락부호)
		lib.F테스트_참임(t, s.M지수_등락폭 >= 0) // 절대값?
	} else {
		lib.F테스트_참임(t, s.M채권지수 > 0)
		lib.F테스트_참임(t, f올바른_등락부호(s.M채권지수_등락부호), s.M채권지수_등락부호)
		lib.F테스트_참임(t, s.M채권지수_등락폭 >= 0) // 절대값?
	}
}

func f지수_이름_검색(종목_코드 string) string {
	url := `http://finance.naver.com/item/main.nhn?code=` + 종목_코드

	본문, 에러 := lib.F_HTTP회신_본문_CP949(url)

	if 에러 != nil {
		lib.F문자열_출력("지수 이름 검색 중 HTTP 에러 발생.\n%v", 에러)
		return ""
	}

	검색_결과 := lib.F정규식_검색(본문, []string{
		`<th scope="row">기초지수</th>` +
			`(.|\r|\n){0,300}` +
			`<th scope="row">유형</th>`,
		`<span title=".{0,100}">.{0,100}</span>`,
		`>.{0,100}<`})

	검색_결과 = 검색_결과[1 : len(검색_결과)-1]

	return 검색_결과
}

func f지수_이름_정리(지수_이름 string) string {
	지수_이름 = strings.ToUpper(지수_이름)
	변환할_문자열_모음 := []([]string){
		[]string{" IDX(", ""},
		[]string{"INDEX", ""},
		[]string{"지수", ""},
		[]string{" MOM ", "MOMENTUM"},
		[]string{"TOTALRETURN", "TR"},
		[]string{"(TR)", "(총수익)"},
		[]string{"(ER)", ""},
		[]string{"/", ""},
		[]string{"(", ""},
		[]string{")", ""},
		[]string{" ", ""},
		[]string{"F-KOSPI200", "코스피200선물"},
		[]string{"KOSPI", "코스피"},
		[]string{"KP100", "코스피100"},
		[]string{"KP200", "코스피200"},
		[]string{"FUTURES", "선물"},
		[]string{"KTB", "국고채"},
		[]string{"TB", "국고채"},
		[]string{"COVCALL", "커버드콜"},
		[]string{"KOSTAR", "코스닥스타"},
		[]string{"MKMSB", "MK통안채"},
		[]string{"KISMSB", "KIS통안채"},
		[]string{"KOBI", "KIS"},
		[]string{"3M", "3개월"},
		[]string{"3Y", "3년"},
		[]string{"10Y", "10년"},
		[]string{"GROUP", "그룹"},
		[]string{"KOREA", "코리아"},
		[]string{"MOMENTUM", "모멘텀"},
		[]string{"FN-", "FNGUIDE-"},
		[]string{"SELECT", ""},
		[]string{"동일가중코스피100", "코스피100동일가중"},
		[]string{"MKFHMG", "MKF현대차그룹"},
		[]string{"WISESB", "WISE스마트베타"},
		[]string{"스마트베타QTY", "스마트베타QUALITY"},
		[]string{"FN-RAFIKOREALARGE", "FNGUIDE-RAFI코리아대형"},
		[]string{"LATINAME35", "LATINAMERICA35"},
		[]string{"LATINAMEREICA", "LATINAMERICA"},
		[]string{"HSMAINLAND", "HANGSENGMAINLAND"},
		[]string{"S&PCSCI", "S&PGSCI"},
		[]string{"GSCINACITR", "GSCINORTHAMERICANCOPPER"},
		[]string{"GSCIIMSI", "GSCIINDUSTRIALMETALSSELECT"},
		[]string{"GSCIPMI", "GSCIPRECIOUSMETALS"},
		[]string{"GSCICCI", "GSCICASHCOPPER"},
		[]string{"CRUDEOIL", "OIL"}}

	for _, 변환 := range 변환할_문자열_모음 {
		지수_이름 = strings.Replace(지수_이름, 변환[0], 변환[1], -1)
	}

	return 지수_이름
}

func f지수_이름_같음(원본, 비교값 string) bool {
	원본 = f지수_이름_정리(원본)
	비교값 = f지수_이름_정리(비교값)

	if strings.Contains(원본, 비교값) ||
		strings.Contains(비교값, 원본) {
		return true
	}

	원본_길이 := int(math.Min(float64(len(원본)), 10))
	비교값_길이 := int(math.Min(float64(len(비교값)), 10))

	if 원본_길이 == 0 || 비교값_길이 == 0 {
		return false
	}

	원본 = 원본[0:원본_길이]
	비교값 = 비교값[0:비교값_길이]

	if strings.Contains(원본, 비교값) ||
		strings.Contains(비교값, 원본) {
		return true
	}

	return false
}
