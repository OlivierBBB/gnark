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
	r1cs, err := frontend.Compile(gurvy.BN256, &f)

	assert.NoError(err)

	{
		var g FoldingCircuit

		for i := 0; i < 1<<(bG+bN); i++ {
			g.Table[i].Assign(i)
		}

		r1 := 2487186132541724
		r2 := 1287571234672148

		g.Q[0].Assign(r1)
		g.QPrime[0].Assign(r2)

		g.Claimed.Assign(2*r1 + r2)

		assert.ProverSucceeded(r1cs, &g)
	}
}
