package base

import (
	"fmt"
	"strings"
)

// GetterPlugin Getter插件
type GetterPlugin interface {
	Exec(s *Setter, data any, expression []string, param []any) any
	Verify(param []any) ([]any, error)
	Name() string
}

type SetterPlugin interface {
	Exec(s *Setter, src, dst any, phase map[string]any) (any, map[string]any)
	Name() string
}

type Setter struct {
	getterPlugins map[string]GetterPlugin
	setterPlugins map[string]SetterPlugin
	pluginPrefix  string
	segmentation  string
}

func (s *Setter) Verify(expression any) error {
	expressionAnyArray, ok := expression.([]any)
	if !ok {
		return fmt.Errorf("expression must be a []any:%v", expression)
	}
	for _, expressionAny := range expressionAnyArray {
		expressionMany, ok := expressionAny.(map[string]any)
		if !ok {
			return fmt.Errorf("expression must be a map[string]any:%v", expression)
		}
		for _, getterExpression := range expressionMany {
			switch getterExpression := getterExpression.(type) {
			case map[string]any:
				routerString, ok := getterExpression["router"].(string)
				if !ok {
					return fmt.Errorf(
						"field 'router' of plugin expression must be a string:%v", getterExpression)
				}
				param, ok := getterExpression["param"].([]any)
				if !ok {
					return fmt.Errorf(
						"field 'param' of plugin expression must be a []any:%v", getterExpression)
				}
				routers := strings.Split(routerString, s.segmentation)
				for _, router := range routers {
					if plugin, ok := s.getGetterPlugin(router); ok {
						var err error
						param, err = plugin.Verify(param)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func (s *Setter) Get(data, expression any) any {
	switch expression := expression.(type) {
	case string:
		return s.GetByRouter(data, strings.Split(expression, s.segmentation), nil)
	case []any:
		return s.GetBySlice(data, expression)
	case map[string]any:
		return s.GetByObject(data, expression)
	default:
		return nil
	}
}

func (s *Setter) GetByRouter(data any, expressions []string, param []any) any {
	if len(expressions) == 0 {
		return data
	}
	if plugin, ok := s.getGetterPlugin(expressions[0]); ok {
		return plugin.Exec(s, data, expressions[1:], param)
	}
	m, ok := data.(map[string]any)
	if !ok {
		return nil
	}
	return s.GetByRouter(m[expressions[0]], expressions[1:], param)
}

func (s *Setter) GetBySlice(data any, expressions []any) []any {
	var values []any
	for _, expression := range expressions {
		values = append(values, s.Get(data, expression))
	}
	return values
}

func (s *Setter) GetByObject(data any, expressions map[string]any) any {
	router := strings.Split(expressions["router"].(string), s.segmentation)
	param := expressions["param"].([]any)
	return s.GetByRouter(data, router, param)
}

func (s *Setter) getGetterPlugin(expression string) (GetterPlugin, bool) {
	if !strings.HasPrefix(expression, s.pluginPrefix) {
		return nil, false
	}
	name := strings.TrimPrefix(expression, s.pluginPrefix)
	plugin, ok := s.getterPlugins[name]
	return plugin, ok
}

func (s *Setter) GetPluginName(name string) string {
	return s.pluginPrefix + name
}

func (s *Setter) GetSegmentation() string {
	return s.segmentation
}

func (s *Setter) Set(src any, dst any, phases []map[string]any) (any, map[string]any) {
	info := map[string]any{}

	modeField := s.GetPluginName("mode")
	for _, phase := range phases {
		mode, ok := phase[modeField]
		if !ok {
			mode = "router"
		}
		switch mode {
		case "router":
			for k, v := range phase {
				if k == modeField {
					continue
				}
				dst = s.SetByRouter(dst, strings.Split(k, s.segmentation), s.Get(src, v))
			}
		case "literal":
			for k, v := range phase {
				if k == modeField {
					continue
				}
				dst = s.SetByRouter(dst, strings.Split(k, s.segmentation), v)
			}
		default:
			var pluginInfo map[string]any
			dst, pluginInfo = s.setterPlugins[mode.(string)].Exec(s, src, dst, phase)
			for k, v := range pluginInfo {
				info[k] = v
			}
		}
	}
	return dst, info
}

func (s *Setter) SetByRouter(dst any, router []string, data any) any {
	if data == nil {
		return dst
	}
	if router == nil || router[0] == s.GetPluginName("this") {
		return deepCopy(data)
	}

	if router[0] == s.GetPluginName("array") {
		data, ok := data.([]any)
		if !ok {
			return nil
		}
		if dst == nil {
			dst = []any{}
		}
		dst, ok := dst.([]any)
		if !ok {
			return dst
		}
		if len(dst) != 0 {
			for i, datum := range dst {
				if i < len(data) {
					dst[i] = s.SetByRouter(dst[i], rest(router), datum)
				}
			}
			return dst
		} else {
			var nv []any
			for _, datum := range data {
				nv = append(nv, s.SetByRouter(nil, rest(router), datum))
			}
			return nv
		}
	}

	if dst == nil {
		dst = map[string]any{}
	}
	if dst, ok := dst.(map[string]any); ok {
		dst[router[0]] = s.SetByRouter(dst[router[0]], rest(router), data)
	}
	return dst
}

func (s *Setter) SetPlugins(plugins []GetterPlugin) {
	s.getterPlugins = make(map[string]GetterPlugin, len(plugins))
	for _, plugin := range plugins {
		s.getterPlugins[plugin.Name()] = plugin
	}
}

func (s *Setter) SetPluginPrefix(prefix string) {
	s.pluginPrefix = prefix
}

func (s *Setter) SetSegmentation(segmentation string) {
	s.segmentation = segmentation
}

func (s *Setter) SetPhases(phases []SetterPlugin) {
	s.setterPlugins = make(map[string]SetterPlugin, len(phases))
	for _, phase := range phases {
		s.setterPlugins[phase.Name()] = phase
	}
}
