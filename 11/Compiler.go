package main

import "fmt"

const (
	Declared = iota
	Used
)

type Token struct {
	typ   int
	value string
}

type SymbolTable struct {
	name     string
	category string
	index    int
}

type Compiler struct {
	currentToken          *Token
	nextToken             func() *Token
	symbolTable           []SymbolTable
	currentFieldIndex     int
	currentStaticIndex    int
	classNames            []string // classNames
	currentClass          string
	subroutineSymbolTable []SymbolTable
	currentArgIndex       int
	currentLocalIndex     int
	xmlOutput             string
}

func newCompiler(next func() *Token) Compiler {
	return Compiler{
		nextToken: next,
	}
}

func (c *Compiler) useNextToken() {
	c.currentToken = c.nextToken()
}

func (c *Compiler) identifier(name string, index int, category string, usage int) {
	c.xmlOutput += "<identifier>\n"
	c.xmlOutput += "<name> " + name + " </name>\n"
	c.xmlOutput += "<index> " + fmt.Sprint(index) + " </index>\n"
	c.xmlOutput += "<category> " + category + " </category>\n"
	var usageStr string
	if usage == 0 {
		usageStr = "declared"
	} else {
		usageStr = "used"
	}
	c.xmlOutput += "<usage> " + usageStr + " </usage>\n"
	c.xmlOutput += "</identifier>\n"
}

func (c *Compiler) symbol() {
	c.xmlOutput += "<symbol> " + c.currentToken.value + " </symbol>\n"
}

func (c *Compiler) keyword() {
	c.xmlOutput += "<keyword> " + c.currentToken.value + " </keyword>\n"
}

func (c *Compiler) expectSymbol(sym string) {
	if c.currentToken.typ == Symbol && c.currentToken.value == sym {
		c.symbol()
	} else {
		panic("expected symbol " + "'" + sym + "'")
	}
}

func (c *Compiler) class() {
	c.useNextToken()
	if c.currentToken.value == "class" {
		c.currentFieldIndex = 0
		c.currentStaticIndex = 0
		c.symbolTable = []SymbolTable{}
		c.xmlOutput += "<class>\n"
		c.keyword()
		// identifier, add to slice
		c.useNextToken()
		if c.currentToken.typ == Identifier {
			c.identifier(c.currentToken.value, 0, "class", Declared)
			c.classNames = append(c.classNames, c.currentToken.value)
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
			if c.currentToken.typ == Keyword {
				switch c.currentToken.value {
				case "static", "field":
					c.classVarDec()
				case "constructor", "function", "method":
					c.subroutineDec()
				default:
					panic("expected variable or subroutine definition")
				}
			} else if c.currentToken.typ == Symbol && c.currentToken.value == "}" {
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

	currentCategory := c.currentToken.value
	var currentIndex int

	if currentCategory == "field" {
		currentIndex = c.currentFieldIndex
		c.currentFieldIndex++
	} else {
		currentIndex = c.currentStaticIndex
		c.currentStaticIndex++
	}

	c.useNextToken()
	c.typ()
	c.useNextToken()

	if c.currentToken.typ == Identifier {
		c.symbolTable = append(c.symbolTable, SymbolTable{
			name:     c.currentToken.value,
			index:    currentIndex,
			category: currentCategory,
		})
		c.identifier(c.currentToken.value, currentIndex, currentCategory, Declared)
		currentIndex++
	} else {
		panic("expected identifier")
	}

	for {
		c.useNextToken()
		if c.currentToken.typ == Symbol {
			if c.currentToken.value == ";" {
				c.symbol()
				break
			} else if c.currentToken.value == "," {
				c.symbol()
				c.useNextToken()
				if c.currentToken.typ != Identifier {
					panic("expected identifier")
				}
				c.symbolTable = append(c.symbolTable, SymbolTable{
					name:     c.currentToken.value,
					index:    currentIndex,
					category: currentCategory,
				})
				c.identifier(c.currentToken.value, currentIndex, currentCategory, Declared)
				currentIndex++
				if currentCategory == "field" {
					c.currentFieldIndex++
				} else {
					c.currentStaticIndex++
				}
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
	c.currentArgIndex = 0
	c.currentLocalIndex = 0
	c.subroutineSymbolTable = []SymbolTable{}
	c.xmlOutput += "<subroutineDec>\n"
	c.keyword()
	c.useNextToken()

	if c.currentToken.typ == Keyword && c.currentToken.value == "void" {
		c.keyword()
	} else if c.currentToken.typ == Keyword || c.currentToken.typ == Identifier {
		c.typ()
	} else {
		panic("expected type or void")
	}

	c.useNextToken()
	if c.currentToken.typ != Identifier {
		panic("expected identifier")
	}
	c.identifier(c.currentToken.value, 0, "subroutine", Declared)

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
	case Keyword:
		switch c.currentToken.value {
		case "int", "boolean", "char":
			c.keyword()
		default:
			panic("invalid keyword")
		}
	case Identifier:
		c.identifier(c.currentToken.value, 0, "class", Used)
	default:
		panic("expected type")
	}
}

func (c *Compiler) parameterList() {
	c.xmlOutput += "<parameterList>\n"

	for {
		c.useNextToken()

		if c.currentToken.typ == Keyword {
			c.typ()
			c.useNextToken()
			if c.currentToken.typ != Identifier {
				panic("expected identifier")
			}
			c.subroutineSymbolTable = append(c.subroutineSymbolTable, SymbolTable{
				name:     c.currentToken.value,
				index:    c.currentArgIndex,
				category: "argument",
			})
			c.identifier(c.currentToken.value, c.currentArgIndex, "argument", Declared)
			c.currentArgIndex++
		} else if c.currentToken.typ == Symbol {
			if c.currentToken.value == ")" {
				break
			} else if c.currentToken.value == "," {
				c.symbol()
				c.useNextToken()
				c.typ()
				c.useNextToken()
				if c.currentToken.typ != Identifier {
					panic("expected identifier")
				}
				c.subroutineSymbolTable = append(c.subroutineSymbolTable, SymbolTable{
					name:     c.currentToken.value,
					index:    c.currentArgIndex,
					category: "argument",
				})
				c.identifier(c.currentToken.value, c.currentArgIndex, "argument", Declared)
				c.currentArgIndex++
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
		if c.currentToken.typ == Keyword && c.currentToken.value == "var" {
			c.xmlOutput += "<varDec>\n"
			c.keyword()
			c.useNextToken()
			c.typ()
			c.useNextToken()
			if c.currentToken.typ != Identifier {
				panic("expected identifier")
			}
			c.subroutineSymbolTable = append(c.subroutineSymbolTable, SymbolTable{
				name:     c.currentToken.value,
				index:    c.currentLocalIndex,
				category: "local",
			})
			c.identifier(c.currentToken.value, c.currentLocalIndex, "local", Declared)
			c.currentLocalIndex++

			for {
				c.useNextToken()

				if c.currentToken.typ == Symbol {
					if c.currentToken.value == ";" {
						c.symbol()
						break
					} else if c.currentToken.value == "," {
						c.symbol()
						c.useNextToken()
						if c.currentToken.typ != Identifier {
							panic("expected identifier")
						}
						c.subroutineSymbolTable = append(c.subroutineSymbolTable, SymbolTable{
							name:     c.currentToken.value,
							index:    c.currentLocalIndex,
							category: "local",
						})
						c.identifier(c.currentToken.value, c.currentLocalIndex, "local", Declared)
						c.currentLocalIndex++
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

	for c.currentToken.typ != Symbol || c.currentToken.value != "}" {

		if c.currentToken.typ != Keyword {
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
	if c.currentToken.typ != Identifier {
		panic("expected identifier")
	}
	category, index := c.findSymbol(c.currentToken.value)
	if category == "" {
		panic("symbol not found")
	}
	c.identifier(c.currentToken.value, index, category, Used)
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

	if c.currentToken.typ == Keyword && c.currentToken.value == "else" {
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
	if c.currentToken.typ == Symbol && c.currentToken.value == ";" {
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
		if c.currentToken.typ == Symbol {
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
	case Keyword:
		c.keywordConstant()
		c.useNextToken()
	case IntegerConstant:
		c.xmlOutput += "<integerConstant> " + c.currentToken.value + " </integerConstant>\n"
		c.useNextToken()
	case StringConstant:
		c.xmlOutput += "<stringConstant> " + c.currentToken.value + " </stringConstant>\n"
		c.useNextToken()
	case Symbol:
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
	case Identifier:
		prevToken := c.currentToken
		c.useNextToken()

		switch c.currentToken.value {
		case ".", "(":
			c.subroutineCall(prevToken)
			c.useNextToken()
		case "[":
			if prevToken.typ == Identifier {
				category, index := c.findSymbol(prevToken.value)
				if category == "" {
					panic("symbol not found")
				}
				c.identifier(prevToken.value, index, category, Used)
			} else {
				panic("expected identifier")
			}

			c.expectSymbol("[")
			c.useNextToken()
			c.expression()
			c.expectSymbol("]")
			c.useNextToken()
		default:
			category, index := c.findSymbol(prevToken.value)
			if category == "" {
				panic("symbol not found")
			}
			c.identifier(prevToken.value, index, category, Used)
		}
	default:
		panic("invalid term")
	}
	c.xmlOutput += "</term>\n"
}

func (c *Compiler) subroutineCall(prev *Token) {
	if prev.typ == Identifier {
		category, index := c.findSymbol(prev.value)
		if category == "" {
			category = "class"
		}
		c.identifier(prev.value, index, category, Used)
	} else {
		panic("expected identifier")
	}

	if c.currentToken.typ != Symbol {
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
		if c.currentToken.typ != Identifier {
			panic("expected identifier")
		}
		c.identifier(c.currentToken.value, 0, "subroutine", Used)
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
		if c.currentToken.typ == Symbol {
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

func (c *Compiler) findSymbol(name string) (string, int) {
	for _, i := range c.subroutineSymbolTable {
		if i.name == name {
			return i.category, i.index
		}
	}

	for _, i := range c.symbolTable {
		if i.name == name {
			return i.category, i.index
		}
	}

	return "", 0
}
