package luavm

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestLua(t *testing.T) {
	if err := RunLuaStr(`
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
