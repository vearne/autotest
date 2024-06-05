package rule

import (
	"github.com/stretchr/testify/assert"
	"github.com/vearne/autotest/internal/model"
	"testing"
)

func TestGrpcLuaRule(t *testing.T) {
	luaStr := `
function verify(r)
	local json = require "json";
	local person = json.decode(r:body());
	return person.age == 10 and person.name == "John";
end
`
	var resp model.GrpcResp
	resp.Body = `{"age": 10, "name": "John"}`
	rule := GrpcLuaRule{LuaStr: luaStr}
	assert.True(t, rule.Verify(&resp))
}
