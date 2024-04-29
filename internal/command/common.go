package command

import (
	"github.com/antchfx/jsonquery"
	"github.com/flosch/pongo2/v6"
	"github.com/vearne/autotest/internal/config"
	"github.com/vearne/autotest/internal/model"
	"github.com/vearne/autotest/internal/resource"
	"reflect"
	"strings"
)

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
		if reflect.TypeOf(value).String() == "float64" && export.Type == model.TypeInteger {
			return int(value.(float64)), nil
		}
	}
	return nil, nil
}
