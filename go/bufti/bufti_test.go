package bufti

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	city := NewModel("city",
		NewField(0, "name", StringType),
		NewField(1, "population", Int32Type),
	)

	person := NewModel("person",
		NewField(0, "name", StringType),
		NewField(1, "age", Int8Type),
		NewField(2, "hight", Float32Type),
		NewField(3, "active", BoolType),
		NewField(4, "hobbies", NewListType(StringType)),
		NewField(5, "city", NewModelType(city)),
	)

	bu := map[string]any{
		"name":    "alice",
		"age":     33,
		"hight":   6.6,
		"active":  true,
		"hobbies": []string{"swimming", "singing", "painting"},
		"city": map[string]any{
			"name":       "Cairo",
			"population": 10000000,
		},
	}

	b, err := person.Encode(bu)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)

	bu2, err := person.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bu2)
}

func TestEdge(t *testing.T) {
	model2 := NewModel("b",
		NewField(0, "a", NewListType(StringType)),
	)

	model := NewModel("a",
		NewField(3, "a", NewListType(NewListType(BoolType))),
		NewField(4, "b", NewListType(NewModelType(model2))),
	)

	bu := map[string]any{
		"a": [][]bool{
			{true, false},
			{false, true},
		},
		"b": []map[string]any{
			{
				"a": []string{"aaa", "bbb", "ccc"},
			},
		},
	}

	fmt.Println(model)

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
