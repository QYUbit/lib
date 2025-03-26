package butil

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Butil struct {
	buf          []byte
	cursor       *int
	littleEndian bool
}

func NewButil(b []byte, littleEndian bool) *Butil {
	cursor := 0
	return &Butil{
		buf:          b,
		cursor:       &cursor,
		littleEndian: littleEndian,
	}
}

func (b *Butil) String() string {
	return fmt.Sprintf("%s", b.buf)
}

func (b *Butil) Bytes() []byte {
	return b.buf
}

func (b *Butil) IsAtEnd() bool {
	return *b.cursor == len(b.buf)
}

func (b *Butil) ReadBytes(size int) ([]byte, error) {
	if *b.cursor+size > len(b.buf) {
		return nil, io.EOF
	}
	p := b.buf[*b.cursor : *b.cursor+size]
	*b.cursor += size
	return p, nil
}

func (b *Butil) ReadUint8() (uint8, error) {
	p, err := b.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return uint8(p[0]), nil
}

func (b *Butil) ReadUint16() (uint16, error) {
	p, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	if b.littleEndian {
		return binary.LittleEndian.Uint16(p), nil
	} else {
		return binary.BigEndian.Uint16(p), nil
	}
}

func (b *Butil) ReadUint32() (uint32, error) {
	p, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	if b.littleEndian {
		return binary.LittleEndian.Uint32(p), nil
	} else {
		return binary.BigEndian.Uint32(p), nil
	}
}

func (b *Butil) ReadUint64() (uint64, error) {
	p, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	if b.littleEndian {
		return binary.LittleEndian.Uint64(p), nil
	} else {
		return binary.BigEndian.Uint64(p), nil
	}
}

func (b *Butil) ReadChars(length int) (string, error) {
	p, err := b.ReadBytes(length)
	return string(p), err
}

func (b *Butil) ReadBool() (bool, error) {
	p, err := b.ReadBytes(1)
	if err != nil {
		return false, err
	}
	return p[0] == 1, nil
}

func (b *Butil) WriteUint8(v uint8) {
	b.buf = append(b.buf, v)
}

func (b *Butil) WriteUint16(v uint16) {
	if b.littleEndian {
		binary.LittleEndian.AppendUint16(b.buf, v)
	} else {
		binary.BigEndian.AppendUint16(b.buf, v)
	}
}

func (b *Butil) WriteUint32(v uint32) {
	if b.littleEndian {
		binary.LittleEndian.AppendUint32(b.buf, v)
	} else {
		binary.BigEndian.AppendUint32(b.buf, v)
	}
}

func (b *Butil) WriteUint64(v uint64) {
	if b.littleEndian {
		binary.LittleEndian.AppendUint64(b.buf, v)
	} else {
		binary.BigEndian.AppendUint64(b.buf, v)
	}
}

func (b *Butil) WriteString(a string) {
	b.buf = append(b.buf, []byte(a)...)
}

func (b *Butil) WriteBool(v bool) {
	val := 0
	if v {
		val = 1
	}
	b.buf = append(b.buf, byte(val))
}
