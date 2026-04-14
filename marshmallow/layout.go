package marshmallow

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Layout interface {
	Encode(value any) ([]byte, error)
	Decode(data []byte, output any) error
}

type BinaryLayout struct{}

func NewBinaryLayout() BinaryLayout {
	return BinaryLayout{}
}

func (b BinaryLayout) Encode(value any) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("encode failed: %w", err)
	}
	return buf.Bytes(), nil
}

func (b BinaryLayout) Decode(data []byte, output any) error {
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, output); err != nil {
		return fmt.Errorf("decode failed: %w", err)
	}
	return nil
}
