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
	HLPoly       Polynomial     `gnark:",public"` // deg = 2 wrt hL => 3 coeffs; bG = 1 => only one poly required
	HRPoly       Polynomial     `gnark:",public"` // deg = 8 wrt hR => 9 coeffs; bG = 1 => only one poly required
	HPrimePolys  [bN]Polynomial `gnark:",public"` // deg = 8 wrt h' => 9 coeffs; bN polys required
}

// Solve verifies a sumcheck instance EXCEPT FOR THE FINAL VERIFICATION.
func (sc *Sumcheck) Solve(curveID gurvy.ID, cs *frontend.ConstraintSystem, mimc *mimc.MiMC) (
	hL, hR frontend.Variable,
	hPrime [bN]frontend.Variable,
	lastClaim frontend.Variable,
) {
	// initialize current claim:
	claimCurr := sc.InitialClaim

	// Elimination of hL:
	zeroAndOne := sc.HLPoly.zeroAndOne(cs)
	cs.AssertIsEqual(zeroAndOne, claimCurr)       // claim == HLPoly(0) + HLPoly(1)
	hL = mimc.Hash(cs, sc.HLPoly.Coefficients...) // Hash the polynomial
	claimCurr = sc.HLPoly.eval(cs, hL)            // Get new current claim

	// Elimination of hR:
	zeroAndOne = sc.HRPoly.zeroAndOne(cs)
	cs.AssertIsEqual(zeroAndOne, claimCurr)       // claim == HRPoly(0) + HRPoly(1)
	hR = mimc.Hash(cs, sc.HRPoly.Coefficients...) // Hash the polynomial
	claimCurr = sc.HRPoly.eval(cs, hR)            // Get new current claim

	// elimination of the variables in h':
	for round := 0; round < bN; round++ {

		// elimination of h'_round:
		zeroAndOne = sc.HPrimePolys[round].zeroAndOne(cs)
		cs.AssertIsEqual(zeroAndOne, claimCurr)                              // claim == PCurr(0) + PCurr(1) where PCurr = HPrimePolys[round]
		hPrime[round] = mimc.Hash(cs, sc.HPrimePolys[round].Coefficients...) // Hash the polynomial
		claimCurr = sc.HPrimePolys[round].eval(cs, hPrime[round])            // Get new current claim
	}

	lastClaim = claimCurr

	return hL, hR, hPrime, lastClaim
}

// Combinator combines the previously computed folded values of Eq, Copy, Cipher
// and the two foldings (VR & VL) of V_i into the evalution of the polynomial
// being summed over.
func Combinator(cs *frontend.ConstraintSystem, eq, copy, cipher, VL, VR, roundConstant frontend.Variable) (computedClaim frontend.Variable) {

	computedClaim = cs.Add(VR, roundConstant)     // VR + C
	aux := cs.Mul(computedClaim, computedClaim)   // (VR + C)^2
	computedClaim = cs.Mul(computedClaim, aux)    // (VR + C)^3
	aux = cs.Mul(aux, aux)                        // (VR + C)^4
	computedClaim = cs.Mul(computedClaim, aux)    // (VR + C)^7
	computedClaim = cs.Add(computedClaim, VL)     // VL + (VR + C)^7
	computedClaim = cs.Mul(computedClaim, cipher) // cipher * (VL + (VR+C)^7)
	aux = cs.Mul(copy, VL)                        // copy * VL
	computedClaim = cs.Add(computedClaim, aux)    // [ copy * VL + cipher * (VL + (VR+C)^7) ]
	computedClaim = cs.Mul(computedClaim, eq)     // eq * [ copy * VL + cipher * (VL + (VR+C)^7) ]

	return computedClaim
}
