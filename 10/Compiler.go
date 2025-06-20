/*
TODO:

square.jack to fail line 67 expected ; got -
*/

package main


type Token struct {
	typ   int
	value string
}

type Compiler struct {
	currentToken    *Token
	nextToken       func() *Token
	varNames        map[string][]string // className -> varName
	classNames      []string            // classNames
	currentClass    string
	subroutineNames map[string][]string // className -> subroutineName
	xmlOutput       string
}

func newCompiler(next func() *Token) Compiler {
	return Compiler{
		nextToken:       next,
		varNames:        map[string][]string{},
		subroutineNames: map[string][]string{},
	}
}

func (c *Compiler) useNextToken() {
	c.currentToken = c.nextToken()
}

func (c *Compiler) identifier() {
	c.xmlOutput += "<identifier> " + c.currentToken.value + " </identifier>\n"
}

func (c *Compiler) symbol() {
	c.xmlOutput += "<symbol> " + c.currentToken.value + " </symbol>\n"
}

func (c *Compiler) keyword() {
	c.xmlOutput += "<keyword> " + c.currentToken.value + " </keyword>\n"
}

func (c *Compiler) expectSymbol(sym string) {
	if c.currentToken.typ == SYMBOL && c.currentToken.value == sym {
		c.symbol()
	} else {
		panic("expected symbol " + "'" + sym + "'")
	}
}

func (c *Compiler) expectIdentifier() {
	if c.currentToken.typ == IDENTIFIER {
		c.identifier()
	} else {
		panic("expected identifier")
	}
}

func (c *Compiler) class() {
	c.useNextToken()
	if c.currentToken.value == "class" {
		c.xmlOutput += "<class>\n"
		c.keyword()
		// identifier, add to slice
		c.useNextToken()
		if c.currentToken.typ == IDENTIFIER {
			c.identifier()
			c.classNames = append(c.classNames, c.currentToken.value)
			c.varNames[c.currentToken.value] = []string{}
			c.subroutineNames[c.currentToken.value] = []string{}
			c.currentClass = c.currentToken.value
		} else {
			panic("expected identifier")
		}
		// {
		c.useNextToken()
		c.expectSymbol("{")
		// var decs
		// subroutine decs
		for {
			c.useNextToken()
			if c.currentToken.typ == KEYWORD {
				switch c.currentToken.value {
				case "static", "field":
					c.classVarDec()
				case "constructor", "function", "method":
					c.subroutineDec()
				default:
					panic("expected variable or subroutine definition")
				}
			} else if c.currentToken.typ == SYMBOL && c.currentToken.value == "}" {
				c.xmlOutput += "<symbol> } </symbol>\n"
				break
			}
		}

		c.xmlOutput += "</class>"
	}
}

func (c *Compiler) classVarDec() {
	c.xmlOutput += "<classVarDec>\n"
	c.keyword()

	c.useNextToken()
	c.typ()
	c.useNextToken()

	if c.currentToken.typ == IDENTIFIER {
		c.identifier()
		c.varNames[c.currentClass] = append(c.varNames[c.currentClass], c.currentToken.value)
	} else {
		panic("expected identifier")
	}

	for {
		c.useNextToken()
		if c.currentToken.typ == SYMBOL {
			if c.currentToken.value == ";" {
				c.symbol()
				break
			} else if c.currentToken.value == "," {
				c.symbol()
				c.useNextToken()
				c.expectIdentifier()
			} else {
				panic("expected symbol ',' or ';'")
			}
		} else {
			panic("expected symbol ',' or ';'")
		}
	}

	c.xmlOutput += "</classVarDec>\n"
}

func (c *Compiler) subroutineDec() {
	c.xmlOutput += "<subroutineDec>\n"
	c.keyword()
	c.useNextToken()

	if c.currentToken.typ == KEYWORD && c.currentToken.value == "void" {
		c.keyword()
	} else if c.currentToken.typ == KEYWORD || c.currentToken.typ == IDENTIFIER {
		c.typ()
	} else {
		panic("expected type or void")
	}

	c.useNextToken()
	c.expectIdentifier()

	c.useNextToken()

	c.expectSymbol("(")
	c.parameterList()
	c.expectSymbol(")")

	c.useNextToken()

	c.subroutineBody()

	c.xmlOutput += "</subroutineDec>\n"
}

func (c *Compiler) typ() {
	switch c.currentToken.typ {
	case KEYWORD:
		switch c.currentToken.value {
		case "int", "boolean", "char":
			c.keyword()
		default:
			panic("invalid keyword")
		}
	case IDENTIFIER:
		c.identifier()
	default:
		panic("expected type")
	}
}

func (c *Compiler) parameterList() {
	c.xmlOutput += "<parameterList>\n"

	for {
		c.useNextToken()

		if c.currentToken.typ == KEYWORD {
			c.typ()
			c.useNextToken()
			c.expectIdentifier()
		} else if c.currentToken.typ == SYMBOL {
			if c.currentToken.value == ")" {
				break
			} else if c.currentToken.value == "," {
				c.symbol()
				c.useNextToken()
				c.typ()
				c.useNextToken()
				c.expectIdentifier()
			} else {
				panic("expected symbol ')' or ','")
			}
		}
	}

	c.xmlOutput += "</parameterList>\n"
}

func (c *Compiler) subroutineBody() {
	c.xmlOutput += "<subroutineBody>\n"
	c.expectSymbol("{")

	c.varDec()
	c.statements()

	c.expectSymbol("}")
	c.xmlOutput += "</subroutineBody>\n"
}

func (c *Compiler) varDec() {
	for {
		c.useNextToken()
		if c.currentToken.typ == KEYWORD && c.currentToken.value == "var" {
			c.xmlOutput += "<varDec>\n"
			c.keyword()
			c.useNextToken()
			c.typ()
			c.useNextToken()
			c.expectIdentifier()

			for {
				c.useNextToken()

				if c.currentToken.typ == SYMBOL {
					if c.currentToken.value == ";" {
						c.symbol()
						break
					} else if c.currentToken.value == "," {
						c.symbol()
						c.useNextToken()
						c.expectIdentifier()
					}
				} else {
					panic("expected symbol ';' or ','")
				}
			}

			c.xmlOutput += "</varDec>\n"
		} else {
			break
		}
	}
}

func (c *Compiler) statements() {
	c.xmlOutput += "<statements>\n"

	for {
		if c.currentToken.typ == SYMBOL && c.currentToken.value == "}" {
			break
		}

		if c.currentToken.typ != KEYWORD {
			panic("expected keyword: 'let', 'if', 'while', 'do' or 'return'")
		}

		switch c.currentToken.value {
		case "let":
			c.letStatement()
			c.useNextToken()
		case "if":
			c.ifStatement()
		case "while":
			c.whileStatement()
			c.useNextToken()
		case "do":
			c.doStatement()
			c.useNextToken()
		case "return":
			c.returnStatement()
			c.useNextToken()
		default:
			panic("expected keyword: 'let', 'if', 'while', 'do' or 'return'")
		}
	}
	c.xmlOutput += "</statements>\n"
}

func (c *Compiler) letStatement() {
	c.xmlOutput += "<letStatement>\n"
	c.keyword()
	c.useNextToken()
	c.expectIdentifier()
	c.useNextToken()

	if c.currentToken.value == "[" {
		c.expectSymbol("[")
		c.useNextToken()
		c.expression()
		c.expectSymbol("]")
		c.useNextToken()
	}

	c.expectSymbol("=")
	c.useNextToken()

	c.expression()
	c.expectSymbol(";")
	c.xmlOutput += "</letStatement>\n"
}

func (c *Compiler) ifStatement() {
	c.xmlOutput += "<ifStatement>\n"
	c.keyword()
	c.useNextToken()
	c.expectSymbol("(")
	c.useNextToken()
	c.expression()
	c.expectSymbol(")")
	c.useNextToken()
	c.expectSymbol("{")
	c.useNextToken()
	c.statements()
	c.expectSymbol("}")
	c.useNextToken()

	if c.currentToken.typ == KEYWORD && c.currentToken.value == "else" {
		c.keyword()
		c.useNextToken()
		c.expectSymbol("{")
		c.useNextToken()
		c.statements()
		c.expectSymbol("}")
		c.useNextToken()
	}
	c.xmlOutput += "</ifStatement>\n"
}

func (c *Compiler) whileStatement() {
	c.xmlOutput += "<whileStatement>\n"
	c.keyword()
	c.useNextToken()
	c.expectSymbol("(")
	c.useNextToken()
	c.expression()
	c.expectSymbol(")")
	c.useNextToken()
	c.expectSymbol("{")
	c.useNextToken()
	c.statements()
	c.expectSymbol("}")
	c.xmlOutput += "</whileStatement>\n"
}

func (c *Compiler) doStatement() {
	c.xmlOutput += "<doStatement>\n"
	c.keyword()
	c.useNextToken()
	prevToken := c.currentToken
	c.useNextToken()
	c.subroutineCall(prevToken)
	c.useNextToken()
	c.expectSymbol(";")
	c.xmlOutput += "</doStatement>\n"
}

func (c *Compiler) returnStatement() {
	c.xmlOutput += "<returnStatement>\n"
	c.keyword()
	c.useNextToken()
	if c.currentToken.typ == SYMBOL && c.currentToken.value == ";" {
		c.symbol()
	} else {
		c.expression()
		c.expectSymbol(";")
	}
	c.xmlOutput += "</returnStatement>\n"
}

func (c *Compiler) expression() {
	c.xmlOutput += "<expression>\n"
	c.term()

Loop:
	for {
		if c.currentToken.typ == SYMBOL {
			switch c.currentToken.value {
			case "+", "-", "*", "/", "&", "|", "<", ">", "=":
				c.op()
				c.useNextToken()
				c.term()
			default:
				break Loop
			}
		} else {
			break
		}
	}
	c.xmlOutput += "</expression>\n"
}

func (c *Compiler) term() {
	c.xmlOutput += "<term>\n"
	switch c.currentToken.typ {
	case KEYWORD:
		c.keywordConstant()
		c.useNextToken()
	case INTEGER_CONSTANT:
		c.xmlOutput += "<integerConstant> " + c.currentToken.value + " </integerConstant>\n"
		c.useNextToken()
	case STRING_CONSTANT:
		c.xmlOutput += "<stringConstant> " + c.currentToken.value + " </stringConstant>\n"
		c.useNextToken()
	case SYMBOL:
		switch c.currentToken.value {
		case "(":
			c.xmlOutput += "<symbol> ( </symbol>\n"
			c.useNextToken()
			c.expression()
			if c.currentToken.value != ")" {
				panic("expected ')'")
			}
			c.xmlOutput += "<symbol> ) </symbol>\n"
			c.useNextToken()
		default:
			c.unaryOp()
			c.useNextToken()
			c.term()
		}
	case IDENTIFIER:
		prevToken := c.currentToken
		c.useNextToken()

		switch c.currentToken.value {
		case ".", "(":
			c.subroutineCall(prevToken)
			c.useNextToken()
		case "[":
			if prevToken.typ == IDENTIFIER {
				c.xmlOutput += "<identifier> " + prevToken.value + " </identifier>\n"
			} else {
				panic("expected identifier")
			}

			c.expectSymbol("[")
			c.useNextToken()
			c.expression()
			c.expectSymbol("]")
			c.useNextToken()
		default:
			c.xmlOutput += "<identifier> " + prevToken.value + " </identifier>\n"
		}
	default:
		panic("invalid term")
	}
	c.xmlOutput += "</term>\n"
}

func (c *Compiler) subroutineCall(prev *Token) {
	if prev.typ == IDENTIFIER {
		c.xmlOutput += "<identifier> " + prev.value + " </identifier>\n"
	} else {
		panic("expected identifier")
	}

	if c.currentToken.typ != SYMBOL {
		panic("expected symbol '(' or '.'")
	}

	switch c.currentToken.value {
	case "(":
		c.symbol()
		c.useNextToken()
		c.expressionList()
		c.expectSymbol(")")
	case ".":
		c.symbol()
		c.useNextToken()
		c.expectIdentifier()
		c.useNextToken()
		c.expectSymbol("(")
		c.useNextToken()
		c.expressionList()
		c.expectSymbol(")")
	}
}

func (c *Compiler) expressionList() {
	c.xmlOutput += "<expressionList>\n"
	for {
		if c.currentToken.typ == SYMBOL {
			if c.currentToken.value == ")" {
				break
			} else if c.currentToken.value == "," {
				c.xmlOutput += "<symbol> , </symbol>\n"
			} else if c.currentToken.value == "(" {
				c.expression()
				continue
			}
			c.useNextToken()
		} else {
			c.expression()
		}
	}
	c.xmlOutput += "</expressionList>\n"
}

func (c *Compiler) op() {
	switch c.currentToken.value {
	case "+", "-", "*", "/", "|", "=":
		c.xmlOutput += "<symbol> " + c.currentToken.value + " </symbol>\n"
	case "&":
		c.xmlOutput += "<symbol> &amp; </symbol>\n"
	case "<":
		c.xmlOutput += "<symbol> &lt; </symbol>\n"
	case ">":
		c.xmlOutput += "<symbol> &gt; </symbol>\n"
	default:
		panic("expected op: '+', '-', '*', '/', '&', '|', '<', '>' or '='")
	}
}

func (c *Compiler) unaryOp() {
	switch c.currentToken.value {
	case "-", "~":
		c.xmlOutput += "<symbol> " + c.currentToken.value + " </symbol>\n"
	default:
		panic("expected unary op: '-' or '~'")
	}
}

func (c *Compiler) keywordConstant() {
	c.xmlOutput += "<keyword> "
	switch c.currentToken.value {
	case "true", "false", "this", "null":
		c.xmlOutput += c.currentToken.value
	default:
		panic("expected keyword: 'true', 'false', 'this' or 'null'")
	}
	c.xmlOutput += " </keyword>\n"
}
