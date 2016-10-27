package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/layneson/rowsofb/env"
	"github.com/layneson/rowsofb/matrix"
)

//CommandHandler represents a command handler.
//A handler takes an environment and a slice of arguments and executes, returning an error if
//the arguments are incorrect.
type CommandHandler func(e *env.E, args []string) error

var errIncorrectNumberOfArgs = errors.New("incorrect number of arguments")

var commands = map[string]CommandHandler{

	//Unary operations

	"inv": func(e *env.E, args []string) error {
		if len(args) != 1 {
			return errIncorrectNumberOfArgs
		}

		m, err := e.Get(args[0])
		if err != nil {
			return err
		}

		rm, err := matrix.Inverse(m)
		if err != nil {
			return err
		}

		e.SetResult(rm)

		return nil
	},

	"trans": func(e *env.E, args []string) error {
		if len(args) != 1 {
			return errIncorrectNumberOfArgs
		}

		m, err := e.Get(args[0])
		if err != nil {
			return err
		}

		rm := matrix.Transpose(m)

		e.SetResult(rm)

		return nil
	},

	"ref": func(e *env.E, args []string) error {
		if len(args) != 1 {
			return errIncorrectNumberOfArgs
		}

		m, err := e.Get(args[0])
		if err != nil {
			return err
		}

		rm := matrix.Ref(m)

		e.SetResult(rm)

		return nil
	},

	"rref": func(e *env.E, args []string) error {
		if len(args) != 1 {
			return errIncorrectNumberOfArgs
		}

		m, err := e.Get(args[0])
		if err != nil {
			return err
		}

		rm := matrix.Rref(m)

		e.SetResult(rm)

		return nil
	},

	//Binary operations

	"add": func(e *env.E, args []string) error {
		if len(args) != 2 {
			return errIncorrectNumberOfArgs
		}

		a, err := e.Get(args[0])
		if err != nil {
			return err
		}

		b, err := e.Get(args[1])
		if err != nil {
			return err
		}

		rm, err := matrix.Add(a, b)
		if err != nil {
			return err
		}

		e.SetResult(rm)

		return nil
	},

	"mul": func(e *env.E, args []string) error {
		if len(args) != 2 {
			return errIncorrectNumberOfArgs
		}

		a, err := e.Get(args[0])
		if err != nil {
			return err
		}

		b, err := e.Get(args[1])
		if err != nil {
			return err
		}

		rm, err := matrix.Multiply(a, b)
		if err != nil {
			return err
		}

		e.SetResult(rm)

		return nil
	},

	"scl": func(e *env.E, args []string) error {
		if len(args) != 2 {
			return errIncorrectNumberOfArgs
		}

		m, err := e.Get(args[1])
		if err != nil {
			return err
		}

		fparts := strings.Split(args[0], "/")

		if len(fparts) == 1 {
			fparts = append(fparts, "1")
		}

		n, err := strconv.Atoi(fparts[0])
		if err != nil {
			return errors.New("invalid fraction")
		}

		d, err := strconv.Atoi(fparts[1])
		if err != nil {
			return errors.New("invalid fraction")
		}

		s := matrix.NewFrac(n, d)

		rm := matrix.Scale(s, m)

		e.SetResult(rm)

		return nil
	},

	"aug": func(e *env.E, args []string) error {
		if len(args) != 2 {
			return errIncorrectNumberOfArgs
		}

		a, err := e.Get(args[0])
		if err != nil {
			return err
		}

		b, err := e.Get(args[1])
		if err != nil {
			return err
		}

		rm, err := matrix.Augment(a, b)
		if err != nil {
			return err
		}

		e.SetResult(rm)

		return nil
	},

	//Other commands

	"def": func(e *env.E, args []string) error {
		if len(args) != 1 {
			return errIncorrectNumberOfArgs
		}

		if _, err := e.IsDefined(args[0]); err != nil { // is the variable a real variable?
			return err
		}

		var mvals []matrix.Frac
		var r, c int

		fmt.Println()

		scanner := bufio.NewScanner(os.Stdin)

		matInputColor.Set()

		scanner.Scan()

		firstLine := scanner.Text()

		if firstLine == "" {
			return errors.New("matrix definition canceled")
		}

		fields := strings.Split(firstLine, "\t")

		c = len(fields)

		for _, s := range fields {
			f, err := matrix.ParseFrac(s)
			if err != nil {
				return err
			}

			mvals = append(mvals, f.Reduce())
		}

		for {
			//fmt.Print(":\t")

			scanner.Scan()
			line := scanner.Text()

			if line == "" {
				break
			}

			fields = strings.Split(line, "\t")

			if len(fields) != c {
				return errors.New("uneven columns")
			}

			for _, s := range fields {
				f, err := matrix.ParseFrac(s)
				if err != nil {
					return err
				}

				mvals = append(mvals, f)
			}
		}

		r = len(mvals) / c

		rm := matrix.NewWithValues(r, c, mvals)

		e.Set(args[0], rm)
		e.SetResult(rm)

		return nil
	},

	"print": func(e *env.E, args []string) error {
		if len(args) != 1 {
			return errIncorrectNumberOfArgs
		}

		m, err := e.Get(args[0])
		if err != nil {
			return err
		}

		e.SetResult(m)

		return nil
	},

	"set": func(e *env.E, args []string) error {
		if len(args) != 1 && len(args) != 2 {
			return errIncorrectNumberOfArgs
		}

		var m matrix.M
		if len(args) == 1 {
			m = e.GetResult()
		} else {
			var err error
			m, err = e.Get(args[1])
			if err != nil {
				return err
			}
		}

		return e.Set(args[0], m)
	},

	"zero": func(e *env.E, args []string) error {
		if len(args) != 1 {
			return errIncorrectNumberOfArgs
		}

		m, err := e.Get(args[0])
		if err != nil {
			return err
		}

		for r := 1; r <= m.Rows(); r++ {
			for c := 1; c <= m.Cols(); c++ {
				m.Set(r, c, matrix.NewScalarFrac(0))
			}
		}

		err = e.Set(args[0], m)
		if err != nil {
			return err
		}

		return nil
	},

	"del": func(e *env.E, args []string) error {
		if len(args) != 1 {
			return errIncorrectNumberOfArgs
		}

		return e.Delete(args[0])
	},

	"clr": func(e *env.E, args []string) error {
		if len(args) != 0 {
			return errIncorrectNumberOfArgs
		}

		e.Clear()

		return nil
	},

	"help": func(e *env.E, args []string) error {
		fmt.Println()

		for cmd, hlp := range commandHelp {
			printHelp(cmd, hlp)
		}
		return nil
	},
}

var commandHelp = map[string]string{
	"inv [A]":     "Calculates the inverse of 'A', or indicates that no such inverse exists.",
	"trans [A]":   "Calculates the transpose of 'A'.",
	"ref [A]":     "Uses Gauss-Jordan elimination to put 'A' into row echelon form.",
	"rref [A]":    "Uses Gauss-Jordan elimination to put 'A' into reduced row echelon form.",
	"add [A] [B]": "Adds matrix 'A' to 'B'. The matrices must have identical dimensions.",
	"mul [A] [B]": "Multiplies matrix 'A' by 'B'. The number of columns of A must equal the number of rows of B.",
	"scl c [A]":   "Multiplies matrix 'A' by the scalar 'c'. The scalar can be a fraction or an integer.",
	"aug [A] [B]": "Augments 'A' with 'B'. A must have the same number of rows as B.",
	"def [A]":     "Opens an interactive process to define the matrix 'A'.",
	"set [A] [B]": "Sets matrix 'A' to matrix 'B'. In essence, it copies 'B' into 'A'. B is optional; if B is not specified, the result is stored into A.",
	"zero [A]":    "Sets every element of matrix 'A' to 0.",
	"del [A]":     "Deletes matrix 'A'. A deleted matrix has no size and no entries, and any operation using a deleted matrix raises an error.",
	"clr":         "Deletes all matrices. This is equivalent to restarting RowsOfB.",
	"exit":        "Exits RowsOfB.",
}

func printHelp(cmd, desc string) {
	cfields := strings.Fields(cmd)
	color.New(color.FgBlue).Print(cfields[0])
	for _, ce := range cfields[1:] {
		color.New(color.FgHiBlue).Printf(" %s", ce)
	}
	color.New(color.FgMagenta).Printf("    -    %s\n", desc)
	fmt.Println()
}
