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

// #cgo CFLAGS: -m32 -Wall
// #include <stdlib.h>
// #include "./c_func.h"
import "C"

import (
	"github.com/ghts/lib"

	"unsafe"
)

func f일반TR_실행(TR식별번호 int64, TR코드 string, c데이터 unsafe.Pointer, 길이, 계좌_인덱스 int) bool {
	cTR식별번호 := C.int(TR식별번호)
	cTR코드 := C.CString(TR코드)
	c길이 := C.int(길이)
	c계좌_인덱스 := C.int(계좌_인덱스)

	defer C.free(unsafe.Pointer(cTR코드))
	//defer C.free(unsafe.Pointer(c데이터))	// 생성한 곳에서 free()하도록 한다.

	반환값 := C.wmcaQuery(cTR식별번호, cTR코드, (*C.char)(c데이터), c길이, c계좌_인덱스)

	return bool(반환값)
}

func f실시간_정보_구독(RT코드 string, 전체_종목코드 string, 단위_길이 int) bool {
	cRT코드 := C.CString(RT코드)
	c코드_모음 := C.CString(전체_종목코드)
	c단위_길이 := C.int(단위_길이)
	c전체_길이 := C.int(len(전체_종목코드))

	defer func() {
		C.free(unsafe.Pointer(cRT코드))
		C.free(unsafe.Pointer(c코드_모음))
	}()

	반환값 := C.wmcaAttach(cRT코드, c코드_모음, c단위_길이, c전체_길이)

	return bool(반환값)
}

func f실시간_정보_해지(타입 string, 전체_종목코드 string, 단위_길이 int) bool {
	c타입 := C.CString(타입)
	c코드_모음 := C.CString(전체_종목코드)
	c단위_길이 := C.int(단위_길이)
	c전체_길이 := C.int(len(전체_종목코드))

	defer func() {
		C.free(unsafe.Pointer(c타입))
		C.free(unsafe.Pointer(c코드_모음))
	}()

	반환값 := C.wmcaDetach(c타입, c코드_모음, c단위_길이, c전체_길이)

	return bool(반환값)
}

func f접속(아이디, 암호, 공인인증서_암호 string) bool {
	f자원_정리()

	c아이디 := C.CString(아이디)
	c암호 := C.CString(암호)
	c공인인증서_암호 := C.CString(공인인증서_암호)
	c웹_공지 := C.CString("NOTICEINTRO")
	c긴급_공지 := C.CString("NOTICEURGENT")
	cN := C.CString("N")
	c서버_이름 := C.CString(lib.F서버명_NH())
	c포트_번호 := C.int(lib.F포트번호_NH())

	defer func() {
		C.free(unsafe.Pointer(c아이디))
		C.free(unsafe.Pointer(c암호))
		C.free(unsafe.Pointer(c공인인증서_암호))
		C.free(unsafe.Pointer(c웹_공지))
		C.free(unsafe.Pointer(c긴급_공지))
		C.free(unsafe.Pointer(cN))
		C.free(unsafe.Pointer(c서버_이름))
	}()

	로드_성공 := bool(C.wmcaLoad())
	if !로드_성공 {
		lib.F문자열_출력("로드 실패")
		return false
	}

	서버_설정_성공 := bool(C.wmcaSetServer(c서버_이름))
	if !서버_설정_성공 {
		lib.F문자열_출력("서버 설정 실패")
		return false
	}

	포트_설정_성공 := bool(C.wmcaSetPort(c포트_번호))
	if !포트_설정_성공 {
		lib.F문자열_출력("포트 설정 실패")
		return false
	}

	//공지(웹)을 사용하지 않음.
	웹_공지_비활성화_성공 := bool(C.wmcaSetOption(c웹_공지, cN))
	if !웹_공지_비활성화_성공 {
		lib.F문자열_출력("웹 공지 비활성화 실패")
		return false
	}

	//긴급공지를  사용하지 않음.
	긴급_공지_비활성화_성공 := bool(C.wmcaSetOption(c긴급_공지, cN))
	if !긴급_공지_비활성화_성공 {
		lib.F문자열_출력("긴급 공지 비활성화 실패")
		return false
	}

	return bool(C.wmcaConnect(c아이디, c암호, c공인인증서_암호))
}

func f접속_해제() bool {
	return f호출("wmcaDisconnect")
}

func f실시간_서비스_모두_해지() bool {
	return f호출("wmcaDetachAll")
}

func f접속됨() bool {
	return f호출("wmcaIsConnected")
}

func f자원_정리() {
	// cgo의 버그로 인해서 인수가 없으면 '사용하지 않는 변수' 컴파일 경고 발생.
	// 컴파일 경고를 없애기 위해서 사용하지 않는 인수를 추가함.
	C.wmcaFreeResource(C.int(1))
}
