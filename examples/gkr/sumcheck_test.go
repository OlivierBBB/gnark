package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestSumcheck(t *testing.T) {
	assert := groth16.NewAssert(t)

	var f SingleFold
	var verifierCircuit SumcheckVerfier

	verifierCircuit.finalFold = f

	_, err := frontend.Compile(gurvy.BN256, &verifierCircuit)
	assert.NoError(err)
}
