package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestSumcheck(t *testing.T) {
	assert := groth16.NewAssert(t)

	degrees := [2]int{2, 3}
	var claim frontend.Variable
	sumcheckVerifierCircuit := NewSumcheckVerifier(claim, degrees)

	// var sumcheckVerifierCircuit SumcheckVerifier

	_, err := frontend.Compile(gurvy.BN256, &sumcheckVerifierCircuit)

	assert.NoError(err)
}
