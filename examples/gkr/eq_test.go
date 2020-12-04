package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

type EqCircuit struct {
	Eq          Eq
	ActualValue frontend.Variable
}

func (eq *EqCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	foldedEqValue := eq.Eq.Fold(cs)

	cs.AssertIsEqual(eq.ActualValue, foldedEqValue)

	return nil
}

// for this test to work set bN = 1
func TestEq(t *testing.T) {

	var eq EqCircuit

	r1cs, err := frontend.Compile(gurvy.BN256, &eq)

	assert := groth16.NewAssert(t)

	assert.NoError(err)

	{
		var witness EqCircuit

		var x, y int = 13, 88
		e := 1 - x - y + 2*x*y
		witness.Eq.QPrime[0].Assign(x)
		witness.Eq.HPrime[0].Assign(y)

		witness.ActualValue.Assign(e)

		assert.ProverSucceeded(r1cs, &witness)
	}
}
