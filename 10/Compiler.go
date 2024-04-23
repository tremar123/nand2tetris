package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

const (
	KEYWORD = iota + 1
	SYMBOL
	INTEGER_CONSTANT
	STRING_CONSTANT
	IDENTIFIER
)

// comment type
const (
	C_NONE = iota + 1
	C_ONE_LINE
	C_MULTI_LINE
)

type Token struct {
	typ   int
	value string
}

func main() {
	info, err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if info.IsDir() {
		dir, err := os.ReadDir(os.Args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, file := range dir {
			nameAndExtension := strings.Split(file.Name(), ".")
			if nameAndExtension[1] != "jack" {
				continue
			}
			compile(filepath.Join(os.Args[1], file.Name()))
		}

	} else {
		compile(os.Args[1])
	}
}

func compile(filename string) {
	tokens := tokenizer(filename)

	for _, token := range tokens {
		var typ string
		switch token.typ {
        case 1:
            typ = "KEYWORD"
        case 2:
            typ = "SYMBOL"
        case 3:
            typ = "INTEGER_CONSTANT"
        case 4:
            typ = "STRING_CONSTANT"
        case 5:
            typ = "IDENTIFIER"
		}
		fmt.Printf("%s - %q\n", typ, token.value)
	}
}

func tokenizer(filename string) []Token {
	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	inString := false
	comment := C_NONE
	tokens := []Token{}
	currentToken := ""

	for i, ch := range string(code) {
		if comment != C_NONE {
			if comment == C_ONE_LINE && ch == '\n' {
				comment = C_NONE
			} else { // C_MULTI_LINE
				if currentToken == "*" {
					if ch == '/' {
						comment = C_NONE
						currentToken = ""
						continue
					} else {
						currentToken = ""
					}
				}

				if ch == '*' {
					currentToken += "*"
				}
			}
			continue
		}

		if unicode.IsControl(ch) {
			continue
		}

		// never touch this
		if (ch != ' ' && inString == false) || (ch != ' ' && inString == true) || (ch == ' ' && inString == true) {
			currentToken += string(ch)
		}

		switch currentToken {
		case "/":
			if code[i+1] == '/' {
				comment = C_ONE_LINE
				currentToken = ""
			} else if code[i+1] == '*' {
				comment = C_MULTI_LINE
				currentToken = ""
			} else {
				tokens = append(tokens, Token{typ: SYMBOL, value: currentToken})
				currentToken = ""
			}
			// keywords
		case "class", "method", "field", "constructor", "do", "let", "function", "static", "var", "boolean", "int", "char", "void", "return", "this", "null", "true", "false", "if", "else", "while":
			if !unicode.IsLetter(rune(code[i+1])) && !unicode.IsDigit(rune(code[i+1])) {
				tokens = append(tokens, Token{typ: KEYWORD, value: currentToken})
				currentToken = ""
			}
		case "(", ")", "{", "}", "[", "]", ",", ".", ";", "=", "+", "-", "*", "&", "|", "~", ">", "<":
			tokens = append(tokens, Token{typ: SYMBOL, value: currentToken})
			currentToken = ""
		default:
			// check if current token is string literal, if there is '\' fount before '"' we are not at the end
			if strings.HasPrefix(currentToken, "\"") {
				inString = true
				if strings.HasSuffix(currentToken, "\"") && len(currentToken) > 1 && currentToken[len(currentToken)-2] != '\\' {
					// TODO: here maybe strip '"'
					str := strings.TrimPrefix(currentToken, "\"")
					str = strings.TrimSuffix(str, "\"")
					tokens = append(tokens, Token{typ: STRING_CONSTANT, value: str})
					currentToken = ""
					inString = false
					continue
				}
				continue
			}

			// should be identifier or integer constant after this
			if !unicode.IsLetter(rune(code[i+1])) && !unicode.IsDigit(rune(code[i+1])) && len(currentToken) > 0 {
				integer, err := strconv.Atoi(currentToken)

				if err != nil {
					// not integer, it's a identifier
					tokens = append(tokens, Token{typ: IDENTIFIER, value: currentToken})
					currentToken = ""
					continue
				}

				tokens = append(tokens, Token{typ: INTEGER_CONSTANT, value: fmt.Sprint(integer)})
				currentToken = ""
			}
		}
	}

	return tokens
}
