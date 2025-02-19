package main

type Token struct {
	typ   int
	value string
}

type Compiler struct {
	nextToken func() *Token
    varNames map[string]string // className -> varName
    classNames []string // classNames
    subroutineNames map[string]string // className -> subroutineName
    xmlOutput string
}

func newCompiler(next func() *Token) Compiler {
	return Compiler{nextToken: next}
}

func (c* Compiler) class() {
    if c.nextToken().value == "class" {
        // TODO: identifier, add to slice
        // {
        // var decs
        // subroutine decs
        // }
    }
    panic("not implemented!")
}

func (c* Compiler) classVarDec() {
    panic("not implemented!")
}

func (c* Compiler) typ() {
    panic("not implemented!")
}

func (c* Compiler) subroutineDec() {
    panic("not implemented!")
}

func (c* Compiler) parameterList() {
    panic("not implemented!")
}

func (c* Compiler) subroutineBody() {
    panic("not implemented!")
}

func (c* Compiler) varDev() {
    panic("not implemented!")
}

func (c* Compiler) className() {
    panic("not implemented!")
}

func (c* Compiler) subroutineName() {
    panic("not implemented!")
}

func (c* Compiler) varName() {
    panic("not implemented!")
}

func (c* Compiler) statement() {
    panic("not implemented!")
}

func (c* Compiler) letStatement() {
    panic("not implemented!")
}

func (c* Compiler) ifStatement() {
    panic("not implemented!")
}

func (c* Compiler) whileStatement() {
    panic("not implemented!")
}

func (c* Compiler) doStatement() {
    panic("not implemented!")
}

func (c* Compiler) returnStatement() {
    panic("not implemented!")
}

func (c* Compiler) expression() {
    panic("not implemented!")
}

func (c* Compiler) term() {
    panic("not implemented!")
}

func (c* Compiler) subroutineCall() {
    panic("not implemented!")
}

func (c* Compiler) expressionList() {
    panic("not implemented!")
}

func (c* Compiler) op() {
    panic("not implemented!")
}

func (c* Compiler) unaryOp() {
    panic("not implemented!")
}

func (c* Compiler) keywordConstant() {
    panic("not implemented!")
}
