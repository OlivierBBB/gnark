package gkr

import (
	"fmt"
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestGKR(t *testing.T) {

	var gkr CircuitGKR

	// set size of gkr.HLPolynomials
	for layer := range gkr.HLPolynomials {
		gkr.HLPolynomials[layer].Coefficients = make([]frontend.Variable, degHL+1)
	}

	// set size of gkr.HRPolynomials
	for layer := range gkr.HRPolynomials {
		gkr.HRPolynomials[layer].Coefficients = make([]frontend.Variable, degHR+1)
	}

	// set size of gkr.HPrimePolynomials
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

		// Set public values
		witness.setPublicInputs()

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

		// go test -timeout 5m -run ^TestGKR$ github.com/consensys/gnark/examples/gkr
	}
}

// J'aimerais faire un calcul du nombre de contraintes en fonction de bN ...
// mais pour l'instant bN est une constante ...
func TestNumberConstraints(t *testing.T) {

	var gkr CircuitGKR

	// for _, poly := range gkr.HRPolynomials (etc ...) will produce an error.

	// fix size of gkr.HRPolynomials
	for i := range gkr.HRPolynomials {
		// hR has degree 2
		gkr.HRPolynomials[i].Coefficients = make([]frontend.Variable, degHR+1)
	}

	// fix size of gkr.HLPolynomials
	for i := range gkr.HLPolynomials {
		// hL has degree 8
		gkr.HLPolynomials[i].Coefficients = make([]frontend.Variable, degHL+1)
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
	fmt.Printf("bN = %v,\nNumber of Constraints %v\n", bN, r1cs.GetNbConstraints())
	assert.Error(err)
}

// BenchmarkGKR benchmarks the GKR prover & verifier
func BenchmarkGKR(b *testing.B) {

	var gkr CircuitGKR

	// set size of gkr.HLPolynomials
	for layer := range gkr.HLPolynomials {
		gkr.HLPolynomials[layer].Coefficients = make([]frontend.Variable, degHL+1)
	}

	// set size of gkr.HRPolynomials
	for layer := range gkr.HRPolynomials {
		gkr.HRPolynomials[layer].Coefficients = make([]frontend.Variable, degHR+1)
	}

	// set size of gkr.HPrimePolynomials
	for layer := range gkr.HPrimePolynomials {
		for varIndex := range gkr.HPrimePolynomials[layer] {
			gkr.HPrimePolynomials[layer][varIndex].Coefficients = make([]frontend.Variable, degHPrime+1)
		}
	}

	r1cs, _ := frontend.Compile(gurvy.BN256, &gkr)

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

	// Set public values
	witness.setPublicInputs()

	// Prover:
	// assert.ProverSucceeded(r1cs, &witness)
	// assert.ProverFailed(r1cs, &witness)

	// Solver:
	// assert.SolvingSucceeded(r1cs, &witness)

	pk := groth16.DummySetup(r1cs)
	// proof, err := groth16.Prove(r1cs, pk, witness)
	// assert.NoError(err)
	// err = groth16.Verify(proof, vk, witness)
	// assert.NoError(err)

	// go test -timeout 5m -run ^TestGKR$ github.com/consensys/gnark/examples/gkr

	b.ResetTimer()

	for _c := 0; _c < b.N; _c++ {

		// b.StartTimer()
		_, err := groth16.Prove(r1cs, pk, &witness)
		if err != nil {
			b.Fatal(err)
		}
		// b.StopTimer()

		// groth16.Verify(proof, vk, witness)
	}
}
