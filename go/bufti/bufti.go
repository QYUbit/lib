package bufti

import (
	"errors"
	"fmt"
)

var (
	ErrNilSlice = errors.New("bytes slice is nil")
	ErrFormat   = errors.New("unexpected buffer format")
	ErrModel    = errors.New("invalid bufti model")
	ErrBufti    = errors.New("unexpected bufti map format")
)

var registeredModels = make(map[string]*Model)

type BuftiType string

const (
	Int8Type    BuftiType = "int8"
	Int16Type   BuftiType = "int16"
	Int32Type   BuftiType = "int32"
	Int64Type   BuftiType = "int64"
	Float32Type BuftiType = "float32"
	Float64Type BuftiType = "float64"
	BoolType    BuftiType = "bool"
	StringType  BuftiType = "string"
)

// Creates a new BuftiListType based on the given element type.
func NewListType(elementType BuftiType) BuftiType {
	return BuftiType(fmt.Sprintf("[%s]", elementType))
}

// Creates a new BuftiModelType with the specified reference model.
func NewModelType(model *Model) BuftiType {
	return BuftiType(fmt.Sprintf("*%s", model.name))
}

type Field struct {
	index     byte
	label     string
	fieldType BuftiType
}

// Creates a new model field based on index, label and type. Index has to be between 0 and 255. Panics when given unexpected inputs.
func NewField(index int, label string, fieldType BuftiType) Field {
	if !isInRange(float64(index), 0, 255) {
		panic("index has to be between 0 and 255")
	}
	if label == "" {
		panic("label must not be empty")
	}
	return Field{
		index:     byte(index),
		label:     label,
		fieldType: fieldType,
	}
}

type fieldSchema struct {
	label     string
	fieldType BuftiType
}

type Model struct {
	name   string
	schema map[byte]fieldSchema
	labels map[string]byte
}

// Creates a new model which represents the way data gets en/decoded. Model name has to be unique. Panics when given unexpected inputs.
func NewModel(name string, fields ...Field) *Model {
	if name == "" {
		panic("model name must not be empty")
	}
	if _, exists := registeredModels[name]; exists {
		panic(fmt.Sprintf("model with name \"%s\" already exists", name))
	}

	m := &Model{
		name:   name,
		schema: make(map[byte]fieldSchema),
		labels: make(map[string]byte),
	}
	registeredModels[name] = m

	for _, f := range fields {
		if _, exists := m.labels[f.label]; exists {
			panic(fmt.Sprintf("duplicate label %s in model %s", f.label, m.name))
		}
		if _, exists := m.schema[f.index]; exists {
			panic(fmt.Sprintf("duplicate index %d in model %s", f.index, m.name))
		}

		m.labels[f.label] = f.index
		m.schema[f.index] = fieldSchema{
			label:     f.label,
			fieldType: f.fieldType,
		}
	}
	return m
}

// Returns the string representation of the model.
func (m *Model) String() string {
	return fmt.Sprintf("model %s %v", m.name, m.schema)
}
