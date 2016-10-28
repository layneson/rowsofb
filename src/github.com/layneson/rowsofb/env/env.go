package env

import (
	"fmt"
	"strings"

	"github.com/layneson/rowsofb/matrix"
)

var varnames = strings.Split("A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z", ",")

func getVarOffset(s string) int {
	for i, v := range varnames {
		if s == v {
			return i
		}
	}

	return -1
}

//mvar represents a matrix variable value. It either exists or doesn't.
type mvar struct {
	exists bool // defaults to false

	value matrix.M
}

//E represents an environment with 26 matrix variables, A-Z, and a result variable.
type E struct {
	res mvar

	vars []mvar
}

//New initializes an enviroment and returns a pointer to it.
func New() *E {
	e := &E{}

	e.res = mvar{exists: true, value: matrix.New(3, 3)} // initialize result with a 3x3 zero matrix

	e.vars = make([]mvar, len(varnames))

	for i := range e.vars {
		e.vars[i] = mvar{false, matrix.M{}}
	}

	return e
}

//IsDefined returns true if the variable in v has a value in the environment.
//It returns an error if the variable is unrecognized.
func (e *E) IsDefined(v string) (bool, error) {
	if v == "$" {
		return true, nil
	}

	voff := getVarOffset(v)
	if voff < 0 {
		return false, InvalidVariableError{v}
	}

	return e.vars[voff].exists, nil
}

//Get returns the value of the matrix at the given variable name.
//It returns an error if the given variable name is invalid or if the matrix is undefined.
func (e *E) Get(v string) (matrix.M, error) {
	if v == "$" {
		return e.res.value, nil
	}

	voff := getVarOffset(v)
	if voff < 0 {
		return matrix.M{}, InvalidVariableError{v}
	}

	if !e.vars[voff].exists {
		return matrix.M{}, UndefinedVariableError{v}
	}

	return e.vars[voff].value, nil
}

//Set sets the value at the given variable to the given matrix.
//It returns an error if no such variable exists.
func (e *E) Set(v string, m matrix.M) error {
	if v == "$" {
		return nil
	}

	m = matrix.CopyMatrix(m) // enforce copy

	voff := getVarOffset(v)
	if voff < 0 {
		return InvalidVariableError{v}
	}

	e.vars[voff].exists = true
	e.vars[voff].value = m

	return nil
}

//GetResult returns the value of the result variable.
func (e *E) GetResult() matrix.M {
	return e.res.value
}

//SetResult sets the value of the result variable.
func (e *E) SetResult(m matrix.M) {
	e.res.value = m
}

//Delete sets a variable to be undefined.
func (e *E) Delete(v string) error {
	voff := getVarOffset(v)
	if voff < 0 {
		return InvalidVariableError{v}
	}

	e.vars[voff].exists = false

	return nil
}

//Clear deletes all variables and resets the result variable to the 3x3 zero matrix.
func (e *E) Clear() {
	e.res.value = matrix.New(3, 3)

	for i := range varnames {
		e.vars[i].exists = false
	}
}

//InvalidVariableError represents an error where an invalid variable name was passed to an env function.
type InvalidVariableError struct {
	vname string
}

func (i InvalidVariableError) Error() string {
	return fmt.Sprintf("%q is not a valid variable name", i.vname)
}

//UndefinedVariableError represents an error where an undefined variable is accessed.
type UndefinedVariableError struct {
	vname string
}

func (u UndefinedVariableError) Error() string {
	return fmt.Sprintf("%q is undefined", u.vname)
}
