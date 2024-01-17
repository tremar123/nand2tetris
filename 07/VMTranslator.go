package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	C_ARITHMETIC = iota
	C_PUSH
	C_POP
)

const (
	O_ADD = "add"
	O_SUB = "sub"
	O_NEG = "neg"
	O_EQ  = "eq"
	O_GT  = "gt"
	O_LT  = "lt"
	O_AND = "and"
	O_OR  = "or"
	O_NOT = "not"

	SEG_CONSTANT = "constant"
	SEG_ARGUMENT = "argument"
	SEG_LOCAL    = "local"
	SEG_STATIC   = "static"
	SEG_THIS     = "this"
	SEG_THAT     = "that"
	SEG_POINTER  = "pointer"
	SEG_TEMP     = "temp"
)

var jump = 0

func main() {
	vmCode, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	filepath := strings.Split(os.Args[1], ".")[0]
	pathSplit := strings.Split(filepath, "/")
	filename := pathSplit[len(pathSplit)-1]
	lines := strings.Split(string(vmCode), "\n")
	asm := ""

	for _, line := range lines {
		commandType, arg1, arg2 := parse(line)
		if commandType == -1 {
			continue
		}

		code := generateCode(filename, line, commandType, arg1, arg2)
		asm += code
	}

	err = os.WriteFile(filepath+".asm", []byte(asm), 0600)
	if err != nil {
		fmt.Println(err.Error())
	}
}

/* Returns command type, arg1, arg2 (value for POP and PUSH commands) */
func parse(command string) (int, string, int) {
	command = strings.TrimSpace(command)

	if command == "" || command[:2] == "//" {
		return -1, "", 0
	}

	split := strings.Split(command, " ")

	if len(split) == 1 {
		return C_ARITHMETIC, split[0], 0
	}

	arg2, err := strconv.Atoi(split[2])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch split[0] {
	case "push":
		return C_PUSH, split[1], arg2
	case "pop":
		return C_POP, split[1], arg2
	default:
		return -1, "", 0
	}
}

/* Returns asm code for instruction */
func generateCode(filename string, commandLine string, commandType int, arg1 string, arg2 int) string {
	code := "// " + commandLine + "\n"
	switch commandType {
	case C_PUSH, C_POP:
		code += writePushPop(commandType, filename, arg1, arg2)
	case C_ARITHMETIC:
		switch arg1 {
		case O_ADD:
			code += binaryOperation("+")
		case O_SUB:
			code += binaryOperation("-")
		case O_NEG:
			code += unaryOperation("-")
		case O_EQ:
			code += compareOperation("JEQ")
		case O_GT:
			code += compareOperation("JGT")
		case O_LT:
			code += compareOperation("JLT")
		case O_AND:
			code += binaryOperation("&")
		case O_OR:
			code += binaryOperation("|")
		case O_NOT:
			code += unaryOperation("!")
		}
	}

	return code
}

func writePushPop(commandType int, filename string, segment string, arg2 int) string {
	code := ""
	switch segment {
	case SEG_ARGUMENT, SEG_LOCAL, SEG_THIS, SEG_THAT:
		segmentSymbol := ""
		switch segment {
		case SEG_ARGUMENT:
			segmentSymbol = "ARG"
		case SEG_LOCAL:
			segmentSymbol = "LCL"
		case SEG_THIS:
			segmentSymbol = "THIS"
		case SEG_THAT:
			segmentSymbol = "THAT"
		}
		if commandType == C_PUSH {
			code = A(segmentSymbol) + C("D=M") + A(arg2) + C("A=A+D") + C("D=M") + insertToStack() + incSP()
		} else {
			code = A(segmentSymbol) + C("D=M") + A(arg2) + C("D=D+A") + A("R13") + C("M=D") + decSP() + getFromStack() + A("R13") + C("A=M") + C("M=D")
		}
	case SEG_TEMP:
		// 5 is where temp segment starts
		if commandType == C_PUSH {
			code = A(5+arg2) + C("D=M") + insertToStack() + incSP()
		} else {
			code = decSP() + getFromStack() + A(5+arg2) + C("M=D")
		}
	case SEG_CONSTANT:
		code = A(arg2) + C("D=A") + insertToStack() + incSP()
	case SEG_STATIC:
		if commandType == C_PUSH {
			code = A(filename, ".", arg2) + C("D=M") + insertToStack() + incSP()
		} else {
			code = decSP() + getFromStack() + A(filename, ".", arg2) + C("M=D")
		}
	case SEG_POINTER:
		segmentSymbol := ""
		if arg2 == 0 {
			segmentSymbol = "THIS"
		} else if arg2 == 1 {
			segmentSymbol = "THAT"
		} else {
			panic("push/pop pointer only accepts 0 or 1")
		}
		if commandType == C_PUSH {
			code = A(segmentSymbol) + C("D=M") + insertToStack() + incSP()
		} else {
			code = decSP() + getFromStack() + A(segmentSymbol) + C("M=D")
		}
	}

	return code
}

func A(value ...any) string {
	value = append([]any{"@"}, value...)
	value = append(value, "\n")
	return fmt.Sprint(value...)
}

func C(str string) string {
	return str + "\n"
}

func insertToStack() string {
	return A("SP") + C("A=M") + C("M=D")
}

func getFromStack() string {
	return A("SP") + C("A=M") + C("D=M")
}

func incSP() string {
	return A("SP") + C("M=M+1")
}

func decSP() string {
	return A("SP") + C("M=M-1")
}

func unaryOperation(operation string) string {
	return A("SP") + C("A=M-1") + C("M="+operation+"M")
}

func binaryOperation(operation string) string {
	return decSP() + C("A=M") + C("D=M") + C("A=A-1") + C("M=M"+operation+"D")
}

func compareOperation(operation string) string {
	label := fmt.Sprintf("END_%s_%d", operation, jump)
	jump++
	return decSP() + C("A=M") + C("D=M") + C("A=A-1") + C("D=M-D") + C("M=-1") + A(label) + C("D;"+operation) + A("SP") + C("A=M-1") + C("M=0") + fmt.Sprintf("(%s)\n", label)
}
