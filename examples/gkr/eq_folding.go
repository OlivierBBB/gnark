package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// EqFoldingCircuit contains the data to fold an Eq Table completely
// along with the value obtained through folding.
type EqFoldingCircuit struct {
	QPrime  []frontend.Variable `gnark:",public"`
	HPrime  []frontend.Variable `gnark:",public"`
	EqValue frontend.Variable   `gnark:",public"`
}

// Define declares the circuit constraints of a single folding circuit c
func (e *EqFoldingCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	// initialize res to 1
	res := cs.Constant(1)

	// multiply all the Eq's into res
	for i := range e.QPrime {
		term := cs.Mul(e.QPrime[i], e.HPrime[i])
		term = cs.Add(1, term, term)
		term = cs.Sub(term, e.QPrime[i])
		term = cs.Sub(term, e.HPrime[i])
		res = cs.Mul(res, term)
	}

	e.EqValue.Assign(res)

	return nil
}
