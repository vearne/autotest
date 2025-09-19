package rule

import "github.com/vearne/autotest/internal/luavm"

func init() {
	registerHttpRespType(luavm.L)
	registerGrocRespType(luavm.L)
}
