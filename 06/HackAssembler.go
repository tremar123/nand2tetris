package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var freeMemory = 16

var symbolMap map[string]int = map[string]int{
	"SP":   0,
	"LCL":  1,
	"ARG":  2,
	"THIS": 3,
	"THAT": 4,

	"R0":  0,
	"R1":  1,
	"R2":  2,
	"R3":  3,
	"R4":  4,
	"R5":  5,
	"R6":  6,
	"R7":  7,
	"R8":  8,
	"R9":  9,
	"R10": 10,
	"R11": 11,
	"R12": 12,
	"R13": 13,
	"R14": 14,
	"R15": 15,

	"SCREEN": 16384,
	"KBD":    24576,
}

var compMap map[string]string = map[string]string{
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"D+A": "0000010",
	"D-A": "0010011",
	"A-D": "0000111",
	"D&A": "0000000",
	"D|A": "0010101",

	"M":   "1110000",
	"!M":  "1110001",
	"-M":  "1110011",
	"M+1": "1110111",
	"M-1": "1110010",
	"D+M": "1000010",
	"D-M": "1010011",
	"M-D": "1000111",
	"D&M": "1000000",
	"D|M": "1010101",
}

var destMap map[string]string = map[string]string{
	"":    "000",
	"M":   "001",
	"D":   "010",
	"MD":  "011",
	"A":   "100",
	"AM":  "101",
	"AD":  "110",
	"AMD": "111",
}

var jumpMap map[string]string = map[string]string{
	"":    "000",
	"JGT": "001",
	"JEQ": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Argument expected: file name")
	}

	asm, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
	}

	lines := strings.Split(string(asm), "\n")
	lines = lines[:len(lines)-1]

	// find all the labels and add them to symbol map
	instructionIndex := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || line[:2] == "//" {
			continue
		}

		if len(line) > 0 && line[0] == '(' {
			label := strings.TrimSuffix(line[1:], ")")
			symbolMap[label] = instructionIndex
		} else {
			instructionIndex++
		}
	}

	machineCode := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || line[0] == '(' || line[:2] == "//" { // skip whitespace and labels
			continue
		} else if line[0] == '@' { // A
			instruction := handleAInstruction(line)
			machineCode += instruction
		} else { // C
			instruction := handleCInstruction(line)
			machineCode += instruction
		}
	}

	newFileName := strings.Split(os.Args[1], ".")[0] + ".hack"

	err = os.WriteFile(newFileName, []byte(machineCode), 0600)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func handleAInstruction(line string) string {
	value := strings.TrimPrefix(line, "@")
	value = strings.Split(value, "//")[0] // if there's a comment at the end of the line delete it
	value = strings.TrimSpace(value)
	number, err := strconv.Atoi(value)
	if err != nil {
		symbolValue, exist := symbolMap[value]
		if exist {
			number = symbolValue
		} else {
			symbolMap[value] = freeMemory
			number = freeMemory
			freeMemory++
		}
	}

	return fmt.Sprintf("0%015b\n", number)
}

func handleCInstruction(line string) string {
	command := strings.Split(line, "//")[0] // get rid of the comment
	command = strings.TrimSpace(command)

	comp := parseComp(line)
	dest := parseDest(line)
	jump := parseJump(line)

	compBinary := translateComp(comp)
	destBinary := translateDest(dest)
	jumpBinary := translateJump(jump)

	return fmt.Sprintf("111%s%s%s\n", compBinary, destBinary, jumpBinary)
}

// initially used spliting strings but this seems to improve execution time by 0.1s on pong.asm :)
func parseComp(line string) string {
	startIndex := 0
	endIndex := 0
	for i, c := range line {
		endIndex = i
		if c == '=' {
			startIndex = i + 1
		} else if c == ';' {
			endIndex--
			break
		}
	}

	comp := line[startIndex : endIndex+1]
	return strings.TrimSpace(comp)
}

func parseDest(line string) string {
	for i, c := range line {
		if c == '=' {
			return strings.TrimSpace(line[0:i])
		}
	}

	return ""
}

func parseJump(line string) string {
	for i := len(line) - 1; i > 0; i-- {
		if line[i] == ';' {
			jump := line[i+1:]
			return strings.TrimSpace(jump)
		}
	}
	return ""
}

func translateComp(comp string) string {
	return compMap[comp]
}

func translateDest(dest string) string {
	return destMap[dest]
}

func translateJump(jump string) string {
	return jumpMap[jump]
}
