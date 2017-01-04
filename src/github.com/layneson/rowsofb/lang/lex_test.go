package lang

import "testing"

func TestLex(t *testing.T) {
	tmap := map[string][]TokenType{
		"()":                  []TokenType{TTLParen, TTRParen},
		"a":                   []TokenType{TTSVar},
		"F":                   []TokenType{TTMVar},
		"5+   7":              []TokenType{TTNum, TTPlus, TTNum},
		"sin(f*X) /   4":      []TokenType{TTFunc, TTLParen, TTSVar, TTMult, TTMVar, TTRParen, TTDiv, TTNum},
		"cos(5) + $a*$Z - $$": []TokenType{TTFunc, TTLParen, TTNum, TTRParen, TTPlus, TTDSVar, TTMult, TTDMVar, TTMinus, TTDAMVar},
		"5 - 7 -> A":          []TokenType{TTNum, TTMinus, TTNum, TTArrow, TTMVar},
		"blub(5, 6, A)":       []TokenType{TTFunc, TTLParen, TTNum, TTComma, TTNum, TTComma, TTMVar, TTRParen},
	}

	for input, expected := range tmap {
		output, err := Lex(input)
		if err != nil {
			t.Fatalf("lex test failed with error: %v", err)
			return
		}

		for i, ett := range expected {
			if output[i].TType != ett {
				outTypes := make([]TokenType, len(output))
				for ii, tok := range output {
					outTypes[ii] = tok.TType
				}
				t.Fatalf("lex test failed to match expected output:\n\tinput: %q\n\texpected output: %v\n\toutput: %v", input, expected, outTypes)
				break
			}
		}
	}
}
