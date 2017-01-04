package env

import "testing"
import "github.com/layneson/rowsofb/lang"
import "github.com/layneson/rowsofb/matrix"

func TestEvaluate(t *testing.T) {
	tinputs := []*lang.ExprNode{
		buildExpr(
			buildTerm(buildNumFactor("5")).
				div(buildNumFactor("4")).
				mult(buildNumFactor("10")).term,
		).expr,
		buildExpr(
			buildTerm(buildParenFactor(
				buildExpr(
					buildTerm(buildNumFactor("5")).
						div(buildNumFactor("4")).term,
				).expr,
			)).mult(buildNumFactor("10")).term,
		).expr,
		buildExpr(
			buildTerm(buildNumFactor("2")).
				mult(buildNumFactor("3")).
				div(buildNumFactor("4")).
				mult(buildNumFactor("6")).
				div(buildNumFactor("10")).term,
		).expr,
	}

	toutputs := []*Value{
		&Value{
			VType:  SVar,
			SValue: matrix.NewFrac(1, 8),
		},
		&Value{
			VType:  SVar,
			SValue: matrix.NewFrac(25, 2),
		},
		&Value{
			VType:  SVar,
			SValue: matrix.NewFrac(6, 240),
		},
	}

	for i, input := range tinputs {
		output, err := Evaluate(input, New(nil, nil, nil))
		if err != nil {
			t.Fatalf("call to Evaluate failed with error: %v", err)
			return
		}

		if output.VType != toutputs[i].VType {
			t.Fatalf("expected output vtype %s but got %s", toutputs[i].VType, output.VType)
			return
		}

		if !output.SValue.Equals(toutputs[i].SValue) {
			t.Fatalf("expected output %s but got %s", toutputs[i].SValue, output.SValue)
		}
	}
}

type exprbuilder struct {
	expr *lang.ExprNode
}

func buildExpr(first *lang.TermNode) exprbuilder {
	return exprbuilder{&lang.ExprNode{First: first}}
}

func (eb exprbuilder) add(term *lang.TermNode) exprbuilder {
	eb.expr.Operators = append(eb.expr.Operators, &lang.Token{
		Literal: "+",
		TType:   lang.TTPlus,
	})

	eb.expr.Terms = append(eb.expr.Terms, term)

	return eb
}

func (eb exprbuilder) sub(term *lang.TermNode) exprbuilder {
	eb.expr.Operators = append(eb.expr.Operators, &lang.Token{
		Literal: "-",
		TType:   lang.TTMinus,
	})

	eb.expr.Terms = append(eb.expr.Terms, term)

	return eb
}

type termbuilder struct {
	term *lang.TermNode
}

func buildTerm(first *lang.FactorNode) termbuilder {
	return termbuilder{term: &lang.TermNode{First: first}}
}

func (tb termbuilder) mult(factor *lang.FactorNode) termbuilder {
	tb.term.Operators = append(tb.term.Operators, &lang.Token{
		Literal: "*",
		TType:   lang.TTMult,
	})

	tb.term.Factors = append(tb.term.Factors, factor)

	return tb
}

func (tb termbuilder) div(factor *lang.FactorNode) termbuilder {
	tb.term.Operators = append(tb.term.Operators, &lang.Token{
		Literal: "/",
		TType:   lang.TTDiv,
	})

	tb.term.Factors = append(tb.term.Factors, factor)

	return tb
}

func buildNumFactor(num string) *lang.FactorNode {
	return &lang.FactorNode{
		FType: lang.NumFactor,
		Num: &lang.Token{
			Literal: num,
			TType:   lang.TTNum,
		},
	}
}

func buildParenFactor(expr *lang.ExprNode) *lang.FactorNode {
	return &lang.FactorNode{
		FType:     lang.ParenFactor,
		ParenExpr: expr,
	}
}
