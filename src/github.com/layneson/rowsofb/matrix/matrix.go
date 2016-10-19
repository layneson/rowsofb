package matrix

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

//Frac represents a fractional number.
type Frac struct {
	//Numerator and denominator.
	n, d int
}

//NewFrac creates a new fraction with the specified numerator and denominator.
func NewFrac(n, d int) Frac {
	if d == 0 {
		panic("Zero denominator not acceptable!")
	}

	f := Frac{n: n, d: d}

	f.normalizeSignage()

	return f
}

//NewScalarFrac returns a fraction that represents a whole number.
func NewScalarFrac(s int) Frac {
	return Frac{n: s, d: 1}
}

//ParseFrac parses a string fraction, returning an error if the fraction is invalid.
func ParseFrac(s string) (Frac, error) {
	fields := strings.Split(s, "/")

	if len(fields) == 1 {
		fields = append(fields, "1")
	}

	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return Frac{}, err
	}

	d, err := strconv.Atoi(fields[1])
	if err != nil {
		return Frac{}, err
	}

	return NewFrac(n, d), nil
}

//String returns a string representation of the fraction.
func (f Frac) String() string {
	if f.d == 1 {
		return strconv.Itoa(f.n)
	}
	return strconv.Itoa(f.n) + "/" + strconv.Itoa(f.d)
}

//IsZero returns true if the fraction is equal to zero.
func (f Frac) IsZero() bool {
	return f.n == 0
}

//normalizeSignage is the only method that mutates Frac. It fixes negative n and d, and also moves the negative from d to n if it exists.
//The end result is a fraction with no negative or a fraction with a negative n.
//Must only be called from another Frac method.
func (f *Frac) normalizeSignage() {
	//Switch signs if n and d are both negative, or if d is negative but n is not (to always keep negative on top)
	if f.n < 0 && f.d < 0 || f.d < 0 {
		f.n *= -1
		f.d *= -1
	}
}

//Numerator returns the fraction's numerator.
func (f Frac) Numerator() int {
	return f.n
}

//Denominator returns the fraction's denominator.
func (f Frac) Denominator() int {
	return f.d
}

//Mul multiplies two fractions and returns the result as a new fraction.
func (f1 Frac) Mul(f2 Frac) Frac {
	f := Frac{n: f1.n * f2.n, d: f1.d * f2.d}

	f.normalizeSignage()

	return f
}

//Div divides the fraction by another and returns the result as a new fraction.
func (f1 Frac) Div(f2 Frac) Frac {
	return f1.Mul(f2.Reciprocal())
}

//Add adds two fractions and returns the result.
//This method does not compute the LCD.
func (f1 Frac) Add(f2 Frac) Frac {
	f := Frac{n: f1.n*f2.d + f2.n*f1.d, d: f1.d * f2.d}

	f.normalizeSignage()

	return f
}

//Reciprocal returns the reciprocal (multiplicative inverse) of the fraction.
func (f Frac) Reciprocal() Frac {
	f1 := Frac{n: f.d, d: f.n}

	f1.normalizeSignage()

	return f1
}

//Neg negates the fraction (multiplies it by -1).
func (f Frac) Neg() Frac {
	return Frac{n: f.n * -1, d: f.d}
}

func gcd(i1, i2 int) int {
	for i2 != 0 {
		t := i2
		i2 = i1 % i2
		i1 = t
	}

	return i1
}

//Reduce reduces the fraction and returns the result.
func (f Frac) Reduce() Frac {
	d := gcd(f.n, f.d)

	f1 := Frac{n: f.n / d, d: f.d / d}

	f1.normalizeSignage()

	return f1
}

//M represents a matrix.
type M struct {
	//Number of rows and columns.
	r, c int

	//The values of the matrix, in row order.
	values []Frac
}

//NewWithValues returns a new matrix with initial values vals.
func NewWithValues(r, c int, vals []Frac) M {
	return M{r: r, c: c, values: vals}
}

//New returns a zero matrix of size r,c.
func New(r, c int) M {
	vals := make([]Frac, r*c)
	for i := range vals {
		vals[i] = NewScalarFrac(0)
	}

	return M{r: r, c: c, values: vals}
}

//String returns a string representation of the matrix.
func (m M) String() string {
	var buff bytes.Buffer

	buff.WriteString("┌  ")
	for c := 1; c < m.Cols(); c++ {
		buff.WriteString(" \t")
	}
	buff.WriteString("  ┐\n")
	for r := 1; r <= m.Rows(); r++ {
		buff.WriteString("│ ")
		buff.WriteString(m.Get(r, 1).String())
		for c := 2; c <= m.Cols(); c++ {
			buff.WriteString("\t")
			buff.WriteString(m.Get(r, c).String())
		}
		buff.WriteString(" │\n")
	}
	buff.WriteString("└  ")
	for c := 1; c < m.Cols(); c++ {
		buff.WriteString(" \t")
	}
	buff.WriteString("  ┘")

	return buff.String()
}

//Rows returns the number of rows in the matrix.
func (m M) Rows() int {
	return m.r
}

//Cols returns the number of columns in the matrix.
func (m M) Cols() int {
	return m.c
}

//Get returns the value at the specified row and column.
func (m M) Get(r, c int) Frac {
	r, c = r-1, c-1
	return m.values[r*m.c+c]
}

//Set sets the value at the specified row and column to the given fraction.
func (m *M) Set(r, c int, v Frac) {
	r, c = r-1, c-1
	m.values[r*m.c+c] = v
}

//SwitchRows switches two rows. It is an elementary row operation.
func (m *M) SwitchRows(r1, r2 int) {
	r1, r2 = r1-1, r2-1
	for c := 0; c < m.c; c++ {
		tmp := m.values[r1*m.c+c]
		m.values[r1*m.c+c] = m.values[r2*m.c+c]
		m.values[r2*m.c+c] = tmp
	}
}

//MultiplyRow multiplies row r by the scalar s. It is an elementary row operation.
func (m *M) MultiplyRow(r int, s Frac) {
	r = r - 1
	for c := 0; c < m.c; c++ {
		m.values[r*m.c+c] = m.values[r*m.c+c].Mul(s).Reduce()
	}
}

//MultiplyAndAddRow adds row r1 multiplied by scalar s to row r2. It is an elementary row operation.
func (m *M) MultiplyAndAddRow(r1 int, s Frac, r2 int) {
	r1, r2 = r1-1, r2-1
	for c := 0; c < m.c; c++ {
		m.values[r2*m.c+c] = m.values[r2*m.c+c].Add(m.values[r1*m.c+c].Mul(s)).Reduce()
	}
}

//Creates a copy of the matrix, with same contents and size.
func copyMatrix(m M) M {
	mm := M{r: m.r, c: m.c, values: make([]Frac, len(m.values))}
	copy(mm.values, m.values)

	return mm
}

//Transpose takes a copy of a matrix and returns its transpose.
func Transpose(m M) M {
	rm := M{r: m.c, c: m.r, values: make([]Frac, m.c*m.r)}

	for r := 1; r <= m.r; r++ {
		for c := 1; c <= m.c; c++ {
			rm.Set(c, r, m.Get(r, c))
		}
	}

	return rm
}

func isLeadingEntry(m M, r, c int) bool {
	if m.Get(r, c).IsZero() {
		return false // must not be zero to be a leading entry
	}

	for cc := c - 1; cc > 0; cc-- {
		if !m.Get(r, cc).IsZero() {
			return false
		}
	}

	return true
}

//Ref takes a copy of a matrix and returns itself in row echelon form.
func Ref(m M) M {
	m = copyMatrix(m)

	startr := 1
	for c := 1; c <= m.Cols(); c++ { // find a leading entry in this column
		found := false
		for r := startr; r <= m.Rows(); r++ {
			if isLeadingEntry(m, r, c) {
				found = true
				m.SwitchRows(startr, r) // move it to the top, son!
				break
			}
		}

		if !found { // no leading entry, next column please
			continue
		}

		m.MultiplyRow(startr, m.Get(startr, c).Reciprocal()) // make first entry one

		for r := startr + 1; r <= m.Rows(); r++ {
			if isLeadingEntry(m, r, c) {
				m.MultiplyAndAddRow(startr, m.Get(startr, c).Reciprocal().Mul(m.Get(r, c).Neg()), r) // zero first column entry
			}
		}

		startr++ // row is now in ref
	}

	return m
}

//Rref takes a copy of a matrix and returns it in rrrrrrreduced rrrrow echelon-a forrrrm-a!
func Rref(m M) M {
	m = copyMatrix(m)

	m = Ref(m)

	for c := 1; c <= m.Cols(); c++ {
		for r := 1; r <= m.Rows(); r++ {
			if isLeadingEntry(m, r, c) {
				m.MultiplyRow(r, m.Get(r, c).Reciprocal()) // make the leading entry 1
				for rr := r - 1; rr > 0; rr-- {            // for each row above the current row...
					if !m.Get(rr, c).IsZero() {
						m.MultiplyAndAddRow(r, m.Get(rr, c).Neg().Mul(m.Get(r, c).Reciprocal()), rr) // clear entry above leading entry
					}
				}
			}
		}
	}

	return m
}

//Inverse takes a copy of a matrix and returns its inverse.
//An error is returned if the matrix has no inverse.
func Inverse(m M) (M, error) {
	m = copyMatrix(m)

	if m.Rows() != m.Cols() {
		return m, errors.New("non-square matrices have no inverse")
	}

	m, _ = Augment(m, Identity(m.r)) // ignore error because Identity will always match m row size

	m = Rref(m)

	for c := 1; c <= m.Rows(); c++ { // must be square, so m.Rows() works here
		found := false
		for r := 1; r <= m.Rows(); r++ {
			if isLeadingEntry(m, r, c) && m.Get(r, c).Numerator() == 1 && m.Get(r, c).Denominator() == 1 { //Leading entry that is one
				found = true
				break
			}
		}

		if !found { //no leading entry that is one... not the identity, not invertible :(
			return m, errors.New("matrix has no inverse")
		}
	}

	rm := M{r: m.r, c: m.r, values: make([]Frac, m.r*m.r)}

	for r := 1; r <= rm.Rows(); r++ {
		for c := 1; c <= rm.Cols(); c++ {
			rm.Set(r, c, m.Get(r, rm.Cols()+c))
		}
	}

	return rm, nil
}

//Identity returns the identity matrix of size i.
func Identity(i int) M {
	rm := M{r: i, c: i, values: make([]Frac, i*i)}

	for r := 1; r <= rm.Rows(); r++ {
		for c := 1; c <= rm.Cols(); c++ {
			if r == c {
				rm.Set(r, c, NewScalarFrac(1))
			} else {
				rm.Set(r, c, NewScalarFrac(0))
			}
		}
	}

	return rm
}

//Augment augments a with b then returns this matrix.
//It returns an error if the two matrices do not have the same number of rows.
func Augment(a, b M) (M, error) {
	if a.r != b.r {
		return a, errors.New("augmented matrices must have equal row counts")
	}

	rm := M{r: a.r, c: a.c + b.c}

	rm.values = make([]Frac, rm.r*rm.c)

	for r := 1; r <= a.Rows(); r++ {
		for c := 1; c <= a.Cols(); c++ {
			rm.Set(r, c, a.Get(r, c))
		}
	}

	for r := 1; r <= a.Rows(); r++ {
		for c := 1; c <= b.Cols(); c++ {
			rm.Set(r, a.Cols()+c, b.Get(r, c))
		}
	}

	return rm, nil
}

//Add adds matrix a to matrix b.
//It returns an error if a and b are not the same size.
func Add(a, b M) (M, error) {
	if a.r != b.r || a.c != b.c {
		return a, errors.New("addition requires two identically-sized matrices")
	}

	rm := copyMatrix(b)

	for r := 1; r <= rm.Rows(); r++ {
		for c := 1; c <= rm.Cols(); c++ {
			rm.Set(r, c, rm.Get(r, c).Add(a.Get(r, c)))
		}
	}

	return rm, nil
}

//Scale scales a matrix by a... you guessed it... scalar. It returns this new matrix.
func Scale(s Frac, m M) M {
	m = copyMatrix(m)

	for r := 1; r <= m.Rows(); r++ {
		m.MultiplyRow(r, s)
	}

	return m
}

//Multiply multiplies a by b. If the matrices cannot be multiplied, an error is returned.
func Multiply(a, b M) (M, error) {
	if a.c != b.r {
		return a, errors.New("multiplication can only be done on matrices A and B if the number of columns of A equals the number of rows of B")
	}

	rm := M{r: a.r, c: b.c, values: make([]Frac, a.r*b.c)}

	for r := 1; r <= rm.Rows(); r++ {
		for c := 1; c <= rm.Cols(); c++ {
			sum := NewScalarFrac(0)
			for count := 1; count <= a.c; count++ {
				f1 := a.Get(r, count)
				f2 := b.Get(count, c)
				res := f1.Mul(f2)
				sum = sum.Add(res)
			}
			rm.Set(r, c, sum)
		}
	}

	return rm, nil
}
