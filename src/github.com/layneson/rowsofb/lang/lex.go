package lang

import (
	"fmt"
)

type lexer struct {
	data        []rune
	next, ipeek int
}

func newLexer(data string) *lexer {
	return &lexer{data: []rune(data), next: 0, ipeek: 0}
}

func (lex *lexer) read() rune {
	r := lex.data[lex.next]
	lex.next++
	return r
}

func (lex *lexer) peek() rune {
	if lex.ipeek == len(lex.data) {
		return rune(0)
	}

	return lex.data[lex.ipeek]
}

func (lex *lexer) peekInc() rune {
	r := lex.peek()

	if r != rune(0) {
		lex.ipeek++
	}

	return r
}

func (lex *lexer) peekReset() {
	lex.ipeek = lex.next
}

func (lex *lexer) canRead() bool {
	return lex.next < len(lex.data)
}

// consume returns a token with a literal of the characters from data[next:ipeek]. lex.next is then set to lex.ipeek.
func (lex *lexer) consume(ttype TokenType) *Token {
	tok := &Token{TType: ttype, Literal: string(lex.data[lex.next:lex.ipeek])}
	lex.next = lex.ipeek
	return tok
}

func (lex *lexer) consumeIgnore() {
	lex.next = lex.ipeek
}

// A Token represents a single token of input.
type Token struct {
	TType   TokenType
	Literal string
}

// A TokenType represents a type of token.
type TokenType int

// tokentype definitions
const (
	TTEOF TokenType = iota

	// operators
	TTPlus
	TTMinus
	TTMult
	TTDiv

	TTArrow
	TTComma

	// parenthesis
	TTLParen
	TTRParen

	// literals
	TTNum
	TTFunc

	// variables
	TTMVar   // matrix variable
	TTSVar   // scalar variable
	TTDMVar  // define matrix variable
	TTDSVar  // define scalar variable
	TTDAMVar // define anonymous matrix variable
)

func (tt TokenType) String() string {
	switch tt {
	case TTEOF:
		return "EOF"
	case TTPlus:
		return "plus"
	case TTMinus:
		return "minus"
	case TTMult:
		return "mult"
	case TTDiv:
		return "div"
	case TTArrow:
		return "arrow"
	case TTComma:
		return "comma"
	case TTLParen:
		return "lparen"
	case TTRParen:
		return "rparen"
	case TTNum:
		return "num"
	case TTFunc:
		return "func"
	case TTMVar:
		return "mvar"
	case TTSVar:
		return "svar"
	case TTDMVar:
		return "dmvar"
	case TTDSVar:
		return "dsvar"
	case TTDAMVar:
		return "damvar"
	}

	return "unknown"
}

// Lex takes a line of text and parses it into tokens.
func Lex(data string) ([]*Token, error) {
	toks := []*Token{}

	lex := newLexer(data)

	for lex.canRead() {
		if runeMatchWhitespace(lex.peek()) {
			lex.peekInc()
			lex.consumeIgnore()
			continue
		}

		switch lex.peek() {
		case '+':
			lex.peekInc()
			toks = append(toks, lex.consume(TTPlus))
			continue
		case '-':
			lex.peekInc()
			if lex.peek() == '>' {
				lex.peekInc()
				toks = append(toks, lex.consume(TTArrow))
			} else {
				toks = append(toks, lex.consume(TTMinus))
			}
			continue
		case '*':
			lex.peekInc()
			toks = append(toks, lex.consume(TTMult))
			continue
		case '/':
			lex.peekInc()
			toks = append(toks, lex.consume(TTDiv))
			continue
		case '(':
			lex.peekInc()
			toks = append(toks, lex.consume(TTLParen))
			continue
		case ')':
			lex.peekInc()
			toks = append(toks, lex.consume(TTRParen))
			continue
		case ',':
			lex.peekInc()
			toks = append(toks, lex.consume(TTComma))
			continue
		}

		if matchNumber(lex) {
			toks = append(toks, lex.consume(TTNum))
			continue
		}

		if matchFunction(lex) {
			toks = append(toks, lex.consume(TTFunc))
			continue
		} else {
			lex.peekReset()
		}

		if lex.peek() == '$' {
			lex.peekInc()

			if runeMatchUppercase(lex.peek()) {
				lex.peekInc()

				toks = append(toks, lex.consume(TTDMVar))
				continue
			}

			if runeMatchLowercase(lex.peek()) {
				lex.peekInc()

				toks = append(toks, lex.consume(TTDSVar))
				continue
			}

			if lex.peek() == '$' {
				lex.peekInc()

				toks = append(toks, lex.consume(TTDAMVar))
				continue
			}

			lex.peekReset()
		}

		if runeMatchUppercase(lex.peek()) {
			lex.peekInc()

			toks = append(toks, lex.consume(TTMVar))
			continue
		}

		if runeMatchLowercase(lex.peek()) {
			lex.peekInc()

			toks = append(toks, lex.consume(TTSVar))
			continue
		}

		return toks, fmt.Errorf("unrecognized token starting at %q", lex.peek())
	}

	toks = append(toks, &Token{TType: TTEOF})

	return toks, nil
}

func runeMatchNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func runeMatchUppercase(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func runeMatchLowercase(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func runeMatchLetter(r rune) bool {
	return runeMatchLowercase(r) || runeMatchUppercase(r)
}

func runeMatchWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func matchNumber(lex *lexer) bool {
	if !runeMatchNumber(lex.peek()) {
		return false
	}

	lex.peekInc()

	for runeMatchNumber(lex.peek()) {
		lex.peekInc()
	}

	return true
}

func matchFunction(lex *lexer) bool {
	if !runeMatchLetter(lex.peek()) {
		return false
	}

	lex.peekInc()

	if !runeMatchLetter(lex.peek()) {
		return false
	}

	lex.peekInc()

	for runeMatchLetter(lex.peek()) {
		lex.peekInc()
	}

	return true
}
