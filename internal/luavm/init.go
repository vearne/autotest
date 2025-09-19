package luavm

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

var L *lua.LState

/*
Lua virtual machine is not thread-safe
so we need LuaVMLock to protect L
*/
var LuaVMLock sync.Mutex

func init() {
	L = lua.NewState()
	luajson.Preload(L)
}

func ExecuteLuaWithGlobals(globals map[string]lua.LValue, source string) (lua.LValue, error) {
	LuaVMLock.Lock()
	defer LuaVMLock.Unlock()

	// 记录执行前的栈大小
	stackSizeBefore := L.GetTop()

	// 原子性地设置所有全局变量
	for name, value := range globals {
		L.SetGlobal(name, value)
	}

	// 执行Lua代码
	err := L.DoString(source)
	if err != nil {
		// 出错时也要清理栈
		L.SetTop(stackSizeBefore)
		return nil, err
	}

	// 获取结果
	result := L.Get(-1)

	// 清理栈，恢复到执行前的大小
	L.SetTop(stackSizeBefore)

	return result, nil
}
