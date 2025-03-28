package bufti

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// Decode decodes the specified byte array into a map.
func (m *Model) Decode(b []byte) (map[string]any, error) {
	if b == nil {
		return nil, ErrNilSlice
	}
	buf := bytes.NewBuffer(b)
	return m.decode(buf, len(b))
}

func (m *Model) decode(buf *bytes.Buffer, limit int) (map[string]any, error) {
	bufti := make(map[string]any)

	for range limit {
		var index byte
		if err := binary.Read(buf, binary.BigEndian, &index); err != nil {
			break
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return nil, fmt.Errorf("%w: index not found (%d)", ErrFormat, index)
		}
		valType := schemaField.fieldType
		label := schemaField.label

		value, err := decodeValue(buf, valType)
		if err != nil {
			return nil, err
		}

		bufti[label] = value
	}

	return bufti, nil
}

func decodeValue(buf *bytes.Buffer, valType BuftiType) (any, error) {
	var size int
	if valType == "string" || (strings.HasPrefix(string(valType), "[") && strings.HasSuffix(string(valType), "]")) || strings.HasPrefix(string(valType), "*") {
		var v uint16
		if err := binary.Read(buf, binary.BigEndian, &v); err != nil {
			return nil, err
		}
		size = int(v)
	}

	switch valType {
	case Int8Type:
		var v int8
		err := binary.Read(buf, binary.BigEndian, &v)
		return v, err
	case Int16Type:
		var v int16
		err := binary.Read(buf, binary.BigEndian, &v)
		return v, err
	case Int32Type:
		var v int32
		err := binary.Read(buf, binary.BigEndian, &v)
		return v, err
	case Int64Type:
		var v int64
		err := binary.Read(buf, binary.BigEndian, &v)
		return v, err
	case Float32Type:
		var v float32
		err := binary.Read(buf, binary.BigEndian, &v)
		return v, err
	case Float64Type:
		var v float64
		err := binary.Read(buf, binary.BigEndian, &v)
		return v, err
	case BoolType:
		var v bool
		err := binary.Read(buf, binary.BigEndian, &v)
		return v, err
	case StringType:
		p := make([]byte, size)
		if _, err := buf.Read(p); err != nil {
			return nil, err
		}
		return string(p), nil
	default:
		listType, isList := isListType(valType)
		if isList {
			var list []any
			for range size {
				item, err := decodeValue(buf, listType)
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

			bu, err := model.decode(buf, size)
			if err != nil {
				return nil, err
			}
			return bu, nil
		}

		return nil, fmt.Errorf("%w: invalid type (%s)", ErrFormat, valType)
	}
}
