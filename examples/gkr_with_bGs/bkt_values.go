package gkr

import (
	"github.com/consensys/gnark/frontend"
)

// // ValuesBKT contains the data of a circuit that folds a bookkeeping-table depending on (q, q')
// // In GKR ValuesBKT intervene twice: V_Input and V_Output; both are public.
// type ValuesBKT struct {
// 	Table [1 << (bG + bN)]frontend.Variable `gnark:",public"` // Table of values of a function V on a cube
// }

// OutputValuesBKT contains a bookkeeping-table depending on q' only.
// In GKR OutputValuesBKT intervenes once: V_Output (public).
type OutputValuesBKT struct {
	Table [1 << (bG - 1 + bN)]frontend.Variable `gnark:",public"` // Table of values of a function V on a cube
}

// InputValuesBKT contains a bookkeeping-table depending on (q, q').
// In GKR InputValuesBKT intervenes once: V_Input (public).
type InputValuesBKT struct {
	Table [1 << (bG + bN)]frontend.Variable `gnark:",public"` // Table of values of a function V on a cube
}

// SingleFold returns V(q, q') where V is the function represented by bkt.Table; recall: bG = 1.
func (bkt *OutputValuesBKT) SingleFold(cs *frontend.ConstraintSystem, QPrime [bN]frontend.Variable) frontend.Variable {

	// folding q' into the mix
	for j, x := range QPrime {
		J := 1 << (bG - 1 + bN - 1 - j)
		for k := 0; k < J; k++ {
			bkt.Table[k+J] = cs.Sub(bkt.Table[k+J], bkt.Table[k])
			bkt.Table[k+J] = cs.Mul(x, bkt.Table[k+J])
			bkt.Table[k] = cs.Add(bkt.Table[k], bkt.Table[k+J])
		}
	}

	return bkt.Table[0]
}

// DoubleFold returns V(q1, q') and V(q2, q') where V is the function represented by bkt.Table
// It first folds q' into V and then does the final folding
func (bkt *InputValuesBKT) DoubleFold(cs *frontend.ConstraintSystem, Q1, Q2 frontend.Variable, QPrime [bN]frontend.Variable) (v1, v2 frontend.Variable) {

	// folding q' into the mix
	for j, x := range QPrime {
		J := 1 << (bN + bG - 1 - j)
		TwoK := 0
		TwoKPlus1 := 1
		for k := 0; k < J; k++ {
			if k != 0 {
				TwoK += 2
				TwoKPlus1 += 2
			}
			bkt.Table[TwoKPlus1] = cs.Sub(bkt.Table[TwoKPlus1], bkt.Table[TwoK])
			bkt.Table[TwoKPlus1] = cs.Mul(x, bkt.Table[TwoKPlus1])
			bkt.Table[TwoK] = cs.Add(bkt.Table[TwoK], bkt.Table[TwoKPlus1])
		}
	}

	Folded0 := bkt.Table[0]
	Folded1 := bkt.Table[1]

	delta := cs.Sub(Folded1, Folded0)

	Q1delta := cs.Mul(Q1, delta)
	Q2delta := cs.Mul(Q2, delta)

	v1 = cs.Add(Folded0, Q1delta)
	v2 = cs.Add(Folded0, Q2delta)

	return v1, v2
}

// // DoubleFold returns V(q1, q') and V(q2, q') where V is the function represented by bkt.Table
// func (bkt *ValuesBKT) DoubleFold(cs *frontend.ConstraintSystem, Q1, Q2 frontend.Variable, QPrime [bN]frontend.Variable) (v1, v2 frontend.Variable) {

// 	bktCopy := ValuesBKT{
// 		Table: bkt.Table,
// 	}

// 	// wasteful !
// 	return bkt.SingleFold(cs, Q1, QPrime), bktCopy.SingleFold(cs, Q2, QPrime)

// }
