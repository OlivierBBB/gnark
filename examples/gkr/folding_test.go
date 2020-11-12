package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestFolding(t *testing.T) {
	assert := groth16.NewAssert(t)

	var f FoldingCircuit
	_, err := frontend.Compile(gurvy.BN256, &f)

	assert.NoError(err)
}
