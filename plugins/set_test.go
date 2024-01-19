package plugins

import (
	"encoding/json"
	"github.com/Asterism12/many/base"
	"reflect"
	"strings"
	"testing"
)

func deepEqual(v1, v2 any) bool {
	switch v1 := v1.(type) {
	case []any:
		v2, ok := v2.([]any)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		valid := true
		for i := range v1 {
			valid = valid && deepEqual(v1[i], v2[i])
		}
		return valid
	case map[string]any:
		v2, ok := v2.(map[string]any)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		valid := true
		for i := range v1 {
			valid = valid && deepEqual(v1[i], v2[i])
		}
		return valid
	default:
		return reflect.DeepEqual(v1, v2)
	}
}

func Test_schemaVerify_Exec(t *testing.T) {
	type args struct {
		s     *base.Setter
		src   any
		dst   string
		phase string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		want1     string
		unmarshal func(data string, t *testing.T) any
	}{
		{
			name: "json.unmarshal-valid",
			args: args{
				s:   &base.Setter{},
				src: nil,
				dst: `{"a_string":"string","a_int":123,"a_float":1.23,"a_array":["1",2],"a_object":{"a_int":123}}`,
				phase: `{"a_string":"string","a_int":"number","a_float":"number","a_array":"array",
"a_object":"object","a_object.a_int":"number"}`,
			},
			want:  `{"a_string":"string","a_int":123,"a_float":1.23,"a_array":["1",2],"a_object":{"a_int":123}}`,
			want1: `{"schema_valid":true,"schema_invalid_info":[]}`,
			unmarshal: func(data string, t *testing.T) any {
				var d any
				if err := json.Unmarshal([]byte(data), &d); err != nil {
					t.Errorf("unmarshal err %v", err)
				}
				return d
			},
		},
		{
			name: "json.unmarshal-invalid",
			args: args{
				s:     &base.Setter{},
				src:   nil,
				dst:   `{"a_string":"string","a_int":123,"a_float":1.23,"a_array":["1",2],"a_object":{"a_int":123}}`,
				phase: `{"a_object.a_int":"string"}`,
			},
			want: `{"a_string":"string","a_int":123,"a_float":1.23,"a_array":["1",2],"a_object":{"a_int":123}}`,
			want1: `{"schema_valid":false,
"schema_invalid_info":[{"field":"a_object.a_int","want":"string","value":123}]}`,
			unmarshal: func(data string, t *testing.T) any {
				var d any
				if err := json.Unmarshal([]byte(data), &d); err != nil {
					t.Errorf("unmarshal err %v", err)
				}
				return d
			},
		},
		{
			name: "json.useNumber-valid",
			args: args{
				s:     &base.Setter{},
				src:   nil,
				dst:   `{"a_string":"string","a_int":123,"a_float":1.23,"a_array":["1",2],"a_object":{"a_int":123}}`,
				phase: `{"a_int":"int"}`,
			},
			want:  `{"a_string":"string","a_int":123,"a_float":1.23,"a_array":["1",2],"a_object":{"a_int":123}}`,
			want1: `{"schema_valid":true,"schema_invalid_info":[]}`,
			unmarshal: func(data string, t *testing.T) any {
				decoder := json.NewDecoder(strings.NewReader(data))
				decoder.UseNumber()
				var d any
				if err := decoder.Decode(&d); err != nil {
					t.Errorf("unmarshal err %v", err)
				}
				return d
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.s.SetSegmentation(".")
			p := schemaVerify{}
			got, info := p.Exec(tt.args.s, tt.args.src,
				tt.unmarshal(tt.args.dst, t), tt.unmarshal(tt.args.phase, t).(map[string]any))
			if !deepEqual(got, tt.unmarshal(tt.want, t)) {
				t.Errorf("Exec() got = %v, want %v", got, tt.want)
			}
			if !deepEqual(info, tt.unmarshal(tt.want1, t)) {
				t.Errorf("Exec() info = %v, want %v", info, tt.want1)
			}
		})
	}
}
