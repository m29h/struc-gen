package test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

var boRefBytes = []byte{0x21, 0x74, 0x43, 0x6d}
var boRef = &DefaultBO{L: 0x7421, B: 0x436d}

func TestDefaultEncode(t *testing.T) {
	buf := make([]byte, len(boRefBytes))
	if l := boRef.MarshalBinary(buf); l != len(boRefBytes) {
		t.Fatalf("got different number of bytes as expected %d <> %d", l, len(boRefBytes))
	}
	if !bytes.Equal(buf, boRefBytes) {
		fmt.Printf("got: %#v\nwant: %#v\n", buf, boRefBytes)
		t.Fatal("encode failed")
	}
}

func TestDefaultDecode(t *testing.T) {

	out := &DefaultBO{}
	if l := out.UnmarshalBinary(boRefBytes); l != len(boRefBytes) {
		t.Fatalf("got different number of bytes as expected %d %d", l, len(boRefBytes))
	}
	if !reflect.DeepEqual(boRef, out) {
		fmt.Printf("got: %#v\nwant: %#v\n", out, boRef)
		t.Fatal("decode failed")
	}
}
