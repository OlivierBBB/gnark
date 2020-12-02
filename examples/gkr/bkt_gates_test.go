package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

type GateBKTCircuit struct {
	PrefoldedBKT PrefoldedGateBKT
	Expected     frontend.Variable
}

func (g *GateBKTCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	hL := cs.Constant(3)
	hR := cs.Constant(5)

	actual := g.PrefoldedBKT.Fold(cs, hL, hR)
	cs.AssertIsEqual(g.Expected, actual)
	return nil
}

func TestGateBKT(t *testing.T) {

	var g GateBKTCircuit

	r1cs, err := frontend.Compile(gurvy.BN256, &g)

	assert := groth16.NewAssert(t)

	assert.NoError(err)

	{
		var witness GateBKTCircuit

		witness.PrefoldedBKT.Table[0].Assign(0)
		witness.PrefoldedBKT.Table[1].Assign(1)
		witness.PrefoldedBKT.Table[2].Assign(2)
		witness.PrefoldedBKT.Table[3].Assign(3)
		witness.Expected.Assign(11)
		assert.ProverSucceeded(r1cs, &witness)
	}
}
