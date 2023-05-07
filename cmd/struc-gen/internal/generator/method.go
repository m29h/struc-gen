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
	lastWasBit     bool
}

// checks the mb local variable and introduces necessary code to byte align the next statement
func (m *MethodBuilder) byteAlign() *jen.Statement {
	st := jen.If(jen.Id("m").Op("%").Lit(8).Op(">").Lit(0)).Block(jen.Id("m").Op("+=").Lit(8).Op("-").Id("m").Op("%").Lit(8)).Line()
	return st
}

func NewMethodBuilder(TypeName string, Type *types.Struct) *MethodBuilder {
	return &MethodBuilder{
		structType:     Type,
		sourceTypeName: TypeName,
		lastWasBit:     false,
	}
}

func (m *MethodBuilder) MakeMethodCode(gen structag.TagPrompter) (*jen.Statement, error) {
	var funcBlock []jen.Code
	anyBitFields := false

	// 2. Build "m := make(map[string]interface{})"
	funcBlock = append(funcBlock, jen.Id("m").Op(":=").Lit(0))

	for i := 0; i < m.structType.NumFields(); i++ {

		tag := structag.NewStrucTag(m.structType.Field(i), m.structType.Tag(i))
		// Generate code for each changeset field
		p, err := gen(tag)
		if err != nil {
			return nil, fmt.Errorf("unable to serialize Field %s: %w", m.structType.Field(i), err)
		}
		if tag.GetBitLen() < 8 {
			anyBitFields = true
			m.lastWasBit = true
		} else {
			if m.lastWasBit {
				// prepend byte align statement
				// to make sure that the tag is written to the next free byte aligned position
				p = m.byteAlign().Add(p.Clone())
			}
			m.lastWasBit = false
		}
		funcBlock = append(funcBlock, p)

	}
	if anyBitFields {
		//prepend the mb variable definition
		//funcBlock = append([]jen.Code{jen.Id("mb").Op(":=").Lit(0)}, funcBlock...)
		//append byteAlign statement to make sure the m variable is up to date before returning
		funcBlock = append(funcBlock, m.byteAlign())
	}
	// 4. Build return statement
	funcBlock = append(funcBlock, jen.Return(jen.Id("m").Op("/").Lit(8)))
	return jen.Block(funcBlock...), nil
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

		// Build method
		result.Add(definition).Add(block).Line().Line()

	}
	return result
}
