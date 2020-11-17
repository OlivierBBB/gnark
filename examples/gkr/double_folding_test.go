package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestDoubleFolding(t *testing.T) {

	assert := groth16.NewAssert(t)

	var df DoubleFoldingCircuit

	_, err := frontend.Compile(gurvy.BN256, &df)

	assert.NoError(err)
}
