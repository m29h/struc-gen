package structag

import (
	"encoding/binary"
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
	return tag.render(j), err

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
		return jen.Null()
	}
	return valueReceiver.Clone().Op("=").Add(v)
}
func (tag *StrucTag) BinaryToType(goType string, n *jen.Statement) *jen.Statement {
	if tag.method == UnmarshalBinary {
		*(tag.context.checkBound) = *(jen.If(tag.bitPosStart().Op("+").Add(tag.fieldLenStatement(false)).Op(">").Len(jen.Id("b")).Op("*").Lit(8)).Block(jen.Return(jen.Lit(0))).Line())
	}

	switch goType { //wrap incoming gotype or treat separately if necessary
	case "bool":
		return tag.genBinToTypeStatement("uint8", tag.Type, jen.Op(">").Lit(0), n)
	case "string":
		return tag.binaryToType(goType, n)
	}

	switch tag.Type { //only allow  valid tag types
	case "bool":
		return tag.genBinToTypeStatement(goType, "uint8", jen.Op("&").Lit(1), n)
	case "byte", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "int", "uint", "float32", "float64", "uint1", "uint2", "uint3", "uint4", "uint5", "uint6", "uint7":
		return tag.genBinToTypeStatement(goType, tag.Type, jen.Null(), n)

	default:
		panic("Unable to handle unknown struc:" + tag.Type + " for field " + tag.field.Name())
	}
}

func (tag *StrucTag) binaryToType(goType string, n *jen.Statement) *jen.Statement {
	switch goType {
	case "byte", "uint8", "bool":
		n.Lit(tag.ElementBitLen())
		return jen.Id("b").Index(tag.bytePos(0))
	case "string":
		pos := tag.bytePos(0)
		_, exactlen := tag.GetNewLoopVar()
		j := jen.Add(n.Clone()).Op("=").String().Call(jen.Id("b").Index(pos.Clone().Op(":").Add(pos.Clone()).Op("+").Add(exactlen)))
		if tag.method == Size {
			j = jen.Null()
		}

		return j
		//}
	case "int8", "int16", "int32", "int64":
		unsignedType := "u" + goType
		return jen.Id(goType).Call(tag.binaryToType(unsignedType, n))
	case "int", "uint":
		unsignedType := "uint32"
		return jen.Id(goType + "32").Call(tag.binaryToType(unsignedType, n))
	case "uint16", "uint32", "uint64":
		len := tag.ElementBitLen() / 8
		st := jen.Null()
		for i := 0; i < len; i++ {
			if i > 0 {
				st.Op("|")
			}
			if tag.Order == binary.LittleEndian {
				st.Id(goType).Call(jen.Id("b").Index(tag.bytePos(i))).Op("<<").Lit(8 * i)
			} else {
				st.Id(goType).Call(jen.Id("b").Index(tag.bytePos(i))).Op("<<").Lit(8 * (len - 1 - i))
			}
		}
		n.Lit(len * 8)
		return st
	case "uint1", "uint2", "uint3", "uint4", "uint5", "uint6", "uint7": // dense packed bitfield types
		len := tag.ElementBitLen()
		st := jen.Null()
		st.Id("tmp").Op(":=").Uint16().Call(jen.Id("b").Index(tag.bytePos(0))).Line()
		st.If(tag.bitOffset(len).Op(">").Lit(8)).Block( // this is the case where the bitfield wraps across byte boundary
			jen.Id("tmp").Op("|=").Uint16().Call(jen.Id("b").Index(tag.bytePos(1))).Op("<<").Lit(8)).Line()
		st.Return(jen.Parens(jen.Id("tmp").Op(">>").Add(tag.bitOffset(0))).Op("&").Parens(jen.Lit((1 << tag.ElementBitLen()) - 1)))
		n.Lit(len)
		// wrap result into a anonymous function
		return jen.Func().Params().Uint16().Block(st).Call()

	case "float32", "float64":
		return jen.Qual("math", "F"+goType[1:]+"frombits").Call(tag.binaryToType("uint"+goType[5:], n))

	default:
		log.Printf("type serialization unknown %s", goType)
		return jen.Qual("fmt", "Println").Call(jen.Params(jen.Id(goType)).Params(n.Clone()))

	}

}
