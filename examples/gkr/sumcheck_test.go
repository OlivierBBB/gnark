package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gurvy"
)

type SumcheckCircuit struct {
	SC       Sumcheck
	Expected frontend.Variable
}

func (scc *SumcheckCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	mimc, _ := mimc.NewMiMC("seed", curveID)
	hR, hL, hPrime, lastClaim := scc.SC.Solve(curveID, cs, &mimc)
	cs.AssertIsEqual(hL, cs.Constant(0))
	cs.AssertIsEqual(hR, cs.Constant(0))
	cs.AssertIsEqual(hPrime[0], cs.Constant(0))
	cs.AssertIsEqual(lastClaim, cs.Constant(0))

	//cs.AssertIsEqual(combinator(cs, eq, copy, cipher, VL, vR, roundconstant), scc.Expected)
	return nil
}

func TestSumcheckCircuit(t *testing.T) {

	var scc SumcheckCircuit

	// initialize the polynomials for eliminating hL, hR and the bN variables for h'
	scc.SC.HLPoly.Coefficients = make([]frontend.Variable, degHL+1)
	scc.SC.HRPoly.Coefficients = make([]frontend.Variable, degHR+1)
	for varIndex := range scc.SC.HPrimePolys {
		scc.SC.HPrimePolys[varIndex].Coefficients = make([]frontend.Variable, degHPrime+1)
	}

	r1cs, err := frontend.Compile(gurvy.BN256, &scc)

	assert := groth16.NewAssert(t)

	assert.NoError(err)

	{
		var witness SumcheckCircuit

		// initialize the HRPoly and fill it up
		witness.SC.HRPoly.Coefficients = make([]frontend.Variable, degHR+1)
		for i := range witness.SC.HRPoly.Coefficients {
			witness.SC.HRPoly.Coefficients[i].Assign(i)
		}

		// initialize the HLPoly and fill it up
		witness.SC.HLPoly.Coefficients = make([]frontend.Variable, degHL+1)
		for i := range witness.SC.HLPoly.Coefficients {
			witness.SC.HLPoly.Coefficients[i].Assign(i + 100)
		}

		// initialize the HPrimePolys and fill them up
		for varIndex := range witness.SC.HPrimePolys {
			witness.SC.HPrimePolys[varIndex].Coefficients = make([]frontend.Variable, degHPrime+1)
			for i := range witness.SC.HPrimePolys[varIndex].Coefficients {
				witness.SC.HPrimePolys[varIndex].Coefficients[i].Assign(varIndex + i + 200)
			}
		}

		// finish the initialization fo witness.
		witness.SC.InitialClaim.Assign(3)
		witness.Expected.Assign(400)

		assert.ProverSucceeded(r1cs, &witness)
	}
}
