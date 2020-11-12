package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// DoubleFold contains the data of a circuit that folds a bookkeeping table
// DoubleFold occurs twice at the input level (level 0) to compute V_0(qL, q') and V_0(qR, q')
type DoubleFold struct {
	table     [1 << (bG + bN)]frontend.Variable // table of values of a function V on a cube
	varValues [bN + 2*bG]frontend.Variable      // array variable values where to compute V
	claimed   frontend.Variable                 // claimed value of V(q, q')
}

// Define declares the circuit constraints of a single double-folding circuit c
func (c *DoubleFold) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	// For simplicity we assume that at the bottom layer the variables appear in
	// the order h', hR, hL.
	for varIndex := 0; varIndex < bN; varIndex++ {
		for i := 0; i < 1<<(bN-1-varIndex); i++ {
			c.table[i+1<<(bN-1-varIndex)] = cs.Sub(c.table[i+1<<(bG+bN-1)], c.table[i])
			c.table[i+1<<(bN-1-varIndex)] = cs.Mul(c.varValues[varIndex], c.table[i+1<<(bG+bN-1-varIndex)])
			c.table[i] = cs.Add(c.table[i], c.table[i+1<<(bG+bN-1-varIndex)])
		}
	}

	// at this point the relevant values in table are concentrated among the first 2**bG indices
	for varIndex := 0; varIndex < bG; varIndex++ {
		for i := 0; i < 1<<(bG-1-varIndex); i++ {
			c.table[i+1<<(bG-1-varIndex)] = cs.Sub(c.table[i+1<<(bG+bN-1)], c.table[i])
			c.table[i+1<<(2*bG-1-varIndex)] = cs.Mul(c.varValues[bN+bG+varIndex], c.table[i+1<<(bG+bN-1-varIndex)])
			c.table[i+1<<(bG-1-varIndex)] = cs.Mul(c.varValues[bN+varIndex], c.table[i+1<<(bG+bN-1-varIndex)])
			c.table[i] = cs.Add(c.table[i], c.table[i+1<<(bG+bN-1-varIndex)])
			c.table[i] = cs.Add(c.table[i], c.table[i+1<<(bG+bN-1-varIndex)])
		}
	}

	cs.AssertIsEqual(c.claimed, c.table[0])

	return nil
}
