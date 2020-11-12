package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gurvy"
)

// SumcheckVerifier contains the circuit data for the verification of a sumcheck run.
type SumcheckVerifier struct {
	nRounds     int                   // number of rounds
	claim       frontend.Variable     // initial claim
	polynomials [][]frontend.Variable // polynomials[r][i] = coeff of X^i in P_r
	finalFold   FoldingCircuit        // final evaluation
}

// NewSumcheckVerifier generates a sumcheck verifier circuit from a claim and degrees
func NewSumcheckVerifier(claim frontend.Variable, degrees [bN + bG]int) SumcheckVerifier {

	polynomials := make([][]frontend.Variable, bG+bN)
	for i := 0; i < bG+bN; i++ {
		for i, d := range degrees {
			polynomials[i] = make([]frontend.Variable, d)
		}
	}

	return SumcheckVerifier{
		nRounds:     len(degrees),
		claim:       claim,
		polynomials: polynomials,
		finalFold:   NewFoldingCircuit(),
	}
}

// Define contains the circuit data for an nRounds long sumcheck verfier
func (c *SumcheckVerifier) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	mimc, _ := mimc.NewMiMC("seed", curveID)

	var r, eval frontend.Variable
	currClaim := c.claim

	for round := 1; round <= c.nRounds; round++ {

		// eval = P_i(0) + P_i(1)
		eval = cs.Add(c.polynomials[round-1][0], c.polynomials[round-1][0], c.polynomials[round-1][1:])
		cs.AssertIsEqual(eval, currClaim)

		// deduce randomness from P_i (and the claim in the first round)
		if round == 1 {
			toHash := make([]frontend.Variable, len(c.polynomials[0])+1)
			toHash[0] = currClaim
			for i := 0; i < len(c.polynomials[round-1]); i++ {
				toHash[i+1] = c.polynomials[round-1][i]
			}
			r = mimc.Hash(cs, toHash...)
		} else {
			r = mimc.Hash(cs, c.polynomials[round-1]...)
		}

		// remember the current randomness for the final folding
		c.finalFold.varValues[round-1] = r

		// compute the next claim P_i(r_i)
		currClaim = c.polynomials[round-1][len(c.polynomials[round-1])-1]
		for i := len(c.polynomials[round-1]) - 2; i >= 0; i-- {
			cs.Mul(r, currClaim)
			cs.Add(currClaim, c.polynomials[round-1][i])
		}
	}

	// at this point: currClaim = P(r_1, r_2, ... , r_N)
	c.finalFold.claimed = currClaim
	c.finalFold.Define(curveID, cs)

	return nil
}
