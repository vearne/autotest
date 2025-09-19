package rule

import (
	"github.com/stretchr/testify/assert"
	"github.com/vearne/autotest/internal/luavm"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestLua(t *testing.T) {
	source := `
		function verify(r)
			local json = require "json";
			local car = json.decode(r:body());
			return r:code() == "200" and car.age == 10 and car.name == "buick";
		end
        r = HttpResp.new("200", "{\"age\": 10,\"name\": \"buick\"}")
		return verify(r)
    `
	value, err := luavm.ExecuteLuaWithGlobals(nil, source)
	if err != nil {
		panic(err)
	}

	result := false
	if value == lua.LTrue {
		result = true
	}
	assert.True(t, result)
}
