package main

import (
	"os"
	"strconv"
)

const PUSH = "@SP\n" +
	"A=M\n" +
	"M=D\n" +
	"@SP\n" +
	"M=M+1\n"

const POP = "@R13\n" +
	"M=D\n" +
	"@SP\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"@R13\n" +
	"A=M\n" +
	"M=D\n"

const ADD = "@SP\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"A=A-1\n" +
	"M=D+M\n"

const SUB = "@SP\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"A=A-1\n" +
	"M=M-D\n"

const NEG = "@SP\n" +
	"A=M-1\n" +
	"M=-M\n"

const AND = "@SP\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"A=A-1\n" +
	"M=D&M\n"

const OR = "@SP\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"A=A-1\n" +
	"M=D|M\n"

const NOT = "@SP\n" +
	"A=M-1\n" +
	"M=!M\n"

const RETURN = "@LCL\n" +
	"D=M\n" +
	"@5\n" +
	"A=D-A\n" +
	"D=M\n" +
	"@R13\n" +
	"M=D\n" +
	"@SP\n" +
	"A=M-1\n" +
	"D=M\n" +
	"@ARG\n" +
	"A=M\n" +
	"M=D \n" +
	"D=A+1\n" +
	"@SP\n" +
	"M=D\n" +
	"@LCL\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"@THAT\n" +
	"M=D\n" +
	"@LCL\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"@THIS\n" +
	"M=D\n" +
	"@LCL\n" +
	"AM=M-1\n" +
	"D=M\n" +
	"@ARG\n" +
	"M=D\n" +
	"@LCL\n" +
	"A=M-1\n" +
	"D=M\n" +
	"@LCL\n" +
	"M=D\n" +
	"@R13\n" +
	"A=M\n" +
	"0;JMP\n"

var counter int

// Initializes counter
func init() {
	counter = 0
}

type CodeWriter struct {
	outputFile      *os.File
	currentFileName string
}

func CreateWriter() *CodeWriter {
	w := &CodeWriter{}
	return w
}

// Creates a new output file from the input file
func (w *CodeWriter) setFileName(fileName string) {
	outputFile, err := os.Create(fileName + ".asm") // Creates the .asm output file

	//checks if there was an error while creating the output file
	if err != nil {
		panic(err)
	}
	w.outputFile = outputFile
}

// Returns the next counter for defining a new symbol
func GetNextSymbolCount() string {
	counter++
	return strconv.Itoa(counter)
}

// Writes to the output file the assembly code that implements the given arithmetric-logical command
func (w *CodeWriter) WriteArithmetic(command string) {
	asm := ""
	switch command {
	case "add":
		{
			asm = ADD
		}
	case "sub":
		{
			asm = SUB
		}
	case "neg":
		{
			asm = NEG
		}
	case "and":
		{
			asm = AND
		}
	case "or":
		{
			asm = OR
		}
	case "not":
		{
			asm = NOT
		}
	case "eq":
		{
			asm = EQ()
		}
	case "gt":
		{
			asm = GT()
		}
	case "lt":
		{
			asm = LT()
		}
	}
	w.outputFile.WriteString(asm)
}

// Writes to the output file the assembly code that implements the given push or pop command
func (w *CodeWriter) WritePushPop(command int, segment string, index int) {
	asm := ""
	strIndex := strconv.Itoa(index)
	switch command {
	case C_PUSH:
		{
			switch segment {
			case "constant":
				{
					asm = "@" + strIndex + "\n" + "D=A\n" + PUSH
				}
			case "local":
				{
					asm = "@LCL\n" +
						"D=M\n" +
						"@" + strIndex + "\n" +
						"A=D+A\n" +
						"D=M\n" + PUSH
				}
			case "this":
				{
					asm = "@THIS\n" +
						"D=M\n" +
						"@" + strIndex + "\n" +
						"A=D+A\n" +
						"D=M\n" +
						PUSH
				}
			case "that":
				{
					asm = "@THAT\n" +
						"D=M\n" +
						"@" + strIndex + "\n" +
						"A=D+A\n" +
						"D=M\n" +
						PUSH
				}
			case "argument":
				{
					asm = "@ARG\n" +
						"D=M\n" +
						"@" + strIndex + "\n" +
						"A=D+A\n" +
						"D=M\n" +
						PUSH
				}
			case "temp":
				{
					asm = "@R5\n" +
						"D=A\n" +
						"@" + strIndex + "\n" +
						"A=D+A\n" +
						"D=M\n" +
						PUSH
				}
			case "pointer":
				{
					if index == 0 {
						asm = "@THIS\n" +
							"D=M\n" +
							PUSH
					} else {
						asm = "@THAT\n" +
							"D=M\n" +
							PUSH
					}
				}
			case "static":
				{
					asm = "@" + w.currentFileName + "." + strIndex + "\n" +
						"D=M\n" +
						PUSH
				}
			}
		}
	case C_POP:
		{
			switch segment {
			case "local":
				{
					asm = "@LCL\n" +
						"D=M\n" +
						"@" + strIndex + "\n" +
						"D=D+A\n" +
						POP
				}
			case "this":
				{
					asm = "@THIS\n" +
						"D=M\n" +
						"@" + strIndex + "\n" +
						"D=D+A\n" +
						POP
				}
			case "that":
				{
					asm = "@THAT\n" +
						"D=M\n" +
						"@" + strIndex + "\n" +
						"D=D+A\n" +
						POP
				}
			case "argument":
				{
					asm = "@ARG\n" +
						"D=M\n" +
						"@" + strIndex + "\n" +
						"D=D+A\n" +
						POP
				}
			case "pointer":
				{
					if index == 0 {
						asm = "@THIS\n" +
							"D=A\n" +
							POP
					} else {
						asm = "@THAT\n" +
							"D=A\n" +
							POP
					}
				}
			case "static":
				{
					asm = "@" + w.currentFileName + "." + strIndex + "\n" +
						"D=A\n" +
						POP
				}
			case "temp":
				{
					asm = "@R5\n" +
						"D=A\n" +
						"@" + strIndex + "\n" +
						"D=D+A\n" +
						POP
				}
			}
		}
	}
	w.outputFile.WriteString(asm)
}

// Writes to the output file the assembly code that effects the label command
func (w *CodeWriter) WriteLabel(label string) {
	asm := LABEL(label)
	w.outputFile.WriteString(asm)
}

// Writes to the output file the assembly code that effects the goto command
func (w *CodeWriter) WriteGoTo(label string) {
	asm := GOTO(label)
	w.outputFile.WriteString(asm)
}

// Writes to the output file the assembly code that effects the if-goto command
func (w *CodeWriter) WriteIf(label string) {
	asm := IFGOTO(label)
	w.outputFile.WriteString(asm)
}

// Writes to the output file the assembly code that effects the function command
func (w *CodeWriter) WriteFunction(functionName string, nArgs int) {
	asm := FUNCTION(functionName, nArgs)
	w.outputFile.WriteString(asm)
}

// Writes to the output file the assembly code that effects the call command
func (w *CodeWriter) WriteCall(functionName string, nArgs int) {
	asm := CALL(functionName, nArgs)
	w.outputFile.WriteString(asm)
}

// Writes to the output file the assembly code that effects the return command
func (w *CodeWriter) WriteReturn() {
	w.outputFile.WriteString(RETURN)
}

// Writes to the output file the Sys.init assembly code
func (w *CodeWriter) Write_Sys_Init() {
	asm := "@256\n" +
		"D=A\n" +
		"@SP\n" +
		"M=D\n" +
		CALL("Sys.init", 0) +
		"0;JMP\n"
	w.outputFile.WriteString(asm)
}

// Returns the assembly code for the eq command
func EQ() string {
	i := GetNextSymbolCount()
	asm :=
		"@SP\n" +
			"AM=M-1\n" +
			"D=M\n" +
			"A=A-1\n" +
			"D=M-D\n" +
			"@EQ.a" + i + "\n" +
			"D;JEQ\n" +
			"@SP\n" +
			"A=M-1\n" +
			"M=0\n" +
			"@EQ.b" + i + "\n" +
			"0;JMP\n" +
			"(EQ.a" + i + ")\n" +
			"@SP\n" +
			"A=M-1\n" +
			"M=-1\n" +
			"(EQ.b" + i + ")\n"
	return asm
}

// Returns the assembly code for the gt command
func GT() string {
	i := GetNextSymbolCount()
	asm :=
		"@SP\n" +
			"AM=M-1\n" +
			"D=M\n" +
			"A=A-1\n" +
			"D=M-D\n" +
			"@GT.a" + i + "\n" +
			"D;JGT\n" +
			"@SP\n" +
			"A=M-1\n" +
			"M=0\n" +
			"@GT.b" + i + "\n" +
			"0;JMP\n" +
			"(GT.a" + i + ")\n" +
			"@SP\n" +
			"A=M-1\n" +
			"M=-1\n" +
			"(GT.b" + i + ")\n"
	return asm
}

// Returns the assembly code for the lt command
func LT() string {
	i := GetNextSymbolCount()
	asm :=
		"@SP\n" +
			"AM=M-1\n" +
			"D=M\n" +
			"A=A-1\n" +
			"D=M-D\n" +
			"@LT.a" + i + "\n" +
			"D;JLT\n" +
			"@SP\n" +
			"A=M-1\n" +
			"M=0\n" +
			"@LT.b" + i + "\n" +
			"0;JMP\n" +
			"(LT.a" + i + ")\n" +
			"@SP\n" +
			"A=M-1\n" +
			"M=-1\n" +
			"(LT.b" + i + ")\n"
	return asm
}

// Returns the assembly code for the label command
func LABEL(label string) string {
	asm := "(" + funcName + "$" + label + ")\n"
	return asm
}

// Returns the assembly code for the goto command
func GOTO(label string) string {
	asm :=
		"@" + funcName + "$" + label + "\n" +
			"0;JMP\n"
	return asm
}

// Returns the assembly code for the if-goto command
func IFGOTO(label string) string {
	asm :=
		"@SP\n" +
			"AM=M-1\n" +
			"D=M\n" +
			"@" + funcName + "$" + label + "\n" +
			"D;JNE\n"
	return asm
}

// Returns the assembly code for the function command
func FUNCTION(funcName string, nArgs int) string {
	asm :=
		"(" + funcName + ")\n" +
			"@SP\n" +
			"A=M\n"
	for i := 0; i < nArgs; i++ {
		asm += "M=0\n" +
			"A=A+1\n"
	}
	asm += "D=A\n" +
		"@SP\n" +
		"M=D\n"
	return asm
}

// Returns the assembly code for the call command
func CALL(funcName string, nArgs int) string {
	i := GetNextSymbolCount()
	n := strconv.Itoa(nArgs)
	asm :=
		"@SP\n" +
			"D=M\n" +
			"@R13\n" +
			"M=D\n" +
			"@RET." + i + "\n" +
			"D=A\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n" +
			"@LCL\n" +
			"D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n" +
			"@ARG\n" +
			"D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n" +
			"@THIS\n" +
			"D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n" +
			"@THAT\n" +
			"D=M\n" +
			"@SP\n" +
			"A=M\n" +
			"M=D\n" +
			"@SP\n" +
			"M=M+1\n" +
			"@R13\n" +
			"D=M\n" +
			"@" + n + "\n" +
			"D=D-A\n" +
			"@ARG\n" +
			"M=D\n" +
			"@SP\n" +
			"D=M\n" +
			"@LCL\n" +
			"M=D\n" +
			"@" + funcName + "\n" +
			"0;JMP\n" +
			"(RET." + i + ")\n"

	return asm
}

// Closes the output file
func (w *CodeWriter) Close() {
	err := w.outputFile.Close()
	checkWriterErr(err)
}

// Terminates the program if there was an error while closing the output file
func checkWriterErr(e error) {
	if e != nil {
		panic(e)
	}
}
