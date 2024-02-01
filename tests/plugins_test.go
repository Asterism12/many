package tests

import (
	"encoding/json"
	"fmt"
	"github.com/Asterism12/many"
	"github.com/Asterism12/many/base"
	"testing"
)

func mustUnmarshal[T any](s string) T {
	var v T
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		panic(fmt.Errorf("json err %w", err))
	}
	return v
}

func TestSwitch(t *testing.T) {
	type args struct {
		data       string
		expression string
	}
	tests := []struct {
		name      string
		args      args
		want      any
		verifyErr bool
	}{
		{
			name: "literal-standard",
			args: args{
				data:       `{"type":"apple"}`,
				expression: `["#switch.type",["literal",{"apple":"fruit","tomato":"vegetable"},"no_idea"]]`,
			},
			want: "fruit",
		},
		{
			name: "literal-omit-parameter-default",
			args: args{
				data:       `{"type":"meat"}`,
				expression: `["#switch.type",[{"apple":"fruit","tomato":"vegetable"},"no_idea"]]`,
			},
			want: "no_idea",
		},
		{
			name: "string-standard",
			args: args{
				data:       `{"type":"apple","word":{"vegetable":"good","fruit":"better"}}`,
				expression: `["#switch.type",["string",{"apple":"word.fruit","tomato":"word.vegetable"},"no_idea"]]`,
			},
			want: "better",
		},
		{
			name: "deep-standard",
			args: args{
				data:       `{"type":"apple","word":{"vegetable":"good","fruit":"better"}}`,
				expression: `["#switch.type",["deep",[{"case":"apple","router":"word.fruit"}]]]`,
			},
			want: "better",
		},
		{
			name: "deep-default",
			args: args{
				data:       `{"type":"meat"}`,
				expression: `["#switch.type",["deep",[{"case":"apple","router":"word.fruit"}],{"literal":"no_idea"}]]`,
			},
			want: "no_idea",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := mustUnmarshal[any](tt.args.data)
			expression := mustUnmarshal[any](tt.args.expression)
			m := many.New()
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
