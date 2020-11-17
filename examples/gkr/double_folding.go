package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// DoubleFoldingCircuit contains the data of a circuit that folds a bookkeeping Table V twice
// to get VL := V(hL, h') and VR := V(hR, h')
type DoubleFoldingCircuit struct {
	Table  [1 << (bG + bN)]frontend.Variable // Table of values of a function V on a cube
	HPrime [bN]frontend.Variable             // values of h' to substitude into V
	HL     [bG]frontend.Variable             // values of hL to substitude into V
	HR     [bG]frontend.Variable             // values of hR to substitude into V
	VL     frontend.Variable                 // claimed value of V(hL, h')
	VR     frontend.Variable                 // claimed value of V(hR, h')
}

// Define declares the circuit constraints of a single double-folding circuit c
func (c *DoubleFoldingCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	// fold h' into Table
	for varIndex := 0; varIndex < bN; varIndex++ {
		for i := 0; i < 1<<(bN-1-varIndex); i++ {
			c.Table[i+1<<(bN-1-varIndex)] = cs.Sub(c.Table[i+1<<(bG+bN-1)], c.Table[i])
			c.Table[i+1<<(bN-1-varIndex)] = cs.Mul(c.HPrime[varIndex], c.Table[i+1<<(bG+bN-1-varIndex)])
			c.Table[i] = cs.Add(c.Table[i], c.Table[i+1<<(bG+bN-1-varIndex)])
		}
	}

	// enough to put RightCopy = c.Table[:1<<bG] ?
	var RightCopy [1 << bG]frontend.Variable
	for i := 0; i < 1<<bG; i++ {
		RightCopy[i] = c.Table[i]
	}

	// at this point the relevant values in Table are concentrated among the first 2**bG indices
	for varIndex := 0; varIndex < bG; varIndex++ {
		for i := 0; i < 1<<(bG-1-varIndex); i++ {
			// VL computation
			c.Table[i+1<<(bG-1-varIndex)] = cs.Sub(c.Table[i+1<<(bG-1-varIndex)], c.Table[i])
			c.Table[i+1<<(bG-1-varIndex)] = cs.Mul(c.HL[varIndex], c.Table[i+1<<(bG-1-varIndex)])
			c.Table[i] = cs.Add(c.Table[i], c.Table[i+1<<(bG-1-varIndex)])
			// VR computation
			RightCopy[i+1<<(bG-1-varIndex)] = cs.Sub(RightCopy[i+1<<(bG-1-varIndex)], RightCopy[i])
			RightCopy[i+1<<(bG-1-varIndex)] = cs.Mul(c.HR[varIndex], RightCopy[i+1<<(bG-1-varIndex)])
			RightCopy[i] = cs.Add(RightCopy[i], RightCopy[i+1<<(bG-1-varIndex)])
		}
	}

	cs.AssertIsEqual(c.VL, c.Table[0])
	cs.AssertIsEqual(c.VR, RightCopy[0])

	return nil
}
