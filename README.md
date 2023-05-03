# struc-gen
Struc-gen is a code generator for Go that generates methods for binary Marshaling and Unmarshaling. The behaviour can be configured by adding struct tags. The configuration options are heavily inspired and substantially compatible with those used by [`lunixbochs/struc`](https://github.com/lunixbochs/struc)

This code is currently EXPERIMENTAL. The API may be changed at any time without notice. In favour of performance there is currently no error handling. The byte slice for marshaling must be pre-allocated to sufficient size. Also a sufficient binary slice with valid data must be provided for unmarshaling. If not, method will read out of bounds and panic(). Use with care.

## Useage
```go
// the go generate expression will run code generator for all structs in this file.
// put the struct in a seperate source code file to limit scope of code generation and avoid syntax errors while parsing file for code generation

//go:generate go run github.com/m29h/struc-gen/cmd/struc-gen
type Example struct {
	A int     `struc:"uint64,big,sizeof=B"` //encode in big endian, automatically set to length of slice B
	B []int64 `struc:"[]int16,little"`      //encode values in little endian
	C int     `struc:"int8,little,sizeof=D"`
	D string
}
```

```go
func main() {
	t := &Example{B: []int64{1, 2, 30000, 4, 5, 6}}
	buf := make([]byte, t.SizeOf())
	t.MarshalBinary(buf)
	o := &Example{B: make([]int64, t.A)}  
	o.UnmarshalBinary(buf)
	fmt.Printf("t=%v,o=%v\n", t, o)
	//t=&{6 [1 2 30000 4 5 6] 0 },o=&{6 [1 2 30000 4 5 6] 0 }

}
```

See `testcase/example.go` for a more extensive example including recursively serializing structs.

## Benchmark
Thanks to code generation struc-gen generated marshaling and unmarshaling methods do not require reflection and does not require allocations. This speeds up the methods by a factor of around 30x. In extremely simple scenarios where the struct size can be known at compile time the performance gain can be even higher.

```
go test --bench=. github.com/m29h/struc-gen/testcase/
goos: linux
goarch: amd64
pkg: github.com/m29h/struc-gen/testcase
cpu: AMD Ryzen 9 3900X 12-Core Processor            
BenchmarkMarshal__strucgen-24                    9780751               127.9 ns/op
BenchmarkUnmarshal__strucgen-24                  5429912               241.4 ns/op
BenchmarkMarshal__lunixbochs_struc-24             219531              4780 ns/op
BenchmarkUnmarshal__lunixbochs_struc-24           165632              6616 ns/op
PASS
ok      github.com/m29h/struc-gen/testcase      5.192s
```