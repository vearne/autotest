package resource

import "sync"

var EnvVars map[string]string
var AllVars sync.Map

func init() {
	EnvVars = make(map[string]string, 10)
}
