package structag

import (
	"encoding/binary"
	"log"
	"strings"

	"github.com/dave/jennifer/jen"
)

func (tag *StrucTag) Pack() (*jen.Statement, error) {
	tag.method = MarshalBinary
	if tag.Skip {
		return jen.Null(), nil
	}

	j, err := tag.unroll(tag.field.Type(), RootStructName().Dot(tag.field.Name()))
	if tag.sizeof != "" {
		//log.Printf("%s", tag.field.Type().String())
		targetType := strings.Replace(tag.field.Type().String(), "command-line-arguments.", "", 1)
		j = RootStructName().Dot(tag.field.Name()).Op("=").Id(targetType).Call(jen.Len(RootStructName().Dot(tag.sizeof))).Line().Add(j)
	}
	tag.loopvars = nil
	return tag.render(j), err

}

func (tag *StrucTag) TypeToBinary(goType string, n *jen.Statement) *jen.Statement {
	switch goType { //wrap incoming gotype or treat separately if necessary
	case "bool":
		//convert the bool to an uint8
		//uint8(*(*uint8)(unsafe.Pointer(&s.BoolT)) & 1 & 1)
		n = jen.Op("*").Parens(jen.Op("*").Uint8()).Parens(jen.Qual("unsafe", "Pointer").Call(jen.Op("&").Add(n))).Op("&").Lit(1)
		return tag.TypeToBinary("uint8", n) //Note: this nust call the exported "TypeToBinary" again for recursion
	case "string":
		return (tag.typeToBinary(goType, n))
	}

	switch tag.Type { //only allow  valid tag types
	case "bool":
		return (tag.typeToBinary("uint8", jen.Uint8().Call(n.Op("&").Lit(1))))
	case "byte", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "int", "uint", "float32", "float64":
		return (tag.typeToBinary(tag.Type, jen.Id(tag.Type).Call(n.Clone())))
	case "uint1", "uint2", "uint3", "uint4", "uint5", "uint6", "uint7":
		return tag.typeToBinary(tag.Type, jen.Uint16().Call(n.Clone().Op("&").Parens(jen.Lit((1<<tag.ElementBitLen())-1))))
	default:
		panic("Unable to handle unknown struc:" + tag.Type + " for field " + tag.field.Name())
	}
}

func (tag *StrucTag) typeToBinary(goType string, n *jen.Statement) *jen.Statement {
	switch goType {
	case "byte", "uint8":
		return jen.Id("b").Index(tag.bytePos(0)).Op("=").Add(n)
	case "string":
		pos := tag.bytePos(0)
		lvar, exactlen := tag.GetNewLoopVar()
		j := jen.For(
			lvar.Clone().Op(":=").Copy(jen.Id("b").Index(pos.Clone().Op(":").Add(pos.Clone()).Op("+").Add(exactlen)), jen.Index().Byte().Parens(n)),
			lvar.Clone().Op("<").Add(exactlen), lvar.Clone().Op("++"))
		j.Block(jen.Id("b").Index(tag.bytePos(0)).Op("=").Lit(0))
		return j

	case "int8", "int16", "int32", "int64":
		unsignedType := "u" + goType
		return tag.typeToBinary(unsignedType, jen.Parens(jen.Id(unsignedType)).Parens(n))
	case "int", "uint":
		unsignedType := "uint32"
		return tag.typeToBinary(unsignedType, jen.Parens(jen.Id(unsignedType)).Parens(n))
	case "uint16", "uint32", "uint64":
		len := tag.ElementBitLen() / 8
		st := jen.Id("_").Op("=").Id("b").Index(tag.bytePos(len - 1)).Comment("bounds check hint to compiler; see golang.org/issue/14808").Line()
		for i := 0; i < len; i++ {
			if tag.Order == binary.LittleEndian {
				st.Id("b").Index(tag.bytePos(i)).Op("=").Byte().Call(n.Clone().Op(">>").Lit(8 * i)).Line()
			} else {
				st.Id("b").Index(tag.bytePos(i)).Op("=").Byte().Call(n.Clone().Op(">>").Lit(8 * (len - 1 - i))).Line()
			}
		}
		return st
	case "uint1", "uint2", "uint3", "uint4", "uint5", "uint6", "uint7": // dense packed bitfield types
		len := tag.ElementBitLen()
		st := jen.Null()

		st.If(
			tag.bitOffset(len).Op("<=").Lit(8),
		).Block( // This is the case where the bitfield fits in a single byte
			jen.Id("b").Index(tag.bytePos(0)).Op("^=").Byte().Call(n.Clone()).Op("<<").Add(tag.bitOffset(0)),
		).Else().Block( // this is the case where the bitfield wraps across byte boundary
			jen.Id("tmp").Op(":=").Add(n.Clone()).Op("<<").Add(tag.bitOffset(0)),
			jen.Id("b").Index(tag.bytePos(0)).Op("^=").Byte().Call(jen.Id("tmp")),
			jen.Id("b").Index(tag.bytePos(1)).Op("^=").Byte().Call(jen.Id("tmp").Op(">>").Lit(8)),
		).Line()
		return st
	case "float32", "float64":
		return jen.Block(jen.Id("tmp").Op(":=").Qual("math", "F"+goType[1:]+"bits").Call(jen.Id(goType).Call(n)).Line().Add(tag.typeToBinary("uint"+goType[5:], jen.Id("tmp"))))

	default:
		log.Printf("type serialization unknown %s", goType)
		return jen.Qual("fmt", "Println").Call(jen.Params(jen.Id(goType)).Params(n.Clone()))

	}

}
