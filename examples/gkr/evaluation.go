package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// PolyEval describes a circuit to evaluate a
// polynomial and compare the value against a claimed
// value.
type PolyEval struct {
	claimed      frontend.Variable
	varValue     frontend.Variable
	coefficients []frontend.Variable // [a0, a1, ... , ad] <-> P = a0 + a1X + ... + adX^d
}

// Define defines the circuit's constraints
// Claim == P(0) + P(1)
func (circuit *PolyEval) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	for i := 1; i < len(circuit.coefficients); i++ {
		circuit.coefficients[len(circuit.coefficients)-1] = cs.Mul(circuit.coefficients[len(circuit.coefficients)-1], circuit.varValue)
		circuit.coefficients[len(circuit.coefficients)-1] = cs.Add(circuit.coefficients[len(circuit.coefficients)-1-i], circuit.coefficients[len(circuit.coefficients)-1])
	}

	cs.AssertIsEqual(circuit.claimed, circuit.coefficients[len(circuit.coefficients)-1])

	return nil
}
