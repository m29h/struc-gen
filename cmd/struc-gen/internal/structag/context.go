package structag

import "github.com/dave/jennifer/jen"

type Context struct {
	checkBound    *jen.Statement
	staticBitPos  int
	dynamicBitPos *jen.Statement
}

// inserts a new Null CheckBound statement placeholder
// moves the global checkBound pointer to the new placeholder statement
// returns the placeholder
func (c *Context) insertNewCheckBound() *jen.Statement {
	newCheckBound := jen.Null()
	c.checkBound = newCheckBound
	return jen.Line().Add(newCheckBound)
}

// flush the content of context state into "m" variable
// after flush "m" points to the next free byte-aligned position and
// the context is empty again (static + dynamic bit positions)
func (c *Context) Flush() *jen.Statement {
	bitP := jen.Lit(c.staticBitPos)
	if c.dynamicBitPos != nil {
		bitP.Op("+").Add(c.dynamicBitPos.Clone())
	}
	bitP = jen.Parens(bitP)
	c.staticBitPos = 0
	c.dynamicBitPos = nil

	result := jen.Id("m").Op("+=").Add(bitP.Clone().Op("/").Lit(8)).Line()
	result.If(bitP.Clone().Op("%").Lit(8).Op("!=").Lit(0)).Block(
		jen.Id("m").Op("++"),
	)
	result.Add(c.insertNewCheckBound())
	return result
}
