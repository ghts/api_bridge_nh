''' Copyright (C) 2015-2018 김운하(UnHa Kim)  unha.kim@kuh.pe.kr

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
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. '''

import random
import sys
import time
import zmq

P메시지_GET = "G"
P메시지_종료 = "Q"
P회신_OK = "O"
UTF8 = 'utf-8'

def 테스트용_종목정보_요청_모듈(P주소_종목정보, P주소_테스트_결과, 테스트_반복횟수):
    #print("종목정보 요청 : 시작")
    
    context = zmq.Context()
    
    종목정보_REQ = context.socket(zmq.REQ)
    종목정보_REQ.connect(P주소_종목정보)
    
    테스트_결과_REQ = context.socket(zmq.REQ)
    테스트_결과_REQ.connect(P주소_테스트_결과)
    
    #print("종목정보 요청 : 초기화 완료")
    
    질의예상값_모음 = []
    질의예상값_모음.append(["000020", "동화약품"])
    질의예상값_모음.append(["000030", "우리은행"])
    질의예상값_모음.append(["000040", "KR모터스"])
    질의예상값_모음.append(["000050", "경방"])
    질의예상값_모음.append(["000060", "메리츠화재"])    
    
    테스트_결과 = True
    
    for 반복횟수 in range(테스트_반복횟수):
        질의예상값 = 질의예상값_모음[random.randint(0, len(질의예상값_모음) - 1)]
        질의값 = 질의예상값[0]
        
        송신_메시지 = [P메시지_GET.encode(UTF8), 질의값.encode(UTF8)]
        
        #print("제공 : send_multipart() 시작")
        종목정보_REQ.send_multipart(송신_메시지)
        #print("제공 : send_multipart() 완료")
        
        #print("제공 : recv_multipart() 시작")
        수신_메시지 = 종목정보_REQ.recv_multipart()
        #print("제공 : recv_multipart() 완료")
        
        구분 = 수신_메시지[0].decode(UTF8)
        
        if 구분 != P회신_OK:
            에러_메시지 = 수신_메시지[1].decode(UTF8)
            print("제공 : 에러 메시지 회신.", 구분, 에러_메시지)
            테스트_결과 = False
            break
        
        종목코드 = 수신_메시지[1].decode(UTF8)
        종목이름 = 수신_메시지[2].decode(UTF8)
        
        if 종목코드 == 질의예상값[0] and 종목이름 == 질의예상값[1]:            
            continue       
        
        print("제공 : 회신값이 예상과 다름.", 질의값, 종목코드, 종목이름)
        테스트_결과 = False
        break
    
    #print("종목정보 요청 : 테스트 결과 전송 시작")
    메시지 = [P메시지_GET.encode(UTF8), str(테스트_결과).encode(UTF8)]
    테스트_결과_REQ.send_multipart(메시지)
    테스트_결과_REQ.recv_multipart()
    #print("종목정보 요청 : 테스트 결과 전송 완료")
    
    # 리소스 정리 후 종료
    테스트_결과_REQ.close()
    context.destroy()
    
    #print("종목정보 요청 : 종료")

if __name__ == "__main__":
    P주소_종목정보 = sys.argv[1]
    P주소_테스트_결과 = sys.argv[2]
    테스트_반복횟수 = int(sys.argv[3])
    테스트용_종목정보_요청_모듈(P주소_종목정보, P주소_테스트_결과, 테스트_반복횟수)