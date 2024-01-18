# many
A converter of schema of map[string]any and []any in Go

## Usage
src:
```json
[
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
]
```
phases:
```json
[
  {
    "#mode": "router",
    "#array.result": "#strict.#array.#array.a"
  }
]
```
code:
```go
	setter := many.New()
	var expression []map[string]any
	var src any
	_ = json.Unmarshal(expressionJSON, &expression)
	_ = json.Unmarshal(srcBys, &src)

	dst, info := setter.Set(src, nil, expression)
	bys, err := json.Marshal(dst)
	fmt.Println(string(bys), err, info)
```
output:
```
[{"result":["a1","a2"]},{"result":["a3","a4"]}] <nil> map[]
```
