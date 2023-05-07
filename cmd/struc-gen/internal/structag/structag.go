package structag

import (
	"encoding/binary"
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
	field     *types.Var
	Type      string
	Order     binary.ByteOrder
	sizeof    string
	Skip      bool
	sizefrom  string
	loopvars  []string
	forceSize []*jen.Statement
	method    directionT
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
	if len(t.forceSize) < len(t.loopvars) {
		panic("Syntax error for " + t.field.Name() + ": Insufficient definition for slice/array sizeof")
	}
	return jen.Id(t.loopvars[len(t.loopvars)-1]), jen.Int().Call(t.forceSize[len(t.loopvars)-1])
}

func (t *StrucTag) GetBitLen() int {
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
		return 8
	}
}

var (
	sizeofMap map[string]string
	//sizefromMap map[string]string
)

func init() {
	//sizefromMap = make(map[string]string)
	sizeofMap = make(map[string]string)
}

func RootStructName() *jen.Statement {
	return jen.Id("s")
}
func NewStrucTag(fieldVal *types.Var, tagVal string) *StrucTag {
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
			//sizefromMap[t.sizeof] = t.field.Name()
		} else if strings.HasPrefix(s, "sizefrom=") {
			tmp := strings.SplitN(s, "=", 2)
			t.sizefrom = tmp[1]
			//sizefromMap[t.field.Name()] = t.sizefrom
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
	tmp := strings.Split(t.Type, "]")
	if len(tmp) < len(strings.Split(t.field.Type().String(), "]")) {
		panic("Go type must match struc type dimensions in " + t.field.Name())
	}
	t.forceSize = make([]*jen.Statement, 0)
	for _, v := range tmp[:len(tmp)-1] {
		if val, err := strconv.Atoi(v[1:]); err == nil {
			t.forceSize = append(t.forceSize, jen.Lit(val))
		} else if v == "[" {
			t.forceSize = append(t.forceSize, RootStructName().Dot(sizeofMap[t.field.Name()]))
		} else if len(v) > 1 {
			sizeofMap[t.field.Name()] = v[1:] //allow [sizeOfField]Type syntax to set the sizeof field in square brackets
			t.forceSize = append(t.forceSize, RootStructName().Dot(sizeofMap[t.field.Name()]))
		} else {
			panic("syntax error in struct tag type statement" + t.field.Name() + t.Type)
		}
	}
	t.Type = tmp[len(tmp)-1]
	if t.Type == "string" && len(t.forceSize) == 0 {
		t.forceSize = append(t.forceSize, RootStructName().Dot(sizeofMap[t.field.Name()]))
	}
	return t
}
