package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/layneson/rowsofb/env"
)

var errorColor = color.New(color.FgRed)
var promptColor = color.New(color.FgGreen)
var cmdInputColor = color.New(color.FgCyan)
var matInputColor = color.New(color.FgHiBlue)
var resultColor = color.New(color.FgMagenta)

//Run starts the CLI loop.
func Run() {
	environment := env.New()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		promptColor.Print("> ")
		cmdInputColor.Set()

		scanner.Scan()

		line := scanner.Text()

		if line == "" {
			resultColor.Set()

			fmt.Println()
			fmt.Println(renderMatrix(environment.GetResult()))
			fmt.Println()

			continue
		}

		fields := strings.Fields(line)

		if fields[0] == "exit" {
			break
		}

		lastEntry := fields[len(fields)-1]
		hasRedirect := false
		if strings.HasPrefix(lastEntry, ">") && len(fields) > 1 {
			if _, err := environment.IsDefined(lastEntry[1:]); err != nil {
				errorColor.Printf("[!] %v.\n", err)
				continue
			}
			hasRedirect = true
		}

		cmd, ok := commands[fields[0]]
		if !ok {
			errorColor.Println("[!] That command does not exist.")
			continue
		}

		if hasRedirect {
			err := cmd(environment, fields[1:len(fields)-1])
			if err != nil {
				errorColor.Printf("[!] %v.\n", err)
				continue
			}

			environment.Set(lastEntry[1:], environment.GetResult())
		} else {
			err := cmd(environment, fields[1:])
			if err != nil {
				errorColor.Printf("[!] %v.\n", err)
				continue
			}
		}

		if fields[0] != "help" {
			fmt.Println()

			resultColor.Println(renderMatrix(environment.GetResult()))
			fmt.Println()
		}
	}
}
