package marshmallow

import "testing"

type sample struct {
	A uint16
	B uint32
}

func TestBinaryLayoutEncodeDecode(t *testing.T) {
	layout := NewBinaryLayout()
	encoded, err := layout.Encode(sample{A: 1, B: 2})
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	out := sample{}
	if err := layout.Decode(encoded, &out); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if out.A != 1 || out.B != 2 {
		t.Fatalf("decoded mismatch")
	}
	if _, err := layout.Encode(map[string]int{"a": 1}); err == nil {
		t.Fatalf("expected encode error")
	}
	if err := layout.Decode([]byte{0x01}, &out); err == nil {
		t.Fatalf("expected decode error")
	}
}
