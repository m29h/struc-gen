package structag

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/dave/jennifer/jen"
)

func (tag *StrucTag) Unpack() (*jen.Statement, error) {
	tag.method = UnmarshalBinary
	if tag.Skip {
		return jen.Null(), nil
	}

	j, err := tag.unroll(tag.field.Type(), RootStructName().Dot(tag.field.Name()))
	tag.loopvars = nil
	return jen.Comment(tag.field.Name()).Line().Add(j), err

}

func (tag *StrucTag) BinaryToType(goType string, n *jen.Statement) *jen.Statement {
	switch goType { //wrap incoming gotype or treat seperately if necessary
	case "bool":
		plus := jen.Null()
		v := tag.binaryToType("uint8", plus).Op(">").Lit(0)
		if tag.method == Size {
			return plus
		}
		return n.Clone().Op("=").Add(v).Line().Add(plus)
	case "string":
		return tag.binaryToType(goType, n)
	}

	switch tag.Type { //only allow  valid tag types
	case "bool":
		plus := jen.Null()
		v := tag.binaryToType("uint8", plus)
		if tag.method == Size {
			return plus
		}
		return n.Clone().Op("=").Id(goType).Call(v).Op("&").Lit(1).Line().Add(plus)
	case "byte", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "int", "uint", "float32", "float64":
		plus := jen.Null()
		v := tag.binaryToType(tag.Type, plus)
		if tag.method == Size {
			return plus
		}
		return n.Clone().Op("=").Id(goType).Call(v).Line().Add(plus)
	default:
		panic("Unable to handle unknown struc:" + tag.Type + " for field " + tag.field.Name())
	}
}

func (tag *StrucTag) binaryToType(goType string, n *jen.Statement) *jen.Statement {
	switch goType {
	case "byte", "uint8":
		n.Id("m").Op("+=").Lit(1)
		return jen.Id("b").Index(jen.Id("m"))
	case "string":
		_, exactlen := tag.GetNewLoopVar()
		j := n.Clone().Op("=").String().Call(jen.Id("b").Index(jen.Id("m").Op(":").Id("m").Op("+").Add(exactlen)))
		if tag.method == Size {
			j = jen.Null()
		}
		j.Line().Id("m").Op("+=").Add(exactlen)
		return j
		//}
	case "int8", "int16", "int32", "int64":
		unsignedType := "u" + goType
		return jen.Id(goType).Call(tag.binaryToType(unsignedType, n))
	case "int", "uint":
		unsignedType := "uint32"
		return jen.Id(goType + "32").Call(tag.binaryToType(unsignedType, n))
	case "uint16", "uint32", "uint64":
		len := 0
		fmt.Sscanf(goType, "uint%2d", &len)
		len /= 8
		st := jen.Null()
		for i := 0; i < len; i++ {
			if i > 0 {
				st.Op("|")
			}
			if tag.Order == binary.LittleEndian {
				st.Id(goType).Call(jen.Id("b").Index(jen.Id("m").Op("+").Lit(i))).Op("<<").Lit(8 * i)
			} else {
				st.Id(goType).Call(jen.Id("b").Index(jen.Id("m").Op("+").Lit(i))).Op("<<").Lit(8 * (len - 1 - i))
			}
		}
		n.Id("m").Op("+=").Lit(len)
		return st
	case "float32", "float64":
		return jen.Qual("math", "F"+goType[1:]+"frombits").Call(tag.binaryToType("uint"+goType[5:], n))

	default:
		log.Printf("type serialization unknown %s", goType)
		return jen.Qual("fmt", "Println").Call(jen.Params(jen.Id(goType)).Params(n.Clone()))

	}

}
