package luavm

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	slog "github.com/vearne/simplelog"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

func resetPreloadedState() {
	preloadedFilesMutex.Lock()
	loadedLuaFiles = make(map[string]string)
	preloadedGlobals = nil
	preloadedFilesMutex.Unlock()

	GlobalLuaVMPool = &LuaVMPool{
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
}

func TestLoadPreloadLuaFiles(t *testing.T) {
	resetPreloadedState()
	tmpDir := t.TempDir()

	luaFile := filepath.Join(tmpDir, "utils.lua")
	err := os.WriteFile(luaFile, []byte(`
		function helperFuncA()
			return "a"
		end

		function helperFuncB(v)
			return string.match(v, "[^@]+@[^@]+") ~= nil
		end
	`), 0644)
	assert.NoError(t, err)

	err = LoadPreloadLuaFiles([]string{luaFile})
	assert.NoError(t, err)
	assert.NotEmpty(t, loadedLuaFiles)
	assert.NotEmpty(t, preloadedGlobals)
}

func TestLoadPreloadLuaFiles_SyntaxError(t *testing.T) {
	resetPreloadedState()
	tmpDir := t.TempDir()

	luaFile := filepath.Join(tmpDir, "invalid.lua")
	err := os.WriteFile(luaFile, []byte(`
		function bad()
			this is invalid syntax
		end
	`), 0644)
	assert.NoError(t, err)

	err = LoadPreloadLuaFiles([]string{luaFile})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "syntax error")
}

func TestLoadPreloadLuaFiles_FileNotFound(t *testing.T) {
	err := LoadPreloadLuaFiles([]string{"nonexistent.lua"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read")
}

func TestIsPreloadedGlobal(t *testing.T) {
	resetPreloadedState()
	tmpDir := t.TempDir()

	luaFile := filepath.Join(tmpDir, "utils.lua")
	err := os.WriteFile(luaFile, []byte(`
		function preloadedCustomFunc()
			return "hello"
		end
	`), 0644)
	assert.NoError(t, err)

	err = LoadPreloadLuaFiles([]string{luaFile})
	assert.NoError(t, err)

	assert.True(t, isPreloadedGlobal("preloadedCustomFunc"))
	assert.False(t, isPreloadedGlobal("nonExistent"))
}

func TestPreloadedGlobalsProtection(t *testing.T) {
	resetPreloadedState()
	tmpDir := t.TempDir()

	luaFile := filepath.Join(tmpDir, "utils.lua")
	err := os.WriteFile(luaFile, []byte(`
		function protectedFunc(ts)
			return "protected"
		end
	`), 0644)
	assert.NoError(t, err)

	err = LoadPreloadLuaFiles([]string{luaFile})
	assert.NoError(t, err)

	L := GlobalLuaVMPool.GetLuaState()
	GlobalLuaVMPool.PutLuaState(L)

	L.SetGlobal("temporaryGlobal", lua.LString("temporary"))

	GlobalLuaVMPool.PutLuaState(L)

	L2 := GlobalLuaVMPool.GetLuaState()
	defer GlobalLuaVMPool.PutLuaState(L2)

	preloaded := L2.GetGlobal("protectedFunc")
	assert.NotEqual(t, lua.LNil, preloaded)

	temporary := L2.GetGlobal("temporaryGlobal")
	assert.Equal(t, lua.LNil, temporary)
}

func TestExecuteLuaWithPreloadedFunctions(t *testing.T) {
	resetPreloadedState()
	tmpDir := t.TempDir()

	luaFile := filepath.Join(tmpDir, "utils.lua")
	err := os.WriteFile(luaFile, []byte(`
		function formatTimestamp(ts)
			return "formatted_" .. tostring(ts)
		end

		function validateEmail(email)
			return string.match(email, "[^@]+@[^@]+") ~= nil
		end
	`), 0644)
	assert.NoError(t, err)

	err = LoadPreloadLuaFiles([]string{luaFile})
	assert.NoError(t, err)

	value, err := ExecuteLuaWithGlobalsPool(nil, nil, `
		return formatTimestamp(12345)
	`)
	assert.NoError(t, err)
	assert.NotEqual(t, lua.LNil, value)
	assert.Equal(t, lua.LString("formatted_12345"), value)

	value2, err := ExecuteLuaWithGlobalsPool(nil, nil, `
		return validateEmail("test@example.com")
	`)
	assert.NoError(t, err)
	assert.Equal(t, lua.LTrue, value2)

	value3, err := ExecuteLuaWithGlobalsPool(nil, nil, `
		return validateEmail("invalid-email")
	`)
	assert.NoError(t, err)
	assert.Equal(t, lua.LFalse, value3)
}

func TestLoadPreloadLuaFiles_EmptyList(t *testing.T) {
	err := LoadPreloadLuaFiles([]string{})
	assert.NoError(t, err)
}

func TestIsBuiltinGlobal(t *testing.T) {
	assert.True(t, isBuiltinGlobal("print"))
	assert.True(t, isBuiltinGlobal("json"))
	assert.True(t, isBuiltinGlobal("os"))
	assert.False(t, isBuiltinGlobal("customFunction"))
}
