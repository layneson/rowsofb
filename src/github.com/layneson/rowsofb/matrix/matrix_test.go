package matrix

import (
	"strconv"
	"strings"
	"testing"
)

func manualMatrix(mat [][]string) M {
	m := New(len(mat), len(mat[0]))

	for r := 0; r < len(mat); r++ {
		for c := 0; c < len(mat[r]); c++ {
			spl := strings.Split(mat[r][c], "/")
			var n, d int

			if len(spl) == 1 {
				n, _ = strconv.Atoi(spl[0])
				d = 1
			} else {
				n, _ = strconv.Atoi(spl[0])
				d, _ = strconv.Atoi(spl[1])
			}

			m.Set(r+1, c+1, NewFrac(n, d))
		}
	}

	return m
}

func fractionEquals(f1, f2 Frac) bool {
	f1.normalizeSignage()
	f2.normalizeSignage()
	return f1.n == f2.n && f1.d == f2.d
}

func matrixEquals(m1, m2 M) bool {
	if m1.r != m2.r || m1.c != m2.c {
		return false
	}

	for r := 1; r <= m1.r; r++ {
		for c := 1; c <= m1.c; c++ {
			plc1 := m1.Get(r, c).Reduce()
			plc2 := m2.Get(r, c).Reduce()
			if !fractionEquals(plc1, plc2) {
				return false
			}
		}
	}

	return true
}

func TestFracAdd(t *testing.T) {
	res := NewFrac(-15, 2).Add(NewFrac(7, 8))
	if !fractionEquals(res, NewFrac(-53, 8)) {
		t.Errorf("Fraction addition failed: expected %v but got %v", NewFrac(-53, 8), res)
	}
}

func TestReduce(t *testing.T) {
	tests := [][]Frac{
		{NewFrac(21, 3), NewScalarFrac(7)},
		{NewFrac(3, 3), NewScalarFrac(1)},
		{NewFrac(0, 124), NewFrac(0, 1)},
		{NewFrac(0, 13286025), NewFrac(0, 1)},
		{NewFrac(-135, 405), NewFrac(-1, 3)},
	}

	for _, tst := range tests {
		res := tst[0].Reduce()
		if !fractionEquals(res, tst[1]) {
			t.Errorf("Reduced fraction %d/%d should be %d/%d but was %d/%d!", tst[0].n, tst[0].d, tst[1].n, tst[1].d, res.n, res.d)
		}
	}
}

func TestGCD(t *testing.T) {
	tests := [][]int{
		{4, 8, 4},
		{48, 60, 12},
		{3348384, 901488, 128784},
	}

	for _, tst := range tests {
		res := gcd(tst[0], tst[1])
		if res != tst[2] {
			t.Errorf("gcd(%d, %d) should be %d, but %d was found instead!", tst[0], tst[1], tst[2], res)
		}
	}
}

func TestIsLeadingEntry(t *testing.T) {
	input := manualMatrix([][]string{
		{"2", "3", "1", "-1"},
		{"0", "2", "1", "2"},
		{"0", "0", "0", "1"},
		{"0", "0", "0", "0"},
	})

	if !isLeadingEntry(input, 1, 1) {
		t.Error("1,1 must be a leading entry!")
	}
	if !isLeadingEntry(input, 2, 2) {
		t.Error("2, 2 must be a leading entry!")
	}
	if !isLeadingEntry(input, 3, 4) {
		t.Error("3, 4 must be a leading entry!")
	}
	if isLeadingEntry(input, 4, 4) {
		t.Error("4, 4, must NOT be a leading entry!")
	}
}

func TestRef(t *testing.T) {
	tests := [][]M{
		{manualMatrix([][]string{
			{"3", "0", "-5"},
			{"1", "-5", "0"},
			{"1", "1", "-2"},
		}), manualMatrix([][]string{
			{"1", "0", "-5/3"},
			{"0", "1", "-1/3"},
			{"0", "0", "0"},
		})},

		{New(5, 5), New(5, 5)},

		{Identity(3), Identity(3)},
	}

	for _, tst := range tests {
		res := Ref(tst[0])
		if !matrixEquals(res, tst[1]) {
			t.Error("Incorrect matrix reduction result!")
		}
	}

}

func TestRref(t *testing.T) {
	tests := [][]M{
		{manualMatrix([][]string{
			{"-12", "2", "-6"},
			{"18", "-3", "9"},
			{"-2", "1/3", "-1"},
		}), manualMatrix([][]string{
			{"1", "-1/6", "1/2"},
			{"0", "0", "0"},
			{"0", "0", "0"},
		})},

		{manualMatrix([][]string{
			{"-1", "0", "1"},
			{"-1", "3", "0"},
			{"-4", "12", "-1"},
		}), manualMatrix([][]string{
			{"1", "0", "0"},
			{"0", "1", "0"},
			{"0", "0", "1"},
		})},

		{manualMatrix([][]string{
			{"1", "2", "3"},
			{"2", "3", "4"},
			{"5/2", "6/3", "7/8"},
		}), manualMatrix([][]string{
			{"1", "0", "-1"},
			{"0", "1", "2"},
			{"0", "0", "-5/8"},
		})},

		{New(5, 5), New(5, 5)},
	}

	for _, tst := range tests {
		res := Rref(tst[0])
		if !matrixEquals(res, tst[1]) {
			t.Error("Incorrect matrix row reduction result!")
		}
	}
}

func TestAugment(t *testing.T) {
	tests := [][]M{
		{manualMatrix([][]string{
			{"1", "3", "5"},
			{"2", "4", "6"},
			{"-2", "-4", "-6"},
		}), manualMatrix([][]string{
			{"7", "9"},
			{"8", "10"},
			{"-8", "-10"},
		}), manualMatrix([][]string{
			{"1", "3", "5", "7", "9"},
			{"2", "4", "6", "8", "10"},
			{"-2", "-4", "-6", "-8", "-10"},
		})},
	}

	for _, tst := range tests {
		res, _ := Augment(tst[0], tst[1])
		if !matrixEquals(res, tst[2]) {
			t.Error("Incorrect matrix augmentation!")
		}
	}
}

func TestInverse(t *testing.T) {
	tests := [][]M{
		{manualMatrix([][]string{
			{"2", "6", "8"},
			{"6", "18", "25"},
			{"6", "17", "32"},
		}), manualMatrix([][]string{
			{"151/2", "-28", "3"},
			{"-21", "8", "-1"},
			{"-3", "1", "0"},
		})},
	}

	for _, tst := range tests {
		res, err := Inverse(tst[0])
		if err != nil {
			t.Errorf("Got error during matrix inverse calculation: %v", err)
		}
		if !matrixEquals(res, tst[1]) {
			t.Error("Incorrect matrix inverse!")
		}
	}
}

func TestAdd(t *testing.T) {
	tests := [][]M{
		{manualMatrix([][]string{
			{"1", "2", "3"},
			{"4", "5", "6"},
		}), manualMatrix([][]string{
			{"4", "5", "6"},
			{"2", "2", "6"},
		}), manualMatrix([][]string{
			{"5", "7", "9"},
			{"6", "7", "12"},
		})},
	}

	for _, tst := range tests {
		res, err := Add(tst[0], tst[1])
		if err != nil {
			t.Errorf("Got error during matrix addition: %v", err)
		}
		if !matrixEquals(res, tst[2]) {
			t.Error("Incorrect matrix addition!")
		}
	}
}

func TestScale(t *testing.T) {
	tests := [][]interface{}{
		{manualMatrix([][]string{
			{"1", "2", "2"},
			{"7", "1", "-1"},
			{"2", "3/2", "0"},
		}), NewScalarFrac(2), manualMatrix([][]string{
			{"2", "4", "4"},
			{"14", "2", "-2"},
			{"4", "3", "0"},
		})},
	}

	for _, tst := range tests {
		m1 := tst[0].(M)
		s := tst[1].(Frac)
		m2 := tst[2].(M)

		res := Scale(s, m1)
		if !matrixEquals(res, m2) {
			t.Error("Incorrect matrix scalar multiplication!")
		}
	}
}

func TestMult(t *testing.T) {
	tests := [][]M{
		{manualMatrix([][]string{
			{"1", "2", "1", "1"},
			{"7", "1", "2", "0"},
			{"3", "-1", "1", "0"},
		}), manualMatrix([][]string{
			{"1", "0", "1"},
			{"0", "2", "0"},
			{"1", "7", "0"},
			{"0", "0", "-1"},
		}), manualMatrix([][]string{
			{"2", "11", "0"},
			{"9", "16", "7"},
			{"4", "5", "3"},
		})},
	}

	for _, tst := range tests {
		res, err := Multiply(tst[0], tst[1])
		if err != nil {
			t.Errorf("Matrix multiplication failed with error: %v", err)
		}

		if !matrixEquals(res, tst[2]) {
			t.Error("Incorrect result of matrix multiplication!")
			t.Errorf("Expected \n%v but got \n%v", tst[2], res)
		}
	}
}
