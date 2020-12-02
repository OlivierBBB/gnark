package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gurvy"
)

// Sumcheck contains the circuit data of a sumcheck run
// EXCEPT WHAT IS REQUIRED FOR THE FINAL CHECK.
type Sumcheck struct {
	InitialClaim frontend.Variable
	HRPoly       Polynomial     `gnark:",public"` // deg = 2 wrt hR => 3 coeffs; bG = 1 => only one poly required
	HLPoly       Polynomial     `gnark:",public"` // deg = 8 wrt hL => 9 coeffs; bG = 1 => only one poly required
	HPrimePolys  [bN]Polynomial `gnark:",public"` // deg = 8 wrt h' => 9 coeffs; bN polys required
}

// Solve verifies a sumcheck instance EXCEPT FOR THE FINAL VERIFICATION.
func (sc *Sumcheck) Solve(curveID gurvy.ID, cs *frontend.ConstraintSystem, mimc *mimc.MiMC) (
	hR, hL frontend.Variable,
	hPrime [bN]frontend.Variable,
	lastClaim frontend.Variable,
) {
	// initialize current claim:
	claimCurr := sc.InitialClaim

	// Elimination of hR:
	zeroAndOne := sc.HRPoly.zeroAndOne(cs)
	cs.AssertIsEqual(zeroAndOne, claimCurr)       // claim == P_r(0) + P_r(1)
	hR = mimc.Hash(cs, sc.HRPoly.Coefficients...) // Hash the polynomial
	claimCurr = sc.HRPoly.eval(cs, hR)            // Get new current claim

	// Elimination of hL:
	zeroAndOne = sc.HLPoly.zeroAndOne(cs)
	cs.AssertIsEqual(zeroAndOne, claimCurr)       // claim == P_l(0) + P_l(1)
	hL = mimc.Hash(cs, sc.HLPoly.Coefficients...) // Hash the polynomial
	claimCurr = sc.HLPoly.eval(cs, hL)            // Get new current claim

	// elimination of the variables in h':
	for round := 0; round < bN; round++ {

		// elimination of h'_round:
		zeroAndOne = sc.HPrimePolys[round].zeroAndOne(cs)
		cs.AssertIsEqual(zeroAndOne, claimCurr)                              // claim == P_l(0) + P_l(1)
		hPrime[round] = mimc.Hash(cs, sc.HPrimePolys[round].Coefficients...) // Hash the polynomial
		claimCurr = sc.HPrimePolys[round].eval(cs, hPrime[round])            // Get new current claim
	}

	lastClaim = claimCurr

	return hR, hL, hPrime, lastClaim
}

// Combinator combines the previously computed folded values of Eq, Copy, Cipher
// and the two foldings (VL & VR) of V_i into the evalution of the polynomial being summed over.
func Combinator(cs *frontend.ConstraintSystem, eq, copy, cipher, VL, VR, roundConstant frontend.Variable) (computedClaim frontend.Variable) {

	// compute eq * [ copy * VL + cipher * (VR + (VL+C)^7) ]
	computedClaim = cs.Add(VL, roundConstant)     // VL + C
	aux := cs.Mul(computedClaim, computedClaim)   // (VL + C)^2
	computedClaim = cs.Mul(computedClaim, aux)    // (VL + C)^3
	aux = cs.Mul(computedClaim, computedClaim)    // (VL + C)^4
	computedClaim = cs.Mul(computedClaim, aux)    // (VL + C)^7
	computedClaim = cs.Add(computedClaim, VR)     // VR + (VL + C)^7
	computedClaim = cs.Mul(computedClaim, cipher) // cipher * (VR + (VL+C)^7)
	aux = cs.Mul(copy, VL)                        // copy * VL
	computedClaim = cs.Add(computedClaim, aux)    // [ copy * VL + cipher * (VR + (VL+C)^7) ]
	computedClaim = cs.Mul(computedClaim, eq)     // eq * [ copy * VL + cipher * (VR + (VL+C)^7) ]

	return computedClaim
}
