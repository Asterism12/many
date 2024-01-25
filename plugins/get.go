package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Asterism12/many/base"
	"unsafe"
)

type getPluginFor struct {
}

func (p getPluginFor) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
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

func (p getPluginFor) Verify(param []any) ([]any, error) {
	return param, nil
}

func (p getPluginFor) Name() string {
	return "for"
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

const (
	Case       = "case"
	ModeString = "string"
	ModeDeep   = "deep"
)

type getPluginSwitch struct {
}

func (g getPluginSwitch) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	value := s.GetByRouter(root, data, expression, base.Rest(param))
	param = param[0].([]any)
	mode, ok := param[0].(string)
	if !ok {
		return g.switchByString(s, root, data, value, param)
	}
	switch mode {
	case ModeDeep:
		return g.switchByDeep(s, root, data, value, param[1:])
	case ModeString:
		return g.switchByString(s, root, data, value, param[1:])
	default:
		return g.switchByString(s, root, data, value, param[1:])
	}
}

func (g getPluginSwitch) switchByString(s *base.Setter, root, data, value any, param []any) any {
	cases := param[0].(map[string]any)
	valueAsString := g.getAsString(value)
	for c, router := range cases {
		if c == valueAsString {
			return s.Get(root, data, router)
		}
	}
	if len(param) == 2 {
		return s.Get(root, data, param[1])
	}
	return nil
}

func (g getPluginSwitch) getAsString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%f", v)
	case json.Number:
		return v.String()
	default:
		bys, _ := json.Marshal(&v)
		return string(bys)
	}
}

func (g getPluginSwitch) bytesToString(bys []byte) string {
	return unsafe.String(unsafe.SliceData(bys), len(bys))
}

func (g getPluginSwitch) switchByDeep(s *base.Setter, root, data, value any, param []any) any {
	for _, c := range param {
		c := c.(map[string]any)
		if caseValue, ok := c[Case]; ok {
			if base.DeepEqual(value, caseValue) {
				return g.getValue(s, root, data, c)
			}
		}
	}
	c := param[len(param)-1].(map[string]any)
	if _, ok := c[Case]; !ok {
		return g.getValue(s, root, data, c)
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

type getPluginLiteral struct {
}

func (g getPluginLiteral) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	return param[0]
}

func (g getPluginLiteral) Verify(param []any) ([]any, error) {
	if len(param) == 0 {
		return nil, errors.New("plugin literal needs param")
	}
	return param[1:], nil
}

func (g getPluginLiteral) Name() string {
	return "literal"
}

type getPluginArray struct {
}

func (g getPluginArray) Exec(s *base.Setter, root, data any, expression []string, param []any) any {
	routers := param[0].([]any)
	result := make([]any, len(routers))
	for _, router := range param[0].([]any) {
		result = append(result, s.Get(root, data, router))
	}
	return result
}

func (g getPluginArray) Verify(param []any) ([]any, error) {
	if len(param) == 0 {
		return nil, errors.New("plugin array needs param")
	}
	if _, ok := param[0].([]any); !ok {
		return nil, errors.New("plugin array must be []any")
	}
	return param[1:], nil
}

func (g getPluginArray) Name() string {
	return "array"
}

// DefaultGetterPlugins be set in Many by default
var DefaultGetterPlugins = []base.GetterPlugin{
	getPluginFor{},
	getPluginKey{},
	getPluginStrict{},
	getPluginSelect{},
	getPluginThis{},
	getPluginRoot{},
	getPluginSwitch{},
}
