package luavm

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// LuaVMPool Lua虚拟机池，避免全局锁竞争
type LuaVMPool struct {
	pool sync.Pool
}

var GlobalLuaVMPool = &LuaVMPool{
	pool: sync.Pool{
		New: func() interface{} {
			L := lua.NewState()
			luajson.Preload(L)
			return L
		},
	},
}

// GetLuaState 从池中获取Lua状态机
func (p *LuaVMPool) GetLuaState() *lua.LState {
	return p.pool.Get().(*lua.LState)
}

// PutLuaState 将Lua状态机归还到池中
func (p *LuaVMPool) PutLuaState(L *lua.LState) {
	// 清理状态机
	L.SetTop(0)
	p.pool.Put(L)
}

// 保持向后兼容
var L *lua.LState
var LuaVMLock sync.Mutex

func init() {
	L = lua.NewState()
	luajson.Preload(L)
}

// ExecuteLuaWithGlobalsPool 使用池化的Lua虚拟机执行脚本
func ExecuteLuaWithGlobalsPool(globals map[string]lua.LValue, source string) (lua.LValue, error) {
	L := GlobalLuaVMPool.GetLuaState()
	defer GlobalLuaVMPool.PutLuaState(L)

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

// 保持向后兼容的原函数
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
