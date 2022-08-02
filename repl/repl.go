package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"monkey_lang/evaluator"
	"monkey_lang/lexer"
	"monkey_lang/object"
	"monkey_lang/parser"
)

// PROMPT string is used to let the user know entry point for REPL
const PROMPT = ">> "

func printMonkeyFace() {
	// MonkeyFace is an ASCII image which is displayed when an error occurs
	const MonkeyFace = `            __,__
	
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

	lines := strings.Split(MonkeyFace, "\n")
	for idx, line := range lines {
		if idx != len(lines)-1 {
			fmt.Println(string("\033[31m"), line)
		} else {
			fmt.Println(string("\033[31m") + line + string("\033[0m"))
		}
	}
}

// Start function is used to start the REPL session
func Start(in io.Reader, out io.Writer) {
	// io.WriteString(out, MONKEY_FACE)
	// printMonkeyFace()
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	// io.WriteString(out, MonkeyFace)
	printMonkeyFace()
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, fmt.Sprintf("\t%s\n", msg))
	}
}
