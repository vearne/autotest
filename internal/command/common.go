package command

import (
	"embed"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/flosch/pongo2/v6"
	"github.com/spf13/cast"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/resource"
	"strings"
)

//go:embed template/*.tpl
var mytpl embed.FS

type ResultInfo struct {
	Total        int
	SuccessCount int
	FailedCount  int
}

type CaseShow struct {
	ID          uint64
	Description string
	State       string
	Reason      string
	Link        string
}

func templateRender(tplStr string) (string, error) {
	// Compile the template first (i. e. creating the AST)
	tpl, err := pongo2.FromString(tplStr)
	if err != nil {
		return "", err
	}

	kvs := make(map[string]any)
	for key, value := range resource.EnvVars {
		kvs[key] = value
	}

	resource.CustomerVars.Range(func(key, value any) bool {
		kvs[key.(string)] = value
		return true
	})

	// Now you can render the template with the given
	// pongo2.Context how often you want to.
	return tpl.Execute(pongo2.Context(kvs))
}

func exportTo(jsonStr string, export *config.Export) (any, error) {
	doc, err := jsonquery.Parse(strings.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}
	node := jsonquery.FindOne(doc, export.Xpath)
	if node != nil && node.Value() != nil {
		value := node.Value()
		str := fmt.Sprintf("%v", value)
		switch export.Type {
		case "integer":
			return cast.ToInt(str), nil
		case "float":
			return cast.ToFloat64(str), nil
		default:
			return str, nil
		}
	}
	return nil, nil
}
