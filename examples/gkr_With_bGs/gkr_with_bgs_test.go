package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestGKRWithBGs(t *testing.T) {

	var gkr FullGKRWithBGsCircuit

	// fix size of gkr.HLPolynomials
	for i := range gkr.HLPolynomials {
		// hL has degree 2
		gkr.HLPolynomials[i].Coefficients = make([]frontend.Variable, degHL+1)
	}

	// fix size of gkr.HRPolynomials
	for i := range gkr.HRPolynomials {
		// hR has degree 8
		gkr.HRPolynomials[i].Coefficients = make([]frontend.Variable, degHR+1)
	}

	// fix size of gkr.HPrimePolynomials
	for l := range gkr.HPrimePolynomials {
		for i := range gkr.HPrimePolynomials[l] {
			// h' has degree 8
			gkr.HPrimePolynomials[l][i].Coefficients = make([]frontend.Variable, degHPrime+1)
		}
	}

	assert := groth16.NewAssert(t)
	r1cs, err := frontend.Compile(gurvy.BN256, &gkr)
	// fmt.Printf("\nNumber of Constraints: %v\n", r1cs.GetNbConstraints())
	assert.NoError(err)

	{
		var witness FullGKRWithBGsCircuit

		// fix size of witness.HLPolynomials and initialize values
		for l := range witness.HLPolynomials {
			witness.HLPolynomials[l].Coefficients = make([]frontend.Variable, degHL+1)
		}

		// fix size of witness.HRPolynomials and initialize values
		for l := range witness.HRPolynomials {
			witness.HRPolynomials[l].Coefficients = make([]frontend.Variable, degHR+1)
		}

		// fix size of witness.HPrimePolynomials and initialize values
		for r := range witness.HPrimePolynomials {
			for l := range witness.HPrimePolynomials[r] {
				witness.HPrimePolynomials[r][l].Coefficients = make([]frontend.Variable, degHPrime+1)
			}
		}

		witness.SetRoundConstants()
		witness.setInputs()
		witness.setOutputs()
		witness.setVLAndVR()
		witness.setPolynomials()
		witness.setQInitial()

		// premier hash: 14596316904690824680289447726830409556807701064825759139913524429557892625

		// Prover:
		// assert.ProverSucceeded(r1cs, &witness)
		// assert.ProverFailed(r1cs, &witness)

		// Solver:
		assert.SolvingSucceeded(r1cs, &witness)

		// pk, vk := groth16.Setup(r1cs)
		// proof, err := groth16.Prove(r1cs, pk, witness)
		// assert.NoError(err)
		// err = groth16.Verify(proof, vk, witness)
		// assert.NoError(err)

		// go test -timeout 30s -run ^TestGKRWithBGs$ github.com/consensys/gnark/examples/gkr_With_bGs
	}
}
