package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gurvy"
)

const (
	bN        = 3
	bG        = 1
	nLayers   = 91
	degHL     = 2
	degHR     = 8
	degHPrime = 8
)

// CircuitGKR contains the circuit data for an nLayers deep GKR circuit.
type CircuitGKR struct {
	QPrimeInitial     [bN]frontend.Variable          `gnark:",public"` // initial randomness (of length bN)
	VLClaimed         [nLayers - 1]frontend.Variable `gnark:",public"` // claimed values of VL for all levels except inputs and outputs
	VRClaimed         [nLayers - 1]frontend.Variable `gnark:",public"` // claimed values of VR for all levels except inputs and outputs
	HLPolynomials     [nLayers]Polynomial            `gnark:",public"` // polynomials for eliminating hL; deg = 2 => 3 coeffs
	HRPolynomials     [nLayers]Polynomial            `gnark:",public"` // polynomials for eliminating hR; deg = 8 => 9 coeffs
	HPrimePolynomials [nLayers][bN]Polynomial        `gnark:",public"` // polynomials for eliminating h'; deg = 8 => 9 coeffs
	VOutput           OutputValuesBKT                `gnark:",public"` // table of outputs
	VInput            InputValuesBKT                 `gnark:",public"` // table of inputs
	RoundConstants    [nLayers]frontend.Variable     `gnark:",public"` // round constants IN REVERSE ORDER
}

// Define declares the circuit constraints of a full GKR circuit
func (gkr *CircuitGKR) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	mimc, _ := mimc.NewMiMC("seed", curveID)

	var (
		initialClaimOfTheSumcheck      frontend.Variable
		qPrimeCurr                     [bN]frontend.Variable
		VL, VR                         frontend.Variable
		rho                            frontend.Variable
		eq, copy, cipher               frontend.Variable
		prefoldedCopy, prefoldedCipher PrefoldedGateBKT
	)

	for round := 0; round < nLayers; round++ {

		// Set up current sumcheck:
		// ========================

		if round == 0 {
			// get the initial claim for the first Sumcheck run
			initialClaimOfTheSumcheck = gkr.VOutput.SingleFold(cs, gkr.QPrimeInitial)
			qPrimeCurr = gkr.QPrimeInitial
			var CipherTable, DummyCopyTable [1 << (2 * bG)]frontend.Variable
			// CipherTable = [0, 0, 1, 0]
			CipherTable[2] = cs.Constant(1)
			for i := range DummyCopyTable {
				if i != 2 {
					CipherTable[i] = cs.Constant(0)
				}
			}
			// there is no cipher table at the top level (and no q!)
			// for simplicity's sake we use a dummy "copy" table set to zero
			for i := range DummyCopyTable {
				DummyCopyTable[i] = cs.Constant(0)
			}
			prefoldedCopy = PrefoldedGateBKT{DummyCopyTable}
			prefoldedCipher = PrefoldedGateBKT{CipherTable}
		}

		// assemble current sumcheck instance
		sc := Sumcheck{
			InitialClaim: initialClaimOfTheSumcheck,
			HLPoly:       gkr.HLPolynomials[round],
			HRPoly:       gkr.HRPolynomials[round],
			HPrimePolys:  gkr.HPrimePolynomials[round],
		}

		// Verify this sumcheck instance EXCEPT FOR THE FINAL VERIFICATION:
		// ================================================================

		hL, hR, hPrime, lastClaimOfThisSumcheck := sc.Solve(curveID, cs, &mimc)

		// Finish the sumcheck verification:
		// =================================

		// fold Eq, prefoldedCopy and prefoldedCipher; get VL and VR
		Eq := Eq{QPrime: qPrimeCurr, HPrime: hPrime}
		eq = Eq.Fold(cs)
		copy = prefoldedCopy.Fold(cs, hL, hR)
		cipher = prefoldedCipher.Fold(cs, hL, hR)
		if round != (nLayers - 1) {
			VL = gkr.VLClaimed[round]
			VR = gkr.VRClaimed[round]
		} else {
			VL, VR = gkr.VInput.DoubleFold(cs, hL, hR, hPrime)
		}
		// compute expected value of the final claim of the current Sumcheck; compare
		expectedClaim := Combinator(cs, eq, copy, cipher, VL, VR, gkr.RoundConstants[round])
		cs.AssertIsEqual(lastClaimOfThisSumcheck, expectedClaim)

		// Prepare the next round:
		// =======================

		if round != (nLayers - 1) {
			rho = mimc.Hash(cs, VL, VR)
			// set the next prefoldedCopy and prefoldedCipher
			prefoldedCopy, prefoldedCipher = PrefoldedCopyAndCipherLinComb(cs, rho, hL, hR)
			// set the next initial claim to VL + rho * VR
			initialClaimOfTheSumcheck = cs.Mul(rho, VR)
			initialClaimOfTheSumcheck = cs.Add(VL, initialClaimOfTheSumcheck)
			// set qPrimeCurr to HPrime
			qPrimeCurr = hPrime
		}
	}

	return nil
}
