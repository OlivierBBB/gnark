package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gurvy"
)

// SumcheckCircuit contains the circuit data for the verification of intermediate sumcheck runs.
type SumcheckCircuit struct {
	InitialClaim         frontend.Variable       `gnark:"c,public"` // initial InitialClaim
	HRPoly               []frontend.Variable     // degree = 2 wrt hR -> 3 coefficients; also bG = 1 so only one poly required
	HLPoly               []frontend.Variable     // degree = 8 wrt hL -> 9 coefficients; also bG = 1 so only one poly required
	HPrimePolys          [bN][]frontend.Variable // degree = 8 wrt h' -> 9 coefficients; bN polynomials required
	VRClaimed            frontend.Variable       // for the final verification
	VLClaimed            frontend.Variable       // for the final verification
	Alpha                frontend.Variable       // for the next round
	Beta                 frontend.Variable       // for the next round
	CopyTablePrefolded   [1 << (2 * bG)]frontend.Variable
	CipherTablePrefolded [1 << (2 * bG)]frontend.Variable
	QPrimeCurr           []frontend.Variable // to compute Eq(QPrimeCurr, QPrimeNext); of size bN
	QRNext               frontend.Variable
	QLNext               frontend.Variable
	QPrimeNext           []frontend.Variable // of size bN
}

// Define contains the circuit data for an NRounds long sumcheck verfier
func (c *SumcheckCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem, FinalSumcheckRun bool) error {

	mimc, _ := mimc.NewMiMC("seed", curveID)

	var r, eval frontend.Variable
	r.Assign(0)
	eval.Assign(0)

	// Get current claim
	claimCurr := c.InitialClaim

	// Elimination of hR:
	eval = cs.Add(c.HRPoly[0], c.HRPoly[0], c.HRPoly[1:]) // eval = P_0(0) + P_0(1)
	cs.AssertIsEqual(eval, claimCurr)                     // claim == eval
	r = mimc.Hash(cs, c.HRPoly...)                        // Hash the polynomial
	c.QRNext.Assign(r)                                    // get qR for the next Sumcheck protocol

	// Get current claim
	p := PolyEvalCircuit{X: r, Coefficients: c.HRPoly}
	p.Define(curveID, cs) // compute p(r)
	claimCurr = p.Value   // set it as current claim

	// Elimination of hL:
	eval = cs.Add(c.HLPoly[0], c.HLPoly[0], c.HLPoly[1:]) // eval = P_1(0) + P_1(1)
	cs.AssertIsEqual(eval, claimCurr)                     // claim == eval
	r = mimc.Hash(cs, c.HLPoly...)                        // Hash the polynomial
	c.QLNext.Assign(r)                                    // get qL for the next Sumcheck protocol

	// Get current claim
	p = PolyEvalCircuit{X: r, Coefficients: c.HRPoly}
	p.Define(curveID, cs) // compute p(r)
	claimCurr = p.Value   // set it as current claim

	for round := 0; round < bN; round++ {

		// elimination of h'_round:
		eval = cs.Add(c.HPrimePolys[round][0], c.HPrimePolys[round][0], c.HPrimePolys[round][1:])
		cs.AssertIsEqual(eval, claimCurr)
		r = mimc.Hash(cs, c.HPrimePolys[round]...)
		c.QPrimeNext[round] = r

		// Get current claim
		p = PolyEvalCircuit{X: r, Coefficients: c.HPrimePolys[round]}
		p.Define(curveID, cs)
		claimCurr = p.Value
	}

	eq := EqFoldingCircuit{QPrime: c.QPrimeCurr, HPrime: c.QPrimeNext}
	eq.Define(curveID, cs) // compute Eq(q', h')

	if !FinalSumcheckRun {
		value := eq.EqValue
		// compute this shit
		cs.AssertIsEqual(claimCurr, eq.EqValue)
	} else {

	}
	// at this point: currInitialClaim = P(r_1, r_2, ... , r_N)
	c.FinalFold.InitialClaimed = currInitialClaim
	c.FinalFold.Define(curveID, cs)

	return nil
}
