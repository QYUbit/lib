package bufti

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"slices"
)

var (
	ErrNilSlice = errors.New("bytes slice has the value nil")
	ErrFormat   = errors.New("unexpected binary format")
	ErrModel    = errors.New("invalid bufti model")
	ErrBufti    = errors.New("unexpected bufti format")
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

var variableSizeTypes = []BuftiType{StringType}

type Field struct {
	index     byte
	label     string
	fieldType BuftiType
}

func NewField(index int, label string, fieldType BuftiType) Field {
	if index < 0 || index > 255 {
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
	schema map[byte]FieldSchema
	labels map[string]byte
}

func NewModel(fields ...Field) *Model {
	m := &Model{
		schema: make(map[byte]FieldSchema),
		labels: make(map[string]byte),
	}
	for _, f := range fields {
		if _, exists := m.labels[f.label]; exists {
			panic(fmt.Sprintf("multiple lables with the same value (%s)", f.label))
		}
		if _, exists := m.schema[f.index]; exists {
			panic(fmt.Sprintf("multiple lables with the same value (%d)", f.index))
		}

		m.labels[f.label] = f.index
		fs := FieldSchema{label: f.label, fieldType: f.fieldType}
		m.schema[f.index] = fs
	}
	return m
}

func (m *Model) Encode(bu map[string]any) ([]byte, error) {
	buf := make([]byte, 0)

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

		buf = append(buf, byte(index))

		if err := encodeValue(&buf, value, valType); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func (m *Model) Decode(b []byte) (map[string]any, error) {
	if b == nil {
		return nil, ErrNilSlice
	}

	bufti := make(map[string]any)
	cursor := 0

	for cursor < len(b) {
		index, err := readBytes(b, &cursor, 1)
		if err != nil {
			return nil, err
		}

		schemaField, exists := m.schema[index[0]]
		if !exists {
			return nil, fmt.Errorf("%w: index not found (%d)", ErrFormat, index[0])
		}
		valType := schemaField.fieldType
		label := schemaField.label

		var size int
		if slices.Contains(variableSizeTypes, valType) {
			p, err := readBytes(b, &cursor, 2)
			if err != nil {
				return nil, err
			}
			size = int(binary.BigEndian.Uint16(p))
		}

		value, err := decodeValue(b, &cursor, valType, size)
		if err != nil {
			return nil, err
		}

		bufti[label] = value
	}

	return bufti, nil
}

func encodeValue(b *[]byte, value any, valType BuftiType) error {
	switch valType {
	case Int8Type:
		v, ok := value.(int8)
		if ok {
			*b = append(*b, itob(v)...)
			return nil
		}
		v2, ok := value.(int)
		if ok {
			*b = append(*b, itob(int8(v2))...)
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Int16Type:
		v, ok := value.(int16)
		if ok {
			*b = append(*b, itob(v)...)
			return nil
		}
		v2, ok := value.(int)
		if ok {
			*b = append(*b, itob(int16(v2))...)
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Int32Type:
		v, ok := value.(int32)
		if ok {
			*b = append(*b, itob(v)...)
			return nil
		}
		v2, ok := value.(int)
		if ok {
			*b = append(*b, itob(int32(v2))...)
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Int64Type:
		v, ok := value.(int64)
		if ok {
			*b = append(*b, itob(v)...)
			return nil
		}
		v2, ok := value.(int)
		if ok {
			*b = append(*b, itob(int64(v2))...)
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Float32Type:
		v, ok := value.(float32)
		if ok {
			*b = append(*b, itob(v)...)
			return nil
		}
		v2, ok := value.(float64)
		if ok {
			*b = append(*b, itob(float32(v2))...)
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Float64Type:
		v, ok := value.(float64)
		if ok {
			*b = append(*b, itob(v)...)
			return nil
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case BoolType:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}
		if v {
			*b = append(*b, 1)
		} else {
			*b = append(*b, 0)
		}
	case StringType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}
		*b = binary.BigEndian.AppendUint16(*b, uint16(len(v)))
		*b = append(*b, []byte(v)...)
	default:
		return fmt.Errorf("%w: invalid schema type (%s)", ErrModel, valType)
	}
	return nil
}

func decodeValue(b []byte, cursor *int, valType BuftiType, size int) (any, error) {
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

func itob(v any) []byte {
	byteBuffer := bytes.NewBuffer([]byte{})
	binary.Write(byteBuffer, binary.BigEndian, v)
	return byteBuffer.Bytes()
}
