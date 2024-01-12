package many

import (
	"many/base"
	"many/plugins"
)

type Setter interface {
	Set(src any, dst any, phases []map[string]any) any
	Get(data, expression any) any
	Verify(expression any) error
}

func New() Setter {
	setter := &base.Setter{}
	setter.SetGetterPlugins(plugins.DefaultPlugins)
	setter.SetPluginPrefix("#")
	return setter
}
