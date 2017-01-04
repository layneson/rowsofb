package cli

import (
	"bufio"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/layneson/rowsofb/env"
	"github.com/layneson/rowsofb/lang"
	"github.com/layneson/rowsofb/matrix"
)

var errorColor = color.New(color.FgRed)
var promptColor = color.New(color.FgGreen)
var cmdInputColor = color.New(color.FgCyan)
var matInputColor = color.New(color.FgHiBlue)
var resultColor = color.New(color.FgMagenta)

var scan = bufio.NewScanner(os.Stdin)

// Run starts the CLI loop.
func Run() {
	e := env.New(defineMatrix, defineAnonymousMatrix, defineScalar)

	for {
		promptColor.Print("> ")

		cmdInputColor.Set()

		scan.Scan()
		line := scan.Text()

		toks, err := lang.Lex(line)
		if err != nil {
			reportError("Lexing Error: ", err)
			continue
		}

		enode, err := lang.Parse(toks)
		if err != nil {
			reportError("Parsing Error: ", err)
			continue
		}

		val, err := env.Evaluate(enode, e)
		if err != nil {
			reportError("Evaluation Error: ", err)
			continue
		}

		switch val.VType {
		case env.MVar:
			e.SetMVar('Z', val.MValue)
			resultColor.Println(renderMatrix(val.MValue))
		case env.SVar:
			e.SetSVar('z', val.SValue)
			resultColor.Println(val.SValue)
		}
	}
}

func reportError(prefix string, e error) {
	errorColor.Printf("[!] %s%v.\n", prefix, e)
}

func reportErrorMsg(message string) {
	errorColor.Printf("[!] %s.\n", message)
}

func defineMatrix(v rune) (matrix.M, bool) {
	promptColor.Printf("Define matrix %c:\n", v)
	return defineMatrixAgnostic()
}

func defineAnonymousMatrix() (matrix.M, bool) {
	promptColor.Println("Define anonymous matrix:")
	return defineMatrixAgnostic()
}

func defineScalar(v rune) (matrix.Frac, bool) {
	promptColor.Printf("Define scalar %c: ", v)
	matInputColor.Set()

	scan.Scan()
	input := strings.TrimSpace(scan.Text())

	frac, err := matrix.ParseFrac(input)
	if err != nil {
		reportError("Failed to parse scalar input: ", err)
		return matrix.Frac{}, false
	}

	return frac, true
}

func defineMatrixAgnostic() (matrix.M, bool) {
	matInputColor.Set()

	scan.Scan()
	first := strings.TrimSpace(scan.Text())

	if first == "" {
		return matrix.M{}, false
	}

	firstFields := strings.Split(first, "\t")

	c := len(firstFields)

	values := []matrix.Frac{}

	firstFracs, err := parseMatrixRow(firstFields)
	if err != nil {
		reportError("Could not define matrix: ", err)
		return matrix.M{}, false
	}

	for _, f := range firstFracs {
		values = append(values, f)
	}

	for {
		scan.Scan()
		line := strings.TrimSpace(scan.Text())

		if line == "" {
			break
		}

		fields := strings.Split(line, "\t")

		if len(fields) != c {
			reportErrorMsg("Matrix has uneven rows")
			return matrix.M{}, false
		}

		fracs, err := parseMatrixRow(fields)
		if err != nil {
			reportError("Could not define matrix: ", err)
			return matrix.M{}, false
		}

		for _, f := range fracs {
			values = append(values, f)
		}
	}

	return matrix.NewWithValues(len(values)/c, c, values), true
}

func parseMatrixRow(fields []string) ([]matrix.Frac, error) {
	fracs := []matrix.Frac{}

	for _, f := range fields {
		frac, err := matrix.ParseFrac(f)
		if err != nil {
			return fracs, err
		}

		fracs = append(fracs, frac)
	}

	return fracs, nil
}
