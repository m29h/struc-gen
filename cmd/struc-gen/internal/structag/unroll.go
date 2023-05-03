package structag

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
)

func (tag *StrucTag) unroll(t types.Type, n *jen.Statement) (*jen.Statement, error) {
	if tag.Type == "pad" {
		lvar, lmax := tag.GetNewLoopVar()
		cast := jen.Id("b").Index(jen.Id("m").Op("+").Id("i")).Op("=").Lit(0)
		cast = jen.For(lvar.Clone().Op(":=").Lit(0), lvar.Clone().Op("<").Add(lmax), lvar.Clone().Op("++")).Block(cast).Line()
		if tag.method != MarshalBinary {
			cast = jen.Null()
		}
		return cast.Id("m").Op("+=").Add(lmax), nil
	}
	switch v := t.(type) {
	case *types.Basic:
		switch tag.method {
		case MarshalBinary:
			return tag.TypeToBinary(v.Underlying().Underlying().String(), n.Clone()), nil
		case UnmarshalBinary, Size:
			return tag.BinaryToType(v.Underlying().Underlying().String(), n.Clone()), nil
		default:
			panic("Unimplemented byte packing direction")
		}
	case *types.Named:
		switch tag.method {
		case MarshalBinary:
			return jen.Id("m").Op("+=").Add(n.Clone().Dot("MarshalBinary").Call(jen.Id("b").Index(jen.Id("m").Op(":")))), nil
		case UnmarshalBinary:
			return jen.Id("m").Op("+=").Add(n.Clone().Dot("UnmarshalBinary").Call(jen.Id("b").Index(jen.Id("m").Op(":")))), nil
		case Size:
			return jen.Id("m").Op("+=").Add(n.Clone().Dot("SizeOf").Call()), nil

		default:
			panic("Unimplemented byte packing direction")
		}

	case *types.Slice, *types.Array:
		lvar, lmax := tag.GetNewLoopVar()
		cast, err := tag.unroll(v.(interface{ Elem() types.Type }).Elem(), n.Clone().Index(lvar))
		if err != nil {
			return jen.Comment("Serialization Error"), err
		}
		return jen.For(lvar.Clone().Op(":=").Lit(0), lvar.Clone().Op("<").Add(lmax), lvar.Clone().Op("++")).Block(cast), nil

	case *types.Pointer:
		return tag.unroll(v.Elem(), jen.Parens(jen.Op("*").Add(n)))

	default:
		return nil, fmt.Errorf("struct field type not handled: %T", v)
	}

}
