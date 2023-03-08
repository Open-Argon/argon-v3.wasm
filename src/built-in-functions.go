package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

type builtinFunc struct {
	name string
	FUNC func(...any) (any, ArErr)
}

func ArgonMult(args ...any) (any, ArErr) {
	return reduce(func(x any, y any) any {
		return newNumber().Mul(y.(number), x.(number))
	}, args), ArErr{}
}

func ArgonInput(args ...any) (any, ArErr) {
	output := []any{}
	for i := 0; i < len(args); i++ {
		output = append(output, anyToArgon(args[i], false, true, 3, 0))
	}
	fmt.Print(output...)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	return input, ArErr{}
}

func ArgonNumber(args ...any) (any, ArErr) {
	if len(args) == 0 {
		return newNumber(), ArErr{}
	}
	switch x := args[0].(type) {
	case string:
		if !numberCompile.MatchString(x) {
			return nil, ArErr{TYPE: "Number Error", message: "Cannot convert type '" + x + "' to a number", EXISTS: true}
		}
		N, _ := newNumber().SetString(x)
		return N, ArErr{}
	case number:
		return x, ArErr{}
	case bool:
		if x {
			return newNumber().SetInt64(1), ArErr{}
		}
		return newNumber().SetInt64(0), ArErr{}
	case nil:
		return newNumber(), ArErr{}
	}

	return nil, ArErr{TYPE: "Number Error", message: "Cannot convert " + typeof(args[0]) + " to a number", EXISTS: true}
}

func ArgonSqrt(a ...any) (any, ArErr) {
	if len(a) == 0 {
		return nil, ArErr{TYPE: "sqrt", message: "sqrt takes 1 argument",
			EXISTS: true}
	}
	if typeof(a[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error", message: "sqrt takes a number not a '" + typeof(a[0]) + "'",
			EXISTS: true}
	}

	r := a[0].(number)

	if r.Sign() < 0 {
		return nil, ArErr{TYPE: "sqrt", message: "sqrt takes a positive number",
			EXISTS: true}
	}

	var x big.Float
	x.SetPrec(30) // I didn't figure out the 'Prec' part correctly, read the docs more carefully than I did and experiement
	x.SetRat(r)

	var s big.Float
	s.SetPrec(15)
	s.Sqrt(&x)

	r, _ = s.Rat(nil)
	return r, ArErr{}
}