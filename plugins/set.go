package plugins

import (
	"encoding/json"
	"github.com/Asterism12/many/base"
	"strings"
)

type schemaVerify struct {
}

func (p schemaVerify) Exec(s *base.Setter, src, dst any, phase map[string]any) (any, map[string]any) {
	result := true
	modeField := s.GetPluginName("mode")
	for k, schemaAny := range phase {
		if k == modeField {
			continue
		}
		value := s.GetByRouter(dst, strings.Split(k, s.GetSegmentation()), nil)
		schema := schemaAny.(string)
		switch value := value.(type) {
		case string:
			result = result && schema == "string"
		case float64:
			result = result && (schema == "number")
		case json.Number:
			errNil := func(_ any, err error) bool {
				return err == nil
			}
			result = result && (schema == "number" || (schema == "int" && errNil(value.Int64())))
		case []any:
			result = result && (schema == "array")
		case map[string]any:
			result = result && (schema == "object")
		default:
			result = result && (schema == "null" && value == nil)
		}
	}
	return dst, map[string]any{"schema": result}
}

func (p schemaVerify) Name() string {
	return "schema"
}
