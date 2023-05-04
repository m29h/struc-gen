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

// Helper method generating a bounds check 'b[m]' for the wish of deserializing exactlen bytes
// In the bounds check failure case a "return m" block is inserted

func checkBound(exactlen *jen.Statement) *jen.Statement {
	return jen.If(jen.Id("m").Op("+").Add(exactlen).Op(">").Len(jen.Id("b"))).Block(jen.Return(jen.Id("m"))).Line()
}

// Helper method generating the deserialization statement consisting of:
// 1. source byte slice 'b[m]' bounds check,
// 2. convert deserialized value to goType and append convertOp
// 3. set deserialized value to valueReceiver
// 4. add deserialized bytes size to slice position index variable 'm'
func (tag *StrucTag) genBinToTypeStatement(goType, binType string, convertOp *jen.Statement, valueReceiver *jen.Statement) *jen.Statement {
	plus := jen.Null()
	v := jen.Id(goType).Call(tag.binaryToType(binType, plus)).Add(convertOp)
	if tag.method == Size {
		return jen.Id("m").Op("+=").Add(plus)
	}
	return checkBound(plus).Add(valueReceiver.Clone()).Op("=").Add(v).Line().Id("m").Op("+=").Add(plus)
}
func (tag *StrucTag) BinaryToType(goType string, n *jen.Statement) *jen.Statement {
	switch goType { //wrap incoming gotype or treat seperately if necessary
	case "bool":
		return tag.genBinToTypeStatement("uint8", "uint8", jen.Op(">").Lit(0), n)
	case "string":
		return tag.binaryToType(goType, n)
	}

	switch tag.Type { //only allow  valid tag types
	case "bool":
		return tag.genBinToTypeStatement(goType, "uint8", jen.Op("&").Lit(1), n)
	case "byte", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "int", "uint", "float32", "float64":
		return tag.genBinToTypeStatement(goType, tag.Type, jen.Null(), n)
	default:
		panic("Unable to handle unknown struc:" + tag.Type + " for field " + tag.field.Name())
	}
}

func (tag *StrucTag) binaryToType(goType string, n *jen.Statement) *jen.Statement {
	switch goType {
	case "byte", "uint8":
		n.Lit(1)
		return jen.Id("b").Index(jen.Id("m"))
	case "string":
		_, exactlen := tag.GetNewLoopVar()
		j := checkBound(exactlen).Add(n.Clone()).Op("=").String().Call(jen.Id("b").Index(jen.Id("m").Op(":").Id("m").Op("+").Add(exactlen)))
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
		n.Lit(len)
		return st
	case "float32", "float64":
		return jen.Qual("math", "F"+goType[1:]+"frombits").Call(tag.binaryToType("uint"+goType[5:], n))

	default:
		log.Printf("type serialization unknown %s", goType)
		return jen.Qual("fmt", "Println").Call(jen.Params(jen.Id(goType)).Params(n.Clone()))

	}

}
