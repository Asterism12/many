package base

import (
	"testing"
)

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
			if !DeepEqual(got, tt.want) {
				t.Errorf("deepCopy() = %v, want %v", got, tt.want)
			}
			tt.modify(tt.args.src)
			if !DeepEqual(got, tt.want) {
				t.Errorf("after modify = %v, want %v", got, tt.want)
			}
		})
	}
}
