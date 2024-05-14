package rule

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestLua(t *testing.T) {
	if err := runLuaStr(`
		function verify(r)
			local json = require "json";
			local car = json.decode(r:body());
			return r:code() == "200" and car.age == 10 and car.name == "buick";
		end
        r = HttpResp.new("200", "{\"age\": 10,\"name\": \"buick\"}")
		return verify(r)
    `); err != nil {
		panic(err)
	}

	result := false
	lv := L.Get(-1)
	if lv == lua.LTrue {
		result = true
	}
	assert.True(t, result)
}

func TestHttpLuaRule(t *testing.T) {
	luaStr := `
function verify(r)
	local json = require "json";
	local person = json.decode(r:body());
	return person.age == 10 and person.name == "John";
end
`
	var resp resty.Response
	resp.SetBody([]byte(`{"age": 10, "name": "John"}`))
	rule := HttpLuaRule{LuaStr: luaStr}
	assert.True(t, rule.Verify(&resp))
}
