package oldlang

import (
	"fmt"
	"strconv"
)

type Statement interface {
}

type ConstStatement struct {
	Statement

	name string
}

type NumStatement struct {
	Statement

	val int
}

type MVarStatement struct {
	Statement

	name string
}

type SVarStatement struct {
	Statement

	name string
}

type DeclMVarStatement struct {
	Statement

	name string
}

type DeclTempMVarStatement struct {
	Statement
}

type CallStatement struct {
	Statement

	name      string
	parameter Statement
}

type AddStatement struct {
	Statement

	left, right Statement
}

type MultStatement struct {
	Statement

	left, right Statement
}

type NegStatement struct {
	Statement

	right Statement
}

type DivStatement struct {
	Statement

	left, right Statement
}

type TokenScanner struct {
	tokens []Token
	index  int
}

func (ts TokenScanner) Peek(i int) Token {
	if i+ts.index >= len(ts.tokens) {
		return ts.tokens[len(ts.tokens)-1]
	}

	return ts.tokens[ts.index+i]
}

func (ts *TokenScanner) Consume() Token {
	t := ts.tokens[ts.index]
	ts.index++
	return t
}

func (ts *TokenScanner) Rollback(i int) {
	ts.index -= i

	if ts.index < 0 {
		ts.index = 0
	}
}

func Compile(line string) (Statement, error) {
	tokens, err := Parse(line)
	if err != nil {
		return nil, err
	}

	ts := TokenScanner{tokens, 0}

	stmt, err := compileTop(ts)

	return stmt, err
}

func compileTop(ts TokenScanner) (Statement, error) {
	_, stmt, err := compileStatement(ts)
	return stmt, err
}

func compileStatement(ts TokenScanner) (int, Statement, error) {
	poss := []func(TokenScanner) (int, Statement, error){
		compileAdd, compileDiv, compileMult, compileNeg,
		compileParen, compileCall, compileDeclMVar, compileDeclTempMVar, compileSVar, compileMVar,
		compileConst, compileNum,
	}

	for _, f := range poss {
		m, s, err := f(ts)
		if err != nil {
			return 0, s, err
		}

		if m > 0 {
			return m, s, nil
		}
	}

	return 0, nil, fmt.Errorf("input parsing error")
}

func compileAdd(ts TokenScanner) (int, Statement, error) {
	accum := 0

	m, s1, err := compileStatement(ts)
	accum += m
	if err != nil || m == 0 {
		ts.Rollback(accum)
		return 0, s1, err
	}

	if ts.Peek(1).ttype == ADD {
		ts.Consume()
		accum++

		m, s2, err := compileStatement(ts)
		accum += m
		if err != nil || m == 0 {
			ts.Rollback(accum)
			return 0, s2, err
		}

		return accum, AddStatement{left: s1, right: s2}, nil
	}

	if ts.Peek(1).ttype == NEG {
		ts.Consume()
		accum++

		m, s2, err := compileStatement(ts)
		accum += m
		if err != nil || m == 0 {
			ts.Rollback(accum)
			return 0, s2, err
		}

		return accum, AddStatement{left: s1, right: NegStatement{right: s2}}, nil
	}

	ts.Rollback(accum)
	return 0, nil, fmt.Errorf("expected statement to follow '+'")
}

func compileDiv(ts TokenScanner) (int, Statement, error) {
	accum := 0

	m, s1, err := compileStatement(ts)
	accum += m
	if err != nil || m == 0 {
		ts.Rollback(accum)
		return 0, s1, err
	}

	if ts.Peek(1).ttype != DIV {
		ts.Rollback(accum)
		return 0, s1, nil
	}

	ts.Consume()
	accum++

	m, s2, err := compileStatement(ts)
	accum += m
	if err != nil || m == 0 {
		ts.Rollback(accum)
		return 0, s2, err
	}

	return accum, DivStatement{left: s1, right: s2}, nil
}

func compileMult(ts TokenScanner) (int, Statement, error) {
	accum := 0

	m, s1, err := compileStatement(ts)
	accum += m
	if err != nil || m == 0 {
		ts.Rollback(accum)
		return 0, s1, err
	}

	if ts.Peek(1).ttype != MULT {
		ts.Rollback(accum)
		return 0, s1, nil
	}

	ts.Consume()
	accum++

	m, s2, err := compileStatement(ts)
	accum += m
	if err != nil || m == 0 {
		ts.Rollback(accum)
		return 0, s2, err
	}

	return accum, DivStatement{left: s1, right: s2}, nil
}

func compileNeg(ts TokenScanner) (int, Statement, error) {
	accum := 0

	if ts.Peek(1).ttype != NEG {
		return 0, nil, nil
	}

	ts.Consume()
	accum++

	m, s, err := compileStatement(ts)
	accum += m
	if err != nil || m == 0 {
		ts.Rollback(accum)
		return 0, s, err
	}

	return accum, NegStatement{right: s}, nil
}

func compileParen(ts TokenScanner) (int, Statement, error) {
	accum := 0

	if ts.Peek(1).ttype != LPAREN {
		return 0, nil, nil
	}

	ts.Consume()
	accum++

	m, s, err := compileStatement(ts)
	accum += m
	if err != nil || m == 0 {
		ts.Rollback(accum)
		return 0, s, err
	}

	if ts.Peek(1).ttype != RPAREN {
		ts.Rollback(accum)
		return 0, nil, nil
	}

	ts.Consume()
	accum++

	return accum, s, nil
}

func compileCall(ts TokenScanner) (int, Statement, error) {
	accum := 0

	if ts.Peek(1).ttype != IDENT {
		return 0, nil, nil
	}

	ntoken := ts.Consume()
	accum++

	if ts.Peek(1).ttype != LPAREN {
		ts.Rollback(accum)
		return 0, nil, nil
	}

	ts.Consume()
	accum++

	m, s, err := compileStatement(ts)
	accum += m
	if err != nil || m == 0 {
		ts.Rollback(accum)
		return 0, s, err
	}

	if ts.Peek(1).ttype != RPAREN {
		ts.Rollback(accum)
		return 0, nil, nil
	}

	ts.Consume()
	accum++

	return accum, CallStatement{name: ntoken.value, parameter: s}, nil
}

func compileDeclMVar(ts TokenScanner) (int, Statement, error) {
	accum := 0

	if ts.Peek(1).ttype != DSIGN {
		return 0, nil, nil
	}

	ts.Consume()
	accum++

	if ts.Peek(1).ttype != MVAR {
		ts.Rollback(accum)
		return 0, nil, nil
	}

	mvar := ts.Consume()
	accum++

	return accum, DeclMVarStatement{name: mvar.value}, nil
}

func compileDeclTempMVar(ts TokenScanner) (int, Statement, error) {
	accum := 0

	if ts.Peek(1).ttype != DSIGN || ts.Peek(2).ttype != DSIGN {
		return 0, nil, nil
	}

	ts.Consume()
	ts.Consume()
	accum += 2

	return accum, DeclTempMVarStatement{}, nil
}

func compileMVar(ts TokenScanner) (int, Statement, error) {
	if ts.Peek(1).ttype != MVAR {
		return 0, nil, nil
	}

	return 1, MVarStatement{name: ts.Consume().value}, nil
}

func compileSVar(ts TokenScanner) (int, Statement, error) {
	if ts.Peek(1).ttype != SVAR {
		return 0, nil, nil
	}

	return 1, SVarStatement{name: ts.Consume().value}, nil
}

func compileConst(ts TokenScanner) (int, Statement, error) {
	if ts.Peek(1).ttype != IDENT {
		return 0, nil, nil
	}

	return 1, ConstStatement{name: ts.Consume().value}, nil
}

func compileNum(ts TokenScanner) (int, Statement, error) {
	if ts.Peek(1).ttype != NUM {
		return 0, nil, nil
	}

	val, err := strconv.Atoi(ts.Consume().value)
	if err != nil {
		ts.Rollback(1)
		return 0, nil, err
	}

	return 1, NumStatement{val: val}, nil
}
