package rule

import (
	"github.com/vearne/autotest/internal/luavm"
	"github.com/vearne/autotest/internal/model"
	"github.com/vearne/zaplog"
	"github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

const luaGrpcRespTypeName = "GrpcResp"

// Registers my person type to given L.
func registerGrocRespType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaGrpcRespTypeName)
	L.SetGlobal("GrpcResp", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newGrpcResp))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), grpcRespMethods))
}

// Constructor
func newGrpcResp(L *lua.LState) int {
	resp := &model.GrpcResp{Code: L.CheckString(1), Body: L.CheckString(2)}
	ud := L.NewUserData()
	ud.Value = resp
	L.SetMetatable(ud, L.GetTypeMetatable(luaGrpcRespTypeName))
	L.Push(ud)
	return 1
}

var grpcRespMethods = map[string]lua.LGFunction{
	"code": getSetGrpcRespCode,
	"body": getSetGrpcRespBody,
}

// Getter and setter for the GrpcResp#Code
func getSetGrpcRespCode(L *lua.LState) int {
	p := checkGrpcResp(L)
	if L.GetTop() == 3 {
		p.Code = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(p.Code))
	return 1
}

// Getter and setter for the GrpcResp#Body
func getSetGrpcRespBody(L *lua.LState) int {
	p := checkGrpcResp(L)
	if L.GetTop() == 3 {
		p.Body = L.CheckString(3)
		return 0
	}
	L.Push(lua.LString(p.Body))
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkGrpcResp(L *lua.LState) *model.GrpcResp {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*model.GrpcResp); ok {
		return v
	}
	L.ArgError(1, "GrpcResp expected")
	return nil
}

type GrpcLuaRule struct {
	LuaStr string `json:"lua"`
}

func (r *GrpcLuaRule) Name() string {
	return "GrpcLuaRule"
}

func (r *GrpcLuaRule) Verify(resp *model.GrpcResp) bool {
	globals := map[string]lua.LValue{
		"codeStr": lua.LString(resp.Code),
		"bodyStr": lua.LString(resp.Body),
	}

	source := r.LuaStr +
		`
	r = GrpcResp.new(codeStr, bodyStr);
	return verify(r);
`
	value, err := luavm.ExecuteLuaWithGlobals(globals, source)
	if err != nil {
		zaplog.Error("GrpcLuaRule-Verify",
			zap.String("code", resp.Code),
			zap.String("body", resp.Body),
			zap.String("LuaStr", r.LuaStr),
			zap.Error(err))
		return false
	}
	return value == lua.LTrue
}
