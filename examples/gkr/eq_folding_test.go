package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestEqFold(t *testing.T) {
	assert := groth16.NewAssert(t)

	QPrime := make([]frontend.Variable, 4)
	HPrime := make([]frontend.Variable, 4)

	e := EqFoldingCircuit{
		QPrime: QPrime,
		HPrime: HPrime,
	}

	r1cs, err := frontend.Compile(gurvy.BN256, &e)

	assert.NoError(err)

	{
		var g EqFoldingCircuit

		g.QPrime = make([]frontend.Variable, 4)
		g.HPrime = make([]frontend.Variable, 4)
		g.EqValue.Assign(0)

		for i := 0; i < 4; i++ {
			g.QPrime[i].Assign(i)
			g.HPrime[i].Assign(i + 1)
		}

		assert.ProverSucceeded(r1cs, &g)
	}

	{
		var g EqFoldingCircuit

		g.QPrime = make([]frontend.Variable, 4)
		g.HPrime = make([]frontend.Variable, 4)

		for i := 0; i < 4; i++ {
			g.QPrime[i].Assign(i + 1)
			g.HPrime[i].Assign(i + 1)
		}
		r := 5 * 13 * 25
		g.EqValue.Assign(r)

		assert.ProverSucceeded(r1cs, &g)
	}
}
