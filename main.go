package many

import (
	"github.com/Asterism12/many/base"
	"github.com/Asterism12/many/plugins"
)

type Many interface {
	Set(src any, dst any, phases []map[string]any) (any, map[string]any)
	Get(data, expression any) any
	Verify(expression any) error
}

type Option func(*base.Setter)

func New(opts ...Option) Many {
	setter := &base.Setter{}
	setter.SetGetterPlugins(plugins.DefaultGetterPlugins)
	setter.SetSetterPlugins(plugins.DefaultSetterPlugins)
	setter.SetPluginPrefix("#")
	setter.SetSegmentation(".")
	for _, opt := range opts {
		opt(setter)
	}
	return setter
}

func WithGetterPlugins(ps []base.GetterPlugin) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetGetterPlugins(append(plugins.DefaultGetterPlugins, ps...))
	}
}

func WithSetterPlugins(ps []base.SetterPlugin) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetSetterPlugins(append(plugins.DefaultSetterPlugins, ps...))
	}
}

func WithPluginPrefix(prefix string) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetPluginPrefix(prefix)
	}
}

func WithSegmentation(segmentation string) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetSegmentation(segmentation)
	}
}
