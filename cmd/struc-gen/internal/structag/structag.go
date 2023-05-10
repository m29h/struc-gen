package structag

import (
	"encoding/binary"
	"errors"
	"fmt"
	"go/types"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
)

type directionT uint32

const (
	MarshalBinary directionT = iota
	UnmarshalBinary
	Size
)

type StrucTag struct {
	field           *types.Var
	Type            string
	Order           binary.ByteOrder
	sizeof          string
	Skip            bool
	sizefrom        string
	loopvars        []string
	arraySize       []string
	method          directionT
	beforeStatement *jen.Statement
	context         *Context
}

// context contains reference to information where the field is in relation to the entire struct
func (t *StrucTag) setContext(context *Context) {
	t.context = context

	if t.context.checkBound == nil || t.FieldType() != Const {
		t.beforeStatement.Add(t.context.insertNewCheckBound())
	}

	if t.FieldType() == SelfContained || t.ElementBitLen() >= 8 {
		//align staticBitPos when needed (transition from bitfield to larger types)
		if t.context.dynamicBitPos == nil {
			// byte-align at compile-time for fully constant sized types
			if t.context.staticBitPos%8 != 0 {
				t.context.staticBitPos += 8 - t.context.staticBitPos%8
			}
		} else {
			//byte align by flushing the context
			t.beforeStatement.Add(t.context.Flush())
		}
	}
}
func (t StrucTag) bytePos(offset int) *jen.Statement {
	p := jen.Id("m")
	if t.loopvars == nil && t.context.dynamicBitPos == nil {
		//fast tracked case as all is static
		return jen.Parens(p.Op("+").Lit(t.context.staticBitPos/8 + offset))
	}
	if t.loopvars == nil && t.context.dynamicBitPos != nil {
		//fast tracked case with bitPosDynamicStatement
		return p.Op("+").Parens(jen.Lit(t.context.staticBitPos + offset*8).Op("+").Add(t.context.dynamicBitPos.Clone())).Op("/").Lit(8)
	}

	byteAlignedType := (t.ElementBitLen()%8 == 0) && (t.context.dynamicBitPos == nil)
	if t.loopvars != nil {
		if byteAlignedType {
			p.Op("+").Lit(t.context.staticBitPos/8 + offset) //without final division
		} else {
			p = jen.Lit(t.context.staticBitPos + offset*8) //for final division
		}
		if t.context.dynamicBitPos != nil {
			p.Op("+").Add(t.context.dynamicBitPos.Clone())
		}
		for i, v := range t.loopvars { // for each loop var level
			p.Op("+").Id(v) //add loop position index

			bitlen := t.ElementBitLen()
			if byteAlignedType { //directly divide by 8 for byte aligned types (saves a division later)
				bitlen /= 8
			}
			for _, w := range t.arraySize[i+1:] {
				if Size, err := strconv.Atoi(w); err != nil {
					p.Op("*").Int().Call(RootStructName().Dot(w)) //handle at runtime
				} else {
					bitlen *= Size //handle at compile time
				}
			}
			if bitlen != 1 { // no need to multiply with one
				p.Op("*").Lit(bitlen)
			}
		}
	}
	if !byteAlignedType {
		//m is byte offset so add it outside of the Parentheses so it'not included in the division
		return jen.Id("m").Op("+").Parens(p).Op("/").Lit(8) //finally divide whole statement unless byteAlignedType
	}
	return jen.Parens(p)
}

// gives the bit position of the first bit of the current strucTag
func (t StrucTag) bitPosStart() *jen.Statement {
	p := jen.Lit(t.context.staticBitPos)
	if t.context.dynamicBitPos != nil {
		p.Op("+").Add(t.context.dynamicBitPos.Clone())
	}
	p.Op("+").Id("m").Op("*").Lit(8)
	return p
}

// gives the bit offset of the current struc tag.
// Statement will loopvariables so it is also valid inside loop blocks
func (t StrucTag) bitOffset(offset int) *jen.Statement {
	if t.loopvars == nil && t.context.dynamicBitPos == nil {
		//fast tracked case as all is static
		return jen.Lit(t.context.staticBitPos%8 + offset)
	}
	p := jen.Lit(t.context.staticBitPos)
	if t.context.dynamicBitPos != nil {
		p.Op("+").Add(t.context.dynamicBitPos.Clone())
	}
	if t.loopvars != nil {
		for i, v := range t.loopvars { // for each loop var level
			lOffs := jen.Id(v) //add loop position index
			bitlen := t.ElementBitLen()
			for _, w := range t.arraySize[i+1:] {
				if Size, err := strconv.Atoi(w); err != nil {
					lOffs.Op("*").Int().Call(RootStructName().Dot(w)) //handle at runtime
				} else {
					bitlen *= Size //handle at compile time
				}
			}
			if bitlen%8 != 0 { // only include loop variables that introduce a mod(8) !=0 cycle
				p.Op("+").Add(lOffs).Op("*").Lit(bitlen % 8)
			}
		}
	}
	return jen.Parens(jen.Parens(p).Op("%").Lit(8).Op("+").Lit(offset))
}

func (t *StrucTag) finalize() *jen.Statement {
	// advance staticBitPos where applicable (compile time sized types)
	t.context.staticBitPos += t.FieldBitLen()
	if t.FieldType() == RuntimeSize && t.ElementBitLen()%8 == 0 {
		//advance m on byte aligned types
		return jen.Line().Id("m").Op("+=").Add(t.fieldLenStatement(true)).Add(t.context.insertNewCheckBound())
	}
	if t.FieldType() == RuntimeSize && t.ElementBitLen()%8 != 0 {
		//advance bitPosDynamicStatement on non byte aligned types with runtime size
		if t.context.dynamicBitPos == nil {
			t.context.dynamicBitPos = jen.Parens(t.fieldLenStatement(false))
		} else {
			t.context.dynamicBitPos = jen.Parens(t.context.dynamicBitPos.Clone().Op("+").Add(t.fieldLenStatement(false)))
		}
		return jen.Null()
	}
	return jen.Null()
}

func (t StrucTag) String() string {
	bo := "big"
	if t.Order == binary.LittleEndian {
		bo = "little"
	}
	return fmt.Sprintf("%s %s Sizeof=%s Skip=%t SizeFrom=%s", t.Type, bo, t.sizeof, t.Skip, t.sizefrom)
}

// Gets a new unique named loop variable name for use in the context of serializing this struct
func (t *StrucTag) GetNewLoopVar() (*jen.Statement, *jen.Statement) {
	if t.loopvars == nil {
		t.loopvars = []string{"i"}
	} else {
		t.loopvars = append(t.loopvars, fmt.Sprintf("i%d", len(t.loopvars)))
	}
	if len(t.arraySize) < len(t.loopvars) {
		panic("Syntax error for " + t.field.Name() + ": Insufficient definition for slice/array sizeof")
	}
	var ArraySize *jen.Statement
	size := t.arraySize[len(t.loopvars)-1]
	if val, err := strconv.Atoi(size); err == nil {
		ArraySize = jen.Lit(val)
	} else {
		ArraySize = RootStructName().Dot(size)
	}
	return jen.Id(t.loopvars[len(t.loopvars)-1]), jen.Int().Call(ArraySize)
}

type FieldType int

const (
	Const FieldType = iota
	RuntimeSize
	SelfContained
)

func (t *StrucTag) FieldType() FieldType {
	if t.ElementBitLen() == 0 {
		return SelfContained
	}
	for _, v := range t.arraySize {
		if _, err := strconv.Atoi(v); err != nil {
			return RuntimeSize
		}
	}
	return Const
}

// Returns jen.Statement representing size in Bit that the Field occupies in binary representation
// If give_bytes is true then a statement in Bytes is produced instead (only works for byte aligned types)

func (t *StrucTag) fieldLenStatement(give_bytes bool) *jen.Statement {
	if t.FieldType() == SelfContained {
		return jen.Nil()
	}
	staticSize := t.ElementBitLen()
	//result := new()
	var result *jen.Statement
	for _, v := range t.arraySize {
		if val, err := strconv.Atoi(v); err == nil {
			staticSize *= val
		} else {
			a := jen.Int().Call(RootStructName().Dot(v))
			if result == nil {
				result = a
			} else {
				result.Op("*").Add(a)
			}
		}
	}
	if give_bytes {
		staticSize /= 8
	}
	s := jen.Lit(staticSize)
	if result == nil {
		result = s
	} else {
		result.Op("*").Add(s)
	}
	return result
}

// Returns size in Bit for the Field occupies in binary representation
func (t *StrucTag) FieldBitLen() int {
	if t.FieldType() != Const {
		return 0
	}
	staticSize := t.ElementBitLen()
	for _, v := range t.arraySize {
		if val, err := strconv.Atoi(v); err == nil {
			staticSize *= val
		} else {
			return 0
		}
	}
	return staticSize
}

func (t *StrucTag) ElementBitLen() int {
	if strings.HasPrefix(t.Type, "uint") ||
		strings.HasPrefix(t.Type, "int") ||
		strings.HasPrefix(t.Type, "float") {

		re := regexp.MustCompile("[0-9]+")
		res := re.FindAllString(t.Type, -1)
		if len(res) == 0 {
			return 32 // default size of platform dependent types
		}
		len, err := strconv.Atoi(res[0])
		if err != nil {
			panic(err)
		}
		return len
	}
	switch t.Type {
	case "string", "byte", "pad", "bool":
		return 8
	default:
		return 0
	}
}

var (
	sizeofMap map[string]string
)

func init() {
	sizeofMap = make(map[string]string)
}

func RootStructName() *jen.Statement {
	return jen.Id("s")
}
func NewStrucTag(ctx *Context, fieldVal *types.Var, tagVal string) (*StrucTag, error) {
	t := &StrucTag{
		field: fieldVal,
		Order: binary.BigEndian,
	}
	tag := reflect.StructTag(tagVal)
	tagStr := tag.Get("struc")
	if tagStr == "" {
		// someone's going to typo this (I already did once)
		// sorry if you made a module actually using this tag
		// and you're mad at me now
		tagStr = tag.Get("struct")
	}
	for _, s := range strings.Split(tagStr, ",") {
		if strings.HasPrefix(s, "sizeof=") {
			tmp := strings.SplitN(s, "=", 2)
			t.sizeof = tmp[1]
			sizeofMap[t.sizeof] = t.field.Name()
		} else if strings.HasPrefix(s, "sizefrom=") {
			tmp := strings.SplitN(s, "=", 2)
			t.sizefrom = tmp[1]
			sizeofMap[t.field.Name()] = t.sizefrom
		} else if s == "big" {
			t.Order = binary.BigEndian
		} else if s == "little" {
			t.Order = binary.LittleEndian
		} else if s == "skip" {
			t.Skip = true
		} else {
			t.Type = s
		}
	}
	if t.Type == "" {
		t.Type = t.field.Type().String()
	}
	slStruc := strings.Split(t.Type, "]")
	//extract the basic type
	t.Type = slStruc[len(slStruc)-1]
	if t.Type == "string" {
		//insert another slice like bracket to make string behave more like []byte for consistency with arrays
		slStruc = append(slStruc[:len(slStruc)-1], "[", "string")
	}
	slGoType := strings.Split(t.field.Type().String(), "]")
	if slGoType[len(slGoType)-1] == "string" {
		//insert another slice like bracket to make string behave more like []byte for consistency with arrays
		slGoType = append(slGoType[:len(slGoType)-1], "[", "string")
	}

	if len(slStruc) != len(slGoType) {
		return nil, errors.New("go type must match struc type dimensions in " + t.field.Name())
	}
	t.arraySize = make([]string, 0)
	for i, indexStruc := range slStruc[:len(slStruc)-1] {
		indexGo := slGoType[i][1:]
		if val, err := strconv.Atoi(indexStruc[1:]); err == nil {
			if valG, err := strconv.Atoi(indexGo); err == nil && valG != val {
				return nil, errors.New("fixed size go array type must match struct type array size")
			}
			// Struc array with Go array or Struc array with Go slice
			t.arraySize = append(t.arraySize, fmt.Sprintf("%d", val))
			continue
		}
		if indexStruc == "[" {
			if valG, err := strconv.Atoi(indexGo); err == nil { // case Go array with struc slice
				t.arraySize = append(t.arraySize, fmt.Sprintf("%d", valG))
				continue
			}
			// case struc slice with Go slice
			t.arraySize = append(t.arraySize, sizeofMap[t.field.Name()])
			delete(sizeofMap, t.field.Name())
			continue
		}
		if len(indexStruc) > 1 {
			if indexGo != "[" {
				return nil, errors.New("struc [sizeof] notation is only allowed in conjunction with Go slice type")
			}
			//allow [sizeOfField]Type syntax to set the sizeof field in square brackets
			t.arraySize = append(t.arraySize, indexStruc[1:])
			continue
		}
		return nil, errors.New("syntax error in struct tag type statement" + t.field.Name() + t.Type)

	}
	t.beforeStatement = jen.Comment(t.field.Name() + " " + t.Type).Line()
	t.setContext(ctx)

	return t, nil
}
