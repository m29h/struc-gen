package structag

import (
	"github.com/dave/jennifer/jen"
)

func (tag *StrucTag) render(j *jen.Statement) *jen.Statement {
	if j == nil {
		return nil
	}

	return tag.beforeStatement.Line().Add(j).Add(tag.finalize())
}
