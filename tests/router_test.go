package tests

import (
	"github.com/Asterism12/many"
	"github.com/Asterism12/many/base"
	"testing"
)

func TestGetRouter(t *testing.T) {
	type args struct {
		data       string
		expression string
	}
	tests := []struct {
		name      string
		args      args
		want      any
		verifyErr bool
		opt       []many.Option
	}{
		{
			name: "standard",
			args: args{
				data:       `{"info":{"type":"apple"}}`,
				expression: `"info.type"`,
			},
			want: "apple",
		},
		{
			name: "forArray-object",
			args: args{
				data:       `{"info":{"type":"apple"}}`,
				expression: `"info.type"`,
			},
			want: "apple",
			opt: []many.Option{
				many.WithForArray(true),
			},
		},
		{
			name: "forArray-array",
			args: args{
				data:       `{"info":[{"type":"apple"},{"type":"tomato"}]}`,
				expression: `"info.type"`,
			},
			want: []any{"apple", "tomato"},
			opt: []many.Option{
				many.WithForArray(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := mustUnmarshal[any](tt.args.data)
			expression := mustUnmarshal[any](tt.args.expression)
			m := many.New(tt.opt...)
			if err := m.Verify([]map[string]any{{"res": expression}}); (err != nil) != tt.verifyErr {
				t.Errorf("Verify() = %v, want %v", err, tt.verifyErr)
			}
			if !tt.verifyErr {
				if got := m.Get(data, data, expression); !base.DeepEqual(got, tt.want) {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSetRouter(t *testing.T) {
	type args struct {
		data       string
		expression string
	}
	tests := []struct {
		name      string
		args      args
		want      any
		verifyErr bool
		opt       []many.Option
	}{
		{
			name: "slice",
			args: args{
				data:       `{"info":{"type":[]}}`,
				expression: `[{"#this":"info"}]`,
			},
			want: map[string]any{
				"type": []any{},
			},
		},
		{
			name: "redirect-src",
			args: args{
				data:       `{"type1":"apple","type2":"fruit"}`,
				expression: `[{"type1":"type2"},{"type3":"type1"}]`,
			},
			want: map[string]any{
				"type1": "fruit",
				"type3": "fruit",
			},
		},
		{
			name: "array-for",
			args: args{
				data:       `{"owns":[{"name":"apple","type":"fruit"}]}`,
				expression: `[{"#array.type":"owns.#for.type"}]`,
			},
			want: []any{
				map[string]any{
					"type": "fruit",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := mustUnmarshal[any](tt.args.data)
			expression := mustUnmarshal[[]map[string]any](tt.args.expression)
			m := many.New(tt.opt...)
			if err := m.Verify(expression); (err != nil) != tt.verifyErr {
				t.Errorf("Verify() = %v, want %v", err, tt.verifyErr)
			}
			if !tt.verifyErr {
				if got, _ := m.Set(data, nil, expression); !base.DeepEqual(got, tt.want) {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
