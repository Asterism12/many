package plugins

import (
	"errors"
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

type getPluginSelect struct {
}

func (g getPluginSelect) Exec(s *base.Setter, data any, expression []string, param []any) any {
	data = s.GetByRouter(data, expression, base.Rest(param))
	phase := param[0].(map[string]any)
	dst, _ := s.Set(data, nil, []map[string]any{phase})
	return dst
}

func (g getPluginSelect) Verify(param []any) ([]any, error) {
	if len(param) == 0 {
		return nil, errors.New("plugin select needs param")
	}
	if _, ok := param[0].(map[string]any); !ok {
		return nil, errors.New("param of plugin select is not map[string]any")
	}
	return base.Rest(param), nil
}

func (g getPluginSelect) Name() string {
	return "select"
}

type getPluginThis struct {
}

func (g getPluginThis) Exec(s *base.Setter, data any, expression []string, param []any) any {
	return data
}

func (g getPluginThis) Verify(param []any) ([]any, error) {
	return param, nil
}

func (g getPluginThis) Name() string {
	return "this"
}

// DefaultGetterPlugins be set
var DefaultGetterPlugins = []base.GetterPlugin{
	getPluginArray{},
	getPluginKey{},
	getPluginStrict{},
	getPluginSelect{},
	getPluginThis{},
}
