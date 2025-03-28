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
		NewField(2, "height", Float64Type),
		NewField(3, "active", BoolType),
		NewField(4, "hobbies", NewListType(StringType)),
		NewField(5, "city", NewModelType(city)),
	)

	bu := map[string]any{
		"name":    "alice",
		"age":     33,
		"height":  6.6,
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

func TestListType(t *testing.T) {
	model := NewModel("listTest",
		NewField(0, "numbers", NewListType(Int32Type)),
	)
	bu := map[string]any{
		"numbers": []int32{1, 2, 3},
	}

	b, err := model.Encode(bu)
	if err != nil {
		t.Fatalf("Encoding failed for list type: %v", err)
	}
	t.Log("Encoded list type:", b)

	decoded, err := model.Decode(b)
	if err != nil {
		t.Fatalf("Decoding failed for list type: %v", err)
	}
	t.Log("Decoded list type:", decoded)
}

func TestMapType(t *testing.T) {
	cityPopulationModel := NewModel("cityPopulation",
		NewField(0, "populations", NewMapType(StringType, Int32Type)),
	)

	bu := map[string]any{
		"populations": map[string]int32{
			"New York":    8419000,
			"Los Angeles": 3980000,
			"Chicago":     2716000,
		},
	}

	b, err := cityPopulationModel.Encode(bu)
	if err != nil {
		t.Fatalf("Encoding failed for map type: %v", err)
	}
	t.Log("Encoded map type:", b)

	bu2, err := cityPopulationModel.Decode(b)
	if err != nil {
		t.Fatalf("Decoding failed for map type: %v", err)
	}
	t.Log("Decoded map type:", bu2)
}

func TestModelType(t *testing.T) {
	nested := NewModel("nested",
		NewField(0, "id", Int32Type),
	)
	main := NewModel("main",
		NewField(0, "nested", NewModelType(nested)),
	)

	bu := map[string]any{
		"nested": map[string]any{
			"id": 42,
		},
	}

	b, err := main.Encode(bu)
	if err != nil {
		t.Fatalf("Encoding failed for model type: %v", err)
	}
	t.Log("Encoded model type:", b)

	decoded, err := main.Decode(b)
	if err != nil {
		t.Fatalf("Decoding failed for model type: %v", err)
	}
	t.Log("Decoded model type:", decoded)
}

func TestMatrix(t *testing.T) {

	model := NewModel("a",
		NewField(3, "a", NewListType(NewListType(BoolType))),
	)

	bu := map[string]any{
		"a": [][]bool{
			{true, false},
			{false, true},
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

func TestEmptyFields(t *testing.T) {
	person := NewModel("person",
		NewField(0, "name", StringType),
		NewField(1, "age", Int8Type),
		NewField(2, "height", Float64Type),
	)

	bu := map[string]any{}

	b, err := person.Encode(bu)
	if err != nil {
		t.Fatalf("Encoding failed for empty fields: %v", err)
	}
	t.Log("Encoded empty fields:", b)

	bu2, err := person.Decode(b)
	if err != nil {
		t.Fatalf("Decoding failed for empty fields: %v", err)
	}
	t.Log("Decoded empty fields:", bu2)
}

func TestInvalidFieldType(t *testing.T) {
	person := NewModel("person",
		NewField(0, "name", StringType),
		NewField(1, "age", Int8Type),
	)

	bu := map[string]any{
		"name": 12345, // Invalid type
		"age":  25,
	}

	_, err := person.Encode(bu)
	if err == nil {
		t.Fatal("Expected error for invalid field type, but got none")
	}
	t.Log("Received expected error for invalid field type:", err)
}

func TestDeeplyNestedStructure(t *testing.T) {
	nestedModel := NewModel("nested",
		NewField(0, "id", Int32Type),
		NewField(1, "attributes", NewListType(NewListType(StringType))),
	)

	mainModel := NewModel("main",
		NewField(0, "details", NewModelType(nestedModel)),
	)

	bu := map[string]any{
		"details": map[string]any{
			"id": 42,
			"attributes": [][]string{
				{"attr1", "attr2"},
				{"attr3"},
			},
		},
	}

	b, err := mainModel.Encode(bu)
	if err != nil {
		t.Fatalf("Encoding failed for deeply nested structure: %v", err)
	}
	t.Log("Encoded deeply nested structure:", b)

	bu2, err := mainModel.Decode(b)
	if err != nil {
		t.Fatalf("Decoding failed for deeply nested structure: %v", err)
	}
	t.Log("Decoded deeply nested structure:", bu2)
}

func TestNullValues(t *testing.T) {
	person := NewModel("person",
		NewField(0, "name", StringType),
		NewField(1, "age", Int8Type),
	)

	bu := map[string]any{
		"name": nil, // Null value
		"age":  nil, // Null value
	}

	_, err := person.Encode(bu)
	if err == nil {
		t.Fatal("Expected error for invalid field type, but got none")
	}
	t.Log("Received expected error for invalid field type:", err)
}

func TestListWithMixedTypes(t *testing.T) {
	person := NewModel("person",
		NewField(0, "hobbies", NewListType(StringType)),
	)

	bu := map[string]any{
		"hobbies": []any{"reading", 123, true}, // Mixed types
	}

	_, err := person.Encode(bu)
	if err == nil {
		t.Fatal("Expected error for list with mixed types, but got none")
	}
	t.Log("Received expected error for list with mixed types:", err)
}

func TestMapTypeInvalidKey(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic for invalid map key type, but got none")
		}
		t.Log("Received expected panic for invalid map key type:", r)
	}()

	debtModel := NewModel("debts",
		NewField(0, "debts", NewMapType(NewListType(StringType), Int32Type)), // Invalid key type
	)

	_ = debtModel
}

func TestMapTypeMixedValueTypes(t *testing.T) {
	cityPopulationModel := NewModel("cityPopulation",
		NewField(0, "populations", NewMapType(StringType, Int32Type)),
	)

	bu := map[string]any{
		"populations": map[string]any{
			"New York":    8419000,
			"Los Angeles": "invalid type",
			"Chicago":     2716000,
		},
	}

	_, err := cityPopulationModel.Encode(bu)
	if err == nil {
		t.Fatal("Expected error for mixed value types in map, but got none")
	}
	t.Log("Received expected error for mixed value types in map:", err)
}

func TestNewFieldValid(t *testing.T) {
	field := NewField(0, "name", StringType)
	if field.index != 0 || field.label != "name" || field.fieldType != StringType {
		t.Fatalf("Field not created as expected: %+v", field)
	}
}

func TestNewFieldInvalidIndex(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for invalid index, but got none")
		}
	}()
	NewField(256, "name", StringType) // Out of valid range
}

func TestNewFieldEmptyLabel(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for empty label, but got none")
		}
	}()
	NewField(0, "", StringType) // Empty label
}

func TestNewModelValid(t *testing.T) {
	model := NewModel("person",
		NewField(0, "name", StringType),
		NewField(1, "age", Int8Type),
	)
	if model.name != "person" || len(model.schema) != 2 {
		t.Fatalf("Model not created as expected: %+v", model)
	}
}

func TestNewModelDuplicateLabel(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for duplicate label, but got none")
		}
	}()
	NewModel("person",
		NewField(0, "name", StringType),
		NewField(1, "name", Int8Type), // Duplicate label
	)
}

func TestNewModelDuplicateIndex(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for duplicate index, but got none")
		}
	}()
	NewModel("person",
		NewField(0, "name", StringType),
		NewField(0, "age", Int8Type), // Duplicate index
	)
}

func TestNewModelEmptyName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic for empty model name, but got none")
		}
	}()
	NewModel("",
		NewField(0, "name", StringType),
	)
}
