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
the Free Software Foundation version 2.1 of the License.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. */

#include <windef.h>

//const DWORD ERR_NONE=WM_USER+300;
//const DWORD ERR_DLL_NOT_FOUND=WM_USER+301;
//const DWORD ERR_FUNC_NOT_FOUND=WM_USER+302;

// 여기서부터 아래 부분은 증권사 제공 예제코드를 복사 후 붙여넣기 한 후 약간 수정한 것임.
/* COPIED FROM 'WmcaIntf.h' in NH OpenAPI sample code.
* LICENSING TERM follows that of original code.
*
* NH 오픈API 예제 소스코드를 참조한 후 야간 에서 복사해서 붙여넣기 됨.
* 저작권 관련 규정은 원래 샘플 소스코드의 저작권 규정을 따름.
* (샘플 소스코드에는 저작권 항목을 찾을 수 없었기에
* 자유롭게 사용할 수 있는 Public Domain이 아닐까 추정하지만,
* 그것은 단지 개인적인 추정일 뿐이며 저작권 관련 정확한 사항은
* API를 배포한 증권사 측에 문의해 봐야함.
*/

const DWORD CA_WMCAEVENT		=WM_USER+8400;
const DWORD CA_CONNECTED		=WM_USER+110;	//접속 및 로그인 성공후 수신되며, 서비스 이용이 가능함을 의미합니다.
const DWORD CA_DISCONNECTED		=WM_USER+120;	//통신 연결이 끊겼을 경우 반환되는 메시지입니다.
const DWORD CA_SOCKET_ERROR		=WM_USER+130;	//네트워크 장애등의 이유로 통신 오류 발생할 경우 수신되는 메시지로, 접속환경 점검이 필요합니다.
const DWORD CA_TR_DATA			=WM_USER+210;	//wmcaTransact() 호출에 따른 처리 결과값이 수신됩니다.
const DWORD CA_REALTIME_DATA	=WM_USER+220;	//wmcaAttach() 호출에 따른 실시간 데이터가 수신됩니다.
const DWORD CA_MESSAGE			=WM_USER+230;	//요청한 서비스에 대한 처리상태가 문자열 형태로 수신되며, 정상처리 및 처리실패등의 각 상태를 보여줍니다.
const DWORD CA_COMPLETE			=WM_USER+240;	//요청한 서비스에 대한 처리가 정상 완료될 경우 수신됩니다.
const DWORD CA_ERROR			=WM_USER+250;	//요청한 서비스에 대한 처리가 실패할 경우 수신되며, 사용자가 잘못된 값을 입력하는 등의 이유로 발생합니다.

// Window SDK에 포함되어 있지만, GCC를 사용하기 위해서는 다시 정의해줘야 함.
// these are in the windows SDK, but need to be repeated here for GCC..
#ifndef MSGFLT_ALLOW
typedef struct tagCHANGEFILTERSTRUCT {
	DWORD cbSize;
	DWORD ExtStatus;
} CHANGEFILTERSTRUCT, *PCHANGEFILTERSTRUCT;

typedef BOOL WINAPI ChangeWindowMessageFilterEx(HWND hWnd, UINT message, DWORD action, PCHANGEFILTERSTRUCT pChangeFilterStruct);

const DWORD MSGFLT_ALLOW = 1;
const DWORD MSGFLT_DISALLOW = 2;
const DWORD MSGFLT_RESET = 0;
#endif
