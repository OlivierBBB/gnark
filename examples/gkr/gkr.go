package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

const (
	nLayers = 91
	bN      = 1
	bG      = 1
)

// FullGKR contains the circuit data for an nLayers deep GKR circuit
// Note: the input folding is not optimized.
type FullGKR struct {
	QInitial      [bG]frontend.Variable                 // initial randomness
	QPrimeInitial [bN]frontend.Variable                 // initial randomness
	VLClaimed     [nLayers - 1]frontend.Variable        // claimed values of VL for all levels except inputs and outputs
	VRClaimed     [nLayers - 1]frontend.Variable        // claimed values of VR for all levels except inputs and outputs
	Polynomials   [nLayers][bG + bN][]frontend.Variable // the polynomials of the sumchecks
	VOutput       [1 << (bG + bN)]frontend.Variable     // table of outputs
	VInput        [1 << (bG + bN)]frontend.Variable     // table of inputs
}

// Define declares the circuit constraints of a full GKR circuit
func (f *FullGKR) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) {

	outputFolding := FoldingCircuit{
		Table:  f.VOutput,
		Q:      f.QInitial,
		QPrime: f.QPrimeInitial,
	}

	// Claim of first round of first Sumcheck
	outputFolding.Define(curveID, cs)

	claim := outputFolding.FoldedValue

	for round := 0; round < nLayers-1; round++ {

		for i := 0; i < bN+bG; i++ {
			eval :=
				cs.AssertIsEqual(claim)
		}
	}
}
