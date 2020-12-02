package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

type ValuesBKTCircuit struct {
	Values   OutputValuesBKT
	Expected frontend.Variable
}

const (
	q      = 12348
	qprime = 721647
)

func (g *ValuesBKTCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	// Q := cs.Constant(q)
	QPrime := [bN]frontend.Variable{cs.Constant(qprime)}

	actual := g.Values.SingleFold(cs, QPrime)
	cs.AssertIsEqual(g.Expected, actual)
	return nil
}

func TestValueBKT(t *testing.T) {

	var g ValuesBKTCircuit

	r1cs, err := frontend.Compile(gurvy.BN256, &g)

	assert := groth16.NewAssert(t)

	assert.NoError(err)

	{
		var witness ValuesBKTCircuit

		var t0, t1, t2, t3 int = 1249, 4173, 871, 6541
		witness.Values.Table[0].Assign(t0)
		witness.Values.Table[1].Assign(t1)
		// witness.Values.Table[2].Assign(t2)
		// witness.Values.Table[3].Assign(t3)

		v0 := t0 + q*(t2-t0)
		v1 := t1 + q*(t3-t1)
		e := v0 + qprime*(v1-v0)
		witness.Expected.Assign(e)

		assert.ProverSucceeded(r1cs, &witness)
	}
}
