package luavm

import (
	"fmt"
	"os"
	"sync"

	slog "github.com/vearne/simplelog"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// 定义一个函数类型，用于作为参数传递
type RegisterType func(L *lua.LState)

var (
	loadedLuaFiles      = make(map[string]string)
	preloadedGlobals    []string
	preloadedFilesMutex sync.RWMutex
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

			if len(loadedLuaFiles) > 0 {
				if err := loadPreloadedFilesToState(L); err != nil {
					slog.Error("failed to load preloaded lua files: %v", err)
				}
			}

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
	// 完全清理状态机
	cleanupLuaState(L)
	p.pool.Put(L)
}

// cleanupLuaState 彻底清理Lua状态机
func cleanupLuaState(L *lua.LState) {
	// 1. 清空栈
	L.SetTop(0)

	// 2. 清理全局变量表中的自定义变量
	// 获取全局表
	global := L.GetGlobal("_G")
	if tbl, ok := global.(*lua.LTable); ok {
		// 清理可能的自定义全局变量，保留Lua内置的
		keysToRemove := make([]lua.LValue, 0)
		tbl.ForEach(func(key, value lua.LValue) {
			if keyStr, ok := key.(lua.LString); ok {
				keyName := string(keyStr)
				// 保留Lua内置的全局变量和函数，以及预加载的全局变量
				if !isBuiltinGlobal(keyName) && !isPreloadedGlobal(keyName) {
					keysToRemove = append(keysToRemove, key)
				}
			}
		})

		// 移除自定义全局变量
		for _, key := range keysToRemove {
			tbl.RawSet(key, lua.LNil)
		}
	}

	// 3. 运行垃圾回收
	L.DoString("collectgarbage('collect')") //nolint:errcheck
}

// isBuiltinGlobal 检查是否为Lua内置全局变量
func isBuiltinGlobal(name string) bool {
	builtins := map[string]bool{
		"_G": true, "_VERSION": true, "assert": true, "collectgarbage": true,
		"dofile": true, "error": true, "getmetatable": true, "ipairs": true,
		"load": true, "loadfile": true, "next": true, "pairs": true,
		"pcall": true, "print": true, "rawequal": true, "rawget": true,
		"rawlen": true, "rawset": true, "require": true, "select": true,
		"setmetatable": true, "tonumber": true, "tostring": true, "type": true,
		"xpcall": true, "string": true, "table": true, "math": true,
		"io": true, "os": true, "debug": true, "coroutine": true,
		"utf8": true, "package": true, "json": true, // json是我们预加载的
	}
	return builtins[name]
}

// isPreloadedGlobal 检查是否为预加载的全局变量
func isPreloadedGlobal(name string) bool {
	preloadedFilesMutex.RLock()
	defer preloadedFilesMutex.RUnlock()

	for _, key := range preloadedGlobals {
		if key == name {
			return true
		}
	}
	return false
}

// ExecuteLuaWithGlobalsPool 使用池化的Lua虚拟机执行脚本
func ExecuteLuaWithGlobalsPool(f RegisterType, globals map[string]lua.LValue, source string) (lua.LValue, error) {
	L := GlobalLuaVMPool.GetLuaState()
	defer GlobalLuaVMPool.PutLuaState(L)

	// 定义数据类型和函数
	if f != nil {
		f(L)
	}

	// 记录执行前的栈大小
	stackSizeBefore := L.GetTop()

	// 记录设置的全局变量，用于后续清理
	var globalKeys []string

	// 原子性地设置所有全局变量
	for name, value := range globals {
		L.SetGlobal(name, value)
		globalKeys = append(globalKeys, name)
	}

	// 执行Lua代码
	err := L.DoString(source)
	if err != nil {
		// 出错时清理栈和全局变量
		L.SetTop(stackSizeBefore)
		cleanupGlobals(L, globalKeys)
		return nil, err
	}

	// 获取结果（复制值，避免引用问题）
	result := lua.LNil
	if L.GetTop() > stackSizeBefore {
		result = L.Get(-1)
		// 对于复杂类型，创建副本以避免状态机回收后的引用问题
		if result != lua.LNil {
			result = copyLuaValue(L, result)
		}
	}

	// 清理栈，恢复到执行前的大小
	L.SetTop(stackSizeBefore)

	// 清理本次设置的全局变量
	cleanupGlobals(L, globalKeys)

	return result, nil
}

// cleanupGlobals 清理指定的全局变量
func cleanupGlobals(L *lua.LState, keys []string) {
	for _, key := range keys {
		L.SetGlobal(key, lua.LNil)
	}
}

// copyLuaValue 复制Lua值，避免引用问题
func copyLuaValue(L *lua.LState, value lua.LValue) lua.LValue {
	switch v := value.(type) {
	case lua.LString:
		return v // 字符串是不可变的，可以直接返回
	case lua.LNumber:
		return v // 数字是不可变的，可以直接返回
	case lua.LBool:
		return v // 布尔值是不可变的，可以直接返回
	case *lua.LTable:
		// 对于表，创建一个新的表并复制内容
		newTable := L.NewTable()
		v.ForEach(func(key, val lua.LValue) {
			newTable.RawSet(key, copyLuaValue(L, val))
		})
		return newTable
	default:
		// 对于其他类型，返回字符串表示
		return lua.LString(value.String())
	}
}

// LoadPreloadLuaFiles 加载并验证Lua预加载文件
func LoadPreloadLuaFiles(files []string) error {
	if len(files) == 0 {
		slog.Info("No Lua preload files specified")
		return nil
	}

	slog.Info("Loading %d Lua preload files", len(files))

	// 先加载文件内容到临时map
	tempLoadedFiles := make(map[string]string)
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read lua file %s: %w", file, err)
		}

		tempL := lua.NewState()
		defer tempL.Close()

		if err := tempL.DoString(string(content)); err != nil {
			return fmt.Errorf("syntax error in lua file %s: %w", file, err)
		}

		tempLoadedFiles[file] = string(content)
		slog.Info("Loaded Lua file: %s", file)
	}

	// 跟踪预加载的全局变量（在加锁之前完成，避免死锁）
	if err := trackPreloadedGlobalsFromFiles(tempLoadedFiles); err != nil {
		return err
	}
	slog.Info("preloadedGlobals: %v", preloadedGlobals)

	// 加锁并更新全局变量
	preloadedFilesMutex.Lock()
	for file, content := range tempLoadedFiles {
		loadedLuaFiles[file] = content
	}
	preloadedFilesMutex.Unlock()

	return nil
}

// loadPreloadedFilesToState 将预加载的Lua文件加载到指定的Lua状态机
func loadPreloadedFilesToState(L *lua.LState) error {
	preloadedFilesMutex.RLock()
	defer preloadedFilesMutex.RUnlock()

	for file, content := range loadedLuaFiles {
		if err := L.DoString(content); err != nil {
			return fmt.Errorf("failed to load lua file %s: %w", file, err)
		}
	}
	return nil
}

// trackPreloadedGlobalsFromFiles 跟踪预加载文件定义的全局变量
func trackPreloadedGlobalsFromFiles(files map[string]string) error {
	if len(files) == 0 {
		return nil
	}

	L := lua.NewState()
	defer L.Close()

	luajson.Preload(L)

	// 直接使用传入的文件内容，不需要加锁
	for _, content := range files {
		if err := L.DoString(content); err != nil {
			return fmt.Errorf("failed to load lua file content: %w", err)
		}
	}

	// 加锁更新 preloadedGlobals
	preloadedFilesMutex.Lock()
	defer preloadedFilesMutex.Unlock()

	global := L.GetGlobal("_G")
	if tbl, ok := global.(*lua.LTable); ok {
		tbl.ForEach(func(key, _ lua.LValue) {
			if keyStr, ok := key.(lua.LString); ok {
				keyName := string(keyStr)
				if !isBuiltinGlobal(keyName) {
					// 检查是否已经存在，避免重复添加
					exists := false
					for _, existing := range preloadedGlobals {
						if existing == keyName {
							exists = true
							break
						}
					}
					if !exists {
						preloadedGlobals = append(preloadedGlobals, keyName)
					}
				}
			}
		})
	}

	slog.Info("Tracked %d preloaded global variables", len(preloadedGlobals))
	return nil
}
