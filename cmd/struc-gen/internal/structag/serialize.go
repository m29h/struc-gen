package structag

import (
	"encoding/binary"
	"fmt"
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
	return jen.Comment(tag.field.Name()).Line().Add(j), err

}

func (tag *StrucTag) TypeToBinary(goType string, n *jen.Statement) *jen.Statement {
	switch goType { //wrap incoming gotype or treat seperately if necessary
	case "bool":
		//convert the bool to an uint8
		//uint8(*(*uint8)(unsafe.Pointer(&s.BoolT)) & 1 & 1)
		n = jen.Op("*").Parens(jen.Op("*").Uint8()).Parens(jen.Qual("unsafe", "Pointer").Call(jen.Op("&").Add(n))).Op("&").Lit(1)
		return tag.TypeToBinary("uint8", n) //Note: this nust call the exported "TypeToBinary" again for recursion
	case "string":
		return tag.typeToBinary(goType, n)
	}

	switch tag.Type { //only allow  valid tag types
	case "bool":
		return tag.typeToBinary("uint8", jen.Uint8().Call(n.Op("&").Lit(1)))
	case "byte", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "int", "uint", "float32", "float64":
		return tag.typeToBinary(tag.Type, jen.Id(tag.Type).Call(n.Clone()))
	default:
		panic("Unable to handle unknown struc:" + tag.Type + " for field " + tag.field.Name())
	}
}

func (tag *StrucTag) typeToBinary(goType string, n *jen.Statement) *jen.Statement {
	switch goType {
	case "byte", "uint8":
		return jen.Id("b").Index(jen.Id("m")).Op("=").Add(n).Line().Id("m").Op("++")
	case "string":
		lvar, exactlen := tag.GetNewLoopVar()
		j := jen.For(
			lvar.Clone().Op(":=").Copy(jen.Id("b").Index(jen.Id("m").Op(":").Id("m").Op("+").Add(exactlen)), jen.Index().Byte().Parens(n)),
			lvar.Clone().Op("<").Add(exactlen), lvar.Clone().Op("++"))
		j.Block(jen.Id("b").Index(jen.Id("m").Op("+").Add(lvar)).Op("=").Lit(0))
		j.Line().Id("m").Op("+=").Add(exactlen)
		return j
		//}
	case "int8", "int16", "int32", "int64":
		unsignedType := "u" + goType
		return tag.typeToBinary(unsignedType, jen.Parens(jen.Id(unsignedType)).Parens(n))
	case "int", "uint":
		unsignedType := "uint32"
		return tag.typeToBinary(unsignedType, jen.Parens(jen.Id(unsignedType)).Parens(n))
	case "uint16", "uint32", "uint64":
		len := 0
		fmt.Sscanf(goType, "uint%2d", &len)
		len /= 8
		st := jen.Id("_").Op("=").Id("b").Index(jen.Id("m").Op("+").Lit(len - 1)).Comment("bounds check hint to compiler; see golang.org/issue/14808").Line()
		for i := 0; i < len; i++ {
			if tag.Order == binary.LittleEndian {
				st.Id("b").Index(jen.Id("m").Op("+").Lit(i)).Op("=").Byte().Call(n.Clone().Op(">>").Lit(8 * i)).Line()
			} else {
				st.Id("b").Index(jen.Id("m").Op("+").Lit(i)).Op("=").Byte().Call(n.Clone().Op(">>").Lit(8 * (len - 1 - i))).Line()
			}
		}
		st.Id("m").Op("+=").Lit(len)
		return st
	case "float32", "float64":
		return jen.Block(jen.Id("tmp").Op(":=").Qual("math", "F"+goType[1:]+"bits").Call(jen.Id(goType).Call(n)).Line().Add(tag.typeToBinary("uint"+goType[5:], jen.Id("tmp"))))

	default:
		log.Printf("type serialization unknown %s", goType)
		return jen.Qual("fmt", "Println").Call(jen.Params(jen.Id(goType)).Params(n.Clone()))

	}

}
