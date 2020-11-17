package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// FoldingCircuit contains the data of a circuit that folds a bookkeeping-table depending on (q, q')
type FoldingCircuit struct {
	Table       [1 << (bG + bN)]frontend.Variable `gnark:"Table,public"`     // Table of values of a function V on a cube
	Q           [bG]frontend.Variable             `gnark:"VarValues,public"` // array variable values where to compute V
	QPrime      [bN]frontend.Variable             `gnark:"VarValues,public"` // array variable values where to compute V
	FoldedValue frontend.Variable                 // FoldedValue = V(q, q')
}

// Define declares the circuit constraints of a single folding circuit c
func (c *FoldingCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	for i, x := range c.Q {
		J := 1 << (bG + bN - 1 - i)
		for k := 0; k < J; k++ {
			c.Table[k+J] = cs.Sub(c.Table[k+J], c.Table[k])
			c.Table[k+J] = cs.Mul(x, c.Table[i+J])
			c.Table[k] = cs.Add(c.Table[i], c.Table[i+J])
		}
	}

	for i, x := range c.QPrime {
		J := 1 << (bN - 1 - i)
		for k := 0; k < J; k++ {
			c.Table[k+J] = cs.Sub(c.Table[k+J], c.Table[k])
			c.Table[k+J] = cs.Mul(x, c.Table[i+J])
			c.Table[k] = cs.Add(c.Table[i], c.Table[i+J])
		}
	}

	c.FoldedValue.Assign(c.Table[0])

	return nil
}
