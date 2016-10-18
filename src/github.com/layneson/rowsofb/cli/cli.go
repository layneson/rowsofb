package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/layneson/rowsofb/env"
)

//Run starts the CLI loop.
func Run() {
	environment := env.New()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		scanner.Scan()

		line := scanner.Text()

		if line == "" {
			continue
		}

		fields := strings.Fields(line)

		if fields[0] == "exit" {
			break
		}

		cmd, ok := commands[fields[0]]
		if !ok {
			fmt.Println("[!] That command does not exist.")
			continue
		}

		err := cmd(environment, fields[1:])
		if err != nil {
			fmt.Printf("[!] %v.\n", err)
			continue
		}

		fmt.Println(environment.GetResult())
	}
}
