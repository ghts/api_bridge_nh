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
GNU Lesser General Public License for more details.M

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. */

package internal

// ghts의 bin디렉토리에 있는 sync_ctype.bat에서
// go tool cgo -godefs 를 실행시켜서
// wmca_type.h에 있는 C언어 구조체를 자동으로 Go언어 구조체로 변환시킴.
// 생성된 결과물은 서로 직접 변환(cast)되어도 안전함.
//go:generate ctype_sync.bat

// #include <stdlib.h>
// #include "./c_func.h"
import "C"

import (
	"github.com/ghts/lib"

	"regexp"
	"strings"
	"time"
	"unsafe"
)

//----------------------------------------------------------------------//
// WMCA_CONNECTED 로그인 구조체
//----------------------------------------------------------------------//
func New로그인_정보(c *C.LOGINBLOCK) (s *lib.NH로그인_정보) {
	//defer func() { C.free(unsafe.Pointer(c)) }()

	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		lib.New에러(r)
		s = nil
	}})

	g := (*LoginInfo)(unsafe.Pointer(c.LoginInfo))
	s = new(lib.NH로그인_정보)
	s.M접속_시각 = lib.F2포맷된_시각_단순형("20060102150405", g.Date)
	s.M접속_서버 = lib.F2문자열(g.ServerName)
	s.M접속_ID = lib.F2문자열(g.UserID)

	계좌_수량 := lib.F2정수_단순형(lib.F2문자열(g.AccountCount))
	s.M계좌_목록 = make([]*lib.NH계좌_정보, 계좌_수량)

	for i, 계좌_정보 := range g.Accountlist[:계좌_수량] {
		g계좌_정보 := new계좌_정보(&계좌_정보, i+1) // 계좌 인덱스는 1부터 시작함.

		if g계좌_정보 == nil {
			continue
		}

		s.M계좌_목록[i] = g계좌_정보
	}

	return s
}

func new계좌_정보(g *AccountInfo, 계좌_인덱스 int) (s *lib.NH계좌_정보) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	s = new(lib.NH계좌_정보)
	s.M계좌_번호 = lib.F2문자열(g.AccountNo)
	s.M계좌명 = lib.F2문자열(g.AccountName)
	s.M상품_코드 = lib.F2문자열(g.AccountProductCode)
	s.M관리점_코드 = lib.F2문자열(g.AmnTabCode)
	s.M위임_만기일 = lib.F2포맷된_일자_단순형("20060102", lib.F2문자열(g.ExpirationDate))
	s.M계좌_인덱스 = 계좌_인덱스

	if lib.F2문자열(g.Granted) == "G" {
		s.M일괄주문_허용계좌 = true
	} else {
		s.M일괄주문_허용계좌 = false
	}

	s.M주석 = lib.F2문자열(g.Filler)

	return s
}

//----------------------------------------------------------------------//
// WMCA 문자 message 구조체
//----------------------------------------------------------------------//
func New수신_메시지_블록(c *C.OUTDATABLOCK) (s *NH수신_메시지_블록) {
	defer func() { C.free(unsafe.Pointer(c)) }()

	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*MsgHeader)(unsafe.Pointer(c.DataStruct.DataString))
	s = new(NH수신_메시지_블록)
	s.M식별번호 = int64(c.TrIdNo)
	s.M메시지_코드 = lib.F2문자열(g.MsgCode)
	s.M메시지_내용 = lib.F2문자열_CP949(g.UsrMsg)

	return s
}

//----------------------------------------------------------------------//
// 주식 현재가 조회 (c1101)
//----------------------------------------------------------------------//
func NewTc1101InBlock(종목코드 string) unsafe.Pointer {
	g := new(Tc1101InBlock)
	lib.F바이트_복사_문자열(g.Lang[:], "k")
	lib.F바이트_복사_문자열(g.Code[:], 종목코드)

	return unsafe.Pointer(g)
}

func New주식_현재가_조회_기본_정보(c *C.char) (s *lib.NH주식_현재가_조회_기본_정보) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	지금 := time.Now()
	금일_24시 := time.Date(지금.Year(), 지금.Month(), 지금.Day(),
		23, 59, 59, 0, 지금.Location())

	g := (*Tc1101OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH주식_현재가_조회_기본_정보)
	s.M종목코드 = lib.F2문자열(g.Code)
	s.M종목명 = lib.F2문자열_CP949(g.Title)
	s.M현재가 = lib.F2정수64_단순형(g.MarketPrice)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = lib.F2정수64_단순형(g.Diff)
	s.M등락율 = lib.F2실수_소숫점_추가(g.DiffRate, 2)
	//s.M매도_호가 = lib.F2정수64_단순형(g.OfferPrice)
	//s.M매수_호가 = lib.F2정수64_단순형(g.BidPrice)
	s.M거래량 = lib.F2정수64_단순형(g.Volume)
	s.M전일대비_거래량_비율 = lib.F2실수_소숫점_추가(g.TrVolRate, 2) // (당일 거래량 / 전일 거래량)
	s.M유동주_회전율 = lib.F2실수_소숫점_추가(g.FloatRate, 2)
	s.M거래대금_100만 = lib.F2정수64_단순형(g.TrAmount)
	s.M상한가 = lib.F2정수64_단순형(g.UpLmtPrice)
	s.M고가 = lib.F2정수64_단순형(g.High)
	s.M시가 = lib.F2정수64_단순형(g.Open)
	s.M시가대비_등락부호 = f2등락부호(g.VsOpenSign)
	s.M시가대비_등락폭 = lib.F2정수64_단순형(g.VsOpenDiff)
	s.M저가 = lib.F2정수64_단순형(g.Low)
	s.M하한가 = lib.F2정수64_단순형(g.LowLmtPrice)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)

	매도_호가_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.OfferPrice1,
			g.OfferPrice2,
			g.OfferPrice3,
			g.OfferPrice4,
			g.OfferPrice5,
			g.OfferPrice6,
			g.OfferPrice7,
			g.OfferPrice8,
			g.OfferPrice9,
			g.OfferPrice10})

	매도_잔량_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.OfferVolume1,
			g.OfferVolume2,
			g.OfferVolume3,
			g.OfferVolume4,
			g.OfferVolume5,
			g.OfferVolume6,
			g.OfferVolume7,
			g.OfferVolume8,
			g.OfferVolume9,
			g.OfferVolume10})

	s.M매도_호가_모음 = make([]int64, 0)
	s.M매도_잔량_모음 = make([]int64, 0)

	for i, 매도_호가 := range 매도_호가_모음 {
		if 매도_호가 >= s.M현재가 && 매도_잔량_모음[i] > 0 {
			s.M매도_호가_모음 = append(s.M매도_호가_모음, 매도_호가)
			s.M매도_잔량_모음 = append(s.M매도_잔량_모음, 매도_잔량_모음[i])
		}
	}

	매수_호가_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.BidPrice1,
			g.BidPrice2,
			g.BidPrice3,
			g.BidPrice4,
			g.BidPrice5,
			g.BidPrice6,
			g.BidPrice7,
			g.BidPrice8,
			g.BidPrice9,
			g.BidPrice10})

	매수_잔량_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.BidVolume1,
			g.BidVolume2,
			g.BidVolume3,
			g.BidVolume4,
			g.BidVolume5,
			g.BidVolume6,
			g.BidVolume7,
			g.BidVolume8,
			g.BidVolume9,
			g.BidVolume10})

	s.M매수_호가_모음 = make([]int64, 0)
	s.M매수_잔량_모음 = make([]int64, 0)

	for i, 매수_호가 := range 매수_호가_모음 {
		if 매수_호가 <= s.M현재가 && 매수_잔량_모음[i] > 0 {
			s.M매수_호가_모음 = append(s.M매수_호가_모음, 매수_호가)
			s.M매수_잔량_모음 = append(s.M매수_잔량_모음, 매수_잔량_모음[i])
		}
	}

	s.M매도_잔량_총합 = lib.F2정수64_단순형(g.OfferVolTot)
	s.M매수_잔량_총합 = lib.F2정수64_단순형(g.BidVolTot)
	s.M시간외_매도_잔량 = lib.F2정수64_단순형(g.OfferVolAfterHour)
	s.M시간외_매수_잔량 = lib.F2정수64_단순형(g.BidVolAfterHour)
	s.M피봇_2차_저항 = lib.F2정수64_단순형(g.PivotUp2)
	s.M피봇_1차_저항 = lib.F2정수64_단순형(g.PivotUp1)
	s.M피봇가 = lib.F2정수64_단순형(g.PivotPrice)
	s.M피봇_1차_지지 = lib.F2정수64_단순형(g.PivotDown1)
	s.M피봇_2차_지지 = lib.F2정수64_단순형(g.PivotDown2)
	s.M시장_구분 = lib.F2문자열_CP949(g.Market)
	s.M업종명 = lib.F2문자열_CP949(g.Sector)
	s.M자본금_규모 = lib.F2문자열_CP949(g.CapSize)
	s.M결산월 = lib.F2문자열_CP949(g.SettleMonth)

	s.M추가_정보_모음 = []string{
		lib.F2문자열_CP949(g.MarketAction1),
		lib.F2문자열_CP949(g.MarketAction2),
		lib.F2문자열_CP949(g.MarketAction3),
		lib.F2문자열_CP949(g.MarketAction4),
		lib.F2문자열_CP949(g.MarketAction5),
		lib.F2문자열_CP949(g.MarketAction6)}

	s.M서킷_브레이커_구분 = strings.TrimSpace(lib.F2문자열_CP949(g.CircuitBreaker))
	s.M액면가 = lib.F2정수64_단순형(g.NominalPrice)
	//s.M전일_종가_타이틀 = lib.F2문자열_CP949(g.PrevPriceTitle)
	s.M전일_종가 = lib.F2정수64_단순형(g.PrevPrice)
	s.M대용가 = lib.F2정수64_단순형(g.CollateralValue)
	s.M공모가 = lib.F2정수64_단순형(g.PublicOfferPrice)
	s.M5일_고가 = lib.F2정수64_단순형(g.High5Day)
	s.M5일_저가 = lib.F2정수64_단순형(g.Low5Day)
	s.M20일_고가 = lib.F2정수64_단순형(g.High20Day)
	s.M20일_저가 = lib.F2정수64_단순형(g.Low20Day)
	s.M52주_고가 = lib.F2정수64_단순형(g.High1Year)

	시각, 에러 := lib.F2포맷된_일자("0102", g.High1YearDate)
	if 에러 == nil {
		s.M52주_고가_일자 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	if s.M52주_고가_일자.After(금일_24시) {
		s.M52주_고가_일자 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	s.M52주_저가 = lib.F2정수64_단순형(g.Low1Year)

	시각, 에러 = lib.F2포맷된_일자("0102", g.Low1YearDate)
	if 에러 == nil {
		s.M52주_저가_일자 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	if s.M52주_저가_일자.After(금일_24시) {
		s.M52주_저가_일자 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	s.M유동_주식수_1000주 = lib.F2정수64_단순형(g.FloatVolume)
	//s.M상장_주식수_1000주,_ = lib.F2정수64_단순형(g.ListVolBy1000)
	s.M시가_총액_억 = lib.F2정수64_단순형(g.MarketCapital)

	s.M거래원_정보_수신_시각 = lib.F2한국증시_개장일_시각_단순형("15:04", g.TraderInfoTime)

	s.M매도_거래원_모음 = []string{
		lib.F2문자열_CP949(g.Seller1),
		lib.F2문자열_CP949(g.Seller2),
		lib.F2문자열_CP949(g.Seller3),
		lib.F2문자열_CP949(g.Seller4),
		lib.F2문자열_CP949(g.Seller5)}

	s.M매도_거래량_모음 = lib.F2정수64_모음_단순형(
		[]interface{}{
			g.Seller1Volume,
			g.Seller2Volume,
			g.Seller3Volume,
			g.Seller4Volume,
			g.Seller5Volume})

	s.M매수_거래원_모음 = []string{
		lib.F2문자열_CP949(g.Buyer1),
		lib.F2문자열_CP949(g.Buyer2),
		lib.F2문자열_CP949(g.Buyer3),
		lib.F2문자열_CP949(g.Buyer4),
		lib.F2문자열_CP949(g.Buyer5)}

	s.M매수_거래량_모음 = lib.F2정수64_모음_단순형(
		[]interface{}{
			g.Buyer1Volume,
			g.Buyer2Volume,
			g.Buyer3Volume,
			g.Buyer4Volume,
			g.Buyer5Volume})

	if strings.TrimSpace(lib.F2문자열(g.ForeignSellVolume)) != "" {
		s.M외국인_매도_거래량 = lib.F2정수64_단순형(g.ForeignSellVolume)
	} else {
		s.M외국인_매도_거래량 = 0
	}

	if strings.TrimSpace(lib.F2문자열(g.ForeignBuyVolume)) != "" {
		s.M외국인_매수_거래량 = lib.F2정수64_단순형(g.ForeignBuyVolume)
	} else {
		s.M외국인_매수_거래량 = 0
	}

	s.M외국인_시간 = lib.F2한국증시_개장일_시각_단순형("15:04", g.ForeignTime)
	s.M외국인_지분율 = lib.F2실수_소숫점_추가(g.ForeignHoldingRate, 2)

	시각, 에러 = lib.F2포맷된_시각("0102", g.SettleDate)
	if 에러 == nil {
		s.M신용잔고_기준_결제일 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	if s.M신용잔고_기준_결제일.After(금일_24시) {
		s.M신용잔고_기준_결제일 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	s.M신용잔고율 = lib.F2실수_소숫점_추가(g.DebtPercent, 2)

	/* if strings.TrimSpace(lib.F2문자열(g.RightsIssueDate)) != "" {
		lib.F문자열_출력("유상 기준일 : '%v'", lib.F2문자열(g.RightsIssueDate))
	}

	시각, 에러 = lib.F2포맷된_시각_단순형("0102", g.RightsIssueDate)
	if 에러 == nil {
		s.M유상_기준일 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	if s.M유상_기준일.After(금일_24시) {
		s.M유상_기준일 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	if strings.TrimSpace(lib.F2문자열(g.BonusIssueDate)) != "" {
		lib.F문자열_출력("무상 기준일 : '%v'", lib.F2문자열(g.BonusIssueDate))
	}

	시각, 에러 = lib.F2포맷된_시각_단순형("0102", g.BonusIssueDate)
	if 에러 == nil {
		s.M무상_기준일 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	if s.M무상_기준일.After(금일_24시) {
		s.M무상_기준일 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	} */

	s.M유상_배정_비율 = lib.F2실수_소숫점_추가(g.RightsIssueRate, 2)
	s.M무상_배정_비율 = lib.F2실수_소숫점_추가(g.BonusIssueRate, 2)
	s.M외국인_순매수량 = lib.F2정수64_단순형(g.ForeignFloatVol)
	s.M당일_자사주_신청_여부 = lib.F2참거짓(g.TreasuryStock, "1", true)
	s.M상장일 = lib.F2포맷된_시각_단순형("20060102", g.IpoDate)
	s.M대주주_지분율 = lib.F2실수_소숫점_추가(g.MajorHoldRate, 2)
	s.M대주주_지분율_정보_일자 = lib.F2포맷된_시각_단순형("060102", g.MajorHoldInfoDate)
	s.M네잎클로버_종목_여부 = lib.F2참거짓(g.FourLeafClover, "1", true)

	증거금_비율_위치 := lib.F2문자열(g.MarginRate)
	증거금_비율_문자열 := ""

	switch 증거금_비율_위치 {
	case "1", "2", "3", "4", "5", "6":
		위치 := lib.F2정수_단순형(증거금_비율_위치)
		증거금_비율_문자열 = s.M추가_정보_모음[위치-1]
	default:
		에러 := lib.New에러("예상하지 못한 증거금_비율_위치. %v", 증거금_비율_위치)
		panic(에러)
	}

	증거금_비율_문자열 = strings.Replace(증거금_비율_문자열, "증거금", "", -1)
	증거금_비율_문자열 = strings.Replace(증거금_비율_문자열, "%", "", -1)
	증거금_비율_문자열 = strings.TrimSpace(증거금_비율_문자열)
	s.M증거금_비율 = lib.F2실수_단순형(증거금_비율_문자열)

	s.M자본금 = lib.F2정수64_단순형(g.Capital)
	s.M전체_거래원_매도_합계 = lib.F2정수64_단순형(g.SellTotalSum)
	s.M전체_거래원_매수_합계 = lib.F2정수64_단순형(g.BuyTotalSum)
	//s.M종목명2 = lib.F2문자열_CP949(g.Title2)
	s.M우회_상장_여부 = lib.F2참거짓(g.BackdoorListing, "1", true)
	//s.M유동주_회전율_2 = lib.F2실수_소숫점_추가(g.FloatRate2, 2)
	//s.M코스피_구분_2 = lib.F2문자열_CP949(g.Market2)

	시각, 에러 = lib.F2포맷된_일자("0102", g.DebtTrDate)
	if 에러 == nil {
		s.M공여율_기준일 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	if s.M공여율_기준일.After(금일_24시) {
		s.M공여율_기준일 = time.Date(지금.Year()-1, 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	s.M공여율 = lib.F2실수_소숫점_추가(g.DebtTrPercent, 2)
	s.PER = lib.F2실수_소숫점_추가(g.PER, 2)

	종목별_신용한도_위치 := lib.F2문자열(g.DebtLimit)
	종목별_신용한도_문자열 := ""

	switch 종목별_신용한도_위치 {
	case "0":
		종목별_신용한도_문자열 = ""
	case "1", "2", "3", "4", "5", "6":
		위치 := lib.F2정수_단순형(종목별_신용한도_위치)
		종목별_신용한도_문자열 = s.M추가_정보_모음[위치-1]
	default:
		에러 := lib.New에러("시각 %v, 종목코드 %v\n"+
			"예상하지 못한 종목별_신용한도_위치. %v\n"+
			" 1 : %v\n2 : %v\n3 : %v\n"+
			" 4 : %v\n5 : %v\n6 : %v\n",
			s.M시각, s.M종목코드, 종목별_신용한도_위치,
			s.M추가_정보_모음[0], s.M추가_정보_모음[1], s.M추가_정보_모음[2],
			s.M추가_정보_모음[3], s.M추가_정보_모음[4], s.M추가_정보_모음[5])
		panic(에러)
	}

	종목별_신용한도_문자열 = strings.TrimSpace(종목별_신용한도_문자열)

	if 종목별_신용한도_문자열 == "" {
		s.M종목별_신용한도 = 0
	} else if strings.Contains(종목별_신용한도_문자열, "없음") {
		s.M종목별_신용한도 = 100
	} else if 일치, _ := regexp.MatchString("^신용한도 [0-9]+억$", 종목별_신용한도_문자열); 일치 {
		금액_문자열_억단위 := lib.F정규식_검색(종목별_신용한도_문자열, []string{"[0-9]+"})
		금액_문자열 := 금액_문자열_억단위 + "00000000"
		s.M종목별_신용한도 = lib.F2실수_단순형(금액_문자열)
	} else {
		lib.New에러("예상하지 못한 패턴.\n'%v'", 종목별_신용한도_문자열)
		return nil
	}

	s.M가중_평균_가격 = lib.F2정수64_단순형(g.WeightAvgPrice)
	s.M상장_주식수 = lib.F2정수64_단순형(g.ListedVolume)
	s.M추가_상장_주식수 = lib.F2정수64_단순형(g.AddListing)
	s.M종목_코멘트 = lib.F2문자열_CP949(g.Comment)
	s.M전일_거래량 = lib.F2정수64_단순형(g.PrevVolume)
	s.M전일_등락부호 = f2등락부호(g.VsPrevSign)
	s.M전일_등락폭 = lib.F2정수64_단순형(g.VsPrevDiff)
	s.M연중_최고가 = lib.F2정수64_단순형(g.High1Year2)

	시각, 에러 = lib.F2포맷된_일자("0102", g.High1YearDate2)
	if 에러 == nil {
		s.M연중_최고가_일자 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	s.M연중_최저가 = lib.F2정수64_단순형(g.Low1Year2)

	시각, 에러 = lib.F2포맷된_일자("0102", g.Low1YearDate2)
	if 에러 == nil {
		s.M연중_최저가_일자 = time.Date(지금.Year(), 시각.Month(), 시각.Day(),
			0, 0, 0, 0, 시각.Location())
	}

	s.M외국인_보유_주식수 = lib.F2정수64_단순형(g.ForeignHoldQty)
	s.M외국인_지분_한도 = lib.F2실수_소숫점_추가(g.ForeignLmtPercent, 2)
	s.M매매_수량_단위 = lib.F2정수64_단순형(g.TrUnitVolume)
	s.M대량_매매_방향 = uint8(g.DarkPoolOfferBid[0])             // 0 = 해당없음 1 = 매도 2 = 매수
	s.M대량_매매_존재 = lib.F2참거짓(g.DarkPoolExist[0], 48, false) // 48 => '0', 49=>'1'

	//lib.F문자열_출력("%v", time.Now())
	//lib.F문자열_출력("%v %v", g.DarkPoolOfferBid, g.DarkPoolExist)
	//lib.F문자열_출력("%v %v", s.M대량_매매_방향, s.M대량_매매_존재)

	return s
}

func New주식_현재가_조회_변동_거래량_정보(c *C.Tc1101OutBlock2) (s *lib.NH주식_현재가_조회_변동_거래량_정보) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tc1101OutBlock2)(unsafe.Pointer(c))
	s = new(lib.NH주식_현재가_조회_변동_거래량_정보)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)
	s.M현재가 = lib.F2정수64_단순형(g.MarketPrice)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = lib.F2정수64_단순형(g.Diff)
	s.M매도_호가 = lib.F2정수64_단순형(g.OfferPrice)
	s.M매수_호가 = lib.F2정수64_단순형(g.BidPrice)
	s.M변동_거래량 = lib.F2정수64_단순형(g.DiffVolume)
	s.M거래량 = lib.F2정수64_단순형(g.Volume)

	if s.M매도_호가 < s.M현재가 && s.M매도_호가 > 0 {
		s.M매도_호가 = s.M현재가
	}

	if s.M매수_호가 > s.M현재가 && s.M현재가 > 0 {
		s.M매수_호가 = s.M현재가
	}

	return s
}

func New주식_현재가_조회_동시호가(c *C.char) (s *lib.NH주식_현재가_조회_동시호가_정보) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tc1101OutBlock3)(unsafe.Pointer(c))
	s = new(lib.NH주식_현재가_조회_동시호가_정보)
	s.M동시호가_구분 = uint8(lib.F2정수_단순형(g.SyncOfferBid))
	s.M예상_체결가 = lib.F2정수64_단순형(g.EstmPrice)
	s.M예상_체결_부호 = f2등락부호(g.EstmSign)
	s.M예상_등락폭 = lib.F2정수64_단순형(g.EstmDiff)
	s.M예상_등락율 = lib.F2실수_소숫점_추가(g.EstmDiffRate, 2)
	s.M예상_체결_수량 = lib.F2정수64_단순형(g.EstmVolume)

	return s
}

//----------------------------------------------------------------------//
// ETF 현재가 조회 (c1151)
//----------------------------------------------------------------------//
func NewTc1151InBlock(종목코드 string) unsafe.Pointer {
	g := new(Tc1151InBlock)
	lib.F바이트_복사_문자열(g.Lang[:], "k")
	lib.F바이트_복사_문자열(g.Code[:], 종목코드)

	return unsafe.Pointer(g)
}

func New_ETF_현재가_조회_기본_정보(c *C.char) (s *lib.NH_ETF_현재가_조회_기본_정보) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tc1151OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH_ETF_현재가_조회_기본_정보)
	s.M종목코드 = lib.F2문자열(g.Code)
	s.M종목명 = lib.F2문자열_CP949(g.Title)
	s.M현재가 = lib.F2정수64_단순형(g.MarketPrice)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = lib.F2정수64_단순형(g.Diff)
	s.M등락율 = lib.F2실수_소숫점_추가(g.DiffRate, 2)
	//s.M매도_호가,_ = lib.F2정수64_단순형(g.OfferPrice)
	//s.M매수_호가,_ = lib.F2정수64_단순형(g.BidPrice)
	s.M거래량 = lib.F2정수64_단순형(g.Volume)
	s.M전일대비_거래량_비율 = lib.F2실수_소숫점_추가(g.TrVolRate, 2)
	s.M유동주_회전율 = lib.F2실수_소숫점_추가(g.FloatVolRate, 2)
	s.M거래대금_100만 = lib.F2정수64_단순형(g.TrAmount)
	s.M상한가 = lib.F2정수64_단순형(g.UpLmtPrice)
	s.M고가 = lib.F2정수64_단순형(g.High)
	s.M시가 = lib.F2정수64_단순형(g.Open)
	s.M시가대비_등락부호 = f2등락부호(g.VsOpenSign)
	s.M시가대비_등락폭 = lib.F2정수64_단순형(g.VsOpenDiff)
	s.M저가 = lib.F2정수64_단순형(g.Low)
	s.M하한가 = lib.F2정수64_단순형(g.LowLmtPrice)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)

	매도_호가_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.OfferPrice1,
			g.OfferPrice2,
			g.OfferPrice3,
			g.OfferPrice4,
			g.OfferPrice5,
			g.OfferPrice6,
			g.OfferPrice7,
			g.OfferPrice8,
			g.OfferPrice9,
			g.OfferPrice10})

	매도_잔량_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.OfferVolume1,
			g.OfferVolume2,
			g.OfferVolume3,
			g.OfferVolume4,
			g.OfferVolume5,
			g.OfferVolume6,
			g.OfferVolume7,
			g.OfferVolume8,
			g.OfferVolume9,
			g.OfferVolume10})

	s.M매도_호가_모음 = make([]int64, 0)
	s.M매도_잔량_모음 = make([]int64, 0)

	for i, 매도_호가 := range 매도_호가_모음 {
		if 매도_호가 >= s.M현재가 && 매도_잔량_모음[i] > 0 {
			s.M매도_호가_모음 = append(s.M매도_호가_모음, 매도_호가)
			s.M매도_잔량_모음 = append(s.M매도_잔량_모음, 매도_잔량_모음[i])
		}
	}

	매수_호가_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.BidPrice1,
			g.BidPrice2,
			g.BidPrice3,
			g.BidPrice4,
			g.BidPrice5,
			g.BidPrice6,
			g.BidPrice7,
			g.BidPrice8,
			g.BidPrice9,
			g.BidPrice10})

	매수_잔량_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.BidVolume1,
			g.BidVolume2,
			g.BidVolume3,
			g.BidVolume4,
			g.BidVolume5,
			g.BidVolume6,
			g.BidVolume7,
			g.BidVolume8,
			g.BidVolume9,
			g.BidVolume10})

	s.M매수_호가_모음 = make([]int64, 0)
	s.M매수_잔량_모음 = make([]int64, 0)

	for i, 매수_호가 := range 매수_호가_모음 {
		if 매수_호가 >= s.M현재가 && 매수_잔량_모음[i] > 0 {
			s.M매수_호가_모음 = append(s.M매수_호가_모음, 매수_호가)
			s.M매수_잔량_모음 = append(s.M매수_잔량_모음, 매수_잔량_모음[i])
		}
	}

	//	s.M매도_잔량_총합,_ = lib.F2정수64_단순형(g.OfferVolTot)
	//	s.M매수_잔량_총합,_ = lib.F2정수64_단순형(g.BidVolTot)
	s.M시간외_매도_잔량 = lib.F2정수64_단순형(g.OfferVolAfterHour)
	s.M시간외_매수_잔량 = lib.F2정수64_단순형(g.BidVolAfterHour)
	s.M피봇_2차_저항 = lib.F2정수64_단순형(g.PivotUp2)
	s.M피봇_1차_저항 = lib.F2정수64_단순형(g.PivotUp1)
	s.M피봇_가격 = lib.F2정수64_단순형(g.PivotPrice)
	s.M피봇_1차_지지 = lib.F2정수64_단순형(g.PivotDown1)
	s.M피봇_2차_지지 = lib.F2정수64_단순형(g.PivotDown2)
	s.M시장_구분 = lib.F2문자열_CP949(g.Market)
	s.M업종명 = lib.F2문자열_CP949(g.Sector)
	s.M자본금_규모 = strings.TrimSpace(lib.F2문자열_CP949(g.CapSize))
	s.M결산월 = lib.F2문자열_CP949(g.SettleMonth)
	s.M추가_정보_모음 = []string{
		lib.F2문자열_CP949(g.MarketAction1),
		lib.F2문자열_CP949(g.MarketAction2),
		lib.F2문자열_CP949(g.MarketAction3),
		lib.F2문자열_CP949(g.MarketAction4),
		lib.F2문자열_CP949(g.MarketAction5),
		lib.F2문자열_CP949(g.MarketAction6)}
	s.M서킷_브레이커_구분 = strings.TrimSpace(lib.F2문자열_CP949(g.CircuitBreaker))
	s.M액면가 = lib.F2정수64_단순형(g.NominalPrice)
	//s.M전일_종가_타이틀 = lib.F2문자열_CP949(g.PrevPriceTitle)
	s.M전일_종가 = lib.F2정수64_단순형(g.PrevPrice)
	s.M대용가 = lib.F2정수64_단순형(g.MortgageValue)
	s.M공모가 = lib.F2정수64_단순형(g.PublicOfferPrice)
	s.M5일_고가 = lib.F2정수64_단순형(g.High5Day)
	s.M5일_저가 = lib.F2정수64_단순형(g.Low5Day)
	s.M20일_고가 = lib.F2정수64_단순형(g.High20Day)
	s.M20일_저가 = lib.F2정수64_단순형(g.Low20Day)
	s.M52주_고가 = lib.F2정수64_단순형(g.High1Year)
	s.M52주_저가 = lib.F2정수64_단순형(g.Low1Year)
	s.M유동_주식수_1000주 = lib.F2정수64_단순형(g.FloatVolume)
	//s.M상장_주식수_1000주,_ = lib.F2정수64_단순형(g.ListVolBy1000)
	s.M시가_총액_억 = lib.F2정수64_단순형(g.MarketCapital)
	s.M거래원_정보_수신_시각 = lib.F2한국증시_개장일_시각_단순형("15:04", g.TraderInfoTime)

	s.M매도_거래원_모음 = []string{
		lib.F2문자열_CP949(g.Seller1),
		lib.F2문자열_CP949(g.Seller2),
		lib.F2문자열_CP949(g.Seller3),
		lib.F2문자열_CP949(g.Seller4),
		lib.F2문자열_CP949(g.Seller5)}

	s.M매도_거래량_모음 = lib.F2정수64_모음_단순형(
		[]interface{}{
			g.Seller1Volume,
			g.Seller2Volume,
			g.Seller3Volume,
			g.Seller4Volume,
			g.Seller5Volume})

	s.M매수_거래원_모음 = []string{
		lib.F2문자열_CP949(g.Buyer1),
		lib.F2문자열_CP949(g.Buyer2),
		lib.F2문자열_CP949(g.Buyer3),
		lib.F2문자열_CP949(g.Buyer4),
		lib.F2문자열_CP949(g.Buyer5)}

	s.M매수_거래량_모음 = lib.F2정수64_모음_단순형(
		[]interface{}{
			g.Buyer1Volume,
			g.Buyer2Volume,
			g.Buyer3Volume,
			g.Buyer4Volume,
			g.Buyer5Volume})

	s.M외국인_매도_거래량 = lib.F2정수64_단순형(g.ForeignSellVolume)
	s.M외국인_매수_거래량 = lib.F2정수64_단순형(g.ForeignBuyVolume)
	s.M외국인_시간 = lib.F2한국증시_개장일_시각_단순형("15:04", g.ForeignTime)
	s.M외국인_지분율 = lib.F2실수_소숫점_추가(g.ForeignHoldingRate, 2)
	s.M신용잔고율 = lib.F2실수_소숫점_추가(g.DebtPercent, 2)

	/* lib.F메모("유상 기준일에 왜 연도가 없는 거지? 매년 하는 것이란 말인가? 그럼 유상배정일이 아니잖아!!!")
	if strings.TrimSpace(lib.F2문자열(g.RightsIssueDate)) != "" {
		lib.F문자열_출력("유상 기준일 : '%v'", lib.F2문자열(g.RightsIssueDate))
	}

	일자, 에러 := lib.F2포맷된_시각_단순형("0102", g.RightsIssueDate)
	if 에러 == nil {
		s.M유상_기준일 = time.Date(지금.Year(), 일자.Month(), 일자.Day(),
			0, 0, 0, 0, 지금.Location())
	}

	if strings.TrimSpace(lib.F2문자열(g.BonusIssueDate)) != "" {
		lib.F문자열_출력("무상 기준일 : '%v'", lib.F2문자열(g.BonusIssueDate))
	}

	일자, 에러 = lib.F2포맷된_시각_단순형("0102", g.BonusIssueDate)
	if 에러 == nil {
		s.M무상_기준일 = time.Date(지금.Year(), 일자.Month(), 일자.Day(),
			0, 0, 0, 0, 지금.Location())
	} */

	s.M유상_배정_비율 = lib.F2실수_소숫점_추가(g.RightsIssueRate, 2)
	s.M무상_배정_비율 = lib.F2실수_소숫점_추가(g.BonusIssueRate, 2)
	s.M상장일 = lib.F2포맷된_시각_단순형("20060102", g.IpoDate)
	s.M상장_주식수 = lib.F2정수64_단순형(g.ListedVolume)
	s.M전체_거래원_매도_합계 = lib.F2정수64_단순형(g.SellTotalSum)
	s.M전체_거래원_매수_합계 = lib.F2정수64_단순형(g.BuyTotalSum)

	return s
}

func New_ETF_현재가_조회_변동_거래_정보(c *C.Tc1151OutBlock2) (s *lib.NH_ETF_현재가_조회_변동_거래량_정보) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tc1151OutBlock2)(unsafe.Pointer(c))
	s = new(lib.NH_ETF_현재가_조회_변동_거래량_정보)

	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)
	s.M현재가 = lib.F2정수64_단순형(g.MarketPrice)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = lib.F2정수64_단순형(g.Diff)
	s.M매도_호가 = lib.F2정수64_단순형(g.OfferPrice)
	s.M매수_호가 = lib.F2정수64_단순형(g.BidPrice)
	s.M변동_거래량 = lib.F2정수64_단순형(g.DiffVolume)
	s.M거래량 = lib.F2정수64_단순형(g.Volume)

	if s.M매도_호가 < s.M현재가 && s.M매도_호가 > 0 {
		s.M매도_호가 = s.M현재가
	}

	if s.M매수_호가 > s.M현재가 && s.M현재가 > 0 {
		s.M매수_호가 = s.M현재가
	}

	return s
}

func New_ETF_현재가_조회_동시호가(c *C.char) (s *lib.NH_ETF_현재가_조회_동시호가_정보) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tc1151OutBlock3)(unsafe.Pointer(c))
	s = new(lib.NH_ETF_현재가_조회_동시호가_정보)
	s.M동시호가_구분 = uint8(lib.F2정수_단순형(g.SyncOfferBid))
	s.M예상_체결가 = lib.F2정수64_단순형(g.EstmPrice)
	s.M예상_체결_부호 = f2등락부호(g.EstmSign)
	s.M예상_등락폭 = lib.F2정수64_단순형(g.EstmDiff)
	s.M예상_등락율 = lib.F2실수_소숫점_추가(g.EstmDiffRate, 2)
	s.M예상_체결_수량 = lib.F2정수64_단순형(g.EstmVolume)

	return s
}

func New_ETF_현재가_조회_ETF자료(c *C.char) (s *lib.NH_ETF_현재가_조회_ETF정보) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tc1151OutBlock4)(unsafe.Pointer(c))
	s = new(lib.NH_ETF_현재가_조회_ETF정보)

	switch lib.F2문자열(g.ETF) {
	case "4":
		s.ETF구분 = P코스닥
	case "8":
		s.ETF구분 = P코스피
	default:
		lib.New에러("예상치 못한 구분값. %v", lib.F2문자열(g.ETF))
	}

	s.NAV = lib.F2실수_소숫점_추가(g.NAV, 2)
	s.NAV등락부호 = f2등락부호(g.DiffSign)
	s.NAV등락폭 = lib.F2실수_소숫점_추가(g.Diff, 2)
	s.M전일NAV = lib.F2실수_소숫점_추가(g.PrevNAV, 2)
	s.M괴리율 = lib.F2실수_소숫점_추가(g.DivergeRate, 2)
	s.M괴리율_부호 = f2등락부호(g.DivergeSign)
	s.M설정단위_당_현금_배당액 = lib.F2실수_단순형(g.DividendPerCU)
	s.M구성_종목수 = lib.F2정수64_단순형(g.ConstituentNo)
	s.M순자산_총액_억 = lib.F2정수64_단순형(g.NAVBy100Million)
	s.M추적_오차율 = lib.F2실수_소숫점_추가(g.TrackingErrRate, 2)

	매도_잔량_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.LP_OfferVolume1,
			g.LP_OfferVolume2,
			g.LP_OfferVolume3,
			g.LP_OfferVolume4,
			g.LP_OfferVolume5,
			g.LP_OfferVolume6,
			g.LP_OfferVolume7,
			g.LP_OfferVolume8,
			g.LP_OfferVolume9,
			g.LP_OfferVolume10})

	s.LP_매도_잔량_모음 = make([]int64, 0)

	for _, 매도_잔량 := range 매도_잔량_모음 {
		// 수량이 0인 주문은 없는 것이나 마찬가지이니 걸러낸다.
		if 매도_잔량 == 0 {
			continue
		}

		s.LP_매도_잔량_모음 = append(s.LP_매도_잔량_모음, 매도_잔량)
	}

	매수_잔량_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.LP_BidVolume1,
			g.LP_BidVolume2,
			g.LP_BidVolume3,
			g.LP_BidVolume4,
			g.LP_BidVolume5,
			g.LP_BidVolume6,
			g.LP_BidVolume7,
			g.LP_BidVolume8,
			g.LP_BidVolume9,
			g.LP_BidVolume10})

	s.LP_매수_잔량_모음 = make([]int64, 0)

	for _, 매수_잔량 := range 매수_잔량_모음 {
		// 수량이 0인 주문은 없는 것이나 마찬가지이니 걸러낸다.
		if 매수_잔량 == 0 {
			continue
		}

		s.LP_매수_잔량_모음 = append(s.LP_매수_잔량_모음, 매수_잔량)
	}

	s.ETF_복제_방법_구분_코드 = lib.F2문자열_CP949(g.TrackingMethod)
	s.ETF_상품_유형_코드 = strings.TrimSpace(lib.F2문자열_CP949(g.ETF_Type))

	return s
}

func New_ETF_현재가_조회_지수_정보(c *C.char) (s *lib.NH_ETF_현재가_조회_지수_정보) {
	g := (*Tc1151OutBlock5)(unsafe.Pointer(c))
	s = new(lib.NH_ETF_현재가_조회_지수_정보)
	s.M업종코드 = strings.TrimSpace(lib.F2문자열(g.SectorCode))
	s.KRX지수_코드 = strings.TrimSpace(lib.F2문자열(g.IndexCode))
	s.M지수_이름 = strings.TrimSpace(lib.F2문자열_CP949(g.IndexName))
	s.M지수 = lib.F2실수_소숫점_추가(g.KP200Index, 2)
	s.M지수_등락부호 = f2등락부호(g.KP200Sign)
	s.M지수_등락폭 = lib.F2실수_소숫점_추가(g.KP200Diff, 2)
	s.M채권지수 = lib.F2실수_소숫점_추가(g.BondIndex, 2)
	s.M채권지수_등락부호 = f2등락부호(g.BondSign)
	s.M채권지수_등락폭 = lib.F2실수_소숫점_추가(g.BondDiff, 2)
	s.M해외지수_코드 = strings.TrimSpace(lib.F2문자열(g.ForeignIndexSymbol))
	s.M기타_업종_코드 = strings.TrimSpace(lib.F2문자열(g.EtcSectorCode))
	s.M채권지수_코드 = strings.TrimSpace(lib.F2문자열(g.BondIndexCode))
	s.M채권지수_세부_코드 = strings.TrimSpace(lib.F2문자열(g.BondDetailCode))

	return s
}

//----------------------------------------------------------------------//
// 코스피/코스닥 호가 잔량 (h1/k3)
//----------------------------------------------------------------------//
func NewNH호가_잔량(c *C.char) (s *lib.NH호가_잔량) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Th1k3OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH호가_잔량)
	s.M종목코드 = lib.F2문자열(g.Code)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)

	매도_호가_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.OfferPrice1,
			g.OfferPrice2,
			g.OfferPrice3,
			g.OfferPrice4,
			g.OfferPrice5,
			g.OfferPrice6,
			g.OfferPrice7,
			g.OfferPrice8,
			g.OfferPrice9,
			g.OfferPrice10})

	매도_잔량_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.OfferVolume1,
			g.OfferVolume2,
			g.OfferVolume3,
			g.OfferVolume4,
			g.OfferVolume5,
			g.OfferVolume6,
			g.OfferVolume7,
			g.OfferVolume8,
			g.OfferVolume9,
			g.OfferVolume10})

	s.M매도_호가_모음 = make([]int64, 0)
	s.M매도_잔량_모음 = make([]int64, 0)

	for i, 매도_잔량 := range 매도_잔량_모음 {
		매도_호가 := 매도_호가_모음[i]

		if 매도_잔량 == 0 || 매도_호가 == 0 {
			continue
		}

		s.M매도_호가_모음 = append(s.M매도_호가_모음, 매도_호가)
		s.M매도_잔량_모음 = append(s.M매도_잔량_모음, 매도_잔량)
	}

	매수_호가_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.BidPrice1,
			g.BidPrice2,
			g.BidPrice3,
			g.BidPrice4,
			g.BidPrice5,
			g.BidPrice6,
			g.BidPrice7,
			g.BidPrice8,
			g.BidPrice9,
			g.BidPrice10})

	매수_잔량_모음 := lib.F2정수64_모음_단순형(
		[]interface{}{
			g.BidVolume1,
			g.BidVolume2,
			g.BidVolume3,
			g.BidVolume4,
			g.BidVolume5,
			g.BidVolume6,
			g.BidVolume7,
			g.BidVolume8,
			g.BidVolume9,
			g.BidVolume10})

	s.M매수_호가_모음 = make([]int64, 0)
	s.M매수_잔량_모음 = make([]int64, 0)

	for i, 매수_잔량 := range 매수_잔량_모음 {
		매수_호가 := 매수_호가_모음[i]

		if 매수_잔량 == 0 || 매수_호가 == 0 {
			continue
		}

		s.M매수_호가_모음 = append(s.M매수_호가_모음, 매수_호가)
		s.M매수_잔량_모음 = append(s.M매수_잔량_모음, 매수_잔량)
	}

	s.M누적_거래량 = lib.F2정수64_단순형(g.Volume)

	return s
}

//----------------------------------------------------------------------//
// 시간외 호가 잔량 (코스피 h2, 코스닥 k4)
//----------------------------------------------------------------------//
func NewNH시간외_호가잔량(c *C.char) (s *lib.NH시간외_호가잔량) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Th2k4OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH시간외_호가잔량)
	s.M종목코드 = lib.F2문자열(g.Code)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)
	s.M총_매도호가_잔량 = lib.F2정수64_단순형(g.OfferVolume)
	s.M총_매수호가_잔량 = lib.F2정수64_단순형(g.BidVolume)

	return s
}

//----------------------------------------------------------------------//
// 예상 호가 잔량 (코스피 h3, 코스닥 k5)
//----------------------------------------------------------------------//
func NewNH예상_호가잔량(c *C.char) (s *lib.NH예상_호가잔량) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Th3k5OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH예상_호가잔량)
	s.M종목코드 = lib.F2문자열(g.Code)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)
	s.M동시호가_구분 = uint8(lib.F2정수_단순형(g.SyncOfferBid))
	s.M예상_체결가 = lib.F2정수64_단순형(g.EstmPrice)
	s.M예상_등락부호 = f2등락부호(g.EstmDiffSign)
	s.M예상_등락폭 = lib.F2정수64_단순형(g.EstmDiff)
	s.M예상_등락율 = lib.F2실수_소숫점_추가(g.EstmDiffRate, 2)
	s.M예상_체결수량 = lib.F2정수64_단순형(g.EstmVolume)
	s.M매도_호가 = lib.F2정수64_단순형(g.OfferPrice)
	s.M매수_호가 = lib.F2정수64_단순형(g.BidPrice)
	s.M매도_호가잔량 = lib.F2정수64_단순형(g.OfferVolume)
	s.M매수_호가잔량 = lib.F2정수64_단순형(g.BidVolume)

	return s
}

//----------------------------------------------------------------------//
// 체결 (코스피 j8, 코스닥 k8)
//----------------------------------------------------------------------//
func NewNH체결_코스피(c *C.char) (s *lib.NH체결) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tj8OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH체결)
	s.M종목코드 = lib.F2문자열(g.Code)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = lib.F2정수64_단순형(g.Diff)
	s.M현재가 = lib.F2정수64_단순형(g.MarketPrice)
	s.M등락율 = lib.F2실수_소숫점_추가(g.DiffRate, 2)
	s.M고가 = lib.F2정수64_단순형(g.High)
	s.M저가 = lib.F2정수64_단순형(g.Low)
	s.M매도_호가 = lib.F2정수64_단순형(g.OfferPrice)
	s.M매수_호가 = lib.F2정수64_단순형(g.BidPrice)
	s.M누적_거래량 = lib.F2정수64_단순형(g.Volume)
	s.M전일대비_거래량_비율 = lib.F2실수_소숫점_추가(g.VsPrevVolRate, 2)
	s.M변동_거래량 = lib.F2정수64_단순형(g.DiffVolume)
	s.M거래대금_100만 = lib.F2정수64_단순형(g.TrAmount)
	s.M시가 = lib.F2정수64_단순형(g.Open)
	s.M가중_평균_가격 = lib.F2정수64_단순형(g.WeightAvgPrice)

	switch lib.F2문자열(g.Market) {
	case "0":
		s.M시장구분 = lib.P시장구분_코스피
	case "1":
		s.M시장구분 = lib.P시장구분_코스닥
	default:
		panic(lib.New에러("예상하지 못한 값. %v", lib.F2문자열(g.Market)))
	}

	return s
}

func NewNH체결_코스닥(c *C.char) (s *lib.NH체결) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tk8OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH체결)
	s.M종목코드 = lib.F2문자열(g.Code)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = lib.F2정수64_단순형(g.Diff)
	s.M현재가 = lib.F2정수64_단순형(g.MarketPrice)
	s.M등락율 = lib.F2실수_소숫점_추가(g.DiffRate, 2)
	s.M고가 = lib.F2정수64_단순형(g.High)
	s.M저가 = lib.F2정수64_단순형(g.Low)
	s.M매도_호가 = lib.F2정수64_단순형(g.OfferPrice)
	s.M매수_호가 = lib.F2정수64_단순형(g.BidPrice)
	s.M누적_거래량 = lib.F2정수64_단순형(g.Volume)
	s.M전일대비_거래량_비율 = lib.F2실수_소숫점_추가(g.VsPrevVolRate, 2)
	s.M변동_거래량 = lib.F2정수64_단순형(g.DiffVolume)
	s.M거래대금_100만 = lib.F2정수64_단순형(g.TrAmount)
	s.M시가 = lib.F2정수64_단순형(g.Open)
	s.M가중_평균_가격 = lib.F2정수64_단순형(g.WeightAvgPrice)

	switch lib.F2문자열(g.Market) {
	case "0":
		s.M시장구분 = lib.P시장구분_코스피
	case "1":
		s.M시장구분 = lib.P시장구분_코스닥
	default:
		panic(lib.New에러("예상하지 못한 값. %v", lib.F2문자열(g.Market)))
	}

	return s
}

//--------------------------------------------//
// ETF NAV (코스피 j1, 코스닥j0)
//--------------------------------------------//
func NewNH_ETF_NAV(c *C.char) (s *lib.NH_ETF_NAV) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tj0j1OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH_ETF_NAV)
	s.M종목코드 = lib.F2문자열(g.Code)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = lib.F2실수_소숫점_추가(g.Diff, 2)
	s.M현재가_NAV = lib.F2실수_소숫점_추가(g.Current, 2)
	s.M시가_NAV = lib.F2실수_소숫점_추가(g.Open, 2)
	s.M고가_NAV = lib.F2실수_소숫점_추가(g.High, 2)
	s.M저가_NAV = lib.F2실수_소숫점_추가(g.Low, 2)
	s.M추적오차_부호 = f2등락부호(g.TrackErrSign)
	s.M추적오차 = lib.F2실수_소숫점_추가(g.TrackingError, 2)
	s.M괴리율_부호 = f2등락부호(g.DivergeSign)
	s.M괴리율 = lib.F2실수_소숫점_추가(g.DivergeRate, 2)

	return s
}

//--------------------------------------------//
// 업종 지수 (코스피 u1, 코스닥 k1)
//--------------------------------------------//
func NewNH업종_지수(c *C.char) (s *lib.NH업종지수) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tu1k1OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH업종지수)
	s.M업종코드 = lib.F2문자열(g.SectorCode)
	s.M시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.Time)
	s.M현재값 = lib.F2실수_소숫점_추가(g.IndexValue, 2)
	s.M등락부호 = f2등락부호(g.DiffSign)
	s.M등락폭 = lib.F2실수_소숫점_추가(g.Diff, 2)
	s.M거래량 = lib.F2정수64_단순형(g.Volume)
	s.M거래대금 = lib.F2정수64_단순형(g.TrAmount)
	s.M개장값 = lib.F2실수_소숫점_추가(g.Open, 2)
	s.M최고값 = lib.F2실수_소숫점_추가(g.High, 2)
	s.M최고값_시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.HighTime)
	s.M최저값 = lib.F2실수_소숫점_추가(g.Low, 2)
	s.M최저값_시각 = lib.F2한국증시_개장일_시각_단순형("15:04:05", g.LowTime)
	s.M지수_등락율 = lib.F2실수_소숫점_추가(g.DiffRate, 2)
	s.M거래비중 = lib.F2실수_소숫점_추가(g.TrVolRate, 2)

	return s
}

// 업종 지수 테이블 이후에도 소스코드 있음.
/* 코스피/코스닥 업종코드 참고표
코스피 업종명			코스닥 업종명
00 	KRX 100			01 	코스닥지수
01 	코스피지수			03 	기타서비스
02 	대형주			04 	코스닥 IT
03 	중형주			06 	제조
04 	소형주			07 	건설
05 	음식료품			08 	유통
06 	섬유의복			10 	운송
07 	종이목재			11 	금융
08 	화학				12 	통신방송서비스
09 	의약품			13 	IT S/W & SVC
10 	비금속광물			14 IT H/W
11 	철강금속			15 	음식료담배
12 	기계				16 	섬유의류
13 	전기전자			17 	종이목재
14 	의료정밀			18 	출판매체복제
15 	운수장비			19 	화학
16 	유통업			20 	제약
17 	전기가스업			21 	비금속
18 	건설업			22 	금속
19 	운수창고			23 	기계장비
20 	통신업			24 	일반전기전자
21 	금융업			25 	의료정밀기기
22 	은행				26 	운송장비부품
24 	증권				27 	기타 제조
25 	보험				28 	통신서비스
26 	서비스업			29 	방송서비스
27 	제조업			30 	인터넷
28 	코스피 200		31 	디지털컨텐츠
29 	코스피 100		32 	소프트웨어
30 	코스피 50			33 	컴퓨터서비스
32 	코스피 배당		34 	통신장비
39 	KP200 건설기계		35 	정보기기
40 	KP200 조선운송		36 	반도체
41 	KP200 철강소재		37 	IT부품
42 	KP200 에너지화학	38 	KOSDAQ 100
43 	KP200 정보기술		39 	KOSDAQ MID 300
44 	KP200 금융		40 	KOSDAQ SMALL
45 	KP200 생활소비재	43 	코스닥 스타
46 	KP200 경기소비재	44 	오락문화
47 	동일가중 KP200		45 	프리미어
48 	동일가중 KP100		46 	우량기업부
49 	동일가중 KP50		47 	벤처기업부
　	　				48 	중견기업부
　	　				49 	기술성장기업부  */

// 주문 및 체결 관련 기능 개발 잠정 보류.
// NH투자증권에서 기술지원에 난색 표시.

//----------------------------------------------------------------------//
// 주식 매도(c8101)
//----------------------------------------------------------------------//
func NewTc8101InBlock(주문 *lib.S질의값_정상주문) unsafe.Pointer {
	switch {
	case 주문.M증권사 != lib.P증권사_NH:
		lib.F패닉("잘못된 증권사. '%v'", 주문.M증권사)
	case 주문.M매수_매도 != lib.P매도:
		lib.F패닉("매도 주문이 아님. %v", 주문.M매수_매도.String())
	}

	g := new(Tc8101InBlock)
	g.Pswd_noz8 = f계좌_비밀번호()
	lib.F바이트_복사_문자열(g.Issue_codez6[:], 주문.M종목코드)
	lib.F바이트_복사_정수(g.Order_qtyz12[:], 주문.M주문수량)
	lib.F바이트_복사_정수(g.Order_unit_pricez10[:], 주문.M주문단가)
	lib.F바이트_복사_문자열(g.Trade_typez2[:], string(f2NH주문유형(주문.M호가유형, 주문.M주문조건)))
	lib.F바이트_복사_문자열(g.Shsll_pos_flagz1[:], "0") // Open API는 공매도 안 됨.
	g.Trad_pswd_no_1z8 = f거래_비밀번호()
	g.Trad_pswd_no_2z8 = f거래_비밀번호()

	return unsafe.Pointer(g)
}

func New매도주문_응답(c *C.char) (s *lib.NH주식_정상주문_응답) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tc8101OutBlock)(unsafe.Pointer(c))

	주문번호, 에러 := lib.F2정수64(g.Noz10)
	if 에러 != nil {
		주문번호 = -1
	}

	s = new(lib.NH주식_정상주문_응답)
	s.M주문번호 = 주문번호
	s.M주문_수량 = lib.F2정수64_단순형(g.Qtyz12)
	s.M주문_단가 = lib.F2정수64_단순형(g.Unit_pricez10)
	s.M매수_매도 = lib.P매도

	return s
}

//----------------------------------------------------------------------//
// 주식 매수(c8102)
//----------------------------------------------------------------------//
func NewTc8102InBlock(주문 *lib.S질의값_정상주문) unsafe.Pointer {
	switch {
	case 주문.M증권사 != lib.P증권사_NH:
		lib.F패닉("잘못된 증권사. '%v'", 주문.M증권사)
	case 주문.M매수_매도 != lib.P매수:
		lib.F패닉("매수 주문이 아님.")
	}

	g := new(Tc8102InBlock)
	g.Pswd_noz8 = f계좌_비밀번호()
	lib.F바이트_복사_문자열(g.Issue_codez6[:], 주문.M종목코드)
	lib.F바이트_복사_정수(g.Order_qtyz12[:], 주문.M주문수량)
	lib.F바이트_복사_정수(g.Order_unit_pricez10[:], 주문.M주문단가)
	lib.F바이트_복사_문자열(g.Trade_typez2[:], string(f2NH주문유형(주문.M호가유형, 주문.M주문조건)))
	g.Trad_pswd_no_1z8 = f거래_비밀번호()
	g.Trad_pswd_no_2z8 = f거래_비밀번호()

	return unsafe.Pointer(g)
}

func New매수주문_응답(c *C.char) (s *lib.NH주식_정상주문_응답) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tc8102OutBlock)(unsafe.Pointer(c))

	주문번호, 에러 := lib.F2정수64(g.Noz10)
	if 에러 != nil {
		주문번호 = -1
	}

	s = new(lib.NH주식_정상주문_응답)
	s.M주문번호 = 주문번호
	s.M주문_수량 = lib.F2정수64_단순형(g.Qtyz12)
	s.M주문_단가 = lib.F2정수64_단순형(g.Unit_pricez10)
	s.M매수_매도 = lib.P매수

	return s
}

// ----------------------------------------------------------------------//
// 주식 정정 주문 (c8103)
// ----------------------------------------------------------------------//
func NewTc8103InBlock(주문 *lib.S질의값_정정주문_NH) unsafe.Pointer {
	if 주문.M증권사 != lib.P증권사_NH {
		lib.F패닉("잘못된 증권사. %v", 주문.M증권사)
	}

	주문수량 := 주문.M주문수량
	if 주문.M잔량_일부 == lib.P잔량 {
		주문수량 = 0
	}

	g := new(Tc8103InBlock)
	g.Pswd_noz8 = f계좌_비밀번호()
	lib.F바이트_복사_문자열(g.Issue_codez6[:], 주문.M종목코드)
	lib.F바이트_복사_정수(g.Crctn_qtyz12[:], 주문수량)
	lib.F바이트_복사_정수(g.Crctn_pricez10[:], 주문.M주문단가)
	lib.F바이트_복사_정수(g.Orgnl_order_noz10[:], 주문.M원주문번호)
	lib.F바이트_복사_문자열(g.All_part_typez1[:], f2NH정정구분(주문.M잔량_일부))
	g.Trad_pswd_no_1z8 = f거래_비밀번호()
	g.Trad_pswd_no_2z8 = f거래_비밀번호()

	return unsafe.Pointer(g)
}

func New정정주문_응답(c *C.char) (주문번호 int64) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		주문번호 = -1
		lib.New에러(r)
	}})

	g := (*Tc8103OutBlock)(unsafe.Pointer(c))
	주문번호, 에러 := lib.F2정수64(g.Order_noz10)
	lib.F에러2패닉(에러)

	return 주문번호
}

// ----------------------------------------------------------------------//
// 주식 취소 주문 (c8104)
// ----------------------------------------------------------------------//
func NewTc8104InBlock(주문 *lib.S질의값_취소주문_NH) unsafe.Pointer {
	if 주문.M증권사 != lib.P증권사_NH {
		lib.F패닉("잘못된 증권사. %v", 주문.M증권사)
	}

	주문수량 := 주문.M주문수량
	if 주문.M잔량_일부 == lib.P잔량 {
		주문수량 = 0
	}

	g := new(Tc8104InBlock)
	g.Pswd_noz8 = f계좌_비밀번호()
	lib.F바이트_복사_문자열(g.Issue_codez6[:], 주문.M종목코드)
	lib.F바이트_복사_정수(g.Canc_qtyz12[:], 주문수량)
	lib.F바이트_복사_정수(g.Orgnl_order_noz10[:], 주문.M원주문번호)
	lib.F바이트_복사_문자열(g.All_part_typez1[:], f2NH취소구분(주문.M잔량_일부))
	g.Trad_pswd_no_1z8 = f거래_비밀번호()
	g.Trad_pswd_no_2z8 = f거래_비밀번호()

	return unsafe.Pointer(g)
}

func New취소주문_응답(c *C.char) (s *lib.NH주식_취소주문_응답) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		s = nil
		lib.New에러(r)
	}})

	g := (*Tc8104OutBlock)(unsafe.Pointer(c))

	모주문번호, 에러 := lib.F2정수64(g.Mom_order_noz10)
	if 에러 != nil {
		모주문번호 = -1
	}

	원주문번호, 에러 := lib.F2정수64(g.Orgnl_order_noz10)
	if 에러 != nil {
		원주문번호 = -1
	}

	주문번호, 에러 := lib.F2정수64(g.Order_noz10)
	if 에러 != nil {
		주문번호 = -1
	}

	s = new(lib.NH주식_취소주문_응답)
	s.M모주문번호 = 모주문번호
	s.M원주문번호 = 원주문번호
	s.M주문번호 = 주문번호
	s.M종목코드 = lib.F2문자열(g.Issue_codez6)
	s.M수량 = lib.F2정수64_단순형(g.Canc_qtyz12)

	return s
}

//----------------------------------------------------------------------//
// 주식 잔고 조회 (c8201)
//----------------------------------------------------------------------//
/*
func NewC8201InBlock(잔고_구분 string) *C.char {
	g := new(Tc8201InBlock)
	g.AccountPassword = f계좌_비밀번호()
	lib.F바이트_복사_문자열(g.Type[:], 잔고_구분)

	return (*C.char)(unsafe.Pointer(g))
}

func New주식잔고_개요(c *C.char) (s *NH주식잔고_개요) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리(M함수with패닉내역:
		func(r interface{}) {
			s = nil
			lib.F에러(r)
		}})

	g := (*Tc8201OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH주식잔고_개요)
	s.M예수금 = lib.F2정수64_단순형(g.Deposits)
	s.M신용_융자금 = lib.F2정수64_단순형(g.CreditLoans)
	s.M이자_미납금 = lib.F2정수64_단순형(g.InterestToBePayed)
	s.M출금_가능_금액 = lib.F2정수64_단순형(g.Withdrawable)
	s.M현금_증거금 = lib.F2정수64_단순형(g.CashMargin)
	s.M대용_증거금 = lib.F2정수64_단순형(g.SecuritiesMargin)
	s.M담보율 = lib.F2실수_단순형(g.PledgeRate)
	s.M현금_미수금 = lib.F2정수64_단순형(g.Receivable)
	s.M주문_가능_금액 = lib.F2정수64_단순형(g.OrderLimit)
	s.M미상환금 = lib.F2정수64_단순형(g.NonPerformingLoans)
	s.M기타_대여금 = lib.F2정수64_단순형(g.EtcLoans)
	s.M대용_금액 = lib.F2정수64_단순형(g.CollateralValue)
	s.M대주_담보금 = lib.F2정수64_단순형(g.LendingPledge)
	s.M매입_원가 = lib.F2정수64_단순형(g.PurchaseCost)
	s.M평가_금액 = lib.F2정수64_단순형(g.AssessedValue)
	s.M자산_합계 = lib.F2정수64_단순형(g.TotalAssets)
	s.M활동_유형 = lib.F2문자열(g.Type)
	s.M대출금 = lib.F2정수64_단순형(g.Loans)
	s.M계좌_증거금율 = lib.F2실수_단순형(g.AccountMarginRate)
	s.M매도_증거금 = lib.F2정수64_단순형(g.SellMargin)
	s.M주문_가능_금액_20퍼센트 = lib.F2정수64_단순형(g.OrderAmountLimit20)
	s.M주문_가능_금액_30퍼센트 = lib.F2정수64_단순형(g.OrderAmountLimit30)
	s.M주문_가능_금액_40퍼센트 = lib.F2정수64_단순형(g.OrderAmountLimit40)
	s.M주문_가능_금액_100퍼센트 = lib.F2정수64_단순형(g.OrderAmountLimit100)
	s.M예수금_D1 = lib.F2정수64_단순형(g.DepositD1)
	s.M예수금_D2 = lib.F2정수64_단순형(g.DepositD2)
	s.M총평가손익 = lib.F2정수64_단순형(g.Earning)
	s.M수익율 = lib.F2실수_소숫점_추가(g.EarningRate, 2)

	return s
}

func New주식잔고_내역(c *C.char) (s *NH주식잔고_내역) {
	defer lib.F에러패닉_내역(lib.S에러패닉_내역(M함수with패닉내역:
		func() {
			s = nil
			lib.F에러(r)
		}})

	g := (*Tc8201OutBlock1)(unsafe.Pointer(c))
	s = new(lib.NH주식잔고_내역)
	s.M종목코드 = lib.F2문자열(g.Code)
	s.M종목명 = lib.F2문자열(g.Title)
	s.M잔고_유형 = lib.F2문자열(g.Type)
	s.M대출일 = lib.F2포맷된_일자_단순형("20060102", g.LoansDate)
	s.M잔고_수량 = lib.F2정수64_단순형(g.Quantity)
	s.M미결제량 = lib.F2정수64_단순형(g.UnsettledQty)
	s.M평균_매입가 = lib.F2정수64_단순형(g.AvgPurchasePrice)
	s.M현재가 = lib.F2정수64_단순형(g.MarketPrice)
	s.M수익_천 = lib.F2정수64_단순형(g.Earning)
	s.M수익율 = lib.F2실수_소숫점_추가(g.EarningRate, 2)
	s.M매입_자금_유형 = lib.F2문자열(g.FinancesType)
	s.M잔량 = lib.F2정수64_단순형(g.RemainedQty)
	s.M만기일 = lib.F2포맷된_일자_단순형("20060102", g.MaturityDate)
	s.M평가_금액 = lib.F2정수64_단순형(g.AssessedValue)
	정수 := lib.F2정수64_단순형(g.MarginRate)
	s.M종목_증거금_비율 = float64(정수)
	s.M평균_매도가 = lib.F2정수64_단순형(g.AvgSellPrice)
	s.M매도_수익 = lib.F2정수64_단순형(g.AfterSellEarning)

	return s
}
*/

//----------------------------------------------------------------------//
// 주문/체결 조회 (s8120)
//----------------------------------------------------------------------//
/*
func NewS8120InBlock(주문_일자 time.Time, 주문_번호 int64,
	시장_구분, 종목코드, 매체_구분, 체결_구분, 조회_순서, 정렬_구분,
	매수도_구분, 매입_자금_구분, 계좌_구분, 계속_여부, CTS string) *C.char {
	g := new(Ts8120InBlock)
	lib.F바이트_복사_문자열(g.QueryType[:], "3")	// 계좌별 조회
	g.AccountPassword = f계좌_비밀번호()
	lib.F바이트_복사_문자열(g.GroupNo[:], "0000")
	lib.F바이트_복사_문자열(g.MarketType[:], 시장_구분)	// 0:전체, 1:3일주문, 2:장내채권, 3:제3시장, 4:선물옵션, 5:장외단주,  7:주식옵션현물
	lib.F바이트_복사_문자열(g.OrderDate[:], 주문_일자.Format("20060102"))
	lib.F바이트_복사_문자열(g.Code[:], 종목코드)
	lib.F바이트_복사_문자열(g.Channel[:], 매체_구분)	// CC:전체,  AA:영업, BB:온라인
	lib.F바이트_복사_문자열(g.OfferBidMatched[:], 체결_구분)	// 0:전체, 1:미체결, 2:체결
	lib.F바이트_복사_문자열(g.Sequence[:], 조회_순서)	// 0:번호, 1:모주문번호
	lib.F바이트_복사_문자열(g.Sorting[:], 정렬_구분)	// 0:주문번호순, 1:주문번호 역순
	lib.F바이트_복사_문자열(g.BuySell[:], 매수도_구분)		// 1:매도, 2:매수, 3:전매, 4:환매
	lib.F바이트_복사_문자열(g.Finances[:], 매입_자금_구분)	// 0:보통, 1:신용, 2:대출
	lib.F바이트_복사_문자열(g.AccountType[:], 계좌_구분)	// 계좌구분 (0:전체, 시장구분(mkt_slctz1 = '4')일때 0:전체, 1:지수선물옵션, 2:주식옵션)
	lib.F바이트_복사_정수(g.OrderId[:], 주문_번호)		// 주문번호
	lib.F바이트_복사_문자열(g.CTS[:], CTS)	// 이건 뭐지?
	g.TrPassword1 = f거래_비밀번호()
	g.TrPassword2 = f거래_비밀번호()
	lib.F바이트_복사_문자열(g.Last[:], 계속_여부)	// 다음화면이 있는경우:'N', 없는경우:' '

	return (*C.char)(unsafe.Pointer(g))
}

func New주문_체결_개요(c *C.char) (s *NH주문_체결_개요) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역:
		func(r interface{}) {
			s = nil
			lib.F에러(r)
		}})

	g := (*Ts8120OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH주문_체결_개요)
	s.M담당자_성명 = lib.F2문자열_CP949(g.Name)
	s.M지점명 = lib.F2문자열_CP949(g.Branch)
	s.M매수_체결수량 = lib.F2정수64_단순형(g.BuyQty)
	s.M매수_체결금액 = lib.F2정수64_단순형(g.BuyAmount)
	s.M매도_체결수량 = lib.F2정수64_단순형(g.SellQty)
	s.M매도_체결금액 = lib.F2정수64_단순형( g.SellAmount)

	return s
}
*/

//----------------------------------------------------------------------//
// 매도 가능 수량 (p8101)
//----------------------------------------------------------------------//
/*
func NewS8101InBlock(매도_구분 string) *C.char {
	g := new(Tc8101InBlock)
	g.AccountPassword = f계좌_비밀번호()
	lib.F바이트_복사_문자열(g.Type[:], 매도_구분)
	// 1:현금, 2:융자, 3:채권, 4:대주, 5:대출주식
	// 6:융자주식합계, 7:대출주식합계, 8:융자주식 및 대출주식
	// 9:융자주식합계+대출주식합계, A:전체

	return (*C.char)(unsafe.Pointer(g))
} */

//----------------------------------------------------------------------//
// 주문 접수 (d3)
//----------------------------------------------------------------------//
func NewNH주문_접수(c *C.char) (s *lib.NH주문_응답) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수: func() { s = nil }})

	g := (*Td3OutBlock)(unsafe.Pointer(c))

	var 주문응답_구분 lib.T주문응답_구분

	switch lib.F2정수_단순형(g.Ordercd) {
	case 10, 11:
		주문응답_구분 = lib.P주문응답_정상
	case 12:
		주문응답_구분 = lib.P주문응답_정정
	case 13:
		주문응답_구분 = lib.P주문응답_취소
	default:
		lib.F패닉("예상하지 못한 값. '%v'", string(g.Ordercd[:]))
	}

	var 시각 time.Time
	if lib.F테스트_모드_실행_중() {
		// 게시판 질문 : order_time 값이 주문시간이 아니라 '주문번호 - 1' 값으로 나옵니다.
		// 답변 : d3OutBlock의 주문시간 필드 order_time 값이 모의투자 환경에서만 별도로 처리 되고 있습니다.
		//        실환경에서는 HHMMSS 이렇게 6자리로 내려옵니다.
		시각 = time.Now()
	} else {
		lib.F2금일_시각_단순형("150405", g.Order_time)
	}

	s = new(lib.NH주문_응답)
	s.RT코드 = lib.NH_RT주문_접수
	s.M주문번호 = lib.F2정수64_단순형(g.Orderno)
	s.M원주문번호 = lib.F2정수64_단순형(g.Orgordno)
	s.M종목코드 = strings.TrimSpace(lib.F2문자열(g.Issuecd))
	s.M주문응답_구분 = 주문응답_구분
	s.M매수_매도 = f2매수_매도(g.Slbygb)
	s.M수량 = lib.F2정수64_단순형(g.Ordergty)
	s.M가격 = lib.F2정수64_단순형(g.Orderprc)
	s.M시각 = 시각

	return s
}

//----------------------------------------------------------------------//
// 주문 체결 (d2)
//----------------------------------------------------------------------//
func NewNH주문_체결(c *C.char) (s *lib.NH주문_응답) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수: func() { s = nil }})

	g := (*Td2OutBlock)(unsafe.Pointer(c))
	s = new(lib.NH주문_응답)
	s.RT코드 = lib.NH_RT주문_체결
	s.M주문번호 = lib.F2정수64_단순형(g.Orderno)
	s.M원주문번호 = -1
	s.M종목코드 = strings.TrimSpace(lib.F2문자열(g.Issuecd))

	switch lib.F2정수_단순형(g.Rejgb) {
	case 1:
		s.M주문응답_구분 = lib.P주문응답_거부
	case 0:
		switch lib.F2정수_단순형(g.Ucgb) {
		case 0:
			s.M주문응답_구분 = lib.P주문응답_체결
		case 1:
			s.M주문응답_구분 = lib.P주문응답_정정
		case 2:
			s.M주문응답_구분 = lib.P주문응답_취소
		case 3:
			s.M주문응답_구분 = lib.P주문응답_거부
		case 4:
			s.M주문응답_구분 = lib.P주문응답_IOC취소
		case 5:
			s.M주문응답_구분 = lib.P주문응답_FOK취소
		default:
			lib.F패닉("예상하지 못한 정정취소구분값. '%v'", lib.F2문자열(g.Ucgb))
		}
	default:
		lib.F패닉("예상하지 못한 정상거부 구분값. '%v'", lib.F2문자열(g.Rejgb))
	}

	var 대출일 time.Time
	if lib.F2문자열(g.Loan_date) == "00000000" {
		대출일 = time.Time{}
	} else {
		lib.F2포맷된_일자_단순형("??", g.Loan_date)
	}

	s.M매수_매도 = f2매수_매도(g.Slbygb)
	s.M수량 = lib.F2정수64_단순형(g.Concgty)
	s.M가격 = lib.F2정수64_단순형(g.Concprc)
	s.M시각 = lib.F2금일_시각_단순형("150405", g.Conctime)
	s.M신용거래_구분 = f2신용거래_구분(lib.F2정수_단순형(g.Sin_gb))
	s.M대출일 = 대출일
	//ato_ord_tpe_chg  //선물옵션주문유형변경여부

	return s
}
