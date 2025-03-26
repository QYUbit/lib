package bufti

import (
	"runtime"
	"testing"
)

func TestBasic(t *testing.T) {
	model := Model{
		labels: map[string]byte{
			"name":   0,
			"age":    1,
			"hight":  3,
			"active": 4,
			"score":  5,
		},
		schama: map[byte]Field{
			0: {label: "name", fieldType: StringType},
			1: {label: "age", fieldType: Int8Type},
			3: {label: "hight", fieldType: Float32Type},
			4: {label: "active", fieldType: BoolType},
			5: {label: "score", fieldType: Int32Type},
		},
	}

	bu := map[string]any{
		"name":   "alice",
		"age":    33,
		"hight":  6.6,
		"active": true,
	}

	b, err := model.Encode(bu)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)

	bu2, err := model.Decode(b)
	if err != nil {
		buf := make([]byte, 1024)
		n := runtime.Stack(buf, false)
		t.Fatalf("Error: %v\n%s", err, buf[:n])
	}
	t.Log(bu2)
}
