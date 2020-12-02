package main

import (
	"testing"

	"github.com/consensys/gurvy"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

func TestPreimage(t *testing.T) {
	assert := groth16.NewAssert(t)

	var mimcCircuit MiMCCircuit

	r1cs, err := frontend.Compile(gurvy.BN256, &mimcCircuit)
	assert.NoError(err)

	{
		var witness MiMCCircuit
		witness.Hash.Assign(42)
		witness.PreImage.Assign(42)
		assert.ProverFailed(r1cs, &witness)
	}

	{
		var witness MiMCCircuit
		witness.PreImage.Assign(35)
		witness.Hash.Assign("19226210204356004706765360050059680583735587569269469539941275797408975356275")
		assert.ProverSucceeded(r1cs, &witness)
	}

}

func TestMiMCOfSlice(t *testing.T) {
	assert := groth16.NewAssert(t)

	var c MiMCOfSliceCircuit

	c.PreImage = make([]frontend.Variable, 5)

	r1cs, err := frontend.Compile(gurvy.BN256, &c)
	assert.NoError(err)

	{
		var witness MiMCOfSliceCircuit
		witness.Hash.Assign(42)
		witness.PreImage = make([]frontend.Variable, 5)
		for i := range witness.PreImage {
			witness.PreImage[i].Assign(i)
		}
		assert.ProverFailed(r1cs, &witness)
	}

	{
		var witness MiMCOfSliceCircuit
		witness.Hash.Assign("4963636333172142546393422809563482779238918718924046481417282740564258588384")
		witness.PreImage = make([]frontend.Variable, 5)
		for i := range witness.PreImage {
			witness.PreImage[i].Assign(1234567890 + i)
		}
		assert.ProverSucceeded(r1cs, &witness)
	}

}
