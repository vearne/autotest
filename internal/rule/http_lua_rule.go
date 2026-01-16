package rule

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/vearne/autotest/internal/luavm"
	"github.com/vearne/zaplog"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

const luaHttpRespTypeName = "HttpResp"

type HttpResp struct {
	Code string
	Body string
}

var httpRespMethods = map[string]lua.LGFunction{
	"code": getSetHttpRespCode,
	"body": getSetHttpRespBody,
}

// Getter and setter for the HttpResp#Code
func getSetHttpRespCode(L *lua.LState) int {
	p := checkHttpResp(L)
	if L.GetTop() == 3 {
		p.Code = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(p.Code))
	return 1
}

// Getter and setter for the HttpResp#Body
func getSetHttpRespBody(L *lua.LState) int {
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

type HttpLuaRule struct {
	LuaStr string `json:"lua"`
}

func (r *HttpLuaRule) Name() string {
	return "HttpLuaRule"
}

func (r *HttpLuaRule) Verify(resp *resty.Response) bool {
	globals := map[string]lua.LValue{
		"codeStr": lua.LString(strconv.Itoa(resp.StatusCode())),
		"bodyStr": lua.LString(resp.String()),
	}
	source := r.LuaStr +
		`
r = HttpResp.new(codeStr, bodyStr);
return verify(r);
`
	value, err := luavm.ExecuteLuaWithGlobalsPool(registerHttpRespType, globals, source)
	if err != nil {
		// 1. print in the console
		var b strings.Builder
		fmt.Fprintln(&b, "================== HttpLuaRule-Verify ==================")
		fmt.Fprintln(&b, "status\t:", resp.StatusCode())
		fmt.Fprintln(&b, "body\t:", resp.String())
		fmt.Fprintln(&b, "LuaStr\t:", r.LuaStr)
		fmt.Fprintln(&b, "error\t:", err.Error())
		os.Stderr.WriteString(b.String())

		// 2. output in the log
		zaplog.Error("HttpLuaRule-Verify",
			zap.Int("status", resp.StatusCode()),
			zap.String("body", resp.String()),
			zap.String("LuaStr", r.LuaStr),
			zap.Error(err))
		return false
	}
	return value == lua.LTrue
}
