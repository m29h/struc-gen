package generator

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/m29h/struc-gen/cmd/struc-gen/internal/structag"
)

type MethodBuilder struct {
	structType     *types.Struct
	sourceTypeName string
}

func NewMethodBuilder(TypeName string, Type *types.Struct) *MethodBuilder {
	return &MethodBuilder{
		structType:     Type,
		sourceTypeName: TypeName,
	}
}

func (m *MethodBuilder) MakeMethodCode(gen structag.TagPrompter) ([]jen.Code, error) {
	var funcBlock []jen.Code

	funcBlock = append(funcBlock, jen.Id("m").Op(":=").Lit(0))

	ctx := &structag.Context{}

	for i := 0; i < m.structType.NumFields(); i++ {

		tag, err := structag.NewStrucTag(ctx, m.structType.Field(i), m.structType.Tag(i))
		if err != nil {
			return nil, fmt.Errorf("unable to serialize Field %s: %w", m.structType.Field(i), err)
		}

		// Generate code for each changeset field
		p, err := gen(tag)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize Field %s: %w", m.structType.Field(i), err)
		}
		funcBlock = append(funcBlock, p)

	}

	// flush the remaining context into "m" and return the value

	funcBlock = append(funcBlock, ctx.Flush(), jen.Return(jen.Id("m")))
	return funcBlock, nil
}

func (m *MethodBuilder) MethodInterfaceBuilder(methodname string) (*jen.Statement, error) {
	switch methodname {
	case "MarshalBinary", "UnmarshalBinary":
		return jen.Func().Params(structag.RootStructName().Op("*").Id(m.sourceTypeName)).Id(methodname).Params(jen.Id("b").Index().Byte()).Int(), nil
	case "SizeOf":
		return jen.Func().Params(structag.RootStructName().Op("*").Id(m.sourceTypeName)).Id(methodname).Params().Int(), nil
	default:
		return nil, fmt.Errorf("MethodInterfaceBuilder: Unknown methodname %s", methodname)
	}
}

func (m *MethodBuilder) MethodHeaderBuilder(methodname string) *jen.Statement {
	switch methodname {
	case "MarshalBinary":
		return jen.If(jen.Len(jen.Id("b")).Op("<").Add(structag.RootStructName()).Dot("SizeOf").Call()).Block(
			jen.Return(jen.Lit(0)), //return zero if b[] is not of sufficient length
		)
	default:
		return jen.Null()
	}
}

func (m *MethodBuilder) MakeMethods() *jen.Statement {
	result := jen.Null()
	for methodname, builder := range structag.MethodBuilderList() {
		block, err := m.MakeMethodCode(builder)
		if err != nil {
			fmt.Println(err)
			return result
		}
		definition, err := m.MethodInterfaceBuilder(methodname)
		if err != nil {
			fmt.Println(err)
			return result
		}
		header := m.MethodHeaderBuilder(methodname)
		// Build method
		def := append([]jen.Code{header}, block...)
		result.Add(definition).Block(def...).Line().Line()

	}
	return result
}
