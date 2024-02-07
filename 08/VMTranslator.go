package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	C_ARITHMETIC = iota
	C_PUSH
	C_POP
	C_GOTO
	C_IF_GOTO
	C_LABEL
	C_CALL
	C_RETURN
	C_FUNCTION
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
var returnCount = 0
var currentFunc = "TOP_LEVEL"

func main() {
	info, err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	asm := ""
	newFileName := ""

	if info.IsDir() {
		dir, err := os.ReadDir(os.Args[1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		asm += writeInit()

		for _, file := range dir {
			nameAndExtension := strings.Split(file.Name(), ".")
			if nameAndExtension[1] != "vm" {
				continue
			}

			processFile(filepath.Join(os.Args[1], file.Name()), nameAndExtension[0], &asm)
		}

		newFileName = filepath.Join(os.Args[1], filepath.Base(os.Args[1])+".asm")
	} else {
		filename := strings.Split(filepath.Base(os.Args[1]), ".")[0]
		dir := filepath.Dir(os.Args[1])
		newFileName = filepath.Join(dir, filename+".asm")
		processFile(os.Args[1], filename, &asm)
	}

	err = os.WriteFile(newFileName, []byte(asm), 0600)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func processFile(filePath string, filename string, asm *string) {
	vmCode, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	lines := strings.Split(string(vmCode), "\n")

	for _, line := range lines {
		commandType, arg1, arg2 := parse(line)
		if commandType == -1 {
			continue
		}

		code := generateCode(filename, line, commandType, arg1, arg2)
		*asm += code
	}
}

/* Returns command type, arg1, arg2 (value for POP and PUSH commands) */
func parse(command string) (int, string, int) {
	command = strings.TrimSpace(command)

	if command == "" || command[:2] == "//" {
		return -1, "", 0
	}

	command = strings.Split(command, "//")[0]
	command = strings.TrimSpace(command)

	split := strings.Split(command, " ")

	if len(split) == 1 {
		if split[0] == "return" {
			return C_RETURN, "", 0
		}
		return C_ARITHMETIC, split[0], 0
	}

	if len(split) == 2 {
		switch split[0] {
		case "goto":
			return C_GOTO, split[1], 0
		case "if-goto":
			return C_IF_GOTO, split[1], 0
		case "label":
			return C_LABEL, split[1], 0
		}
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
	case "call":
		return C_CALL, split[1], arg2
	case "function":
		return C_FUNCTION, split[1], arg2
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
	case C_GOTO:
		code += writeGoto(arg1)
	case C_IF_GOTO:
		code += writeIfGoto(arg1)
	case C_LABEL:
		code += writeLabel(arg1)
	case C_CALL:
		code += writeCall(arg1, arg2)
	case C_RETURN:
		code += writeReturn()
	case C_FUNCTION:
		code += writeFunction(arg1, arg2)
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

func writeInit() string {
	return A(256) + C("D=A") + A("SP") + C("M=D") + writeCall("Sys.init", 0)
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
			code = A(segmentSymbol) + C("D=M") + A(arg2) + C("A=A+D") + C("D=M") + insertDToStack() + incSP()
		} else {
			code = A(segmentSymbol) + C("D=M") + A(arg2) + C("D=D+A") + A("R13") + C("M=D") + decSP() + getFromStack() + A("R13") + C("A=M") + C("M=D")
		}
	case SEG_TEMP:
		// 5 is where temp segment starts
		if commandType == C_PUSH {
			code = A(5+arg2) + C("D=M") + insertDToStack() + incSP()
		} else {
			code = decSP() + getFromStack() + A(5+arg2) + C("M=D")
		}
	case SEG_CONSTANT:
		code = A(arg2) + C("D=A") + insertDToStack() + incSP()
	case SEG_STATIC:
		if commandType == C_PUSH {
			code = A(filename, ".", arg2) + C("D=M") + insertDToStack() + incSP()
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
			code = A(segmentSymbol) + C("D=M") + insertDToStack() + incSP()
		} else {
			code = decSP() + getFromStack() + A(segmentSymbol) + C("M=D")
		}
	}

	return code
}

func writeGoto(label string) string {
	return A(label) + C("0;JMP")
}

func writeIfGoto(label string) string {
	return decSP() + getFromStack() + A(label) + C("D;JNE")
}

func writeLabel(label string) string {
	return fmt.Sprintf("(%s)\n", label)
}

func writeCall(function string, nArgs int) string {
	return saveState() + // push retAddr, LCL, ARG, THIS, THAT
		A(5+nArgs) + C("D=A") + A("SP") + C("D=M-D") + A("ARG") + C("M=D") + // ARG = SP-5-nArgs
		A("SP") + C("D=M") + A("LCL") + C("M=D") + // LCL = SP
		writeGoto(function) + writeLabel(returnLabel())
}

func writeReturn() string {
	return A("LCL") + C("D=M") + A("R14") + C("M=D") + // endframe(R14) = LCL
		A(5) + C("A=D-A") + C("D=M") + A("R15") + C("M=D") + // retAddr = *(endframe - 5)
		writePushPop(C_POP, "", SEG_ARGUMENT, 0) + // *ARG = pop()
		A("ARG") + C("D=M+1") + A("SP") + C("M=D") + // SP = ARG+1
		A("R14") + C("A=M-1") + C("D=M") + A("THAT") + C("M=D") + // THAT = *(endframe - 1)
		A("R14") + C("D=M") + A(2) + C("A=D-A") + C("D=M") + A("THIS") + C("M=D") + // THIS = *(endframe - 2)
		A("R14") + C("D=M") + A(3) + C("A=D-A") + C("D=M") + A("ARG") + C("M=D") + // ARG = *(endframe - 3)
		A("R14") + C("D=M") + A(4) + C("A=D-A") + C("D=M") + A("LCL") + C("M=D") + // LCL = *(endframe - 4)
		A("R15") + C("A=M") + C("0;JMP") // goto retAddr
}

func writeFunction(name string, nLocal int) string {
	currentFunc = name
	code := writeLabel(name)
	for i := 0; i < nLocal; i++ {
		code += A("SP") + C("A=M") + C("M=0") + incSP()
	}

	return code
}

func returnLabel() string {
	return fmt.Sprintf("%s$ret.%d", currentFunc, returnCount)
}

// calling push constant "label" would make sense here, but I am lazy to rework that function
func saveState() string {
	returnCount++
	return A(returnLabel()) + C("D=A") + insertDToStack() + incSP() +
		A("LCL") + C("D=M") + insertDToStack() + incSP() +
		A("ARG") + C("D=M") + insertDToStack() + incSP() +
		A("THIS") + C("D=M") + insertDToStack() + incSP() +
		A("THAT") + C("D=M") + insertDToStack() + incSP()
}

func A(value ...any) string {
	value = append([]any{"@"}, value...)
	value = append(value, "\n")
	return fmt.Sprint(value...)
}

func C(str string) string {
	return str + "\n"
}

func insertDToStack() string {
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
	return decSP() + C("A=M") + C("D=M") + C("A=A-1") + C("D=M-D") + C("M=-1") + A(label) + C("D;"+operation) + A("SP") + C("A=M-1") + C("M=0") + writeLabel(label)
}
