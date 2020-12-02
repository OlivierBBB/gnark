package gkr

import (
	"github.com/consensys/gnark/frontend"
)

// PrefoldedGateBKT contains a prefolded book-keeping table associated to a gate type in a layered arithmetic circuit.
// Examples are: Copy and Cipher.
type PrefoldedGateBKT struct {
	Table [1 << (2 * bG)]frontend.Variable // recall bG = 1
}

// Fold takes a 2 x 2 prefolded gate book-keeping table and returns its value at (hL, hR); recall: bG = 1.
func (pre *PrefoldedGateBKT) Fold(cs *frontend.ConstraintSystem, HL, HR frontend.Variable) frontend.Variable {

	// folding HL into the mix
	pre.Table[2] = cs.Sub(pre.Table[2], pre.Table[0])
	pre.Table[2] = cs.Mul(pre.Table[2], HL)
	pre.Table[0] = cs.Add(pre.Table[0], pre.Table[2]) // S[0] <- T[0] + hL * (T[2] - T[0])
	pre.Table[3] = cs.Sub(pre.Table[3], pre.Table[1])
	pre.Table[3] = cs.Mul(pre.Table[3], HL)
	pre.Table[1] = cs.Add(pre.Table[1], pre.Table[3]) // S[1] <- T[1] + hL * (T[3] - T[1])

	// folding HR into the mix
	pre.Table[1] = cs.Sub(pre.Table[1], pre.Table[0])
	pre.Table[1] = cs.Mul(pre.Table[1], HR)
	pre.Table[0] = cs.Add(pre.Table[0], pre.Table[1]) // FoldedValue == S[0] + hR * (S[1] - S[0])

	return pre.Table[0]
}
