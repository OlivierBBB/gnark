package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// evalTest is tasked with computing a polynomial evaluation P(r)
func TestPolyEval(t *testing.T) {

	assert := groth16.NewAssert(t)

	var e PolyEvalCircuit
	// we need to fix the length of e.Coefficients
	e.Coefficients = make([]frontend.Variable, 3)

	r1cs, err := frontend.Compile(gurvy.BN256, &e)

	assert.NoError(err)

	{
		var p PolyEvalCircuit

		p.Coefficients = make([]frontend.Variable, 3)
		// p = X^2 + 2X + 1 = (X+1)^2
		p.Coefficients[0].Assign(1)
		p.Coefficients[1].Assign(2)
		p.Coefficients[2].Assign(1)

		p.VarValue.Assign(3)
		p.Claimed.Assign(16)

		assert.ProverSucceeded(r1cs, &p)
	}
}
