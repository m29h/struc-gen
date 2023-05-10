[![Tests](https://github.com/m29h/struc-gen/actions/workflows/push.yml/badge.svg)](https://github.com/m29h/struc-gen/actions/workflows/push.yml)

# struc-gen
Struc-gen is a code generator for Go that generates methods for binary struct Marshaling and Unmarshaling. The behaviour can be configured by adding struct tags. The configuration options are heavily inspired and substantially compatible with those used by [`lunixbochs/struc`](https://github.com/lunixbochs/struc)

The code supports most basic Go types as well as bitfield types 1 to 7 bit; bitfield types get "tightly packed" on a bit level.

The API may be changed at any time without notice. In favour of performance there is currently limited error handling. The byte slice for marshaling must be pre-allocated to sufficient size or else marshaling will not do anything and return 0. Use the `SizeOf()` method to determine suitable byte slice size to allocate.

Unmarshaling is terminating gracefully if the end of byte slice is reached prematurely. Unmarshaling an invalid/incomplete byte stream may leave the struct partially uninizialized but it will not read out of bounds and most importantly can not crash your application.

Very basic validation can be achieved with e.g.:
```go
o := &Example{}  
if actual := o.UnmarshalBinary(buf); actual != o.SizeOf(){
	return nil, errors.New("Reached EOF while Unmarshaling")
}
return o, nil
```
Slices and pointer receivers are automatically allocated if they are nil in UnmarshalBinary. Slices are resized when necessary in UnmarshalBinary
## Supported Tag Binary Type Specifiers

The following types are supported for binary marshaling
 - unsigned integer types `uint8` `uint16` `uint32` `uint64` 
 - floating Point types `float32` `float64` 
 - Go Strings `string` requires a linked `sizeof=...` field to store the string length
 - unsigned bitfield types `uint1` `uint2` `uint3` `uint4` `uint5` `uint6` `uint7`. These get tightly packed after each other. Padding bits are automatically introduced after bitfields to make the next non-bitfield types byte-aligned again.
 - Dummy type to introduce byte padding `pad`, always marshaled as `0x00`
 - Array types of any of the above by prepending `[len]` syntax, (including arrays of bitfield types get tightly packed! For example an `[4]uint6` packs 4 numbers into three bytes)
 - Slice types of any of the above `[]`, requires another linked field to be tagged with `sizeof=...` to store the slice size

for compatibility [`lunixbochs/struc`](https://github.com/lunixbochs/struc) also the types `bool` and `byte` are supported that both effectively map to a binary `uint8`

Map types are unsupported currently. You can however make it a custom named type and manually implement its `MarshalBinary([]byte) int`,`UnarshalBinary([]byte) int` and `SizeOf() int` methods

## Useage
- Install the struc-gen code generator `go install github.com/m29h/struc-gen/cmd/struc-gen@latest`
- Annotate your go files for which you want to have the Marshaling methods generated with `//go:generate struc-gen`
- Run `go generate ./...` in your go module directory

```go
// the go generate expression will run code generator for all structs in this file.
// put the struct in a seperate source code file to limit scope of code generation and avoid syntax errors while parsing file for code generation

//go:generate struc-gen
type Example struct {
	//unexported fields such as a + c are no problem and treated just like exported fields
	a int     `struc:"uint64,big,sizeof=B"` //encode in big endian, automatically set to length of slice B
	b int     `struc:"skip"`                // If you want an unexported field to not be marshaled just tag it with "skip"
	B []int64 `struc:"[]int16,little"`      //encode values in little endian
	C int     `struc:"uint4"`
	c int     `struc:"uint7,sizeof=D"` //this bitfield is packed without gap after C and wraps across byte boundary
	D string  // don't worry, types larger than 7 bits will always be written byte-aligned
}
```

```go
func main() {
	t := &Example{B: []int64{1, 2, 30000, 4, 5, 6},
		b: 1337, //field marked with `struc:"skip"`, will not be Marshaled
		C: 30,
		D: "Hello World",
	}
	buf := make([]byte, t.SizeOf())
	t.MarshalBinary(buf)
	o := &Example{}
	o.UnmarshalBinary(buf)
	fmt.Printf("t=%v,o=%v\n", t, o)
	//t=&{6 1337 [1 2 30000 4 5 6] 30 11 Hello World},o=&{6 0 [1 2 30000 4 5 6] 14 11 Hello World}
}
```

See `testcase/example.go` for a more extensive example including recursively serializing structs.

## Benchmark
Thanks to code generation struc-gen generated marshaling and unmarshaling methods do not require reflection and does not require allocations. This speeds up the methods by a factor of around 30x. In extremely simple scenarios where the struct size can be known at compile time the performance gain can be even higher.

```
go test --bench=.
goos: linux
goarch: amd64
pkg: github.com/m29h/struc-gen/testcase
cpu: AMD Ryzen 9 3900X 12-Core Processor            
BenchmarkMarshal__BitfieldArray-24               2092998               576.8 ns/op
BenchmarkUnmarshal__BitfieldArray-24             2074000               558.1 ns/op
BenchmarkMarshal__ByteArray-24                   8339070               148.9 ns/op
BenchmarkUnmarshal__ByteArray-24                 6079618               175.1 ns/op
BenchmarkMarshal__strucgen-24                   10114875               114.4 ns/op
BenchmarkUnmarshal__strucgen-24                  3848166               302.3 ns/op
BenchmarkMarshal__lunixbochs_struc-24             225576              5307 ns/op
BenchmarkUnmarshal__lunixbochs_struc-24           252816              4725 ns/op
PASS
ok      github.com/m29h/struc-gen/testcase      13.429s
```