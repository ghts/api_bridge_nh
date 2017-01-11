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
	"github.com/go-mangos/mangos"
)

// TR : TRansaction(트랜잭션)의 줄임말.
// RT : RealTime(실시간)의 줄임말.
// PUB소켓은 수신 기능이 없으며, 'non-blocking'방식으로 지연없이 동작.

var 실시간_정보_중계_중 = lib.New안전한_bool(false)

// SUB-PUB기반 실시간 정보 중계
func Go루틴_API_실시간_정보_중계(ch초기화 chan lib.T신호) (에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M에러: &에러})

	if 실시간_정보_중계_중.G값() {
		ch초기화 <- lib.P신호_초기화
		return nil
	} else if 에러 = 실시간_정보_중계_중.S값(true); 에러 != nil {
		ch초기화 <- lib.P신호_초기화
		return 에러
	}

	defer 실시간_정보_중계_중.S값(false)

	소켓PUB_CBOR, 에러 := lib.New소켓PUB(lib.P주소_NH_실시간_CBOR)
	lib.F에러2패닉(에러)
	defer 소켓PUB_CBOR.Close()

	소켓PUB_MsgPack, 에러 := lib.New소켓PUB(lib.P주소_NH_실시간_MsgPack)
	lib.F에러2패닉(에러)
	defer 소켓PUB_MsgPack.Close()

	소켓PUB_JSON, 에러 := lib.New소켓PUB(lib.P주소_NH_실시간_JSON)
	lib.F에러2패닉(에러)
	defer 소켓PUB_JSON.Close()

	ch종료 := lib.F공통_종료_채널()
	ch초기화 <- lib.P신호_초기화

	변환형식_모음 := []lib.T변환{lib.CBOR, lib.MsgPack, lib.JSON}
	소켓PUB_모음 := []mangos.Socket{소켓PUB_CBOR, 소켓PUB_MsgPack, 소켓PUB_JSON}

	for {
		select {
		case 실시간_정보 := <-ch실시간_정보:
			for i := 0; i < 3; i++ {
				소켓_메시지, 에러 := lib.New소켓_메시지(변환형식_모음[i], 실시간_정보)
				if 에러 != nil {
					lib.F에러_출력(에러)
					continue
				}

				if 에러 := 소켓_메시지.S소켓_송신(소켓PUB_모음[i], lib.P30초); 에러 != nil {
					lib.F문자열_출력("PUB소켓 전송 에러 : %v %T %v", 변환형식_모음[i], 실시간_정보, 실시간_정보)
					lib.F에러_출력(에러)
					continue
				}
			}
		case <-ch종료:
			return nil
		default:
			lib.F실행권한_양보()
		}
	}
}

var TR소켓_중계_중 = lib.New안전한_bool(false)

// NH OpenAPI는 TR쿼터가 4초 내 TR 20회 발생임.
// 평균적으로 1초당 5회이며, TR의 타임아웃이 30초이므로,
// 워커는 200개 정도면 충분할 듯 함.
// (어차피 쿼터 제한 초과로 차단될 바에야 프록시 소켓에서 대기하는 게 낫다.)
const TR처리_도우미_수량 = 200

// go-mangos의 raw모드를 사용하여 병렬처리.
func Go루틴_소켓TR_중계(ch초기화 chan lib.T신호) (에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M에러: &에러})

	if TR소켓_중계_중.G값() {
		ch초기화 <- lib.P신호_초기화
		return nil
	} else if 에러 = TR소켓_중계_중.S값(true); 에러 != nil {
		ch초기화 <- lib.P신호_초기화
		return 에러
	}

	defer TR소켓_중계_중.S값(false)

	// 병렬처리를 위해서 raw모드를 사용함. raw모드 사용법은 go-mangos 예제를 참고.
	소켓REP, 에러 := lib.New소켓REP_raw(lib.P주소_NH_TR)
	ch도우미_초기화 := make(chan lib.T신호, TR처리_도우미_수량)
	ch도우미_종료 := make(chan error, TR처리_도우미_수량)

	for i := 0; i < TR처리_도우미_수량; i++ {
		go f소켓TR_처리(소켓REP, ch도우미_초기화, ch도우미_종료)
	}

	for i := 0; i < TR처리_도우미_수량; i++ {
		if 신호 := <-ch도우미_초기화; 신호 != lib.P신호_초기화 {
			lib.F패닉("%v번째 TR처리 워커 소켓 초기화 실패.", i)
		}
	}

	ch종료 := lib.F공통_종료_채널()
	ch초기화 <- lib.P신호_초기화 // 초기화 완료.

	for {
		select {
		case 에러 = <-ch도우미_종료:
			lib.F에러_출력(에러)

			// 새로운 worker함수 인스턴스 실행.
			select {
			case <-ch종료:
				continue
			default:
				go f소켓TR_처리(소켓REP, ch도우미_초기화, ch도우미_종료)
				if 신호 := <-ch도우미_초기화; 신호 != lib.P신호_초기화 {
					lib.F패닉("TR처리 워커 소켓 초기화 실패.")
				}
			}
		case <-ch종료:
			return nil
		default:
			lib.F실행권한_양보()
		}
	}
}

// 소켓으로 받은 TR을 실제로 처리하는 worker함수
// 여러 개의 worker함수가 동시에 존재하므로, 처리시간으로 인한 지연에 신경쓰지 않아도 됨.
// 에러나 패닉이 발생하지 않는 한 종료할 필요없음.
func f소켓TR_처리(소켓REP mangos.Socket, ch초기화 chan<- lib.T신호, ch종료 chan<- error) {
	var 에러 error
	var raw메시지 *mangos.Message // 최대한 재활용 해야 성능 문제를 걱정할 필요가 없어진다.
	var 회신_메시지 lib.I소켓_메시지

	defer func() { ch종료 <- 에러 }()

	defer lib.F에러패닉_처리(lib.S에러패닉_처리{M함수with패닉내역: func(r interface{}) {
		if raw메시지 != nil {
			메시지, 에러 := lib.New소켓_메시지_에러(r)
			if 에러 == nil {
				메시지.S소켓_회신(소켓REP, lib.P10초, raw메시지)
			}
		}
	}})

	ch공통종료 := lib.F공통_종료_채널()
	ch초기화 <- lib.P신호_초기화

	for {
		raw메시지, 에러 = 소켓REP.RecvMsg()
		lib.F에러2패닉(에러)
		defer raw메시지.Free() // GC부담을 덜어주고, 재활용을 통한 성능 향상을 꾀함.

		바이트_모음 := raw메시지.Body
		수신_메시지 := lib.New소켓_메시지from바이트_모음(바이트_모음)
		lib.F에러2패닉(수신_메시지.G에러())

		회신_메시지, 에러 = f소켓TR_처리_도우미(수신_메시지)
		lib.F에러2패닉(에러)

		if 회신_메시지 != nil {
			에러 = 회신_메시지.S소켓_회신(소켓REP, lib.P30초, raw메시지)
			lib.F에러2패닉(에러)
		}

		select {
		case <-ch공통종료:
			return
		default: // OK
		}
	}
}

func f소켓TR_처리_도우미(수신_메시지 lib.I소켓_메시지) (응답_메시지 lib.I소켓_메시지, 에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{
		M에러:         &에러,
		M함수with패닉내역: func(r interface{}) { 응답_메시지, _ = lib.New소켓_메시지_에러(r) }})

	lib.F조건부_패닉(수신_메시지.G길이() != 1, "잘못된 메시지 길이. 예상값 1, 실제값 %v. %v", 수신_메시지.G길이(), 수신_메시지)

	TR질의값, 에러 := lib.F소켓_메시지_해석(수신_메시지, 0)
	lib.F에러2패닉(에러)

	변환_형식 := 수신_메시지.G변환_형식(0)
	TR구분 := TR질의값.(lib.I질의값).G_TR구분()
	TR코드 := TR질의값.(lib.I질의값).G_TR코드()

	// TR처리
	switch TR구분 {
	case lib.TR조회, lib.TR주문:
		TR응답_모음 := make([]*lib.S바이트_변환_매개체, 0)
		채널_질의 := lib.New채널_질의(ch조회, lib.P30초, 10).S질의(수신_메시지)

		// 데이터가 여러 번에 나누어서 수신되며, 수신된 데이터를 취합해서 1개의 회신 메시지를 생성함.
		for {
			응답 := 채널_질의.G응답()
			lib.F에러2패닉(응답.G에러())

			TR응답_구분, ok := 응답.G값(0).(lib.TR응답_구분)
			lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 응답.G값(0))

			switch TR응답_구분 {
			case lib.TR응답_메시지:
				lib.F문자열_출력("%v : %v", 응답.G값(1), 응답.G값(2))
				lib.F메모("메시지가 에러발생 통보인지 여부는 이후에 다시 생각 해보자.")
			case lib.TR응답_데이터:
				TR응답, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
				lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 응답.G값(1))

				TR응답_모음 = append(TR응답_모음, TR응답)
			case lib.TR응답_완료:
				TR응답, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
				lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 응답.G값(1))

				if !TR응답.IsNil() {
					TR응답_모음 = append(TR응답_모음, TR응답)
				}
			default:
				lib.F문자열_출력("예상하지 못한 경우.")
				lib.F변수값_확인(응답.G값_모음()...)
			}

			// 모든 데이터를 수신했으므로, 응답 메시지를 반환하고 종료.
			// 데이터를 모두 수신했는 지 확인할 때는 데이터가 순서대로 수신된다는 보장이 없기에,
			// 완료 메시지 수신 '이후'에도 데이터 추가 수신 여부를 확인해야 함.
			if TR응답_수신완료(TR코드, TR응답_모음) {
				return f일반TR_응답_메시지_생성(변환_형식, TR질의값, TR응답_모음)
			}
		}
	case lib.TR실시간_정보_구독, lib.TR실시간_정보_해지:
		채널_질의 := lib.New채널_질의(ch실시간_정보_구독_및_해지, lib.P30초, 10).S질의(수신_메시지)

		for {
			응답 := 채널_질의.G응답()
			lib.F에러2패닉(응답.G에러())

			TR응답_구분, ok := 응답.G값(0).(lib.TR응답_구분)
			lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 응답.G값(0))

			switch TR응답_구분 {
			case lib.TR응답_완료:
				return lib.New소켓_메시지(변환_형식, lib.TR응답_완료)
			default:
				lib.F문자열_출력("응답값 확인 후 적절하게 처리하도록 할 것.")
				lib.F변수값_확인(응답.G값_모음()...)
			}
		}
	case lib.TR접속:
		채널_질의 := lib.New채널_질의(ch접속, lib.P30초, 10).S질의(수신_메시지)

		for {
			응답 := 채널_질의.G응답()
			lib.F에러2패닉(응답.G에러())

			TR구분, ok := 응답.G값(0).(lib.TR응답_구분)
			lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 응답.G값(0))

			switch TR구분 {
			case lib.TR응답_메시지: // 인증 실패  메시지
				lib.F패닉("%v : %v", 응답.G값(1), 응답.G값(2))
			case lib.TR응답_완료:
				switch 응답.G길이() {
				case 1: // 이미 접속되어 있음. 로그인 정보가 없는 데 괜찮을려나?
					return lib.New소켓_메시지(변환_형식, lib.TR응답_완료)
				case 2: // OK
					변환값, ok := 응답.G값(1).(*lib.S바이트_변환_매개체)
					lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 응답.G값(1))

					로그인_정보 := new(lib.NH로그인_정보)
					lib.F에러2패닉(변환값.G값(로그인_정보))

					return lib.New소켓_메시지(변환_형식, 로그인_정보)
				default:
					lib.F패닉("예상하지 못한 길이. %v", 응답.G길이())
				}
			default:
				lib.F문자열_출력("응답값 확인 후 적절하게 처리하도록 할 것.")
				lib.F변수값_확인(응답.G값_모음()...)
			}
		}
	case lib.TR접속_해제:
		채널_질의 := lib.New채널_질의(ch접속_해제, lib.P30초, 10).S질의(수신_메시지)

		for {
			응답 := 채널_질의.G응답()
			lib.F에러2패닉(응답.G에러())
			lib.F조건부_패닉(응답.G길이() != 1, "예상하지 못한 길이. %v", 응답.G길이())

			TR응답_구분, ok := 응답.G값(0).(lib.TR응답_구분)
			lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 응답.G값(0))

			switch TR응답_구분 {
			case lib.TR응답_완료:
				return lib.New소켓_메시지(변환_형식, lib.TR응답_완료)
			default:
				lib.F문자열_출력("응답값 확인 후 적절하게 처리하도록 할 것.")
				lib.F변수값_확인(응답.G값_모음()...)
			}
		}
	case lib.TR실시간_정보_일괄_해지:
		채널_질의 := lib.New채널_질의(ch실시간_정보_일괄_해지, lib.P30초, 10).S질의(수신_메시지)

		for {
			응답 := 채널_질의.G응답()
			lib.F에러2패닉(응답.G에러())
			lib.F조건부_패닉(응답.G길이() != 1, "예상하지 못한 길이. %v", 응답.G길이())

			TR응답_구분, ok := 응답.G값(0).(lib.TR응답_구분)
			lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 응답.G값(0))

			switch TR응답_구분 {
			case lib.TR응답_완료:
				return lib.New소켓_메시지(변환_형식, lib.TR응답_완료)
			default:
				lib.F문자열_출력("응답값 확인 후 적절하게 처리하도록 할 것.")
				lib.F변수값_확인(응답.G값_모음()...)
			}
		}
	case lib.TR접속됨:
		응답 := lib.New채널_질의(ch접속됨, lib.P30초, 10).S질의(수신_메시지).G응답()
		lib.F에러2패닉(응답.G에러())
		lib.F조건부_패닉(응답.G길이() != 1, "예상하지 못한 길이. 예상값 1, 실제값 %v", 응답.G길이())

		접속_여부, ok := 응답.G값(0).(bool)
		lib.F조건부_패닉(!ok, "예상하지 못한 자료형. %T", 응답.G값(0))

		return lib.New소켓_메시지(변환_형식, 접속_여부)
	case lib.TR종료:
		close(lib.F공통_종료_채널())
		return nil, nil // 회신 하지 않음.
	default:
		lib.F패닉("예상하지 못한 TR구분. %v", TR구분)
	}

	return
}

var (
	nh주식_현재가_조회_기본   = lib.F자료형_문자열(lib.NH주식_현재가_조회_기본_정보{})
	nh주식_현재가_조회_변동   = lib.F자료형_문자열(make([]*lib.NH주식_현재가_조회_변동_거래량_정보, 0))
	nh주식_현재가_조회_동시호가 = lib.F자료형_문자열(lib.NH주식_현재가_조회_동시호가_정보{})

	nhETF_현재가_조회_기본   = lib.F자료형_문자열(lib.NH_ETF_현재가_조회_기본_정보{})
	nhETF_현재가_조회_변동   = lib.F자료형_문자열(make([]*lib.NH_ETF_현재가_조회_변동_거래량_정보, 0))
	nhETF_현재가_조회_동시호가 = lib.F자료형_문자열(lib.NH_ETF_현재가_조회_동시호가_정보{})
	nhETF_현재가_조회_ETF  = lib.F자료형_문자열(lib.NH_ETF_현재가_조회_ETF정보{})
	nhETF_현재가_조회_지수   = lib.F자료형_문자열(lib.NH_ETF_현재가_조회_지수_정보{})

	nh호가_잔량     = lib.F자료형_문자열(lib.NH호가_잔량{})
	nh시간외_호가_잔량 = lib.F자료형_문자열(lib.NH시간외_호가잔량{})
	nh예상_호가_잔량  = lib.F자료형_문자열(lib.NH예상_호가잔량{})
	nh체결        = lib.F자료형_문자열(lib.NH체결{})
	nhETF_NAV   = lib.F자료형_문자열(lib.NH_ETF_NAV{})
	nh업종_지수     = lib.F자료형_문자열(lib.NH업종지수{})
)

func TR응답_수신완료(TR코드 string, TR응답_모음 []*lib.S바이트_변환_매개체) bool {
	switch TR코드 {
	case lib.NH_TR주식_현재가_조회:
		수신완료_기본정보, 수신완료_변동정보, 수신완료_동시호가정보 := false, false, false

		for _, TR응답 := range TR응답_모음 {
			switch TR응답.G자료형_문자열() {
			case nh주식_현재가_조회_기본:
				수신완료_기본정보 = true
			case nh주식_현재가_조회_변동:
				수신완료_변동정보 = true
			case nh주식_현재가_조회_동시호가:
				수신완료_동시호가정보 = true
			default:
				lib.F패닉("예상하지 못한 경우. %v", TR응답.G자료형_문자열())
			}
		}

		return 수신완료_기본정보 && 수신완료_변동정보 && 수신완료_동시호가정보
	case lib.NH_TR_ETF_현재가_조회:
		수신완료_기본정보, 수신완료_변동정보, 수신완료_동시호가정보 := false, false, false
		수신완료_ETF정보, 수신완료_지수정보 := false, false

		for _, TR응답 := range TR응답_모음 {
			switch TR응답.G자료형_문자열() {
			case nhETF_현재가_조회_기본:
				수신완료_기본정보 = true
			case nhETF_현재가_조회_변동:
				수신완료_변동정보 = true
			case nhETF_현재가_조회_동시호가:
				수신완료_동시호가정보 = true
			case nhETF_현재가_조회_ETF:
				수신완료_ETF정보 = true
			case nhETF_현재가_조회_지수:
				수신완료_지수정보 = true
			default:
				lib.F패닉("예상하지 못한 경우. %v", TR응답.G자료형_문자열())
			}
		}

		return 수신완료_기본정보 && 수신완료_변동정보 && 수신완료_동시호가정보 &&
			수신완료_ETF정보 && 수신완료_지수정보
	default:
		lib.F패닉("예상하지 못한 TR코드. %v", TR코드)
	}

	return false
}

func f일반TR_응답_메시지_생성(변환_형식 lib.T변환, 질의값 interface{},
	응답_정보_모음 []*lib.S바이트_변환_매개체) (응답_메시지 lib.I소켓_메시지, 에러 error) {
	defer lib.F에러패닉_처리(lib.S에러패닉_처리{
		M에러:         &에러,
		M함수with패닉내역: func(r interface{}) { 응답_메시지, _ = lib.New소켓_메시지_에러(r) }})

	TR코드 := 질의값.(lib.I질의값).G_TR코드()

	switch TR코드 {
	case lib.NH_TR주식_현재가_조회:
		응답값 := lib.NewNH주식_현재가_조회_응답()
		응답값.M질의 = 질의값.(*lib.S질의값_단일종목)

		for _, 단일값_바이트_변환 := range 응답_정보_모음 {
			switch 단일값_바이트_변환.G자료형_문자열() {
			case nh주식_현재가_조회_기본:
				에러 = 단일값_바이트_변환.G값(응답값.M기본_정보)
			case nh주식_현재가_조회_변동:
				에러 = 단일값_바이트_변환.G값(&응답값.M변동_거래량_정보)
			case nh주식_현재가_조회_동시호가:
				에러 = 단일값_바이트_변환.G값(응답값.M동시호가_정보)
			default:
				에러 = lib.New에러("에상하지 못한 자료형. %v", 단일값_바이트_변환.G자료형_문자열())
			}

			lib.F에러2패닉(에러)
		}

		lib.F조건부_패닉(응답값.M질의 == nil || 응답값.M기본_정보 == nil ||
			응답값.M변동_거래량_정보 == nil || 응답값.M동시호가_정보 == nil,
			"응답값에 nil 필드 존재함. %v", 응답값)

		return lib.New소켓_메시지(변환_형식, 응답값)
	case lib.NH_TR_ETF_현재가_조회:
		응답값 := lib.NewNH_ETF_현재가_조회_응답()
		응답값.M질의 = 질의값.(*lib.S질의값_단일종목)

		for _, 변환값 := range 응답_정보_모음 {
			switch 변환값.G자료형_문자열() {
			case nhETF_현재가_조회_기본:
				에러 = 변환값.G값(응답값.M기본_정보)
			case nhETF_현재가_조회_변동:
				에러 = 변환값.G값(&응답값.M변동_거래량_정보)
			case nhETF_현재가_조회_동시호가:
				에러 = 변환값.G값(응답값.M동시호가_정보)
			case nhETF_현재가_조회_ETF:
				에러 = 변환값.G값(응답값.ETF_정보)
			case nhETF_현재가_조회_지수:
				에러 = 변환값.G값(응답값.M지수_정보)
			default:
				에러 = lib.New에러("예상하지 못한 자료형. %v", 변환값.G자료형_문자열())
			}

			lib.F에러2패닉(에러)
		}

		lib.F조건부_패닉(응답값.M질의 == nil || 응답값.M기본_정보 == nil ||
			응답값.M변동_거래량_정보 == nil || 응답값.M동시호가_정보 == nil ||
			응답값.ETF_정보 == nil || 응답값.M지수_정보 == nil,
			"응답값에 nil 필드 존재함. %v", 응답값)

		return lib.New소켓_메시지(변환_형식, 응답값)
	default:
		lib.F패닉("예상하지 못한 TR코드. %v", TR코드)
	}

	return
}
