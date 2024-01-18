package many

import (
	"github.com/Asterism12/many/base"
	"github.com/Asterism12/many/plugins"
)

type Setter interface {
	Set(src any, dst any, phases []map[string]any) (any, map[string]any)
	Get(data, expression any) any
	Verify(expression any) error
}

type Option func(*base.Setter)

func New(opts ...Option) Setter {
	setter := &base.Setter{}
	setter.SetPlugins(plugins.DefaultPlugins)
	setter.SetPluginPrefix("#")
	setter.SetSegmentation(".")
	for _, opt := range opts {
		opt(setter)
	}
	return setter
}

func WithGetterPlugins(ps []base.GetterPlugin) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetPlugins(append(plugins.DefaultPlugins, ps...))
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
