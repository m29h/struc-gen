package test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

var bitfieldRefBytes = []byte{0x7a, 0x01, 0x02, 0xed, 0x17}
var bitfieldRef = &BitfieldExample{A: 0xa, B: 0x7, C2: [2]byte{0x01, 0x02}, D: 0xd, E: 0x7e, F: [2]bool{false, true}}

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
