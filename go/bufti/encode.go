package bufti

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

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

func encodeValue(buf *bytes.Buffer, value any, valType BuftiType) error {
	switch valType {
	case Int8Type:
		v, ok := value.(int8)
		if ok {
			return binary.Write(buf, binary.BigEndian, v)
		}
		v2, ok := value.(int)
		if ok && isInRange(float64(v2), -128, 127) {
			return binary.Write(buf, binary.BigEndian, int8(v2))
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Int16Type:
		v, ok := value.(int16)
		if ok {
			return binary.Write(buf, binary.BigEndian, v)
		}
		v2, ok := value.(int)
		if ok && isInRange(float64(v2), -32768, 32767) {
			return binary.Write(buf, binary.BigEndian, int16(v2))
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Int32Type:
		v, ok := value.(int32)
		if ok {
			return binary.Write(buf, binary.BigEndian, v)
		}
		v2, ok := value.(int)
		if ok && isInRange(float64(v2), -2147483648, 2147483647) {
			return binary.Write(buf, binary.BigEndian, int32(v2))
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Int64Type:
		v, ok := value.(int64)
		if ok {
			return binary.Write(buf, binary.BigEndian, v)
		}
		v2, ok := value.(int)
		if ok {
			return binary.Write(buf, binary.BigEndian, int64(v2))
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Float32Type:
		v, ok := value.(float32)
		if ok {
			return binary.Write(buf, binary.BigEndian, v)
		}
		v2, ok := value.(float64)
		if ok && isInRange(v2, -3.4e38, 3.4e38) {
			return binary.Write(buf, binary.BigEndian, float32(v2))
		}
		return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
	case Float64Type:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}
		return binary.Write(buf, binary.BigEndian, v)
	case BoolType:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}
		return binary.Write(buf, binary.BigEndian, v)
	case StringType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}
		if err := binary.Write(buf, binary.BigEndian, uint16(len(v))); err != nil {
			return err
		}
		_, err := buf.Write([]byte(v))
		return err
	}

	listType, isList := isListType(valType)
	if isList {
		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}

		if err := binary.Write(buf, binary.BigEndian, uint16(val.Len())); err != nil {
			return err
		}

		for i := range val.Len() {
			if err := encodeValue(buf, val.Index(i).Interface(), listType); err != nil {
				return err
			}
		}
		return nil
	}

	keyType, valueType, isList := isMapType(valType)
	if isList {
		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Map {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}

		if err := binary.Write(buf, binary.BigEndian, uint16(val.Len())); err != nil {
			return err
		}

		for _, key := range val.MapKeys() {
			if err := encodeValue(buf, key.Interface(), keyType); err != nil {
				return err
			}
			if err := encodeValue(buf, val.MapIndex(key).Interface(), valueType); err != nil {
				return err
			}
		}
		return nil
	}

	modelName, isModel := isModelType(valType)
	if isModel {
		model, exists := registeredModels[modelName]
		if !exists {
			return fmt.Errorf("%w: can not find the model %s", ErrModel, modelName)
		}
		bu, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("%w: can not apply value of type %T to %s", ErrBufti, value, valType)
		}

		if err := binary.Write(buf, binary.BigEndian, uint16(len(bu))); err != nil {
			return err
		}

		p, err := model.Encode(bu)
		if err != nil {
			return err
		}
		buf.Write(p)
		return nil
	}

	return fmt.Errorf("%w: invalid schema type (%s)", ErrModel, valType)
}
