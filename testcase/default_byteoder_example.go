package test

// with the argument -little we are changing the default byteorder
//###go:generate go run github.com/m29h/struc-gen/cmd/struc-gen -little

type DefaultBO struct {
	L int `struc:"uint16"`
	B int `struc:"uint16,big"`
}
