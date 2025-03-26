package bufti

import (
	"runtime"
	"testing"
)

func TestBasic(t *testing.T) {
	model := NewModel(
		NewField(0, "name", StringType),
		NewField(1, "age", Int8Type),
		NewField(3, "hight", Float32Type),
		NewField(4, "active", BoolType),
		NewField(5, "score", Int32Type),
	)

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
