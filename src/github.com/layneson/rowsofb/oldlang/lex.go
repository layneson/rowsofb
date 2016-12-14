package oldlang

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	data  []rune
	index int
}

func NewScanner(data string) Scanner {
	runes := append([]rune(data), rune(0))

	return Scanner{runes, 0}
}

func (s Scanner) Peek() rune {
	if s.index == len(s.data) {
		return rune(0)
	}

	return s.data[s.index]
}

func (s *Scanner) Read() rune {
	if s.index == len(s.data) {
		s.index++
		return rune(0)
	}

	next := s.data[s.index]

	s.index++

	return next
}

func (s *Scanner) ReadStr(am int) string {
	str := ""
	for i := 0; i < am; i++ {
		str += string(s.Read())
	}
	return str
}

func (s *Scanner) Unread() {
	s.index--
}

type View struct {
	scan  *Scanner
	index int
}

func (s Scanner) View() View {
	return View{&s, s.index}
}

func (s *View) Peek() rune {
	return s.scan.Peek()
}

func (s *View) Inc() {
	s.index++

	_ = s.scan.Read()
}

func (s *View) Dec() {
	s.index--
	s.scan.Unread()
}

func (s *View) Reset() int {
	idx := s.index
	for i := idx; i > 0; i-- {
		s.scan.Unread()
	}
	s.index = 0
	return idx
}

func (s *View) ResetZero() int {
	_ = s.Reset()
	return 0
}

func (s *View) Skip(num int) {
	for i := 0; i < num; i++ {
		s.Inc()
	}
}

func (s *View) View() View {
	return s.scan.View()
}

func matchString(sv View, str string) int {
	for i, w := 0, 0; i < len(str); i += w {
		value, width := utf8.DecodeRuneInString(str[i:])

		if value != sv.Peek() {
			return sv.ResetZero()
		}

		sv.Inc()

		w = width
	}

	return sv.Reset()
}

func matchOne(in rune, runes ...rune) bool {
	for _, r := range runes {
		if r == in {
			return true
		}
	}
	return false
}

func matchOneIgnoreCase(in rune, runes ...rune) bool {
	for _, r := range runes {
		var otherCase rune
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				otherCase = unicode.ToLower(r)
			} else {
				otherCase = unicode.ToUpper(r)
			}

			if in == otherCase {
				return true
			}
		}

		if in == r {
			return true
		}
	}
	return false
}

func matchRange(in rune, ranges ...rune) bool {
	for i := 0; i < len(ranges); i += 2 {
		if in >= ranges[i] && in <= ranges[i+1] {
			return true
		}
	}
	return false
}

type Matcher func(View) int

// TT represents a Token Type.
type TT int

func getMatcher(tt TT) Matcher {
	switch tt {
	case EOF:
		return func(v View) int {
			if v.Peek() == rune(0) {
				v.Inc()
				return v.Reset()
			}

			return v.ResetZero()
		}

	case LPAREN:
		return func(v View) int {
			if v.Peek() == '(' {
				v.Inc()
				return v.Reset()
			}

			return v.ResetZero()
		}

	case RPAREN:
		return func(v View) int {
			if v.Peek() == ')' {
				v.Inc()
				return v.Reset()
			}

			return v.ResetZero()
		}

	case MULT:
		return func(v View) int {
			if v.Peek() == '*' {
				v.Inc()
				return v.Reset()
			}

			return v.ResetZero()
		}

	case DIV:
		return func(v View) int {
			if v.Peek() == '/' {
				v.Inc()
				return v.Reset()
			}

			return v.ResetZero()
		}

	case ADD:
		return func(v View) int {
			if v.Peek() == '+' {
				v.Inc()
				return v.Reset()
			}

			return v.ResetZero()
		}

	case NEG:
		return func(v View) int {
			if v.Peek() == '-' {
				v.Inc()
				return v.Reset()
			}

			return v.ResetZero()
		}

	case IDENT:
		return func(v View) int {
			if !matchRange(v.Peek(), 'a', 'z', 'A', 'Z') {
				return v.ResetZero()
			}

			v.Inc()

			if !matchRange(v.Peek(), 'a', 'z', 'A', 'Z') {
				return v.ResetZero()
			}

			v.Inc()

			for {
				if !matchRange(v.Peek(), 'a', 'z', 'A', 'Z') {
					return v.Reset()
				}

				v.Inc()
			}
		}

	case MVAR:
		return func(v View) int {
			if !matchRange(v.Peek(), 'A', 'Z') {
				return v.ResetZero()
			}

			v.Inc()

			return v.Reset()
		}

	case SVAR:
		return func(v View) int {
			if !matchRange(v.Peek(), 'a', 'z') {
				return v.ResetZero()
			}

			v.Inc()

			return v.Reset()
		}

	case DSIGN:
		return func(v View) int {
			if v.Peek() == '$' {
				v.Inc()
				return v.Reset()
			}

			return v.ResetZero()
		}

	case NUM:
		return func(v View) int {
			if !matchRange(v.Peek(), '1', '9') {
				return v.ResetZero()
			}

			v.Inc()

			for {
				if !matchRange(v.Peek(), '0', '9') {
					return v.Reset()
				}

				v.Inc()
			}
		}

	case WS:
		return func(v View) int {
			if v.Peek() != ' ' && v.Peek() != '\t' {
				return v.ResetZero()
			}

			v.Inc()

			for {
				if v.Peek() != ' ' && v.Peek() != '\t' {
					return v.Reset()
				}

				v.Inc()
			}

			return v.Reset()
		}
	}

	return nil
}

const (
	EOF TT = iota
	LPAREN
	RPAREN
	MULT
	DIV
	ADD
	NEG
	IDENT
	MVAR
	SVAR
	DSIGN
	NUM
	WS
)

var ttypes = []TT{
	EOF,
	LPAREN, RPAREN,
	MULT, DIV, ADD, NEG,
	IDENT,
	MVAR, SVAR,
	DSIGN,
	NUM,
	WS,
}

type Token struct {
	ttype TT
	value string
}

func (s *Scanner) Lex() ([]Token, error) {
	tokens := []Token{}

outer:
	for {
		found := false
		for _, tt := range ttypes {
			fmt.Printf("Getting matcher! Next rune %#v\n", s.Peek())
			matcher := getMatcher(tt)
			am := matcher(s.View())
			if am > 0 {
				found = true

				if tt != WS {
					tokens = append(tokens, Token{tt, s.ReadStr(am)})
				}

				if tt == EOF {
					break outer
				}
			}
		}

		if !found {
			return tokens, fmt.Errorf("unrecognized input at character %d", s.index+1)
		}
	}

	return tokens, nil
}

func Parse(line string) ([]Token, error) {
	sc := NewScanner(line)

	return sc.Lex()
}
