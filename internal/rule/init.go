package rule

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	"os"
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
	// 设置自定义的 print 函数
	L.SetGlobal("print", L.NewFunction(customPrint))
	// register json lib
	luajson.Preload(L)
	registerHttpRespType(L)
	registerGrocRespType(L)
}

func customPrint(L *lua.LState) int {
	top := L.GetTop()
	for i := 1; i <= top; i++ {
		fmt.Fprint(os.Stdout, L.ToString(i)) // 输出到 os.Stdout
		if i != top {
			fmt.Fprint(os.Stdout, "\t")
		}
	}
	fmt.Fprintln(os.Stdout) // 换行
	return 0
}
