package lang

import "testing"

func tokentypesToTokens(types []TokenType) []*Token {
	toks := make([]*Token, len(types))

	for i, tt := range types {
		toks[i] = &Token{TType: tt, Literal: "test"}
	}

	return toks
}

func TestParse(t *testing.T) {
	tinputs := [][]TokenType{
		{TTNum, TTPlus, TTNum, TTEOF},
		{TTMinus, TTNum, TTMult, TTMVar, TTEOF},
		{TTNum, TTMult, TTLParen, TTFunc, TTLParen, TTDAMVar, TTPlus, TTMVar, TTRParen, TTMinus, TTNum, TTRParen, TTEOF},
		{TTFunc, TTLParen, TTNum, TTComma, TTNum, TTRParen, TTEOF},
	}

	toutputs := []string{
		"expr(term(factor(numFactor <num>)) <plus> term(factor(numFactor <num>)))",
		"expr(term(factor(-numFactor <num>) <mult> factor(varFactor <mvar>)))",
		"expr(term(factor(numFactor <num>) <mult> factor(parenFactor (expr(term(factor(funcFactor <func>(expr(term(factor(varFactor <damvar>)) <plus> term(factor(varFactor <mvar>)))))) <minus> term(factor(numFactor <num>)))))))",
		"expr(term(factor(funcFactor <func>(expr(term(factor(numFactor <num>))),expr(term(factor(numFactor <num>)))))))",
	}

	for i, types := range tinputs {
		toks := tokentypesToTokens(types)

		expr, err := Parse(toks)
		if err != nil {
			t.Fatalf("Parsing testing failed due to error: %v", err)
			return
		}

		if expr.String() != toutputs[i] {
			t.Fatalf("Parsing testing failed looking for\n\t%s\nbut found\n\t%s", toutputs[i], expr)
		}
	}
}
