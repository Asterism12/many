package example

import (
	"encoding/json"
	"fmt"
	"many"
	"testing"
)

func TestExample(t *testing.T) {
	expressionJSON := []byte(`[
  {
    "#mode": "router",
    "#array.result": "#strict.#array.#array.a"
  }
]`)
	srcBys := []byte(`[
  [
    {
      "a": "a1",
      "b": "b1"
    },
    {
      "a": "a2",
      "b": "b2"
    }
  ],
  [
    {
      "a": "a3",
      "b": "b3"
    },
    {
      "a": "a4",
      "b": {
        "c": "c1",
        "d": "d1"
      }
    }
  ]
]`)

	setter := many.New()
	var expression []map[string]any
	var src any
	_ = json.Unmarshal(expressionJSON, &expression)
	_ = json.Unmarshal(srcBys, &src)

	dst := setter.Set(src, nil, expression)
	fmt.Println(dst)
	fmt.Println(json.Marshal(dst))
}
