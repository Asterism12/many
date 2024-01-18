package many

import (
	"github.com/Asterism12/many/base"
	"github.com/Asterism12/many/plugins"
)

// Many a converter of schema of map[string]any and []any
type Many interface {
	Set(src any, dst any, phases []map[string]any) (any, map[string]any)
	Get(data, expression any) any
	Verify(phases []map[string]any) error
}

// Option to custom Many
type Option func(*base.Setter)

// New create a Many.
// The methods provided by Many interface are safe to exec concurrently.
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

// WithGetterPlugins plugins in ps and DefaultGetterPlugins will be set in this Many.
// The latter overrides the previous if they use a duplicate name.
func WithGetterPlugins(ps []base.GetterPlugin) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetGetterPlugins(append(plugins.DefaultGetterPlugins, ps...))
	}
}

// WithSetterPlugins plugins in ps and DefaultSetterPlugins will be set in this Many.
// The latter overrides the previous if they use a duplicate name.
func WithSetterPlugins(ps []base.SetterPlugin) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetSetterPlugins(append(plugins.DefaultSetterPlugins, ps...))
	}
}

// WithPluginPrefix plugin prefix is a string to mark a router string is a pointer to a plugin.
// Works in both setter and getter.
// Default value is "#"
func WithPluginPrefix(prefix string) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetPluginPrefix(prefix)
	}
}

// WithSegmentation segmentation is a string to split routers.
// Works in both setter and getter.
// Default value is "."
func WithSegmentation(segmentation string) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetSegmentation(segmentation)
	}
}
