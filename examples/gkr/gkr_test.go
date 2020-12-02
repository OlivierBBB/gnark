package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestGKR(t *testing.T) {

	var gkr CircuitGKR

	// fix the size of gkr.HLPolynomials
	for layer := range gkr.HLPolynomials {
		gkr.HLPolynomials[layer].Coefficients = make([]frontend.Variable, degHL+1)
	}

	// fix the size of gkr.HRPolynomials
	for layer := range gkr.HRPolynomials {
		gkr.HRPolynomials[layer].Coefficients = make([]frontend.Variable, degHR+1)
	}

	// fix the size of gkr.HPrimePolynomials
	for layer := range gkr.HPrimePolynomials {
		for varIndex := range gkr.HPrimePolynomials[layer] {
			gkr.HPrimePolynomials[layer][varIndex].Coefficients = make([]frontend.Variable, degHPrime+1)
		}
	}

	assert := groth16.NewAssert(t)
	r1cs, err := frontend.Compile(gurvy.BN256, &gkr)
	// fmt.Printf("\nNumber of Constraints: %v\n", r1cs.GetNbConstraints())
	assert.NoError(err)

	{
		var witness CircuitGKR

		// fix the size of witness.HLPolynomials and initialize values
		for l := range witness.HLPolynomials {
			witness.HLPolynomials[l].Coefficients = make([]frontend.Variable, degHL+1)
		}

		// fix the size of witness.HRPolynomials and initialize values
		for l := range witness.HRPolynomials {
			witness.HRPolynomials[l].Coefficients = make([]frontend.Variable, degHR+1)
		}

		// fix the size of witness.HPrimePolynomials and initialize values
		for r := range witness.HPrimePolynomials {
			for l := range witness.HPrimePolynomials[r] {
				witness.HPrimePolynomials[r][l].Coefficients = make([]frontend.Variable, degHPrime+1)
			}
		}

		// Insert all public values into the witness
		witness.setup()

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
