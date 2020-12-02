package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gurvy"
)

const (
	nLayers   = 91
	bN        = 3 // 2^bN hash computations
	bG        = 1 // base circuit breadth = 2 = 2^bG
	degHL     = 2
	degHR     = 8
	degHPrime = 8
)

// FullGKRCircuit contains the circuit data for an nLayers deep GKR circuit.
type FullGKRCircuit struct {
	QInitial          frontend.Variable              `gnark:",public"` // initial randomness, recall bG = 1
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
func (gkr *FullGKRCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	mimc, _ := mimc.NewMiMC("seed", curveID)

	var (
		initialClaimOfTheSumcheck      frontend.Variable
		qPrimeCurr                     [bN]frontend.Variable
		VL, VR                         frontend.Variable
		lambda, rho                    frontend.Variable
		eq, copy, cipher               frontend.Variable
		prefoldedCopy, prefoldedCipher PrefoldedGateBKT
	)

	for round := 0; round < nLayers; round++ {

		if round == 0 {
			// get the initial claim for the first Sumcheck run
			initialClaimOfTheSumcheck = gkr.VOutput.SingleFold(cs, gkr.QPrimeInitial)
			qPrimeCurr = gkr.QPrimeInitial
			prefoldedCopy, prefoldedCipher = PrefoldedCopyAndCipher(cs, gkr.QInitial)
		}

		// constitute current sumcheck instance
		sc := Sumcheck{
			InitialClaim: initialClaimOfTheSumcheck,
			HLPoly:       gkr.HLPolynomials[round],
			HRPoly:       gkr.HRPolynomials[round],
			HPrimePolys:  gkr.HPrimePolynomials[round],
		}

		hL, hR, hPrime, lastClaimOfThisSumcheck := sc.Solve(curveID, cs, &mimc)

		// get eq(q', h'), prefoldedCopy(hL, hR) and prefoldedCipher(hL, hR)
		Eq := Eq{QPrime: qPrimeCurr, HPrime: hPrime}
		eq = Eq.Fold(cs)
		copy = prefoldedCopy.Fold(cs, hR, hL)
		cipher = prefoldedCipher.Fold(cs, hR, hL)

		// get VL and VR
		if round != (nLayers - 1) {
			VL = gkr.VLClaimed[round]
			VR = gkr.VRClaimed[round]
		} else {
			VL, VR = gkr.VInput.DoubleFold(cs, hR, hL, hPrime)
		}

		// compute expected value of the final claim of the current Sumcheck run
		expectedClaim := Combinator(cs, eq, copy, cipher, VL, VR, gkr.RoundConstants[round])

		// compare expectedClaim to the lastClaimOfThisSumcheck
		cs.AssertIsEqual(lastClaimOfThisSumcheck, expectedClaim)

		// get lambda and rho
		lambda = mimc.Hash(cs, VL)
		rho = mimc.Hash(cs, VR)

		// Preparing the next round:
		// =========================

		// set the next prefoldedCopy and prefoldedCipher
		if round != (nLayers - 1) {
			prefoldedCopy, prefoldedCipher = PrefoldedCopyAndCipherLinComb(cs, lambda, rho, hR, hL)
		}

		// The next initial claim is lambda * VL + rho * VR
		aux := cs.Mul(lambda, VL)
		initialClaimOfTheSumcheck = cs.Mul(rho, VR)
		initialClaimOfTheSumcheck = cs.Add(initialClaimOfTheSumcheck, aux)

		// redefine qPrimeCurr as the previously computed HPrime
		qPrimeCurr = hPrime
	}

	return nil
}
