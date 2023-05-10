package test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

var bitfieldRefBytes = []byte{0x74, 0x21, 0x43, 0x6d, 0xb0, 0x01, 0xff, 0x01}
var bitfieldRef = &BitfieldExample{A: 0x2, B: 0x7, C2: []byte{0x01, 0x02, 0x03, 0x04}, D: 0xd, E: 0x06, F: []bool{false, true, true, false, true, true}, G: 0xff, H: true}

func TestBitfieldTransCodeDynamic(t *testing.T) {
	for a := 0; a < 16; a++ {
		for b := 0; b < 32; b++ {
			//build a bitfield with two dynamic bitfield slices followed by a byte aligned type, followed by single bit
			ref := &BitfieldExample{C2: make([]byte, a), F: make([]bool, b), G: 0x04, H: true}
			if b > 4 { // set an arbitrary bit
				ref.F[b/3] = true
			}
			if a > 4 { // set an arbitrary value in C2
				ref.C2[a/3] = 0xc
			}
			buf := make([]byte, ref.SizeOf())
			//manual length check verification(it's in fact surprisingly hard)
			x := (a*4 + b) - 5 // a * 4bit + b*1bit and there are 5 padding bit in the empty structure available both to a + b
			if x < 0 {
				x = 0
			} else {
				if x%8 != 0 { // at least one bit of F has overflown (needs full extra byte as G is byte aligned)
					x = x/8 + 1
				} else { // exact fit. G smoothly fits byte-aligned right after F ends
					x = x / 8
				}
			}
			if ref.SizeOf() != 5+x {
				fmt.Printf("got: %#v\n", ref)
				t.Fatalf("wrong len %d %d", ref.SizeOf(), 5+x)
			}

			ref.MarshalBinary(buf)
			if buf[len(buf)-1] != 0x01 {
				//check that the true bool at the end is in fact there in the binary buffer
				t.Fatal("decode failed buf[len(buf)-1]")
			}
			if buf[len(buf)-2] != 0x04 {
				//check that the 0x04 byte is byte aligned always and second-to-last
				t.Fatal("decode failed buf[len(buf)-2]")
			}

			have := &BitfieldExample{C2: make([]byte, 0), F: make([]bool, 0)}
			have.UnmarshalBinary(buf)
			if !reflect.DeepEqual(ref, have) {
				fmt.Printf("got: %#v\nwant: %#v\n", have, ref)
				t.Fatal("decode failed")
			}
		}

	}

}

func TestBitfieldLengthCheckDynamic(t *testing.T) {
	for a := 0; a < 16; a++ {
		for b := 0; b < 32; b++ {
			//build a bitfield with two dynamic bitfield slices followed by a byte aligned type, followed by single bit
			ref := &BitfieldExample{C2: make([]byte, a), F: make([]bool, b), G: 0x04, H: true}
			if b > 4 { // set an arbitrary bit
				ref.F[b/3] = true
			}
			if a > 4 { // set an arbitrary value in C2
				ref.C2[a/3] = 0xc
			}
			buf := make([]byte, ref.SizeOf())
			ref.MarshalBinary(buf)
			for trunc := 0; trunc < len(buf); trunc++ {
				have := &BitfieldExample{C2: make([]byte, 0), F: make([]bool, 0)}
				if have.UnmarshalBinary(buf[:trunc]) != 0 {
					t.Fatal("truncation detection failed")
				}
			}
		}

	}

}

func TestBitfieldEncode(t *testing.T) {
	buf := make([]byte, len(bitfieldRefBytes))
	if l := bitfieldRef.MarshalBinary(buf); l != len(bitfieldRefBytes) {
		t.Fatalf("got different number of bytes as expected %d <> %d", l, len(bitfieldRefBytes))
	}
	if !bytes.Equal(buf, bitfieldRefBytes) {
		fmt.Printf("got: %#v\nwant: %#v\n", buf, bitfieldRefBytes)
		t.Fatal("encode failed")
	}
}
func TestBitfieldSizeOf(t *testing.T) {
	if l := bitfieldRef.SizeOf(); l != len(bitfieldRefBytes) {
		t.Fatalf("got different number of bytes as expected %d <> %d", l, len(bitfieldRefBytes))
	}
}

func TestBitfieldDecode(t *testing.T) {

	out := &BitfieldExample{}
	if l := out.UnmarshalBinary(bitfieldRefBytes); l != len(bitfieldRefBytes) {
		t.Fatalf("got different number of bytes as expected %d %d", l, len(bitfieldRefBytes))
	}
	if !reflect.DeepEqual(bitfieldRef, out) {
		fmt.Printf("got: %#v\nwant: %#v\n", out, bitfieldRef)
		t.Fatal("decode failed")
	}
}

func BenchmarkMarshal__BitfieldArray(b *testing.B) {
	val := &BitfieldExampleBench{}
	buf := make([]byte, val.SizeOf())
	for i := 0; i < b.N; i++ {
		val.MarshalBinary(buf)
	}
}
func BenchmarkUnmarshal__BitfieldArray(b *testing.B) {
	val := &BitfieldExampleBench{}
	buf := make([]byte, val.SizeOf())
	val.MarshalBinary(buf)
	for i := 0; i < b.N; i++ {
		val.UnmarshalBinary(buf)
	}
}
func BenchmarkMarshal__ByteArray(b *testing.B) {
	val := &ByteExampleBench{}
	buf := make([]byte, val.SizeOf())
	for i := 0; i < b.N; i++ {
		val.MarshalBinary(buf)
	}
}
func BenchmarkUnmarshal__ByteArray(b *testing.B) {
	val := &ByteExampleBench{}
	buf := make([]byte, val.SizeOf())
	val.MarshalBinary(buf)
	for i := 0; i < b.N; i++ {
		val.UnmarshalBinary(buf)
	}
}
