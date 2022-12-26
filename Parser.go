package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const C_ARITHMETIC = 0
const C_PUSH = 1
const C_POP = 2
const C_LABEL = 3
const C_GOTO = 4
const C_IF = 5
const C_FUNCTION = 6
const C_RETURN = 7
const C_CALL = 8

const ARITHMETIC_LOGICAL = "add sub neg eq gt lt and or not"

var commandTypesMap map[string]int
var funcName string

// Initializes commandTypesMap
func init() {
	commandTypesMap = make(map[string]int)
	commandTypesMap["arithmetic"] = C_ARITHMETIC
	commandTypesMap["push"] = C_PUSH
	commandTypesMap["pop"] = C_POP
	commandTypesMap["label"] = C_LABEL
	commandTypesMap["goto"] = C_GOTO
	commandTypesMap["if-goto"] = C_IF
	commandTypesMap["function"] = C_FUNCTION
	commandTypesMap["return"] = C_RETURN
	commandTypesMap["call"] = C_CALL
}

type Parser struct {
	scanner            *bufio.Scanner
	currentCommandType int
	arg1               string
	arg2               int
}

// Opens the input file and gets ready to parse it
func CreateParser(fileName string) *Parser {
	p := &Parser{}
	p.setInputFile(fileName)
	return p
}

// Initializes the scanner with the input file
func (p *Parser) setInputFile(fileName string) {
	f, err := os.Open(fileName)
	checkParserErr(err)
	scanner := bufio.NewScanner(f)
	p.scanner = scanner
}

// Returns true if the input file has more lines to read
func (p *Parser) HasMoreLines() bool {
	return p.scanner.Scan()
}

// Reads the next command from the input file and makes it the current command
func (p *Parser) Advance() bool {
	for p.HasMoreLines() {
		line := strings.TrimSpace(p.scanner.Text())
		if len(p.scanner.Text()) != 0 && !strings.HasPrefix(line, "/") { // skip lines that are empty or comments
			line = strings.Split(line, "/")[0]     // Remove comments from the line
			line = strings.TrimSpace(line)         //Remove remaining white spaces
			command := strings.Split(line, " ")[0] // Gets command
			if command == "return" {
				p.currentCommandType = commandTypesMap["return"]
			} else if strings.Contains(ARITHMETIC_LOGICAL, command) {
				p.currentCommandType = commandTypesMap["arithmetic"]
				p.arg1 = command
			} else {
				p.currentCommandType = commandTypesMap[command]
				p.arg1 = strings.Split(line, " ")[1]
				if command == "push" || command == "pop" || command == "function" || command == "call" {
					arg2, err := strconv.Atoi(strings.Split(line, " ")[2]) //Gets arg2 and converts it into int
					checkParserErr(err)
					p.arg2 = arg2
					if command == "function" {
						funcName = p.arg1
					}
				}
			}
			return true
		}
	}
	return false
}

// Returns a constact representing the type of the current command
func (p *Parser) CommandType() int {
	return p.currentCommandType
}

// Returns the first argument of the current command
func (p *Parser) Arg1() string {
	return p.arg1
}

// Returns the second argument of the current command
func (p *Parser) Arg2() int {
	return p.arg2
}

// Terminates the program if there was an error while opening the input file
func checkParserErr(e error) {
	if e != nil {
		panic(e)
	}
}
