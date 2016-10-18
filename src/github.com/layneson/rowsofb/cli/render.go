package cli

import (
	"bytes"

	"github.com/layneson/rowsofb/matrix"
)

func renderMatrix(m matrix.M) string {
	smat := make([]string, m.Rows()*m.Cols())
	cwidths := make([]int, m.Cols())
	csum := 0

	for c := 1; c <= m.Cols(); c++ {
		mwidth := 0
		for r := 1; r <= m.Rows(); r++ {
			str := m.Get(r, c).String()
			smat[(r-1)*m.Cols()+(c-1)] = str
			if len(str) > mwidth {
				mwidth = len(str)
			}
		}
		cwidths[c-1] = mwidth
		csum += mwidth
	}

	nspace := csum + 4*(m.Cols()-1)

	var buff bytes.Buffer

	buff.WriteString("┌ ")
	for i := 0; i < nspace; i++ {
		buff.WriteString(" ")
	}
	buff.WriteString(" ┐\n")

	for r := 1; r <= m.Rows(); r++ {
		buff.WriteString("│ ")

		for c := 1; c <= m.Cols()-1; c++ {
			str := padRight(smat[(r-1)*m.Cols()+(c-1)], cwidths[c-1])
			buff.WriteString(str)
			buff.WriteString("    ") // 4 spaces
		}

		buff.WriteString(padRight(smat[(r-1)*m.Cols()+(m.Cols()-1)], cwidths[m.Cols()-1]))

		buff.WriteString(" │\n")
	}

	buff.WriteString("└ ")
	for i := 0; i < nspace; i++ {
		buff.WriteString(" ")
	}
	buff.WriteString(" ┘")

	return buff.String()
}

//padRight pads the string s with spaces on the right until s has length l.
func padRight(s string, l int) string {
	for len(s) < l {
		s = s + " "
	}

	return s
}
