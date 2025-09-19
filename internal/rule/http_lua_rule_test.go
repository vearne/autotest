package rule

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestHttpLuaRule2(t *testing.T) {
	luaStr := `
function verify(r)
	local json = require "json";
	local person = json.decode(r:body());
	return person.age == 1 and person.name == "John";
end
`
	var resp resty.Response
	resp.SetBody([]byte(`{"age": 10, "name": "John"}`))
	rule := HttpLuaRule{LuaStr: luaStr}
	assert.False(t, rule.Verify(&resp))
}

func TestHttpLuaRule3(t *testing.T) {
	luaStr := `
function verify(r)
	local json = require "json";
	local person = json.decode(r:body());
	print("age:", person.age)
	return person.age == 1 and person.name == "John";
end
`
	var resp resty.Response
	resp.SetBody([]byte(`{"age": 1, "name": "Lily"}`))
	rule := HttpLuaRule{LuaStr: luaStr}
	assert.False(t, rule.Verify(&resp))
}
