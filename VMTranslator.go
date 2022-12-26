package main

import (
	"fmt"
	"os"
	"strings"
)

var parser *Parser
var writer *CodeWriter

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No provided file or directory")
		os.Exit(1)
	}

	writer = CreateWriter()

	fileOrDir := os.Args[1]

	// This returns an *os.FileInfo type
	info, err := os.Stat(fileOrDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if info.IsDir() { // is directory

		files, err := os.ReadDir(fileOrDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		writer.setFileName(fileOrDir + "/" + info.Name()) // sets the output file name to be dirName.asm
		writer.Write_Sys_Init()                           // writes the bootstrap code to the output file

		// for each .vm file we generate its assembly code to the output file
		for _, file := range files {
			fileStructure := strings.Split(file.Name(), ".")
			if fileStructure[1] == "vm" {
				writer.currentFileName = fileStructure[0]
				generateAsmOutputFromFile(fileOrDir + "/" + file.Name())
			}
		}

	} else { // is file
		outputFile := strings.Split(fileOrDir, ".vm")[0]
		writer.currentFileName = outputFile
		writer.setFileName(outputFile)
		generateAsmOutputFromFile(fileOrDir)
	}
	writer.Close()
}

// Generates the asmmebly code from the .vm inputFile and writes it to the .asm output file
func generateAsmOutputFromFile(inputFile string) {

	parser = CreateParser(inputFile)

	for parser.Advance() {
		switch parser.currentCommandType { // writes the assembly code according the the command
		case C_ARITHMETIC:
			{
				writer.WriteArithmetic(parser.Arg1())
			}
		case C_PUSH, C_POP:
			{
				writer.WritePushPop(parser.currentCommandType, parser.Arg1(), parser.Arg2())
			}
		case C_LABEL:
			{
				writer.WriteLabel(parser.Arg1())
			}
		case C_GOTO:
			{
				writer.WriteGoTo(parser.Arg1())
			}
		case C_IF:
			{
				writer.WriteIf(parser.Arg1())
			}
		case C_CALL:
			{
				writer.WriteCall(parser.Arg1(), parser.Arg2())
			}
		case C_FUNCTION:
			{
				writer.WriteFunction(funcName, parser.Arg2())
			}
		case C_RETURN:
			{
				writer.WriteReturn()
			}
		}
	}
}
