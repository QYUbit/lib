package bufti

import (
	"testing"
)

func TestBasic(t *testing.T) {
	model := NewModel(
		NewField(0, "name", StringType),
		NewField(1, "age", Int8Type),
		NewField(3, "hight", Float32Type),
		NewField(4, "active", BoolType),
		NewField(5, "score", Int32Type),
		NewField(6, "hobbies", NewListType(StringType)),
	)

	bu := map[string]any{
		"name":    "alice",
		"age":     33,
		"hight":   6.6,
		"active":  true,
		"hobbies": []string{"swimming", "singing", "painting"},
	}

	b, err := model.Encode(bu)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)

	bu2, err := model.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bu2)
}

func TestEdge(t *testing.T) {
	model := NewModel(
		NewField(0, "a", NewListType(NewListType(BoolType))),
	)

	bu := map[string]any{
		"a": [][]bool{
			{true, false},
			{false, true},
		},
	}

	b, err := model.Encode(bu)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)

	bu2, err := model.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bu2)
}
