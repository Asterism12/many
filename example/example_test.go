package example

import (
	"encoding/json"
	"fmt"
	"github.com/Asterism12/many"
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

	var expression []map[string]any
	var src any
	_ = json.Unmarshal(expressionJSON, &expression)
	_ = json.Unmarshal(srcBys, &src)

	setter := many.New(many.WithPhases(expression))
	fmt.Println(setter.Verify(nil))
	dst, info := setter.Set(src, nil, nil)
	bys, err := json.Marshal(dst)
	fmt.Println(string(bys), err, info)
}

func TestExample2(t *testing.T) {
	expressionJSON := []byte(`[
  {
    "#this": "#this",
    "a": [
      "a.#select",
      {
        "a1": "a1"
      }
    ]
  }
]`)
	srcBys := []byte(`{"a":{"a1":"a1","b1":"b1"},"b":123}`)

	var expression []map[string]any
	var src any
	_ = json.Unmarshal(expressionJSON, &expression)
	_ = json.Unmarshal(srcBys, &src)

	setter := many.New(many.WithPhases(expression))
	fmt.Println(setter.Verify(nil))
	dst, info := setter.Set(src, nil, nil)
	bys, err := json.Marshal(dst)
	fmt.Println(string(bys), err, info)
}

func TestExample3(t *testing.T) {
	expressionJSON := []byte(`[
  {
    "res": [
      "a.#array.#switch.id",
      [
        {
          "case": 1,
          "router": "#root.b"
        },
        {
          "case": 2,
          "literal": "s2"
        },
        {
          "literal": [
            "default",
            "literal"
          ]
        }
      ]
    ]
  }
]`)
	srcBys := []byte(`{"a":[{"id":1},{"id":2},{"id":3}],"b":123}`)

	var expression []map[string]any
	var src any
	_ = json.Unmarshal(expressionJSON, &expression)
	_ = json.Unmarshal(srcBys, &src)

	setter := many.New(many.WithPhases(expression))
	fmt.Println(setter.Verify(nil))
	dst, info := setter.Set(src, nil, nil)
	bys, err := json.Marshal(dst)
	fmt.Println(string(bys), err, info)
}
