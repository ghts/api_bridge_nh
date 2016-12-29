package main

import (
	"github.com/ghts/lib"
	"github.com/ghts/api_bridge_nh/internal"
)

func main() {
	lib.F에러2패닉(internal.F초기화())

	<-lib.F공통_종료_채널()
}
