package base

import (
	"fmt"
	"strings"
)

const (
	Router  = "router"
	Literal = "literal"
	Mode    = "mode"
	This    = "this"
	Array   = "array"
)

// GetterPlugin execute when a router string is equal to result of Name()
type GetterPlugin interface {
	Exec(s *Setter, root, data any, expression []string, param []any) any
	Verify(param []any) ([]any, error)
	Name() string
}

// SetterPlugin execute when value of field 'mode' of phase is equal to result of Name()
type SetterPlugin interface {
	Exec(s *Setter, src, dst any, phase map[string]any) (any, map[string]any)
	Verify(phase map[string]any) error
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
	forArray      bool
	redirectSrc   bool
}

// Verify return error when phases expression is valid
// use defaultPhases when phases is nil.
func (s *Setter) Verify(phases []map[string]any) error {
	if phases == nil {
		phases = s.defaultPhases
	}
	for _, phase := range phases {
		mode, ok := phase[s.GetPluginName(Mode)]
		if !ok {
			mode = Router
		}
		switch mode {
		case Router:
			for _, expression := range phase {
				switch expression := expression.(type) {
				case map[string]any:
					return fmt.Errorf("expression of mode router cannnot be map[string]any:%v", expression)
				case []any:
					if len(expression) == 0 {
						return fmt.Errorf("len of expression must greater than 0:%v", expression)
					}
					routerString, ok := expression[0].(string)
					if !ok {
						return fmt.Errorf(
							"'router' of expression must be a string:%v", expression)
					}
					param := expression[1:]
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
		case Literal:
			continue
		default:
			name, ok := mode.(string)
			if !ok {
				return fmt.Errorf("field %s must be a string:%s", s.GetPluginName(Mode), mode)
			}
			plugin, ok := s.setterPlugins[name]
			if !ok {
				return fmt.Errorf("setter plugin is not exist:%s", name)
			}
			if err := plugin.Verify(phase); err != nil {
				return fmt.Errorf("phase of plugin:%s err:%w", name, err)
			}
		}
	}
	return nil
}

// Get get value data by expression
func (s *Setter) Get(root, data, expression any) any {
	switch expression := expression.(type) {
	case string:
		return s.GetByRouter(root, data, strings.Split(expression, s.segmentation), nil)
	case []any:
		return s.GetBySlice(root, data, expression)
	default:
		return nil
	}
}

// GetByRouter get value by routers and param.
func (s *Setter) GetByRouter(root any, data any, routers []string, param []any) any {
	if len(routers) == 0 || (len(routers) == 1 && routers[0] == "") {
		return data
	}
	if plugin, ok := s.getGetterPlugin(routers[0]); ok {
		return plugin.Exec(s, root, data, routers[1:], param)
	}
	m, ok := data.(map[string]any)
	if ok {
		return s.GetByRouter(root, m[routers[0]], routers[1:], param)
	}
	if s.forArray {
		arr, ok := data.([]any)
		if !ok {
			return nil
		}
		var values []any
		for _, datum := range arr {
			values = append(values, s.GetByRouter(root, datum, routers, param))
		}
		return values
	}
	return nil
}

// GetBySlice get value by expressions
// first element of expressions is declared as router.
// rest elements are declared as params.
func (s *Setter) GetBySlice(root, data any, expressions []any) any {
	router := expressions[0]
	params := expressions[1:]
	return s.GetByRouter(root, data, strings.Split(router.(string), s.segmentation), params)
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
	modeField := s.GetPluginName(Mode)
	thisField := s.GetPluginName(This)
	for _, phase := range phases {
		mode, ok := phase[modeField]
		if !ok {
			mode = Router
		}
		switch mode {
		case Router:
			if v, ok := phase[thisField]; ok {
				dst = s.SetByRouter(dst, strings.Split(thisField, s.segmentation), s.Get(src, src, v))
			}
			for k, v := range phase {
				if k == modeField || k == thisField {
					continue
				}
				dst = s.SetByRouter(dst, strings.Split(k, s.segmentation), s.Get(src, src, v))
			}
		case Literal:
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
		if s.redirectSrc {
			src = dst
		}
	}
	return dst, info
}

// SetByRouter set value of dst by router and data of src
func (s *Setter) SetByRouter(dst any, router []string, data any) any {
	if router == nil || router[0] == s.GetPluginName(This) {
		return deepCopy(data)
	}

	if router[0] == s.GetPluginName(Array) {
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
			for i := range dst {
				if i < len(data) {
					dst[i] = s.SetByRouter(dst[i], Rest(router), data[i])
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

// SetForArray traverse array to get value by rest of routers when meet an array in GetByRouter.
//
// Logic of forArray is as same as plugin 'for'.
// Set forArray to true allowed you to omit router 'for' in routers.
// On the other hand it makes the mean of router ambiguous.
func (s *Setter) SetForArray(forArray bool) {
	s.forArray = forArray
}

// SetRedirectSrc set src to dst while a phase finished in Set
func (s *Setter) SetRedirectSrc(redirectSrc bool) {
	s.redirectSrc = redirectSrc
}
