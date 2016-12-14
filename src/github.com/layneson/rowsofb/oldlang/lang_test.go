package oldlang

import "testing"

func typesMatch(tokens []Token, types []TT) bool {
	if len(tokens) != len(types) {
		return false
	}

	for i, token := range tokens {
		if token.ttype != types[i] {
			return false
		}
	}

	return true
}

func TestLex(t *testing.T) {
	testMap := map[string][]TT{
		"A + B": []TT{MVAR, ADD, MVAR},
	}

	for k, v := range testMap {
		tokens, err := Parse(k)
		if err != nil {
			t.Errorf("failed to parse: %v", err)
		}

		if !typesMatch(tokens, v) {
			t.Errorf("lex input %q does not match token types %v (%v)", k, v, tokens)
		}
	}
}
