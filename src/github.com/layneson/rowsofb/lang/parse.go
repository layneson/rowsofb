package lang

import (
	"fmt"
	"strings"
)

/*
   Parsing grammar:

       expr    -> term ((ttPlus | ttMinus) term)* (arrow (ttMVar | ttSVar))? EOF
       term    -> factor ((ttMult | ttDiv) factor)*
       factor  -> (ttMinus)? ttNum
               -> (ttMinus)? ttFunc ttLParen expr (ttComma expr)* ttRParen
               -> (ttMinus)? ttDMVar | ttDSVar | ttDAMVar | ttMVar | ttSVar
               -> (ttMinus)? ttLParen expr ttRParen
*/

// A ExprNode represents an expression.
type ExprNode struct {
	First *TermNode

	Operators []*Token
	Terms     []*TermNode

	ResultVar *Token
}

func (enode *ExprNode) String() string {
	s := fmt.Sprintf("expr(%s", enode.First)
	for i, op := range enode.Operators {
		s += fmt.Sprintf(" <%s> %s", op.TType, enode.Terms[i])
	}

	return s + ")"
}

// A TermNode represents a term.
type TermNode struct {
	First *FactorNode

	Operators []*Token
	Factors   []*FactorNode
}

func (tnode *TermNode) String() string {
	s := fmt.Sprintf("term(%s", tnode.First)
	for i, op := range tnode.Operators {
		s += fmt.Sprintf(" <%s> %s", op.TType, tnode.Factors[i])
	}

	return s + ")"
}

// FactorType represents a type of factor.
type FactorType int

// FactorType definitions.
const (
	NumFactor FactorType = iota
	FuncFactor
	VarFactor
	ParenFactor
)

func (ft FactorType) String() string {
	switch ft {
	case NumFactor:
		return "numFactor"
	case FuncFactor:
		return "funcFactor"
	case VarFactor:
		return "varFactor"
	case ParenFactor:
		return "parenFactor"
	}

	return "unknown"
}

// A FactorNode represents a factor.
type FactorNode struct {
	FType FactorType // indicates which of the following variables to use

	Neg *Token // nil if there is no ttMinus preceeding the factornode

	Num *Token

	Function *Token
	FuncArgs []*ExprNode

	Variable *Token

	ParenExpr *ExprNode
}

func (fnode *FactorNode) String() string {
	s := "factor("

	if fnode.Neg != nil {
		s += "-"
	}

	s += fnode.FType.String()

	switch fnode.FType {
	case NumFactor:
		s += fmt.Sprintf(" <%s>", fnode.Num.TType)
	case FuncFactor:
		argstrs := []string{}
		for _, e := range fnode.FuncArgs {
			argstrs = append(argstrs, e.String())
		}
		s += fmt.Sprintf(" <func>(%s)", strings.Join(argstrs, ","))
	case VarFactor:
		s += fmt.Sprintf(" <%s>", fnode.Variable.TType)
	case ParenFactor:
		s += fmt.Sprintf(" (%s)", fnode.ParenExpr)
	}

	return s + ")"
}

type parser struct {
	toks []*Token
	pos  int
}

func (p *parser) peek() *Token {
	return p.toks[p.pos]
}

func (p *parser) consume() *Token {
	tok := p.toks[p.pos]
	p.pos++
	return tok
}

// Parse takes a list of token pointers and returns an expression pointer and an error.
func Parse(toks []*Token) (*ExprNode, error) {
	psr := &parser{toks: toks, pos: 0}

	expr, err := parseExpr(psr)
	if err != nil {
		return expr, err
	}

	if psr.peek().TType == TTArrow {
		psr.consume()

		if psr.peek().TType != TTMVar && psr.peek().TType != TTSVar {
			return expr, fmt.Errorf("expected one of (%q, %q) but found %q", TTMVar, TTSVar, psr.peek().TType)
		}

		expr.ResultVar = psr.consume()
	}

	if psr.peek().TType != TTEOF {
		return expr, fmt.Errorf("expected %q but found %q", TTEOF, psr.peek().TType)
	}

	return expr, nil
}

func parseExpr(psr *parser) (*ExprNode, error) {
	enode := &ExprNode{}

	first, err := parseTerm(psr)
	if err != nil {
		return enode, err
	}

	enode.First = first

	for psr.peek().TType == TTPlus || psr.peek().TType == TTMinus {
		op := psr.consume()

		term, err := parseTerm(psr)
		if err != nil {
			return enode, err
		}

		enode.Operators = append(enode.Operators, op)
		enode.Terms = append(enode.Terms, term)
	}

	return enode, nil
}

func parseTerm(psr *parser) (*TermNode, error) {
	tnode := &TermNode{}

	first, err := parseFactor(psr)
	if err != nil {
		return tnode, err
	}

	tnode.First = first

	for psr.peek().TType == TTMult || psr.peek().TType == TTDiv {
		op := psr.consume()

		factor, err := parseFactor(psr)
		if err != nil {
			return tnode, err
		}

		tnode.Operators = append(tnode.Operators, op)
		tnode.Factors = append(tnode.Factors, factor)
	}

	return tnode, nil
}

func parseFactor(psr *parser) (*FactorNode, error) {
	fnode := &FactorNode{}

	if psr.peek().TType == TTMinus {
		fnode.Neg = psr.consume()
	}

	if psr.peek().TType == TTNum {
		fnode.Num = psr.consume()

		fnode.FType = NumFactor
		return fnode, nil
	}

	if psr.peek().TType == TTFunc {
		fnode.Function = psr.consume()

		if psr.peek().TType != TTLParen {
			return fnode, fmt.Errorf("expected %q but found %q", TTLParen, psr.peek().TType)
		}

		psr.consume()

		fexpr, err := parseExpr(psr)
		if err != nil {
			return nil, err
		}

		fnode.FuncArgs = []*ExprNode{fexpr}

		for psr.peek().TType == TTComma {
			psr.consume()

			expr, err := parseExpr(psr)
			if err != nil {
				return nil, err
			}

			fnode.FuncArgs = append(fnode.FuncArgs, expr)
		}

		if psr.peek().TType != TTRParen {
			return fnode, fmt.Errorf("expected %q but found %q", TTRParen, psr.peek().TType)
		}

		psr.consume()

		fnode.FType = FuncFactor
		return fnode, nil
	}

	if psr.peek().TType == TTLParen {
		psr.consume()

		expr, err := parseExpr(psr)
		if err != nil {
			return fnode, err
		}

		fnode.ParenExpr = expr

		if psr.peek().TType != TTRParen {
			return fnode, fmt.Errorf("expected %q but found %q", TTRParen, psr.peek().TType)
		}

		psr.consume()

		fnode.FType = ParenFactor
		return fnode, err
	}

	if psr.peek().TType != TTMVar && psr.peek().TType != TTSVar && psr.peek().TType != TTDMVar && psr.peek().TType != TTDSVar && psr.peek().TType != TTDAMVar {
		return fnode, fmt.Errorf("expected one of (%q, %q, %q, %q, %q) but found %q", TTMVar, TTSVar, TTDMVar, TTDSVar, TTDAMVar, psr.peek().TType)
	}

	fnode.Variable = psr.consume()

	fnode.FType = VarFactor
	return fnode, nil
}
