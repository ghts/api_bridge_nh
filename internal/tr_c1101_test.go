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
)

func TestC1101_주식_현재가(t *testing.T) {
	lib.F메모("주식 현재가 시각값이 비어있는 경우가 종종 발생함.")
	t.SkipNow()

	lib.F대기(lib.P3초)
	lib.F테스트_에러없음(t, f접속_확인())

	기본_정보 := new(lib.NH주식_현재가_조회_기본_정보)
	변동_정보_모음 := make([]*lib.NH주식_현재가_조회_변동_거래량_정보, 0)
	동시호가_정보 := new(lib.NH주식_현재가_조회_동시호가_정보)

	최소_대기_시간 := time.Now().Add(time.Second)
	종목 := lib.F임의_종목_코스피_주식()

	질의값 := new(lib.S질의값_단일종목)
	질의값.TR구분 = lib.TR조회
	질의값.TR코드 = lib.NH_TR주식_현재가_조회
	질의값.M종목코드 = 종목.G코드()

	// 소켓으로 질의를 수신한다는 것을 감안한 변환.
	소켓_메시지, 에러 := lib.New소켓_메시지(lib.CBOR, 질의값)
	lib.F테스트_에러없음(t, 에러)

	채널_질의 := lib.New채널_질의(ch조회, lib.P30초, 1).S질의(소켓_메시지)

	수신완료_기본정보, 수신완료_거래량정보, 수신완료_동시호가정보 := false, false, false

	for {
		if 수신완료_기본정보 && 수신완료_거래량정보 && 수신완료_동시호가정보 {
			break
		}

		응답 := 채널_질의.G응답()
		lib.F테스트_에러없음(t, 응답.G에러())

		구분, ok := 응답.G값(0).(lib.TR응답_구분)
		lib.F테스트_참임(t, ok, 응답)

		lib.F테스트_같음(t, 구분, lib.TR응답_데이터, lib.TR응답_메시지, lib.TR응답_완료)

		switch 구분 {
		case lib.TR응답_데이터:
			lib.F테스트_같음(t, 응답.G길이(), 2)

			변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
			lib.F테스트_참임(t, ok, 응답.G값(1))

			//lib.F테스트_에러없음(t, 변환값.G값(데이터_블록))
			//lib.F테스트_다름(t, 데이터_블록.M데이터, nil)

			switch 변환값.G자료형_문자열() {
			case lib.F자료형_문자열(lib.NH주식_현재가_조회_기본_정보{}):
				수신완료_기본정보 = true
				lib.F테스트_에러없음(t, 변환값.G값(기본_정보))
			case lib.F자료형_문자열([]*lib.NH주식_현재가_조회_변동_거래량_정보{}):
				수신완료_거래량정보 = true
				lib.F테스트_에러없음(t, 변환값.G값(&변동_정보_모음))
			case lib.F자료형_문자열(lib.NH주식_현재가_조회_동시호가_정보{}):
				수신완료_동시호가정보 = true
				lib.F테스트_에러없음(t, 변환값.G값(동시호가_정보))
			default:
				lib.F문자열_출력("예상치 못한 자료형. %v", 변환값.G자료형_문자열())
				t.FailNow()
			}
		case lib.TR응답_메시지:
			lib.F테스트_같음(t, 응답.G길이(), 3)

			_, ok = 응답.G값(1).(string) // 코드
			lib.F테스트_참임(t, ok)

			메시지, ok := 응답.G값(2).(string)
			lib.F테스트_참임(t, ok)

			lib.F테스트_참임(t, strings.Contains(메시지, "조회완료"))
			lib.F체크포인트("완료 메시지")
			continue

			lib.F체크포인트("기타 메시지.\n%v", 메시지)
		case lib.TR응답_완료:
			// 완료 응답이 수신된 이후에도 정보가 수신되므로, break 하면 안 됨.
			lib.F테스트_같음(t, 응답.G길이(), 2)

			변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
			lib.F테스트_참임(t, ok, 응답.G값(1))

			if 변환값.G자료형_문자열() != "<nil>" {
				lib.New에러("예상하지 못한 자료형. %v", 변환값.G자료형_문자열())
			}
		default:
			lib.New에러("예상하지 TR응답 구분값. %v", 구분)
		}
	}

	//lib.F문자열_출력("*** 종목코드 %v ***", 종목.G코드())
	//lib.F문자열_출력("*** 시각 %v ***", 기본_정보.M시각)

	// 기본 자료 테스트
	f주식_현재가_조회_기본_정보_테스트(t, 기본_정보, 종목.G코드())

	// 변동 자료 테스트
	f주식_현재가_조회_변동_거래량_정보_테스트(t, 기본_정보, 변동_정보_모음)

	// 동시호가 자료 테스트
	lib.F테스트_참임(t, 동시호가_정보 != nil, "동시호가 자료를 수신하지 못함.")
	f주식_현재가_조회_동시호가_정보_테스트(t, 기본_정보, 동시호가_정보)

	// 서버 TR 제한을 피하기 위함.
	for time.Now().Before(최소_대기_시간) {
		time.Sleep(100 * time.Millisecond)
	}
}

func f주식_현재가_조회_기본_정보_테스트(t *testing.T, s *lib.NH주식_현재가_조회_기본_정보, 종목코드 string) {
	종목, 에러 := lib.F종목by코드(종목코드)
	lib.F에러2패닉(에러)

	지금 := time.Now()
	십분전 := 지금.Add(-10 * time.Minute)
	십분후 := 지금.Add(10 * time.Minute)
	금일_0시 := time.Date(지금.Year(), 지금.Month(), 지금.Day(), 0, 0, 0, 0, 지금.Location())
	금일_9시 := 금일_0시.Add(9 * time.Hour)
	개장_시각 := 금일_9시
	개장일_0시, 에러 := lib.F한국증시_최근_개장일()
	lib.F에러2패닉(에러)

	삼십일전 := 지금.Add(-30 * 24 * time.Hour)
	연초 := time.Date(지금.Year(), time.January, 1, 0, 0, 0, 0, 지금.Location())
	일년전 := 지금.Add(-1 * 366 * 24 * time.Hour)
	이백년전 := 지금.Add(-200 * 365 * 24 * time.Hour)

	lib.F테스트_참임(t, s != nil, "기본 자료를 수신하지 못함.")
	lib.F테스트_같음(t, s.M종목코드, 종목.G코드())
	lib.F테스트_참임(t, utf8.ValidString(s.M종목명))

	lib.F테스트_다름(t, strings.TrimSpace(s.M종목명), "")
	//lib.F테스트_참임(t, strings.Contains(s.M종목명, 종목.G이름()), s.M종목명, 종목.G이름())

	lib.F테스트_참임(t, s.M등락율 >= 0) // 절대값임.

	f테스트_등락부호(t, s.M등락부호, s.M현재가, s.M전일_종가, s.M상한가, s.M하한가)
	lib.F테스트_같음(t, s.M전일_종가+f등락부호2정수(s.M등락부호)*s.M등락폭, s.M현재가)

	if s.M현재가 != 0 && s.M등락폭 != 0 && s.M등락율 != 0 {
		등락율_근사값 := math.Abs(float64(s.M등락폭)) / float64(s.M현재가) * 100
		lib.F테스트_참임(t, lib.F오차율(등락율_근사값, s.M등락율) < 10)
	}

	lib.F테스트_참임(t, s.M거래량 >= 0)
	lib.F테스트_참임(t, s.M전일대비_거래량_비율 >= 0)

	if s.M거래량 != 0 && s.M전일_거래량 != 0 {
		거래량_비율_근사값 := float64(s.M거래량) / float64(s.M전일_거래량) * 100
		lib.F테스트_참임(t, lib.F오차율(s.M전일대비_거래량_비율, 거래량_비율_근사값) < 10,
			s.M전일대비_거래량_비율, 거래량_비율_근사값, s.M거래량, s.M전일_거래량)
	}

	if s.M유동_주식수_1000주 != 0 {
		유동주_회전율_근사값 := float64(s.M거래량) /
			float64(s.M유동_주식수_1000주*1000) * 100
		유동주_회전율_근사값 = math.Trunc(유동주_회전율_근사값*100) / 100
		lib.F테스트_참임(t, lib.F오차(s.M유동주_회전율, 유동주_회전율_근사값) < 1 ||
			lib.F오차율(s.M유동주_회전율, 유동주_회전율_근사값) < 10,
			s.M유동주_회전율, 유동주_회전율_근사값)
	}

	if s.M거래대금_100만 != 0 && s.M거래량 != 0 && s.M현재가 != 0 {
		거래대금_근사값 := s.M거래량 * s.M현재가 / 1000000
		lib.F테스트_참임(t, lib.F오차율(s.M거래대금_100만, 거래대금_근사값) < 10)
	}

	if s.M거래량 > 0 {
		// 거래량이 0이면 저가, 고가 모두 0임.
		lib.F테스트_참임(t, s.M저가 > 0, s.M저가)
		lib.F테스트_참임(t, s.M고가 > 0, s.M고가)

		lib.F테스트_참임(t, s.M저가 >= s.M하한가, s.M하한가, s.M저가)
		lib.F테스트_참임(t, s.M현재가 <= s.M고가, s.M현재가, s.M고가)
		lib.F테스트_참임(t, s.M상한가 >= s.M고가)
		lib.F테스트_참임(t, s.M고가 >= s.M시가)
		lib.F테스트_참임(t, s.M고가 >= s.M저가)
		lib.F테스트_참임(t, s.M시가 >= s.M저가)
		lib.F테스트_참임(t, s.M저가 >= s.M하한가)
		lib.F테스트_참임(t, s.M현재가 >= s.M저가)
		lib.F테스트_참임(t, s.M현재가 <= s.M고가)
		lib.F테스트_참임(t, s.M가중_평균_가격 >= s.M저가)
		lib.F테스트_참임(t, s.M가중_평균_가격 <= s.M고가)
	}

	lib.F테스트_참임(t, s.M하한가 > 0)
	lib.F테스트_참임(t, s.M연중_최저가 > 0)
	lib.F테스트_참임(t, s.M52주_고가 >= s.M연중_최고가)
	lib.F테스트_참임(t, s.M52주_고가 >= s.M20일_고가)
	lib.F테스트_참임(t, s.M20일_고가 >= s.M5일_고가)
	lib.F테스트_참임(t, s.M5일_고가 >= s.M5일_저가)
	lib.F테스트_참임(t, s.M연중_최저가 >= s.M52주_저가)
	lib.F테스트_참임(t, s.M20일_저가 >= s.M52주_저가)
	lib.F테스트_참임(t, s.M5일_저가 >= s.M20일_저가)
	lib.F테스트_참임(t, s.M연중_최고가 >= s.M연중_최저가)
	f테스트_등락부호(t, s.M시가대비_등락부호, s.M현재가, s.M시가, s.M상한가, s.M하한가)
	lib.F테스트_같음(t, s.M시가+s.M시가대비_등락폭, s.M현재가) // 시가대비_등락폭 자체에 부호가 반영되어 있음.
	lib.F테스트_참임(t, s.M시각.After(개장일_0시.Add(-1*time.Second)))
	lib.F테스트_참임(t, s.M시각.Before(개장일_0시.Add(18*time.Hour)), s.M시각)
	lib.F테스트_참임(t, s.M시각.Before(십분후))

	if lib.F한국증시_정규시장_거래시간임() { // 장중
		lib.F메모("현재가 시각이 상당히 이전 시간이 나오는 경우 발견함.")
		lib.F테스트_참임(t, s.M시각.After(십분전), s.M시각, 십분전)
		lib.F테스트_참임(t, s.M시각.Before(십분후), s.M시각, 십분후)
	} else { // 장중이 아니면 마감 시각 기록.
		lib.F테스트_같음(t, s.M시각.Hour(), 14, 15)
	}

	매도_잔량_합계 := int64(0)
	for i, 매도_잔량 := range s.M매도_잔량_모음 {
		lib.F테스트_참임(t, 매도_잔량 >= 0, i, 매도_잔량)

		if 매도_잔량 == 0 {
			continue
		}

		매도_잔량_합계 += 매도_잔량
		매도_호가 := s.M매도_호가_모음[i]
		lib.F테스트_참임(t, 매도_호가 <= s.M상한가)
		lib.F테스트_참임(t, 매도_호가 >= s.M하한가)

		switch i {
		case 0:
			lib.F테스트_참임(t, 매도_호가 >= s.M현재가)
		default:
			lib.F테스트_참임(t, 매도_호가 > s.M매도_호가_모음[i-1])
		}
	}

	매수_잔량_합계 := int64(0)
	for i, 매수_잔량 := range s.M매수_잔량_모음 {
		lib.F테스트_참임(t, 매수_잔량 >= 0, i, 매수_잔량)

		if 매수_잔량 == 0 {
			continue
		}

		매수_호가 := s.M매수_호가_모음[i]
		lib.F테스트_참임(t, 매수_호가 <= s.M상한가)
		lib.F테스트_참임(t, 매수_호가 >= s.M하한가)

		switch i {
		case 0:
			lib.F테스트_참임(t, 매수_호가 <= s.M현재가)
		default:
			lib.F테스트_참임(t, 매수_호가 < s.M매수_호가_모음[i-1],
				i, 매수_호가, s.M매수_호가_모음[i-1])
		}
	}

	lib.F테스트_참임(t, s.M매도_잔량_총합 >= 매도_잔량_합계)
	lib.F테스트_참임(t, s.M매수_잔량_총합 >= 매수_잔량_합계)
	lib.F테스트_참임(t, s.M시간외_매도_잔량 >= 0)
	lib.F테스트_참임(t, s.M시간외_매수_잔량 >= 0)
	lib.F테스트_참임(t, s.M피봇_2차_저항 >= s.M피봇_1차_저항)
	lib.F테스트_참임(t, s.M피봇_1차_저항 >= s.M피봇가)
	lib.F테스트_참임(t, s.M피봇가 >= s.M피봇_1차_지지)
	lib.F테스트_참임(t, s.M피봇_1차_지지 >= s.M피봇_2차_지지)
	lib.F테스트_참임(t, utf8.ValidString(s.M시장_구분))
	lib.F테스트_같음(t, s.M시장_구분, "코스피", "코스닥")
	lib.F테스트_참임(t, utf8.ValidString(s.M업종명))
	lib.F테스트_참임(t, utf8.ValidString(s.M자본금_규모))
	lib.F테스트_참임(t, strings.Contains(s.M자본금_규모, "형주") ||
		strings.TrimSpace(s.M자본금_규모) == "", s.M자본금_규모)
	lib.F테스트_참임(t, utf8.ValidString(s.M결산월))
	lib.F테스트_참임(t, strings.Contains(s.M결산월, "월 결산"))

	for _, 추가_정보 := range s.M추가_정보_모음 {
		lib.F테스트_참임(t, utf8.ValidString(추가_정보))
	}

	lib.F테스트_참임(t, utf8.ValidString(s.M서킷_브레이커_구분))
	lib.F테스트_같음(t, s.M서킷_브레이커_구분, "", "CB발동", "CB해제", "장종료")
	lib.F테스트_참임(t, s.M액면가 >= 0, s.M종목코드, s.M액면가, s.M종목코드, s.M종목명) // ETN은 액면가가 없음.
	//lib.F테스트_참임(t, strings.Contains(s.M전일_종가_타이틀, "전일종가"))
	lib.F테스트_참임(t, lib.F오차율(s.M상한가, float64(s.M전일_종가)*1.3) < 5)
	lib.F테스트_참임(t, lib.F오차율(s.M하한가, float64(s.M전일_종가)*0.7) < 5)
	lib.F테스트_참임(t, s.M대용가 < s.M전일_종가)
	lib.F테스트_참임(t, s.M대용가 > int64(float64(s.M전일_종가)*0.5))
	lib.F테스트_참임(t, s.M공모가 >= 0, s.M공모가)
	lib.F테스트_참임(t, s.M52주_저가_일자.After(일년전), s.M52주_저가_일자)
	lib.F테스트_참임(t, s.M52주_저가_일자.Before(지금), s.M52주_저가_일자)
	lib.F테스트_참임(t, s.M52주_고가_일자.After(일년전), s.M52주_고가_일자)
	lib.F테스트_참임(t, s.M52주_고가_일자.Before(지금), s.M52주_고가_일자)
	//lib.F테스트_참임(t, lib.F오차(s.M상장_주식수_1000주 - (s.M상장_주식수/1000) <= 1.01 ||
	//	lib.F오차율(s.M상장_주식수_1000주 - (s.M상장_주식수/1000)) < 10)
	lib.F테스트_참임(t, s.M유동_주식수_1000주 >= 0)

	시가총액_근사값 := s.M현재가 * s.M상장_주식수 / 100000000
	lib.F테스트_참임(t, lib.F오차율(s.M시가_총액_억, 시가총액_근사값) < 10)
	lib.F테스트_참임(t, s.M거래원_정보_수신_시각.Before(십분후),
		s.M거래원_정보_수신_시각, 십분후)

	if lib.F한국증시_정규시장_거래시간임() {
		lib.F테스트_참임(t, s.M거래원_정보_수신_시각.After(개장_시각))
		lib.F테스트_참임(t, s.M시각.Before(십분후))
	}

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
	lib.F테스트_참임(t, s.M외국인_시간.After(개장일_0시),
		s.M외국인_시간, 개장일_0시)
	lib.F테스트_참임(t, s.M외국인_시간.Before(개장일_0시.Add(23*time.Hour)),
		s.M외국인_시간, 개장일_0시.Add(23*time.Hour))

	if lib.F한국증시_정규시장_거래시간임() {
		lib.F테스트_참임(t, s.M외국인_시간.Before(십분후))
	}

	lib.F테스트_참임(t, s.M외국인_지분율 >= 0)
	lib.F테스트_참임(t, s.M외국인_지분율 <= 100)
	lib.F테스트_참임(t, s.M신용잔고_기준_결제일.After(삼십일전), s.M신용잔고_기준_결제일)
	lib.F테스트_참임(t, s.M신용잔고_기준_결제일.Before(개장일_0시.Add(18*time.Hour)))

	if lib.F한국증시_정규시장_거래시간임() {
		lib.F테스트_참임(t, s.M신용잔고_기준_결제일.Before(금일_0시))
	}

	lib.F테스트_참임(t, s.M신용잔고율 >= 0)
	lib.F테스트_참임(t, s.M신용잔고율 <= 100)
	//lib.F테스트_참임(t, s.M유상_기준일.After(이백년전) || s.M유상_기준일.IsZero())
	//lib.F테스트_참임(t, s.M무상_기준일.After(이백년전) || s.M무상_기준일.IsZero())
	lib.F테스트_참임(t, s.M유상_배정_비율 >= 0)
	lib.F테스트_참임(t, s.M유상_배정_비율 <= 100)
	//lib.F테스트_참임(t, s.M외국인_순매수량 >= 0, s.M외국인_순매수량)	// 순매도 시 (-) 값을 가질 수 있음.
	lib.F테스트_참임(t, s.M무상_배정_비율 >= 0)
	lib.F테스트_참임(t, s.M무상_배정_비율 <= 100)
	//lib.F변수값_확인(s.M당일_자사주_신청_여부)

	lib.F테스트_참임(t, s.M상장일.After(이백년전))
	lib.F테스트_참임(t, s.M상장일.Before(개장일_0시.Add(1*time.Second)))
	lib.F테스트_참임(t, s.M대주주_지분율 >= 0)
	lib.F테스트_참임(t, s.M대주주_지분율 <= 100)
	lib.F테스트_참임(t, s.M대주주_지분율_정보_일자.Before(십분후),
		s.M대주주_지분율_정보_일자)
	//lib.F변수값_확인(s.M네잎클로버_종목_여부)	// NH투자증권 선정 추천 종목
	lib.F테스트_참임(t, s.M증거금_비율 >= 0)
	lib.F테스트_참임(t, s.M증거금_비율 <= 100)
	lib.F테스트_참임(t, s.M자본금 > 0)
	lib.F테스트_참임(t, s.M전체_거래원_매도_합계 >= 매도_거래량_합계)
	lib.F테스트_참임(t, s.M전체_거래원_매수_합계 >= 매수_거래량_합계)

	//lib.F변수값_확인(s.M종목명2)
	//lib.F테스트_참임(t, utf8.ValidString(s.M종목명2))
	//lib.F변수값_확인(s.M우회_상장_여부)	// 이 항목은 뭐하는 데 필요할까?

	//lib.F테스트_참임(t, s.M코스피_구분_2 == "코스피" || s.M코스피_구분_2 == "코스닥")
	//lib.F테스트_참임(t, utf8.ValidString(s.M코스피_구분_2))   // 앞에 나온 '코스피/코스닥 구분'과 중복 아닌가?

	lib.F테스트_참임(t, s.M공여율_기준일.After(삼십일전))
	lib.F테스트_참임(t, s.M공여율_기준일.Before(지금))        // 공여율은 '신용거래 관련 비율'이라고 함.
	lib.F테스트_참임(t, s.M공여율 >= 0 && s.M공여율 <= 100) // 공여율(%)
	//lib.F테스트_참임(t, math.Abs(float64(s.PER)) < 100, s.PER, 종목.G코드())
	lib.F테스트_참임(t, s.M종목별_신용한도 >= 0)
	//lib.F테스트_참임(t, s.M종목별_신용한도 <= 100)
	lib.F테스트_참임(t, s.M가중_평균_가격 >= s.M저가)
	lib.F테스트_참임(t, s.M가중_평균_가격 <= s.M고가)
	lib.F테스트_참임(t, s.M추가_상장_주식수 >= 0)
	lib.F테스트_참임(t, utf8.ValidString(s.M종목_코멘트))
	lib.F테스트_참임(t, s.M전일_거래량 >= 0)
	lib.F테스트_참임(t, s.M전일_등락폭 >= 0) // 절대값
	lib.F테스트_참임(t, f올바른_등락부호(s.M전일_등락부호))
	lib.F테스트_참임(t, s.M연중_최고가_일자.After(연초))
	lib.F테스트_참임(t, s.M연중_최고가_일자.Before(지금))
	lib.F테스트_참임(t, s.M연중_최저가_일자.After(연초))
	lib.F테스트_참임(t, s.M연중_최저가_일자.Before(지금))
	lib.F테스트_참임(t, s.M외국인_보유_주식수 <= s.M상장_주식수+s.M추가_상장_주식수)
	lib.F테스트_참임(t, s.M외국인_지분_한도 >= 0)
	lib.F테스트_참임(t, s.M외국인_지분_한도 <= 100)
	lib.F테스트_참임(t, s.M매매_수량_단위 == 1, s.M매매_수량_단위)
	lib.F테스트_같음(t, int(s.M대량_매매_방향), 0, 1, 2)

	lib.F문자열_출력("대량 매매 관련 테스트 일시 보류.") //함. API버그로 의심됨.")
	//if s.M대량_매매_방향 == 0 {
	//	lib.F테스트_거짓임(t, s.M대량_매매_존재)
	//} else {
	//	lib.F테스트_참임(t, s.M대량_매매_존재)
	//}
}

func f주식_현재가_조회_변동_거래량_정보_테스트(t *testing.T,
	기본_정보 *lib.NH주식_현재가_조회_기본_정보,
	변동_정보_모음 []*lib.NH주식_현재가_조회_변동_거래량_정보) {
	lib.F테스트_참임(t, len(변동_정보_모음) > 0, "변동 자료를 수신하지 못함.", 기본_정보.M종목코드)

	거래량_잔량 := 기본_정보.M거래량
	지금 := time.Now()
	삼분후 := 지금.Add(3 * time.Minute)
	개장일_0시, 에러 := lib.F한국증시_최근_개장일()
	lib.F에러2패닉(에러)

	for i, s := range 변동_정보_모음 {
		lib.F테스트_참임(t, s.M시각.After(개장일_0시.Add(9*time.Hour)))
		lib.F테스트_참임(t, s.M시각.Before(삼분후), s.M시각)
		lib.F테스트_참임(t, s.M매도_호가 >= 0)
		lib.F테스트_참임(t, s.M매수_호가 >= 0)
		lib.F테스트_참임(t, s.M매도_호가 >= 기본_정보.M하한가 ||
			s.M매도_호가 == 0)
		lib.F테스트_참임(t, s.M매도_호가 <= 기본_정보.M상한가)
		lib.F테스트_참임(t, s.M매수_호가 >= 기본_정보.M하한가 ||
			s.M매수_호가 == 0)
		lib.F테스트_참임(t, s.M매수_호가 <= 기본_정보.M상한가)
		lib.F테스트_참임(t, s.M현재가 <= 기본_정보.M상한가)
		lib.F테스트_참임(t, s.M현재가 >= 기본_정보.M하한가)

		if lib.F한국증시_정규시장_거래시간임() {
			lib.F테스트_참임(t, s.M시각.Before(삼분후), s.M시각)
			lib.F테스트_참임(t, s.M매도_호가 >= s.M현재가 ||
				s.M매도_호가 == 0)
			lib.F테스트_참임(t, s.M매수_호가 <= s.M현재가 ||
				s.M매수_호가 == 0)
		} else {
			lib.F테스트_참임(t, s.M시각.Before(개장일_0시.Add(18*time.Hour)))

			// 장 마감 후 매도호가, 매수호가는 흔히 생각하는 조건을 만족시키지 않음.
			//lib.F테스트_참임(t, s.M매도_호가 >= s.M현재가)
			//lib.F테스트_참임(t, s.M매수_호가 <= s.M현재가)
		}

		lib.F테스트_참임(t, f올바른_등락부호(s.M등락부호))
		lib.F테스트_같음(t, f등락부호2정수(s.M등락부호)*s.M등락폭,
			s.M현재가-기본_정보.M전일_종가)

		// 걸러낸 자료로 인한 오차 수정.
		if i == 0 && s.M거래량 != 거래량_잔량 {
			거래량_잔량 = s.M거래량
		}

		lib.F테스트_같음(t, s.M거래량, 거래량_잔량)
		거래량_잔량 -= s.M변동_거래량
	}
}

func f주식_현재가_조회_동시호가_정보_테스트(t *testing.T,
	기본_정보 *lib.NH주식_현재가_조회_기본_정보,
	s *lib.NH주식_현재가_조회_동시호가_정보) {
	lib.F테스트_다름(t, s, nil)
	lib.F테스트_같음(t, int(s.M동시호가_구분), 0, 1, 2, 3, 4, 5, 6)

	if s.M동시호가_구분 == 0 { // 동시호가 아님.
		return
	}

	lib.F변수값_확인(기본_정보.M시각, 기본_정보.M종목코드, s.M동시호가_구분)
	lib.F테스트_참임(t, f올바른_등락부호(s.M예상_체결_부호), s.M예상_체결_부호)
	lib.F테스트_참임(t, s.M예상_체결가 <= 기본_정보.M상한가)
	lib.F테스트_참임(t, s.M예상_체결가 >= 기본_정보.M하한가)
	lib.F테스트_참임(t, lib.F오차율(s.M예상_체결가, 기본_정보.M현재가) < 10)
	lib.F테스트_같음(t, f등락부호2정수(s.M예상_체결_부호)*s.M예상_등락폭,
		s.M예상_체결가-기본_정보.M전일_종가)

	if s.M예상_등락폭 != 0 && s.M예상_등락율 != 0 {
		예상_등락율_근사값 := math.Abs(float64(s.M예상_등락폭)) /
			float64(s.M예상_체결가) * 100
		lib.F테스트_참임(t, lib.F오차율(s.M예상_등락율, 예상_등락율_근사값) < 10)
	}

	lib.F테스트_참임(t, s.M예상_체결_수량 >= 0)
	lib.F테스트_참임(t, s.M예상_체결_수량 <= 기본_정보.M매도_잔량_총합 ||
		s.M예상_체결_수량 <= 기본_정보.M매수_잔량_총합)
}
