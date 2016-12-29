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

const wmca_dll = "wmca.dll"
const 실행_성공 = "completed successfully"

const (
	P코스피  string = "코스피"
	P코스닥  string = "코스닥"
	P실물복제 string = "실물복제"
	P합성복제 string = "합성복제"
	P일반형  string = "일반형"
	P파생형  string = "파생형"
)

// 주문유형
type NH주문유형 string

const (
	NH주문유형_지정가         NH주문유형 = "00"
	NH주문유형_시장가         NH주문유형 = "03"
	NH주문유형_조건부_지정가     NH주문유형 = "05"
	NH주문유형_최유리_지정가     NH주문유형 = "12"
	NH주문유형_최우선_지정가     NH주문유형 = "13"
	NH주문유형_시간외_단일가     NH주문유형 = "31"
	NH주문유형_장전_시간외_전일종가 NH주문유형 = "61"
	NH주문유형_장후_시간외_금일종가 NH주문유형 = "71"
	NH주문유형_IOC_지정가     NH주문유형 = "C0" // 즉시체결, 잔량취소
	NH주문유형_FOK_지정가     NH주문유형 = "F0" // 즉시전량체결, 전량취소
	NH주문유형_IOC_시장가     NH주문유형 = "C3" // 즉시체결, 잔량취소
	NH주문유형_FOK_시장가     NH주문유형 = "F3" // 즉시전량체결, 전량취소
	NH주문유형_IOC_최유리_지정가 NH주문유형 = "C2" // 즉시체결, 잔량취소
	NH주문유형_FOK_최유리_지정가 NH주문유형 = "F2" // 즉시전량체결, 전량취소
)

// NH 등락부호
const (
	NH상한 uint8 = 0x18 // 24
	NH상승 uint8 = 0x1E // 30
	NH보합 uint8 = 0x20 // 32
	NH하한 uint8 = 0x19 // 25
	NH하락 uint8 = 0x1F // 31
)

// 정정 구분
const (
	NH정정구분_일부 = "1"
	NH정정구분_잔량 = "2"
)

// 취소 구분
const (
	NH취소구분_일부 = "1"
	NH취소구분_잔량 = "2"
)

// 매수 매도 구분
const (
	NH매도 = 1
	NH매수 = 2
)
