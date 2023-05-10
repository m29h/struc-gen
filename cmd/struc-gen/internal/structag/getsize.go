package structag

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

func (tag *StrucTag) GetSize() (*jen.Statement, error) {
	tag.method = Size
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
