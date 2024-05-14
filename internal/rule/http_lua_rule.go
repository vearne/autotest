package rule

import (
	"github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

var L *lua.LState

func init() {
	L = lua.NewState()
	//defer L.Close()
	// 注册 JSON 库
	luajson.Preload(L)
	registerHttpRespType(L)
}

type HttpResp struct {
	Code string
	Body string
}

const luaHttpRespTypeName = "HttpResp"

var httpRespMethods = map[string]lua.LGFunction{
	"code": getSetCode,
	"body": getSetBody,
}

// Getter and setter for the HttpResp#Code
func getSetCode(L *lua.LState) int {
	p := checkHttpResp(L)
	if L.GetTop() == 3 {
		p.Code = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(p.Code))
	return 1
}

// Getter and setter for the HttpResp#Body
func getSetBody(L *lua.LState) int {
	p := checkHttpResp(L)
	if L.GetTop() == 3 {
		p.Body = L.CheckString(3)
		return 0
	}
	L.Push(lua.LString(p.Body))
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkHttpResp(L *lua.LState) *HttpResp {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*HttpResp); ok {
		return v
	}
	L.ArgError(1, "HttpResp expected")
	return nil
}

// Constructor
func newHttpResp(L *lua.LState) int {
	resp := &HttpResp{L.CheckString(1), L.CheckString(2)}
	ud := L.NewUserData()
	ud.Value = resp
	L.SetMetatable(ud, L.GetTypeMetatable(luaHttpRespTypeName))
	L.Push(ud)
	return 1
}

// Registers my person type to given L.
func registerHttpRespType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaHttpRespTypeName)
	L.SetGlobal("HttpResp", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newHttpResp))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), httpRespMethods))
}
