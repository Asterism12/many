package plugins

import (
	"encoding/json"
	"github.com/Asterism12/many/base"
	"strings"
)

const (
	SchemaValid       = "schema_valid"
	SchemaInvalidInfo = "schema_invalid_info"
	Field             = "field"
	Want              = "want"
	Value             = "value"
)

type schemaVerify struct {
}

func (p schemaVerify) Exec(s *base.Setter, _, dst any, phase map[string]any) (any, map[string]any) {
	allValid := true
	var invalidInfo []any
	modeField := s.GetPluginName("mode")
	allowNull := s.GetPluginName("null") == "allow"
	for k, schemaAny := range phase {
		if k == modeField {
			continue
		}
		value := s.GetByRouter(dst, dst, strings.Split(k, s.GetSegmentation()), nil)
		schema := schemaAny.(string)
		valid := true
		switch value := value.(type) {
		case string:
			valid = schema == "string"
		case float64:
			valid = schema == "number"
		case json.Number:
			errNil := func(_ any, err error) bool {
				return err == nil
			}
			valid = schema == "number" || (schema == "int" && errNil(value.Int64()))
		case []any:
			valid = schema == "array"
		case map[string]any:
			valid = schema == "object"
		case bool:
			valid = schema == "bool"
		default:
			valid = allowNull || schema == "null"
		}
		if valid == false {
			allValid = false
			invalidInfo = append(invalidInfo, map[string]any{Field: k, Want: schema, Value: value})
		}
	}
	return dst, map[string]any{SchemaValid: allValid, SchemaInvalidInfo: invalidInfo}
}

func (p schemaVerify) Verify(phase map[string]any) error {
	return nil
}

func (p schemaVerify) Name() string {
	return "schema"
}

type newRouter struct{}

func (p newRouter) Exec(s *base.Setter, src, dst any, phase map[string]any) (any, map[string]any) {
	dst = nil
	thisField := s.GetPluginName(base.This)
	if v, ok := phase[thisField]; ok {
		dst = s.SetByRouter(dst, strings.Split(thisField, s.GetSegmentation()), s.Get(src, src, v))
	}
	for k, v := range phase {
		if k == thisField {
			continue
		}
		dst = s.SetByRouter(dst, strings.Split(k, s.GetSegmentation()), s.Get(src, src, v))
	}
	return dst, nil
}

func (p newRouter) Verify(phase map[string]any) error {
	return nil
}

func (p newRouter) Name() string {
	return "new"
}

var DefaultSetterPlugins = []base.SetterPlugin{
	schemaVerify{},
	newRouter{},
}
