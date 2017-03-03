package main

import (
	"github.com/ghts/api_bridge_nh/internal"
	"github.com/ghts/lib"
)

func main() {
	lib.F테스트_모드_시작() // 이것을 파라메터로 조정할 수 있도록 할 것.

	lib.F에러2패닉(internal.F초기화())
	lib.F문자열_출력("초기화 완료")

	<-lib.F공통_종료_채널()
}
