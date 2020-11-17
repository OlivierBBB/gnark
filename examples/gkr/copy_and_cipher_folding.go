package gkr

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

// CopyAndCipherFoldingCircuit contains the circuit information to compute Copy(q, hL, hR)
// recall that Copy(q, hL, hR) = q * (1-hL) * hR
type CopyAndCipherFoldingCircuit struct {
	Q           frontend.Variable
	HL          frontend.Variable
	HR          frontend.Variable
	CopyValue   frontend.Variable // Copy(q, hL, hR)
	CipherValue frontend.Variable // Cipher(q, hL, hR)
}

// Define
func (c *CopyAndCipherFoldingCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	oneMinusHL := cs.Sub(1, c.HL)

	res := cs.Mul(oneMinusHL, c.HR)
	res2 := cs.Mul(c.Q, res)
	c.CopyValue.Assign(res2)
	res = cs.Sub(res, res2)
	c.CipherValue.Assign(res)

	return nil
}
