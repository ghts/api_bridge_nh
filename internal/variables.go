package internal

import (
	"sync"
	"github.com/ghts/lib"
)

var 계좌_인덱스_맵 = make(map[string]int)
var 계좌_인덱스_맵_잠금 = sync.RWMutex{}
var DLL경로_기본값 = lib.F_GOPATH()  + "/src/github.com/ghts/ghts_dependency/NH_OpenAPI"