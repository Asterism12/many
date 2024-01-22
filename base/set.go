package base

import (
	"fmt"
	"strings"
)

// GetterPlugin execute when a router string is equal to result of Name()
type GetterPlugin interface {
	Exec(s *Setter, data any, expression []string, param []any) any
	Verify(param []any) ([]any, error)
	Name() string
}

// SetterPlugin execute when value of field 'mode' of phase is equal to result of Name()
type SetterPlugin interface {
	Exec(s *Setter, src, dst any, phase map[string]any) (any, map[string]any)
	Name() string
}

// Setter provide method 'set' and 'get'
type Setter struct {
	getterPlugins map[string]GetterPlugin
	setterPlugins map[string]SetterPlugin
	pluginPrefix  string
	segmentation  string
	defaultPhases []map[string]any
	omitempty     bool
}

// Verify return error when phases expression is valid
// use defaultPhases when phases is nil.
func (s *Setter) Verify(phases []map[string]any) error {
	if phases == nil {
		phases = s.defaultPhases
	}
	for _, phase := range phases {
		for _, expression := range phase {
			switch expression := expression.(type) {
			case map[string]any:
				routerString, ok := expression["router"].(string)
				if !ok {
					return fmt.Errorf(
						"field 'router' of plugin expression must be a string:%v", expression)
				}
				param, ok := expression["param"].([]any)
				if !ok {
					return fmt.Errorf(
						"field 'param' of plugin expression must be a []any:%v", expression)
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

// Get get value from data by expression
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

// GetByRouter expression is a slice of router string.
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

// GetBySlice expression is a []any.
// Return an array of results of expressions.
func (s *Setter) GetBySlice(data any, expressions []any) []any {
	var values []any
	for _, expression := range expressions {
		values = append(values, s.Get(data, expression))
	}
	return values
}

// GetByObject expression is a map[string]any.
// Return the result of target plugin.
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

// GetPluginName return plugin name with prefix
func (s *Setter) GetPluginName(name string) string {
	return s.pluginPrefix + name
}

// GetSegmentation return segmentation
func (s *Setter) GetSegmentation() string {
	return s.segmentation
}

// Set dst by phases and data of src.
// Return new value of dst and info of plugins return. The origin dst may be changed.
// The correct way to use it is -- dst,info = s.Set(src,dst,phases)
func (s *Setter) Set(src any, dst any, phases []map[string]any) (any, map[string]any) {
	if phases == nil {
		phases = s.defaultPhases
	}

	info := map[string]any{}
	modeField := s.GetPluginName("mode")
	thisField := s.GetPluginName("this")
	for _, phase := range phases {
		mode, ok := phase[modeField]
		if !ok {
			mode = "router"
		}
		switch mode {
		case "router":
			if v, ok := phase[thisField]; ok {
				dst = s.SetByRouter(dst, strings.Split(thisField, s.segmentation), s.Get(src, v))
			}
			for k, v := range phase {
				if k == modeField || k == thisField {
					continue
				}
				dst = s.SetByRouter(dst, strings.Split(k, s.segmentation), s.Get(src, v))
			}
		case "literal":
			if v, ok := phase[thisField]; ok {
				dst = s.SetByRouter(dst, strings.Split(thisField, s.segmentation), v)
			}
			for k, v := range phase {
				if k == modeField || k == thisField {
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

// SetByRouter set value of dst by router and data of src
func (s *Setter) SetByRouter(dst any, router []string, data any) any {
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
					dst[i] = s.SetByRouter(dst[i], Rest(router), datum)
				}
			}
			return dst
		} else {
			var nv []any
			for _, datum := range data {
				nv = append(nv, s.SetByRouter(nil, Rest(router), datum))
			}
			return nv
		}
	}

	if dst == nil {
		dst = map[string]any{}
	}
	if dst, ok := dst.(map[string]any); ok {
		dst[router[0]] = s.SetByRouter(dst[router[0]], Rest(router), data)
		if s.omitempty && dst[router[0]] == nil {
			delete(dst, router[0])
		}
	}
	return dst
}

// SetGetterPlugins set getter plugins of setter
func (s *Setter) SetGetterPlugins(plugins []GetterPlugin) {
	s.getterPlugins = make(map[string]GetterPlugin, len(plugins))
	for _, plugin := range plugins {
		s.getterPlugins[plugin.Name()] = plugin
	}
}

// SetSetterPlugins set setter plugins of setter
func (s *Setter) SetSetterPlugins(plugins []SetterPlugin) {
	s.setterPlugins = make(map[string]SetterPlugin, len(plugins))
	for _, phase := range plugins {
		s.setterPlugins[phase.Name()] = phase
	}
}

// SetPluginPrefix plugin prefix is a string to mark a router string is a pointer to a plugin.
// Works in both setter and getter.
func (s *Setter) SetPluginPrefix(prefix string) {
	s.pluginPrefix = prefix
}

// SetSegmentation segmentation is a string to split routers.
// Works in both setter and getter.
func (s *Setter) SetSegmentation(segmentation string) {
	s.segmentation = segmentation
}

// SetDefaultPhases defaultPhases is used in Set when phases is nil
func (s *Setter) SetDefaultPhases(phases []map[string]any) {
	s.defaultPhases = phases
}

// SetOmitempty delete null value in map[string]any if omitempty is true
func (s *Setter) SetOmitempty(omitempty bool) {
	s.omitempty = omitempty
}
