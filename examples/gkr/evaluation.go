package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// PolyEvalCircuit describes a circuit to evaluate a
// polynomial and compare the value against a claimed
// value.
type PolyEvalCircuit struct {
	X            frontend.Variable
	Coefficients []frontend.Variable // [a0, a1, ... , ad] <-> P = a0 + a1X + ... + adX^d
	Value        frontend.Variable
}

// Define defines the circuit's constraints
// Claim == P(0) + P(1)
func (c *PolyEvalCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	res := cs.Constant(0)

	for i := len(c.Coefficients) - 1; i >= 0; i-- {
		if i != len(c.Coefficients)-1 {
			res = cs.Mul(res, c.X)
		}
		res = cs.Add(res, c.Coefficients[i])
	}

	c.Value.Assign(res)

	return nil
}
