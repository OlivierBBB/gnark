package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// SingleFold contains the data of a circuit that folds a bookkeeping table
// SingleFold occurs once at the output level (level d) to compute V_d(q, q')
type SingleFold struct {
	table     [1 << (bG + bN)]frontend.Variable // table of values of a function V on a cube
	varValues [bN + bG]frontend.Variable        // array variable values where to compute V
	claimed   frontend.Variable                 // claimed value of V(q, q')
}

// DoubleFold contains the data of a circuit that folds a bookkeeping table
// DoubleFold occurs twice at the input level (level 0) to compute V_0(qL, q') and V_0(qR, q')
type DoubleFold struct {
	table     [1 << (bG + bN)]frontend.Variable // table of values of a function V on a cube
	varValues [bN + 2*bG]frontend.Variable      // array variable values where to compute V
	claimed   frontend.Variable                 // claimed value of V(q, q')
}

// Define declares the circuit constraints of a single folding
func (circuit *SingleFold) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	// compute a + r*(b-a) recursively
	for varIndex := 0; varIndex < bN+bG; varIndex++ {
		for i := 0; i < 1<<(bG+bN-1-varIndex); i++ {
			circuit.table[i+1<<(bG+bN-1-varIndex)] = cs.Sub(circuit.table[i+1<<(bG+bN-1)], circuit.table[i])
			circuit.table[i+1<<(bG+bN-1-varIndex)] = cs.Mul(circuit.varValues[varIndex], circuit.table[i+1<<(bG+bN-1-varIndex)])
			circuit.table[i] = cs.Add(circuit.table[i], circuit.table[i+1<<(bG+bN-1-varIndex)])
		}
	}
	cs.AssertIsEqual(circuit.claimed, circuit.table[0])

	return nil
}

// Define declares the circuit constraints of a single folding
func (circuit *DoubleFold) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	// For simplicity we assume that at the bottom layer the variables appear in
	// the order h', hR, hL.
	for varIndex := 0; varIndex < bN; varIndex++ {
		for i := 0; i < 1<<(bN-1-varIndex); i++ {
			circuit.table[i+1<<(bN-1-varIndex)] = cs.Sub(circuit.table[i+1<<(bG+bN-1)], circuit.table[i])
			circuit.table[i+1<<(bN-1-varIndex)] = cs.Mul(circuit.varValues[varIndex], circuit.table[i+1<<(bG+bN-1-varIndex)])
			circuit.table[i] = cs.Add(circuit.table[i], circuit.table[i+1<<(bG+bN-1-varIndex)])
		}
	}

	// at this point the relevant values in table are concentrated among the first 2**bG indices
	for varIndex := 0; varIndex < bG; varIndex++ {
		for i := 0; i < 1<<(bG-1-varIndex); i++ {
			circuit.table[i+1<<(bG-1-varIndex)] = cs.Sub(circuit.table[i+1<<(bG+bN-1)], circuit.table[i])
			circuit.table[i+1<<(2*bG-1-varIndex)] = cs.Mul(circuit.varValues[bN+bG+varIndex], circuit.table[i+1<<(bG+bN-1-varIndex)])
			circuit.table[i+1<<(bG-1-varIndex)] = cs.Mul(circuit.varValues[bN+varIndex], circuit.table[i+1<<(bG+bN-1-varIndex)])
			circuit.table[i] = cs.Add(circuit.table[i], circuit.table[i+1<<(bG+bN-1-varIndex)])
			circuit.table[i] = cs.Add(circuit.table[i], circuit.table[i+1<<(bG+bN-1-varIndex)])
		}
	}

	cs.AssertIsEqual(circuit.claimed, circuit.table[0])

	return nil
}
