package plugins

import (
	"encoding/json"
	"github.com/Asterism12/many/base"
	"strings"
)

type schemaVerify struct {
}

func (p schemaVerify) Exec(s *base.Setter, _, dst any, phase map[string]any) (any, map[string]any) {
	allValid := true
	var invalidInfo []any
	modeField := s.GetPluginName("mode")
	for k, schemaAny := range phase {
		if k == modeField {
			continue
		}
		value := s.GetByRouter(dst, strings.Split(k, s.GetSegmentation()), nil)
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
			valid = schema == "null"
		}
		if valid == false {
			allValid = false
			invalidInfo = append(invalidInfo, map[string]any{"field": k, "want": schema, "value": value})
		}
	}
	return dst, map[string]any{"schema_valid": allValid, "schema_invalid_info": invalidInfo}
}

func (p schemaVerify) Name() string {
	return "schema"
}

var DefaultSetterPlugins = []base.SetterPlugin{
	schemaVerify{},
}
