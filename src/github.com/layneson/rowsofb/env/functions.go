package env

import (
	"fmt"
	"strings"

	"github.com/layneson/rowsofb/matrix"
)

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

			return valueFromMatrix(matrix.Identity(n)), nil
		},
	},

	"ref": function{
		[]VarType{MVar},
		[]string{"mat"},
		func(vals []*Value) (*Value, error) {
			return valueFromMatrix(matrix.Ref(vals[0].MValue)), nil
		},
	},

	"rref": function{
		[]VarType{MVar},
		[]string{"mat"},
		func(vals []*Value) (*Value, error) {
			return valueFromMatrix(matrix.Rref(vals[0].MValue)), nil
		},
	},

	"invert": function{
		[]VarType{MVar},
		[]string{"mat"},
		func(vals []*Value) (*Value, error) {
			m, err := matrix.Inverse(vals[0].MValue)
			if err != nil {
				return nil, err
			}

			return valueFromMatrix(m), nil
		},
	},

	"augment": function{
		[]VarType{MVar, MVar},
		[]string{"a", "b"},
		func(vals []*Value) (*Value, error) {
			m, err := matrix.Augment(vals[0].MValue, vals[1].MValue)
			if err != nil {
				return nil, err
			}

			return valueFromMatrix(m), nil
		},
	},
}

func valueFromMatrix(m matrix.M) *Value {
	return &Value{
		VType:  MVar,
		MValue: m,
	}
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
