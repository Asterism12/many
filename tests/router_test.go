package tests

import (
	"fmt"
	"github.com/Asterism12/many"
	"github.com/Asterism12/many/base"
	"testing"
)

func TestRouter(t *testing.T) {
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
			data := mustUnmarshal(tt.args.data)
			expression := mustUnmarshal(tt.args.expression)
			m := many.New(tt.opt...)
			if err := m.Verify([]map[string]any{{"res": expression}}); (err != nil) != tt.verifyErr {
				t.Errorf("Verify() = %v, want %v", err, tt.verifyErr)
			}
			if !tt.verifyErr {
				if got := m.Get(data, data, expression); !base.DeepEqual(got, tt.want) {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				} else {
					fmt.Println(got, tt.want)
				}
			}
		})
	}
}
