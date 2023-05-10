package test

//go:generate go run github.com/m29h/struc-gen/cmd/struc-gen

type BitfieldExample struct {
	A  int    `struc:"uint4,sizeof=C2"`
	B  int    `struc:"uint4"`
	C2 []byte `struc:"[]uint4"` //this type must get aligned and a padding bit is inserted
	D  int    `struc:"uint4"`
	E  int    `struc:"uint7,sizeof=F"` //this bitfield wraps across byte boundary
	F  []bool `struc:"[]uint1"`        //this bitfield array dynamically wraps across byte boundary
	G  int    `struc:"uint8"`          //byte aligned type after dynamic sized bitfield type
	H  bool   `struc:"uint1"`
}

type BitfieldExampleBench struct {
	V2 [256]uint  `struc:"[256]uint2"`
	V3 [8][8]uint `struc:"[8][8]uint3"`
}
type ByteExampleBench struct {
	V1 [256]uint  `struc:"[256]uint8"`
	V2 [8][8]uint `struc:"[][]uint8"`
}
