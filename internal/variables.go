package internal

import (
	"sync"
)

var 계좌_인덱스_맵 = make(map[string]int)
var 계좌_인덱스_맵_잠금 = sync.RWMutex{}
