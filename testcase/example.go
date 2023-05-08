package test

import (
	"encoding/binary"
)

//go:generate struc-gen

// Custom Type
type Int3 uint32

func (i *Int3) MarshalBinary(p []byte) int {
	var tmp [4]byte
	binary.BigEndian.PutUint32(tmp[:], uint32(*i))
	copy(p, tmp[1:])
	return 3
}
func (i *Int3) UnmarshalBinary(p []byte) int {
	var tmp [4]byte
	if len(p) < 3 {
		return 0
	}
	tmp[1] = p[0]
	tmp[2] = p[1]
	tmp[3] = p[2]
	*i = Int3(binary.BigEndian.Uint32(tmp[:]))
	return 3
}
func (i *Int3) SizeOf() int {
	return 3
}

type Nested struct {
	Test2 int `struc:"int8"`
}

// test case "struct Example" taken directly from github.com/lunixbochs/struc to verify compatibility
type Example struct {
	Pad    []byte `struc:"[5]pad"`        // 00 00 00 00 00
	I8f    int    `struc:"int8"`          // 01
	I16f   int    `struc:"int16"`         // 00 02
	I32f   int    `struc:"int32"`         // 00 00 00 03
	I64f   int    `struc:"int64"`         // 00 00 00 00 00 00 00 04
	U8f    int    `struc:"uint8,little"`  // 05
	U16f   int    `struc:"uint16,little"` // 06 00
	U32f   int    `struc:"uint32,little"` // 07 00 00 00
	U64f   int    `struc:"uint64,little"` // 08 00 00 00 00 00 00 00
	Boolf  int    `struc:"bool"`          // 01
	Byte4f []byte `struc:"[4]byte"`       // "abcd"

	I8     int8    // 09
	I16    int16   // 00 0a
	I32    int32   // 00 00 00 0b
	I64    int64   // 00 00 00 00 00 00 00 0c
	U8     uint8   `struc:"little"` // 0d
	U16    uint16  `struc:"little"` // 0e 00
	U32    uint32  `struc:"little"` // 0f 00 00 00
	U64    uint64  `struc:"little"` // 10 00 00 00 00 00 00 00
	BoolT  bool    // 01
	BoolF  bool    // 00
	Byte4  [4]byte // "efgh"
	Float1 float32 // 41 a0 00 00
	Float2 float64 // 41 35 00 00 00 00 00 00

	I32f2 int64 `struc:"int32"`  // ff ff ff ff
	U32f2 int64 `struc:"uint32"` // ff ff ff ff

	I32f3 int32 `struc:"int64"` // ff ff ff ff ff ff ff ff

	Size1 int    `struc:"sizeof=Str,little"` // 0a 00 00 00
	Str   string `struc:"[]byte"`            // "ijklmnopqr"
	Strb  string `struc:"[4]byte"`           // "stuv"

	Size2 int    `struc:"uint8,sizeof=Str2"` // 04
	Str2  string // "1234"

	Size3 int    `struc:"uint8,sizeof=Bstr"` // 04
	Bstr  []byte // "5678"

	Size4 int    `struc:"little"`                // 07 00 00 00
	Str4a string `struc:"[]byte,sizefrom=Size4"` // "ijklmno"
	Str4b string `struc:"[]byte,sizefrom=Size4"` // "pqrstuv"

	Size5 int    `struc:"uint8"`          // 04
	Bstr2 []byte `struc:"sizefrom=Size5"` // "5678"

	Nested  Nested  // 00 00 00 01
	NestedP *Nested // 00 00 00 02
	TestP64 *int    `struc:"int64"` // 00 00 00 05

	NestedSize int      `struc:"sizeof=NestedA"` // 00 00 00 02
	NestedA    []Nested // [00 00 00 03, 00 00 00 04]

	Skip int `struc:"skip"`

	CustomTypeSize    Int3   `struc:"sizeof=CustomTypeSizeArr"` // 00 00 00 04
	CustomTypeSizeArr []byte // "ABCD"

}
