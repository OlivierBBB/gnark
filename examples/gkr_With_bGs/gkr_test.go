package gkr

import (
	"fmt"
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func TestGKR(t *testing.T) {

	var gkr FullGKRCircuit

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
	// fmt.Printf("\nNumber of Constraints: %v\n", r1cs.GetNbConstraints())
	assert.NoError(err)

	{
		var witness FullGKRCircuit

		// fix size of witness.HRPolynomials and initialize values
		for l := range witness.HRPolynomials {
			witness.HRPolynomials[l].Coefficients = make([]frontend.Variable, degHR+1)
			for i := range witness.HRPolynomials[l].Coefficients {
				witness.HRPolynomials[l].Coefficients[i].Assign(0)
			}
		}

		// fix size of witness.HLPolynomials and initialize values
		for l := range witness.HLPolynomials {
			witness.HLPolynomials[l].Coefficients = make([]frontend.Variable, degHL+1)
			for i := range witness.HLPolynomials[l].Coefficients {
				witness.HLPolynomials[l].Coefficients[i].Assign(0)
			}
		}

		// fix size of witness.HPrimePolynomials and initialize values
		for r := range witness.HPrimePolynomials {
			for l := range witness.HPrimePolynomials[r] {
				witness.HPrimePolynomials[r][l].Coefficients = make([]frontend.Variable, degHPrime+1)
				for i := range witness.HPrimePolynomials[r][l].Coefficients {
					witness.HPrimePolynomials[r][l].Coefficients[i].Assign(0)
				}
			}
		}

		// Initialize the remaining stuff
		witness.QInitial.Assign(0)
		for i := range witness.QPrimeInitial {
			witness.QPrimeInitial[i].Assign(0)
		}
		for i := range witness.RoundConstants {
			witness.RoundConstants[i].Assign(0)
		}
		for i := range witness.VInput.Table {
			witness.VInput.Table[i].Assign(0)
		}
		for i := range witness.VOutput.Table {
			witness.VOutput.Table[i].Assign(0)
		}
		for i := range witness.VLClaimed {
			witness.VLClaimed[i].Assign(0)
		}
		for i := range witness.VRClaimed {
			witness.VRClaimed[i].Assign(0)
		}

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
	}
}

// J'aimerais faire un calcul du nombre de contraintes en fonction de bN ...
// mais pour l'instant bN est une constante ...
func TestNumberConstraints(t *testing.T) {

	var gkr FullGKRCircuit

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
	assert.NoError(err)
}

// BenchmarkGKR benchmarks the GKR prover & verifier
func BenchmarkGKR(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()

	var witness FullGKRCircuit

	// fix size of gkr.HRPolynomials
	for i := range witness.HRPolynomials {
		// hR has degree 2
		witness.HRPolynomials[i].Coefficients = make([]frontend.Variable, degHR+1)
	}

	// fix size of witness.HLPolynomials
	for i := range witness.HLPolynomials {
		// hL has degree 8
		witness.HLPolynomials[i].Coefficients = make([]frontend.Variable, degHL+1)
	}

	// fix size of witness.HPrimePolynomials
	for l := range witness.HPrimePolynomials {
		for i := range witness.HPrimePolynomials[l] {
			// h' has degree 8
			witness.HPrimePolynomials[l][i].Coefficients = make([]frontend.Variable, degHPrime+1)
		}
	}

	r1cs, _ := frontend.Compile(gurvy.BN256, &witness)
	pk, vk := groth16.Setup(r1cs)

	// fix size of witness.HRPolynomials and initialize values
	for l := range witness.HRPolynomials {
		// witness.HRPolynomials[l].Coefficients = make([]frontend.Variable, degHR+1)
		for i := range witness.HRPolynomials[l].Coefficients {
			witness.HRPolynomials[l].Coefficients[i].Assign(0)
		}
	}

	// fix size of witness.HLPolynomials and initialize values
	for l := range witness.HLPolynomials {
		// witness.HLPolynomials[l].Coefficients = make([]frontend.Variable, degHL+1)
		for i := range witness.HLPolynomials[l].Coefficients {
			witness.HLPolynomials[l].Coefficients[i].Assign(0)
		}
	}

	// fix size of witness.HPrimePolynomials and initialize values
	for r := range witness.HPrimePolynomials {
		for l := range witness.HPrimePolynomials[r] {
			// witness.HPrimePolynomials[r][l].Coefficients = make([]frontend.Variable, degHPrime+1)
			for i := range witness.HPrimePolynomials[r][l].Coefficients {
				witness.HPrimePolynomials[r][l].Coefficients[i].Assign(0)
			}
		}
	}

	// Initialize the remaining stuff
	witness.QInitial.Assign(0)
	for i := range witness.QPrimeInitial {
		witness.QPrimeInitial[i].Assign(0)
	}
	for i := range witness.RoundConstants {
		witness.RoundConstants[i].Assign(0)
	}
	for i := range witness.VInput.Table {
		witness.VInput.Table[i].Assign(0)
	}
	for i := range witness.VOutput.Table {
		witness.VOutput.Table[i].Assign(0)
	}
	for i := range witness.VLClaimed {
		witness.VLClaimed[i].Assign(0)
	}
	for i := range witness.VRClaimed {
		witness.VRClaimed[i].Assign(0)
	}

	for _c := 0; _c < b.N; _c++ {

		b.StartTimer()
		proof, _ := groth16.Prove(r1cs, pk, witness)
		b.StopTimer()

		groth16.Verify(proof, vk, witness)
	}
}
