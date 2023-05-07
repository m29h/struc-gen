package structag

import "github.com/dave/jennifer/jen"

type TagPrompter func(tag *StrucTag) (*jen.Statement, error)

func MethodBuilderList() map[string]TagPrompter {
	return map[string]TagPrompter{
		"MarshalBinary":   func(tag *StrucTag) (*jen.Statement, error) { return tag.Pack() },
		"UnmarshalBinary": func(tag *StrucTag) (*jen.Statement, error) { return tag.Unpack() },
		"SizeOf":          func(tag *StrucTag) (*jen.Statement, error) { return tag.GetSize() },
	}
}
