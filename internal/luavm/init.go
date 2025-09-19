package luavm

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
}

func RunLuaStr(source string) error {
	LuaVMLock.Lock()
	defer LuaVMLock.Unlock()

	return L.DoString(source)
}

func SetGlobal(name string, value lua.LValue) {
	L.SetGlobal(name, value)
}

func Get(idx int) lua.LValue {
	return L.Get(idx)
}
