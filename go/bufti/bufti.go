package bufti

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

var (
	ErrNilSlice = errors.New("bytes slice is nil")
	ErrFormat   = errors.New("invalid binary format")
	ErrModel    = errors.New("invalid bufti model")
	ErrBufti    = errors.New("invalid bufti map format")
)

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

var registeredModels = make(map[string]*Model)

type Field struct {
	index     byte
	label     string
	fieldType BuftiType
}

// Creates a new model field based on index, label and type.
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

type FieldSchema struct {
	label     string
	fieldType BuftiType
}

type Model struct {
	name   string
	schema map[byte]FieldSchema
	labels map[string]byte
}

// Creates a new model which represents the way data gets en/decoded.
func NewModel(name string, fields ...Field) *Model {
	if name == "" {
		panic("model name must not be empty")
	}
	if _, exists := registeredModels[name]; exists {
		panic(fmt.Sprintf("model with name \"%s\" already exists", name))
	}

	m := &Model{
		name:   name,
		schema: make(map[byte]FieldSchema),
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
		m.schema[f.index] = FieldSchema{
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

// Encode encodes the provided map into a byte array.
func (m *Model) Encode(bu map[string]any) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	buf.Grow(2 * len(bu))

	for label, value := range bu {
		index, exists := m.labels[label]
		if !exists {
			return nil, fmt.Errorf("%w: label not found (%s)", ErrModel, label)
		}
		schemaField, exists := m.schema[index]
		if !exists {
			return nil, fmt.Errorf("%w: index not found (%d)", ErrModel, index)
		}
		valType := schemaField.fieldType

		if err := buf.WriteByte(index); err != nil {
			return nil, err
		}

		if err := encodeValue(buf, value, valType); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// Decode decodes the specified byte array into a map.
func (m *Model) Decode(b []byte) (map[string]any, error) {
	cursor := 0
	return m.decode(b, &cursor, len(b))
}

func (m *Model) decode(b []byte, cursor *int, limit int) (map[string]any, error) {
	if b == nil {
		return nil, ErrNilSlice
	}

	bufti := make(map[string]any)

	for range limit {
		if *cursor >= len(b) {
			break
		}

		index, err := readBytes(b, cursor, 1)
		if err != nil {
			return nil, err
		}

		schemaField, exists := m.schema[index[0]]
		if !exists {
			return nil, fmt.Errorf("%w: index not found (%d)", ErrFormat, index[0])
		}
		valType := schemaField.fieldType
		label := schemaField.label

		value, err := decodeValue(b, cursor, valType)
		if err != nil {
			return nil, err
		}

		bufti[label] = value
	}

	return bufti, nil
}

func encodeValue(buf *bytes.Buffer, value any, valType BuftiType) error {
	switch valType {
	case Int8Type:
		v, ok := value.(int8)
		if ok {
			encodeNumber(buf, v)
			return nil
		}
		v2, ok := value.(int)
		if ok && isInRange(float64(v2), -128, 127) {
			encodeNumber(buf, int8(v2))
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Int16Type:
		v, ok := value.(int16)
		if ok {
			encodeNumber(buf, v)
			return nil
		}
		v2, ok := value.(int)
		if ok && isInRange(float64(v2), -32768, 32767) {
			encodeNumber(buf, int16(v2))
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Int32Type:
		v, ok := value.(int32)
		if ok {
			encodeNumber(buf, v)
			return nil
		}
		v2, ok := value.(int)
		if ok && isInRange(float64(v2), -2147483648, 2147483647) {
			encodeNumber(buf, int32(v2))
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Int64Type:
		v, ok := value.(int64)
		if ok {
			encodeNumber(buf, v)
			return nil
		}
		v2, ok := value.(int)
		if ok {
			encodeNumber(buf, int64(v2))
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Float32Type:
		v, ok := value.(float32)
		if ok {
			encodeNumber(buf, v)
			return nil
		}
		v2, ok := value.(float64)
		if ok {
			encodeNumber(buf, float32(v2))
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Float64Type:
		v, ok := value.(float64)
		if ok {
			encodeNumber(buf, v)
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case BoolType:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}
		if v {
			buf.WriteByte(1)
		} else {
			buf.WriteByte(0)
		}
	case StringType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}
		buf.Write(binary.BigEndian.AppendUint16([]byte{}, uint16(len(v))))
		buf.Write([]byte(v))
	default:
		listType, isList := getListType(valType)
		if isList {
			val := reflect.ValueOf(value)
			if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
				return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
			}
			buf.Write(binary.BigEndian.AppendUint16([]byte{}, uint16(val.Len())))

			for i := range val.Len() {
				if err := encodeValue(buf, val.Index(i).Interface(), listType); err != nil {
					return err
				}
			}
			return nil
		}

		modelName, isModel := strings.CutPrefix(string(valType), "*")
		if isModel {
			model, exists := registeredModels[modelName]
			if !exists {
				return fmt.Errorf("%w: can not find the model %s", ErrModel, modelName)
			}
			bu, ok := value.(map[string]any)
			if !ok {
				return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
			}
			buf.Write(binary.BigEndian.AppendUint16([]byte{}, uint16(len(bu))))

			p, err := model.Encode(bu)
			if err != nil {
				return err
			}
			buf.Write(p)
			return nil
		}

		return fmt.Errorf("%w: invalid schema type (%s)", ErrModel, valType)
	}
	return nil
}

func decodeValue(b []byte, cursor *int, valType BuftiType) (any, error) {
	var size int
	if valType == "string" || (strings.HasPrefix(string(valType), "[") && strings.HasSuffix(string(valType), "]")) || strings.HasPrefix(string(valType), "*") {
		p, err := readBytes(b, cursor, 2)
		if err != nil {
			return nil, err
		}
		size = int(binary.BigEndian.Uint16(p))
	}

	switch valType {
	case Int8Type:
		p, err := readBytes(b, cursor, 1)
		if err != nil {
			return nil, err
		}
		return int8(p[0]), nil
	case Int16Type:
		var v int16
		p, err := readBytes(b, cursor, 2)
		if err != nil {
			return nil, err
		}
		btoi(p, &v)
		return v, nil
	case Int32Type:
		var v int32
		p, err := readBytes(b, cursor, 4)
		if err != nil {
			return nil, err
		}
		btoi(p, &v)
		return v, nil
	case Int64Type:
		var v int64
		p, err := readBytes(b, cursor, 8)
		if err != nil {
			return nil, err
		}
		btoi(p, &v)
		return v, nil
	case Float32Type:
		var v float32
		p, err := readBytes(b, cursor, 4)
		if err != nil {
			return nil, err
		}
		btoi(p, &v)
		return v, nil
	case Float64Type:
		var v float64
		p, err := readBytes(b, cursor, 8)
		if err != nil {
			return nil, err
		}
		btoi(p, &v)
		return v, nil
	case BoolType:
		p, err := readBytes(b, cursor, 1)
		if err != nil {
			return nil, err
		}
		if p[0] != 0 && p[0] != 1 {
			return nil, fmt.Errorf("%w: Bool type can only be 0 or 1, instead %d", ErrFormat, p[0])
		}
		return p[0] == 1, nil
	case StringType:
		p, err := readBytes(b, cursor, size)
		if err != nil {
			return nil, err
		}
		return string(p), nil
	default:
		listType, isList := getListType(valType)
		if isList {
			var list []any
			for range size {
				item, err := decodeValue(b, cursor, listType)
				if err != nil {
					return nil, err
				}
				list = append(list, item)
			}
			return list, nil
		}

		modelName, isModel := strings.CutPrefix(string(valType), "*")
		if isModel {
			model, exists := registeredModels[modelName]
			if !exists {
				return nil, fmt.Errorf("%w: can not find the model %s", ErrModel, modelName)
			}

			bu, err := model.decode(b, cursor, size)
			if err != nil {
				return nil, err
			}
			return bu, nil
		}

		return nil, fmt.Errorf("%w: invalid type (%s)", ErrFormat, valType)
	}
}

func readBytes(b []byte, cursor *int, size int) ([]byte, error) {
	if *cursor+size > len(b) {
		return nil, io.EOF
	}
	p := b[*cursor : *cursor+size]
	*cursor += size
	return p, nil
}

func btoi(b []byte, dest any) {
	byteBuffer := bytes.NewBuffer(b)
	binary.Read(byteBuffer, binary.BigEndian, dest)
}

func encodeNumber[T int | int8 | int16 | int32 | int64 | float32 | float64](buf *bytes.Buffer, value T) error {
	return binary.Write(buf, binary.BigEndian, value)
}

func getListType(valType BuftiType) (BuftiType, bool) {
	s, found := strings.CutPrefix(string(valType), "[")
	if !found {
		return "", false
	}
	after, found := strings.CutSuffix(s, "]")
	return BuftiType(after), found
}

func isInRange(v float64, min float64, max float64) bool {
	return v >= min && v <= max
}
