package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestSumcheck(t *testing.T) {
	assert := groth16.NewAssert(t)

	var sumcheckVerifierCircuit SumcheckCircuit

	// var sumcheckVerifierCircuit SumcheckVerifier

	r1cs, err := frontend.Compile(gurvy.BN256, &sumcheckVerifierCircuit)

	{
		degrees := [2]int{2, 3}
		var claim frontend.Variable
		s := NewSumcheckCircuit(claim, degrees)

		assert.ProverSucceeded(r1cs, &s)
	}

	assert.NoError(err)
}
