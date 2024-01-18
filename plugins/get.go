package plugins

import (
	"github.com/Asterism12/many/base"
)

type getPluginArray struct {
}

func (p getPluginArray) Exec(s *base.Setter, data any, expression []string, param []any) any {
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

func (p getPluginArray) Verify(param []any) ([]any, error) {
	return param, nil
}

func (p getPluginArray) Name() string {
	return "array"
}

type getPluginKey struct {
}

func (p getPluginKey) Exec(s *base.Setter, data any, expression []string, param []any) any {
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

func (p getPluginKey) Verify(param []any) ([]any, error) {
	return param, nil
}

func (p getPluginKey) Name() string {
	return "key"
}

type getPluginStrict struct {
}

func (p getPluginStrict) Exec(s *base.Setter, data any, expression []string, param []any) any {
	data = s.GetByRouter(data, expression, param)
	if p.strict(data) {
		return data
	} else {
		return nil
	}
}

func (p getPluginStrict) strict(data any) bool {
	if data == nil {
		return false
	}
	arr, ok := data.([]any)
	if !ok {
		return true
	}
	var values []any
	for _, datum := range arr {
		if p.strict(datum) {
			values = append(values, datum)
		}
	}
	if len(values) == 0 {
		return false
	}
	return true
}

func (p getPluginStrict) Verify(param []any) ([]any, error) {
	return param, nil
}

func (p getPluginStrict) Name() string {
	return "strict"
}

// DefaultGetterPlugins be set
var DefaultGetterPlugins = []base.GetterPlugin{
	getPluginArray{},
	getPluginKey{},
	getPluginStrict{},
}
