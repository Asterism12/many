package many

import (
	"github.com/Asterism12/many/base"
	"github.com/Asterism12/many/plugins"
)

// Many a converter of schema of map[string]any and []any
type Many interface {
	Set(src any, dst any, phases []map[string]any) (any, map[string]any)
	Get(root, data, expression any) any
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
	setter.SetOmitempty(true)
	setter.SetRedirectSrc(true)
	for _, opt := range opts {
		opt(setter)
	}
	return setter
}

// WithGetterPlugins set plugins in ps and DefaultGetterPlugins in this Many.
// The latter overrides the previous if they use a duplicate name.
func WithGetterPlugins(ps []base.GetterPlugin) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetGetterPlugins(append(plugins.DefaultGetterPlugins, ps...))
	}
}

// WithSetterPlugins set plugins in ps and DefaultSetterPlugins in this Many.
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

// WithPhases set default phases in this Many.
// default phases is used in Set when phases is nil
func WithPhases(phases []map[string]any) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetDefaultPhases(phases)
	}
}

// WithOmitempty delete null value in map[string]any if omitempty is true
func WithOmitempty(omitempty bool) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetOmitempty(omitempty)
	}
}

// WithForArray traverse array to get value by rest of routers when meet an array in GetByRouter.
//
// Logic of forArray is as same as plugin 'for'.
// Set forArray to true allowed you to omit router 'for' in routers.
// On the other hand it makes the mean of router ambiguous.
func WithForArray(forArray bool) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetForArray(forArray)
	}
}

// WithRedirectSrc set src to dst while a phase finished in Set
func WithRedirectSrc(redirectSrc bool) func(*base.Setter) {
	return func(setter *base.Setter) {
		setter.SetRedirectSrc(redirectSrc)
	}
}
