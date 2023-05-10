package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"testing"

	"github.com/lunixbochs/struc"
)

var five = 5
var reference = &Example{
	nil,
	1, 2, 3, 4, 5, 6, 7, 8, 0, []byte{'a', 'b', 'c', 'd'},
	9, 10, 11, 12, 13, 14, 15, 16, true, false, [4]byte{'e', 'f', 'g', 'h'},
	20, 21,
	-1,
	4294967295,
	-1,
	10, "ijklmnopqr", "stuv",
	4, "1234",
	4, []byte("5678"),
	7, "ijklmno", "pqrstuv",
	4, []byte("5678"),
	Nested{1}, &Nested{2}, &five,
	6, []Nested{{3}, {4}, {5}, {6}, {7}, {8}},
	0,
	Int3(4), []byte("ABCD"),
}
var referenceBytes = []byte{
	0, 0, 0, 0, 0, // pad(5)
	1, 0, 2, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 4, // fake int8-int64(1-4)
	5, 6, 0, 7, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, // fake little-endian uint8-uint64(5-8)
	0,                  // fake bool(0)
	'a', 'b', 'c', 'd', // fake [4]byte

	9, 0, 10, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0, 12, // real int8-int64(9-12)
	13, 14, 0, 15, 0, 0, 0, 16, 0, 0, 0, 0, 0, 0, 0, // real little-endian uint8-uint64(13-16)
	1, 0, // real bool(1), bool(0)
	'e', 'f', 'g', 'h', // real [4]byte
	65, 160, 0, 0, // real float32(20)
	64, 53, 0, 0, 0, 0, 0, 0, // real float64(21)

	255, 255, 255, 255, // fake int32(-1)
	255, 255, 255, 255, // fake uint32(4294967295)

	255, 255, 255, 255, 255, 255, 255, 255, // fake int64(-1)

	10, 0, 0, 0, // little-endian int32(10) sizeof=Str
	'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', // Str
	's', 't', 'u', 'v', // fake string([4]byte)
	04, '1', '2', '3', '4', // real string
	04, '5', '6', '7', '8', // fake []byte(string)

	7, 0, 0, 0, // little-endian int32(7)
	'i', 'j', 'k', 'l', 'm', 'n', 'o', // Str4a sizefrom=Size4
	'p', 'q', 'r', 's', 't', 'u', 'v', // Str4b sizefrom=Size4
	04, '5', '6', '7', '8', // fake []byte(string)

	1, 2, // Nested{1}, Nested{2}
	0, 0, 0, 0, 0, 0, 0, 5, // &five

	0, 0, 0, 6, // int32(6)
	3, 4, 5, 6, 7, 8, // [Nested{3}, ...Nested{8}]

	0, 0, 4, 'A', 'B', 'C', 'D', // Int3(4), []byte("ABCD")
}

func TestEncode(t *testing.T) {
	buf := make([]byte, len(referenceBytes))
	if l := reference.MarshalBinary(buf); l != len(referenceBytes) {
		t.Fatal("got different number of bytes as expected")
	}
	if !bytes.Equal(buf, referenceBytes) {
		fmt.Printf("got: %#v\nwant: %#v\n", buf, referenceBytes)
		t.Fatal("encode failed")
	}
}

func TestSizeOf(t *testing.T) {
	if l := reference.SizeOf(); l != len(referenceBytes) {
		t.Fatal("got different SizeOf() of bytes as expected")
	}
}
func TestDecode(t *testing.T) {
	out := &Example{}
	if l := out.UnmarshalBinary(referenceBytes); l != len(referenceBytes) {
		t.Fatalf("got different number of bytes as expected %d %d", l, len(referenceBytes))
	}
	if !reflect.DeepEqual(reference, out) {
		fmt.Printf("got: %#v\nwant: %#v\n", out, reference)
		t.Fatal("decode failed")
	}
}
func TestValidate(t *testing.T) {
	out := &Example{}
	//with truncated referencebytes l must be != len(refBytes) (no read out of bounds allowed!)
	for truncate := 0; truncate < len(referenceBytes); truncate++ {
		if l := out.UnmarshalBinary(referenceBytes[:truncate]); l == len(referenceBytes) {
			t.Fatalf("Unmarshal went further than allowed %d %d", l, len(referenceBytes))
		}
	}
}

func BenchmarkMarshal__strucgen(b *testing.B) {
	buf := make([]byte, len(referenceBytes))
	for i := 0; i < b.N; i++ {
		a := reference
		a.MarshalBinary(buf)
		if !bytes.Equal(buf, referenceBytes) {
			fmt.Printf("got: %#v\nwant: %#v\n", buf, referenceBytes)
			b.Fatal("encode failed")
		}
	}
}
func BenchmarkUnmarshal__strucgen(b *testing.B) {
	a := &Example{}
	for i := 0; i < b.N; i++ {
		a.UnmarshalBinary(referenceBytes)
	}
}

func (i *Int3) Pack(p []byte, opt *struc.Options) (int, error) {
	var tmp [4]byte
	binary.BigEndian.PutUint32(tmp[:], uint32(*i))
	copy(p, tmp[1:])
	return 3, nil
}
func (i *Int3) Unpack(r io.Reader, length int, opt *struc.Options) error {
	var tmp [4]byte
	if _, err := r.Read(tmp[1:]); err != nil {
		return err
	}
	*i = Int3(binary.BigEndian.Uint32(tmp[:]))
	return nil
}
func (i *Int3) Size(opt *struc.Options) int {
	return 3
}
func (i *Int3) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}

func BenchmarkMarshal__lunixbochs_struc(b *testing.B) {

	for it := 0; it < b.N; it++ {
		var buf bytes.Buffer
		a := reference
		if err := struc.Pack(&buf, a); err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkUnmarshal__lunixbochs_struc(b *testing.B) {
	out := &Example{}
	for it := 0; it < b.N; it++ {
		buf := bytes.NewReader(referenceBytes)
		if err := struc.Unpack(buf, out); err != nil {
			b.Fatalf("%v,%s", out, err)
		}
	}

}
