package base

import (
	"reflect"
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

func Test_deepCopy(t *testing.T) {
	type args struct {
		src any
	}
	tests := []struct {
		name   string
		args   args
		want   any
		modify func(any)
	}{
		{
			name: "map_change",
			args: args{
				src: map[string]any{
					"i": 1,
					"m": map[string]any{
						"i": 1,
					},
				},
			},
			want: map[string]any{
				"i": 1,
				"m": map[string]any{
					"i": 1,
				},
			},
			modify: func(data any) {
				m := data.(map[string]any)
				m["i"] = 2
				m["m"] = []any{}
			},
		},
		{
			name: "2_layer_map_change",
			args: args{
				src: map[string]any{
					"i": 1,
					"m": map[string]any{
						"i": 1,
					},
				},
			},
			want: map[string]any{
				"i": 1,
				"m": map[string]any{
					"i": 1,
				},
			},
			modify: func(data any) {
				m := data.(map[string]any)
				m["i"] = 2
				m2 := m["m"].(map[string]any)
				m2["i"] = 2
				m2["i2"] = 3
			},
		},
		{
			name: "slices",
			args: args{
				src: map[string]any{
					"slice": []any{
						1,
						"s",
					},
				},
			},
			want: map[string]any{
				"slice": []any{
					1,
					"s",
				},
			},
			modify: func(data any) {
				m := data.(map[string]any)
				s := m["slice"].([]any)
				s[1] = 2
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deepCopy(tt.args.src)
			if !deepEqual(got, tt.want) {
				t.Errorf("deepCopy() = %v, want %v", got, tt.want)
			}
			tt.modify(tt.args.src)
			if !deepEqual(got, tt.want) {
				t.Errorf("after modify = %v, want %v", got, tt.want)
			}
		})
	}
}
