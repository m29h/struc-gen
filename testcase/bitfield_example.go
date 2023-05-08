package test

//go:generate struc-gen

type BitfieldExample struct {
	A  int     `struc:"uint4"`
	B  int     `struc:"uint3"`    // even if a 0xf is in B onld a 0x7 may be written to the uint3
	C2 [2]byte `struc:"[2]uint8"` //this type must get aligned and a padding bit is inserted
	D  int     `struc:"uint4"`
	E  int     `struc:"uint7"`    //this bitfield wraps across byte boundary
	F  [2]bool `struc:"[2]uint1"` //this bitfield wraps across byte boundary
}

type BitfieldExampleBench struct {
	V1  byte `struc:"uint2"`
	V2  byte `struc:"uint2"`
	V3  byte `struc:"uint2"`
	V4  byte `struc:"uint2"`
	V5  byte `struc:"uint2"`
	V6  byte `struc:"uint2"`
	V7  byte `struc:"uint2"`
	V8  byte `struc:"uint2"`
	AV1 byte `struc:"uint2"`
	AV2 byte `struc:"uint2"`
	AV3 byte `struc:"uint2"`
	AV4 byte `struc:"uint2"`
	AV5 byte `struc:"uint2"`
	AV6 byte `struc:"uint2"`
	AV7 byte `struc:"uint2"`
	AV8 byte `struc:"uint2"`
	BV1 byte `struc:"uint2"`
	BV2 byte `struc:"uint2"`
	BV3 byte `struc:"uint2"`
	BV4 byte `struc:"uint2"`
	BV5 byte `struc:"uint2"`
	BV6 byte `struc:"uint2"`
	BV7 byte `struc:"uint2"`
	BV8 byte `struc:"uint2"`
}
type ByteExampleBench struct {
	V1  byte `struc:"uint8"`
	V2  byte `struc:"uint8"`
	V3  byte `struc:"uint8"`
	V4  byte `struc:"uint8"`
	V5  byte `struc:"uint8"`
	V6  byte `struc:"uint8"`
	V7  byte `struc:"uint8"`
	V8  byte `struc:"uint8"`
	AV1 byte `struc:"uint8"`
	AV2 byte `struc:"uint8"`
	AV3 byte `struc:"uint8"`
	AV4 byte `struc:"uint8"`
	AV5 byte `struc:"uint8"`
	AV6 byte `struc:"uint8"`
	AV7 byte `struc:"uint8"`
	AV8 byte `struc:"uint8"`
	BV1 byte `struc:"uint8"`
	BV2 byte `struc:"uint8"`
	BV3 byte `struc:"uint8"`
	BV4 byte `struc:"uint8"`
	BV5 byte `struc:"uint8"`
	BV6 byte `struc:"uint8"`
	BV7 byte `struc:"uint8"`
	BV8 byte `struc:"uint8"`
}
