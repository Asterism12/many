package plugins

import (
	"github.com/Asterism12/many/base"
)

type getPluginArray struct {
}

func (g getPluginArray) Exec(s *base.Setter, data any, expression []string, param []any) any {
	arr, ok := data.([]any)
	if !ok {
		return nil
	}
	var values []any
	for _, datum := range arr {
		values = append(values, s.GetByRouter(datum, expression, param))
	}
	return values
}

func (g getPluginArray) Verify(param []any) ([]any, error) {
	return param, nil
}

func (g getPluginArray) Name() string {
	return "array"
}

type getPluginKey struct {
}

func (g getPluginKey) Exec(s *base.Setter, data any, expression []string, param []any) any {
	m, ok := data.(map[string]any)
	if !ok {
		return nil
	}
	var values []any
	for k := range m {
		values = append(values, s.GetByRouter(k, expression, param))
	}
	return values
}

func (g getPluginKey) Verify(param []any) ([]any, error) {
	return param, nil
}

func (g getPluginKey) Name() string {
	return "key"
}

type getPluginStrict struct {
}

func (g getPluginStrict) Exec(s *base.Setter, data any, expression []string, param []any) any {
	data = s.GetByRouter(data, expression, param)
	if g.strict(data) {
		return data
	} else {
		return nil
	}
}

func (g getPluginStrict) strict(data any) bool {
	if data == nil {
		return false
	}
	arr, ok := data.([]any)
	if !ok {
		return true
	}
	var values []any
	for _, datum := range arr {
		if g.strict(datum) {
			values = append(values, datum)
		}
	}
	if len(values) == 0 {
		return false
	}
	return true
}

func (g getPluginStrict) Verify(param []any) ([]any, error) {
	return param, nil
}

func (g getPluginStrict) Name() string {
	return "strict"
}

var DefaultPlugins = []base.GetterPlugin{
	getPluginArray{},
	getPluginKey{},
	getPluginStrict{},
}
