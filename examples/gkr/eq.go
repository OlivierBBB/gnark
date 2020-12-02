package gkr

import (
	"github.com/consensys/gnark/frontend"
)

// Eq contains two arrays of values to be folded into an Eq table
type Eq struct {
	QPrime [bN]frontend.Variable
	HPrime [bN]frontend.Variable
}

// MonovariateEqEval computes 1 - q - h + 2 * q * h with q = eq.QPrime[i] and h = eq.HPrime[i]
func (eq *Eq) MonovariateEqEval(cs *frontend.ConstraintSystem, i int) frontend.Variable {

	res := cs.Mul(eq.QPrime[i], eq.HPrime[i])
	res = cs.Add(1, res, res)
	res = cs.Sub(res, eq.QPrime[i])
	res = cs.Sub(res, eq.HPrime[i])

	return res
}

// Fold returns Eq(q', h')
func (eq *Eq) Fold(cs *frontend.ConstraintSystem) frontend.Variable {

	res := cs.Constant(1)

	// multiply all the MonovariateEqEval's into res
	for i := range eq.QPrime {
		res = cs.Mul(res, eq.MonovariateEqEval(cs, i))
	}

	return res
}
