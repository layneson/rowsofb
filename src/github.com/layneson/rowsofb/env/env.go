package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/layneson/rowsofb/lang"
	"github.com/layneson/rowsofb/matrix"
)

// MatrixDefiner is a function which takes a matrix variable name and returns the matrix the user defined for it.
// The bool return value is false if the user cancelled the process and true otherwise.
type MatrixDefiner func(rune) (matrix.M, bool)

// AnonymousMatrixDefiner is a function which returns a user-defined anonymous matrix.
// The bool return value is false if the user cancelled the process and true otherwise.
type AnonymousMatrixDefiner func() (matrix.M, bool)

// ScalarDefiner is a function which takes a scalar variable name and returns the scalar the user defined for it.
// The bool return value is false if the user cancelled the process and true otherwise.
type ScalarDefiner func(rune) (matrix.Frac, bool)

// E represents an environment which contains 26 matrix variables (A-Z) and 26 scalar variables (a-z).
// The variables Z and z are set to the results of matrix and scalar-resolving expressions, respectively.
type E struct {
	mvars []matrix.M
	svars []matrix.Frac

	mdef  MatrixDefiner
	amdef AnonymousMatrixDefiner
	sdef  ScalarDefiner
}

// New creates a new environment. Each matrix variable defaults to a 3x3 zero matrix
// and each scalar variable defaults to zero.
func New(mdef MatrixDefiner, amdef AnonymousMatrixDefiner, sdef ScalarDefiner) *E {
	e := &E{mdef: mdef, amdef: amdef, sdef: sdef}

	for r := 'A'; r <= 'Z'; r++ {
		e.mvars = append(e.mvars, matrix.New(3, 3))
	}

	for r := 'a'; r <= 'z'; r++ {
		e.svars = append(e.svars, matrix.NewScalarFrac(0))
	}

	return e
}

// GetMVar returns the value of the given matrix variable.
// It assumes the given rune is a valid matrix variable name.
func (e *E) GetMVar(v rune) matrix.M {
	return e.mvars[v-'A']
}

// SetMVar sets the value of the given matrix variable to the given matrix.
// It assumes the given rune is a valid matrix variable name.
func (e *E) SetMVar(v rune, m matrix.M) {
	e.mvars[v-'A'] = m
}

// GetSVar returns the value of the given scalar variable.
// It assumes the given rune is a valid scalar variable name.
func (e *E) GetSVar(v rune) matrix.Frac {
	return e.svars[v-'a']
}

// SetSVar sets the value of the given scalar variable to the given scalar.
// It assumes the given rune is a valid scalar variable name.
func (e *E) SetSVar(v rune, m matrix.Frac) {
	e.svars[v-'a'] = m.Reduce()
}

// VarType represents the type of a certain variable.
type VarType int

// VarType definitions
const (
	MVar VarType = iota
	SVar
	InvalidVar
)

func (vt VarType) String() string {
	switch vt {
	case MVar:
		return "mvar"
	case SVar:
		return "svar"
	case InvalidVar:
		return "invalid"
	}

	return "unknown"
}

// GetVarType returns the type of variable that the given rune represents.
// Returns InvalidVar if v does not represent a valid variable.
func GetVarType(v rune) VarType {
	if v >= 'A' && v <= 'Z' {
		return MVar
	}

	if v >= 'a' && v <= 'z' {
		return SVar
	}

	return InvalidVar
}

// Value represents either a matrix or scalar value.
type Value struct {
	VType VarType

	MValue matrix.M
	SValue matrix.Frac
}

// Evaluate evaluates a lang.ExprNode within the context of the given environment, returning an error if one occurs.
// It also returns a Value which holds the expression result.
func Evaluate(enode *lang.ExprNode, env *E) (*Value, error) {
	val, err := evalExpr(enode, env)
	if err != nil {
		return nil, err
	}

	if enode.ResultVar != nil {
		if val.VType == MVar && enode.ResultVar.TType == lang.TTSVar {
			return nil, fmt.Errorf("cannot assign a matrix value to a scalar variable")
		}

		if val.VType == SVar && enode.ResultVar.TType == lang.TTMVar {
			return nil, fmt.Errorf("cannot assign a scalar value to a matrix variable")
		}

		v := rune(enode.ResultVar.Literal[0])

		switch val.VType {
		case MVar:
			env.SetMVar(v, val.MValue)
		case SVar:
			env.SetSVar(v, val.SValue)
		}
	}

	return val, nil
}

func evalExpr(enode *lang.ExprNode, env *E) (*Value, error) {
	first, err := evalTerm(enode.First, env)
	if err != nil {
		return nil, err
	}

	for i, tnode := range enode.Terms {
		op := enode.Operators[i]

		tval, err := evalTerm(tnode, env)
		if err != nil {
			return nil, err
		}

		first, err = evalAddition(op.TType == lang.TTMinus, first, tval)
		if err != nil {
			return nil, err
		}
	}

	return first, nil
}

func evalAddition(subtraction bool, left, right *Value) (*Value, error) {
	if left.VType != right.VType {
		return nil, fmt.Errorf("cannot perform addition or subtraction with a scalar and a matrix")
	}

	if left.VType == SVar {
		if subtraction {
			right.SValue = right.SValue.Neg()
		}

		return &Value{VType: SVar, SValue: left.SValue.Add(right.SValue)}, nil
	}

	if subtraction {
		right.MValue = matrix.Scale(matrix.NewScalarFrac(-1), right.MValue)
	}

	if left.MValue.Cols() != right.MValue.Cols() || left.MValue.Rows() != right.MValue.Rows() {
		return nil, fmt.Errorf("cannot perform addition or subtraction on two matrices of different sizes")
	}

	sum, _ := matrix.Add(left.MValue, right.MValue)

	return &Value{VType: MVar, MValue: sum}, nil
}

func evalTerm(tnode *lang.TermNode, env *E) (*Value, error) {
	fstack := vstack{}
	ostack := tstack{}
	divqueue := []*Value{}

	for i := len(tnode.Factors) - 1; i >= 0; i-- {
		val, err := evalFactor(tnode.Factors[i], env)
		if err != nil {
			return nil, err
		}

		fstack.push(val)
	}

	fval, err := evalFactor(tnode.First, env)
	if err != nil {
		return nil, err
	}

	fstack.push(fval)

	for i := len(tnode.Operators) - 1; i >= 0; i-- {
		ostack.push(tnode.Operators[i])
	}

	for ostack.canPop() {
		op := ostack.pop()

		if op.TType == lang.TTDiv {
			divqueue = append(divqueue, fstack.pop())
			continue
		}

		right := fstack.pop()
		left := fstack.pop()

		val, err := evalMultiplication(false, left, right)
		if err != nil {
			return nil, err
		}

		fstack.push(val)
	}

	divqueue = append(divqueue, fstack.pop())

	divaccum := divqueue[0]

	for i := 1; i < len(divqueue); i++ {
		val, err := evalMultiplication(true, divaccum, divqueue[i])
		if err != nil {
			return nil, err
		}

		divaccum = val
	}

	return divaccum, nil
}

func evalMultiplication(division bool, left, right *Value) (*Value, error) {
	if left.VType == SVar && right.VType == SVar {
		rrec := right.SValue
		if division {
			rrec = rrec.Reciprocal()
		}

		return &Value{VType: SVar, SValue: left.SValue.Mul(rrec).Reduce()}, nil
	}

	if left.VType == SVar && right.VType == MVar {
		if division {
			return nil, fmt.Errorf("cannot divide a scalar by a matrix")
		}

		return &Value{VType: MVar, MValue: matrix.Scale(left.SValue, right.MValue)}, nil
	}

	if left.VType == MVar && right.VType == SVar {
		rrec := right.SValue
		if division {
			rrec = rrec.Reciprocal()
		}

		return &Value{VType: MVar, MValue: matrix.Scale(rrec, left.MValue)}, nil
	}

	if left.MValue.Cols() != right.MValue.Rows() {
		return nil, fmt.Errorf("cannot multiply a %dx%d matrix by a %dx%d matrix", left.MValue.Rows(), left.MValue.Cols(), right.MValue.Rows(), right.MValue.Cols())
	}

	if division {
		return nil, fmt.Errorf("cannot divide a matrix by a matrix")
	}

	product, _ := matrix.Multiply(left.MValue, right.MValue)

	return &Value{VType: MVar, MValue: product}, nil
}

func evalFactor(fnode *lang.FactorNode, env *E) (*Value, error) {
	val, err := evalFactorIgnoreNeg(fnode, env)
	if err != nil {
		return nil, err
	}

	if fnode.Neg != nil {
		switch val.VType {
		case MVar:
			val.MValue = matrix.Scale(matrix.NewScalarFrac(-1), val.MValue)
		case SVar:
			val.SValue = val.SValue.Mul(matrix.NewScalarFrac(-1))
		}
	}

	return val, nil
}

func evalFactorIgnoreNeg(fnode *lang.FactorNode, env *E) (*Value, error) {
	switch fnode.FType {
	case lang.NumFactor:
		num, _ := strconv.Atoi(fnode.Num.Literal)
		return &Value{VType: SVar, SValue: matrix.NewScalarFrac(num)}, nil
	case lang.ParenFactor:
		return evalExpr(fnode.ParenExpr, env)
	case lang.FuncFactor:
		return evalFunction(fnode, env)
	}

	switch fnode.Variable.TType {
	case lang.TTDMVar:
		v := rune(fnode.Variable.Literal[1])
		mat, ok := env.mdef(v)
		if !ok {
			return nil, fmt.Errorf("user cancelled matrix input")
		}
		env.SetMVar(v, mat)
		return &Value{VType: MVar, MValue: mat}, nil
	case lang.TTDSVar:
		v := rune(fnode.Variable.Literal[1])
		scal, ok := env.sdef(v)
		if !ok {
			return nil, fmt.Errorf("user cancelled scalar input")
		}
		env.SetSVar(v, scal)
		return &Value{VType: SVar, SValue: scal}, nil
	case lang.TTDAMVar:
		mat, ok := env.amdef()
		if !ok {
			return nil, fmt.Errorf("user cancelled matrix input")
		}
		return &Value{VType: MVar, MValue: mat}, nil
	case lang.TTMVar:
		v := rune(fnode.Variable.Literal[0])
		mat := env.GetMVar(v)
		return &Value{VType: MVar, MValue: mat}, nil
	case lang.TTSVar:
		v := rune(fnode.Variable.Literal[0])
		scal := env.GetSVar(v)
		return &Value{VType: SVar, SValue: scal}, nil
	}

	return nil, fmt.Errorf("unexpected factor")
}

func evalFunction(fnode *lang.FactorNode, env *E) (*Value, error) {
	fname := fnode.Function.Literal

	fn, ok := functions[fname]
	if !ok {
		return nil, fmt.Errorf("%q is not a valid function", fname)
	}

	vals := []*Value{}

	for _, enode := range fnode.FuncArgs {
		val, err := evalExpr(enode, env)
		if err != nil {
			return nil, err
		}

		vals = append(vals, val)
	}

	err := checkFunctionArgs(vals, fname, fn)
	if err != nil {
		return nil, err
	}

	return fn.handler(vals)
}

type vstack struct {
	stack []*Value
}

func (s *vstack) push(v *Value) {
	s.stack = append(s.stack, v)
}

func (s *vstack) canPop() bool {
	return len(s.stack) > 0
}

func (s *vstack) pop() *Value {
	v := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return v
}

type tstack struct {
	stack []*lang.Token
}

func (s *tstack) push(t *lang.Token) {
	s.stack = append(s.stack, t)
}

func (s *tstack) canPop() bool {
	return len(s.stack) > 0
}

func (s *tstack) pop() *lang.Token {
	t := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return t
}

type function struct {
	signature []VarType
	vnames    []string
	handler   func([]*Value) (*Value, error)
}

var functions = map[string]function{
	"identity": function{
		[]VarType{SVar},
		[]string{"size"},
		func(vals []*Value) (*Value, error) {
			if !vals[0].SValue.IsWhole() {
				return nil, fmt.Errorf("size must be an integer")
			}

			n := vals[0].SValue.Integer()

			if n < 0 {
				return nil, fmt.Errorf("size must be positive")
			}

			return &Value{
				VType:  MVar,
				MValue: matrix.Identity(n),
			}, nil
		},
	},

	"ref": function{
		[]VarType{MVar},
		[]string{"mat"},
		func(vals []*Value) (*Value, error) {
			return &Value{
				VType:  MVar,
				MValue: matrix.Ref(vals[0].MValue),
			}, nil
		},
	},

	"rref": function{
		[]VarType{MVar},
		[]string{"mat"},
		func(vals []*Value) (*Value, error) {
			return &Value{
				VType:  MVar,
				MValue: matrix.Rref(vals[0].MValue),
			}, nil
		},
	},
}

func checkFunctionArgs(vals []*Value, fname string, fn function) error {
	if len(vals) != len(fn.signature) {
		return fmt.Errorf("call to %s takes %d arguments, but was supplied %d", fname, len(fn.signature), len(vals))
	}

	correct := true
	for i, vtype := range fn.signature {
		if vtype != vals[i].VType {
			correct = false
			break
		}
	}

	if !correct {
		return fmt.Errorf("call to %s expects arguments (%s) but was supplied (%s)", fname, vartypesToString(fn.signature), valuesVartypesToString(vals))
	}

	return nil
}

func vartypesToString(vtypes []VarType) string {
	strs := []string{}
	for _, vtype := range vtypes {
		var s string
		switch vtype {
		case MVar:
			s = "matrix"
		case SVar:
			s = "scalar"
		}

		strs = append(strs, s)
	}

	return strings.Join(strs, ", ")
}

func valuesVartypesToString(vals []*Value) string {
	strs := []string{}
	for _, val := range vals {
		var s string
		switch val.VType {
		case MVar:
			s = "matrix"
		case SVar:
			s = "scalar"
		}

		strs = append(strs, s)
	}

	return strings.Join(strs, ", ")
}
