package rule

import (
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	"sync"
)

var L *lua.LState

/*
Lua virtual machine is not thread-safe
so we need LuaVMLock to protect L
*/
var LuaVMLock sync.Mutex

func init() {
	L = lua.NewState()
	//defer L.Close()
	// register json lib
	luajson.Preload(L)
	registerHttpRespType(L)
	registerGrocRespType(L)
}
