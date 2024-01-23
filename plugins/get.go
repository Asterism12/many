package plugins

import (
	"errors"
	"github.com/Asterism12/many/base"
)

type getPluginArray struct {
}

func (p getPluginArray) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	arr, ok := data.([]any)
	if !ok {
		return nil
	}
	var values []any
	for _, datum := range arr {
		values = append(values, s.GetByRouter(root, datum, expression, param))
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

func (p getPluginKey) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	m, ok := data.(map[string]any)
	if !ok {
		return nil
	}
	var values []any
	for k := range m {
		values = append(values, s.GetByRouter(root, k, expression, param))
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

func (p getPluginStrict) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	data = s.GetByRouter(root, data, expression, param)
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

func (g getPluginSelect) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	data = s.GetByRouter(root, data, expression, base.Rest(param))
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

func (g getPluginThis) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	return data
}

func (g getPluginThis) Verify(param []any) ([]any, error) {
	return param, nil
}

func (g getPluginThis) Name() string {
	return "this"
}

type getPluginRoot struct {
}

func (g getPluginRoot) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	return s.GetByRouter(root, root, expression, param)
}

func (g getPluginRoot) Verify(param []any) ([]any, error) {
	return param, nil
}

func (g getPluginRoot) Name() string {
	return "root"
}

const Case = "case"

type getPluginSwitch struct {
}

func (g getPluginSwitch) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	data = s.GetByRouter(root, data, expression, base.Rest(param))
	cases := param[0].([]any)
	for _, c := range cases {
		c := c.(map[string]any)
		if caseValue, ok := c[Case]; ok {
			if base.DeepEqual(data, caseValue) {
				return g.getValue(s, root, data, c)
			}
		} else {
			return g.getValue(s, root, data, c)
		}
	}
	return nil
}

func (g getPluginSwitch) getValue(s *base.Setter, root, data any, c map[string]any) any {
	if router, ok := c[base.Router]; ok {
		return s.Get(root, data, router)
	}
	if literal, ok := c[base.Literal]; ok {
		return literal
	}
	return nil
}

func (g getPluginSwitch) Verify(param []any) ([]any, error) {
	if len(param) == 0 {
		return nil, errors.New("plugin switch needs param")
	}
	_, ok := param[0].([]any)
	if !ok {
		return nil, errors.New("param of plugin switch is not []any")
	}
	return base.Rest(param), nil
}

func (g getPluginSwitch) Name() string {
	return "switch"
}

// DefaultGetterPlugins be set in Many by default
var DefaultGetterPlugins = []base.GetterPlugin{
	getPluginArray{},
	getPluginKey{},
	getPluginStrict{},
	getPluginSelect{},
	getPluginThis{},
	getPluginRoot{},
	getPluginSwitch{},
}
