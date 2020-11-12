package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// FoldingCircuit contains the data of a circuit that folds a bookkeeping table
type FoldingCircuit struct {
	table     [1 << (2)]frontend.Variable `gnark:"table,public"`     // table of values of a function V on a cube
	varValues [2]frontend.Variable        `gnark:"varValues,public"` // array variable values where to compute V
	claimed   frontend.Variable           `gnark:"claimed,public"`   // claimed value of \sum_b V(b)
}

// NewFoldingCircuit generates an empty FoldingCircuit circuit
func NewFoldingCircuit() FoldingCircuit {
	var table [1 << (2)]frontend.Variable
	var varValues [2]frontend.Variable
	var claimed frontend.Variable
	return FoldingCircuit{
		table:     table,
		varValues: varValues,
		claimed:   claimed,
	}
}

// Define declares the circuit constraints of a single folding circuit c
func (c *FoldingCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	// compute a + r*(b-a) recursively
	for varIndex := 0; varIndex < 2; varIndex++ {
		j := 2 - 1 - varIndex
		for i := 0; i < 1<<j; i++ {
			c.table[i+1<<j] = cs.Sub(c.table[i+1<<j], c.table[i])
			c.table[i+1<<j] = cs.Mul(c.varValues[varIndex], c.table[i+1<<j])
			c.table[i] = cs.Add(c.table[i], c.table[i+1<<j])
		}
	}

	cs.AssertIsEqual(c.claimed, c.table[0])

	return nil
}
