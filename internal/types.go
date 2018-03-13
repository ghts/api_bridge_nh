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

	"sync"
)

func New대기항목_저장소() *s대기항목_저장소 {
	대기항목_저장소 := new(s대기항목_저장소)
	대기항목_저장소.저장소 = make(map[int64]*lib.S콜백_대기항목)

	return 대기항목_저장소
}

type s대기항목 struct {
	키  int64
	내용 *lib.S콜백_대기항목
}

type s대기항목_저장소 struct {
	sync.RWMutex
	저장소 map[int64]*lib.S콜백_대기항목
}

func (s *s대기항목_저장소) G대기항목(키 int64) *s대기항목 {
	s.RLock()
	defer s.RUnlock()

	내용, 존재함 := s.저장소[키]
	if !존재함 {
		return nil
	}

	대기항목 := new(s대기항목)
	대기항목.키 = 키
	대기항목.내용 = 내용

	return 대기항목
}

func (s *s대기항목_저장소) G대기항목_모음() []*s대기항목 {
	s.RLock()
	defer s.RUnlock()

	대기항목_모음 := make([]*s대기항목, 0)

	for 키, 항목 := range s.저장소 {
		대기항목 := new(s대기항목)
		대기항목.키 = 키
		대기항목.내용 = 항목

		대기항목_모음 = append(대기항목_모음, 대기항목)
	}

	return 대기항목_모음
}

func (s *s대기항목_저장소) S추가(대기항목 *lib.S콜백_대기항목) {
	s.Lock()
	defer s.Unlock()

	s.저장소[대기항목.G식별번호()] = 대기항목
}

func (s *s대기항목_저장소) S삭제(키 int64) {
	s.Lock()
	defer s.Unlock()

	delete(s.저장소, 키)
}


//----------------------------------------------------------------------//
// WMCA 문자 message 구조체
//----------------------------------------------------------------------//
type NH수신_메시지_블록 struct {
	M식별번호   int64
	M메시지_코드 string //00000:정상, 기타:비정상(코드값은 언제든지 변경될 수 있음.)
	M메시지_내용 string
}

//----------------------------------------------------------------------//
// WMCA TR 응답 구조체
//----------------------------------------------------------------------//
type NH수신_데이터_블록 struct {
	M식별번호  int64
	M블록_이름 string
	M데이터   interface{}
}
