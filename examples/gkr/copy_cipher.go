package gkr

import (
	"github.com/consensys/gnark/frontend"
)

// PrefoldedCopyAndCipher generates two prefolded gate BKTs
//		prefoldedCopy(?,?)		:= Copy(Q,?,?)
// 		prefoldedCipher(?,?)	:= Cipher(Q,?,?)
// these BKTs contain a single nonzero value for (?,?) == (1,0) i.e. in (BKT_name).Table[2]
func PrefoldedCopyAndCipher(cs *frontend.ConstraintSystem, Q frontend.Variable) (prefoldedCopy, prefoldedCipher PrefoldedGateBKT) {

	// set prefoldedCopy.Table[2] and prefoldedCipher.Table[2]
	prefoldedCopy.Table[2] = Q
	prefoldedCipher.Table[1] = cs.Sub(1, Q)
	// zeros everywhere else
	for i := 0; i < len(prefoldedCopy.Table); i++ {
		if i != 2 {
			prefoldedCopy.Table[i] = cs.Constant(0)
			prefoldedCipher.Table[i] = cs.Constant(0)
		}
	}

	return prefoldedCopy, prefoldedCipher
}

// PrefoldedCopyAndCipherLinComb generates two linear combinations of prefolded gate BKTs
//		prefoldedCopy(?,?) := lambda * Copy(hL,?,?) + rho * Copy(hR,?,?)
//		prefoldedCipher(?,?) := lambda * Cipher(hL,?,?) + rho * Cipher(hR,?,?)
// these BKTs contain a single nonzero value for (?,?) == (1,0) i.e. in (BKT_name).Table[2]
func PrefoldedCopyAndCipherLinComb(cs *frontend.ConstraintSystem, lambda, rho, hL, hR frontend.Variable) (prefoldedCopy, prefoldedCipher PrefoldedGateBKT) {

	// copy(a,b,c) = a * (1-b) * c
	// prefoldedCopy.Table[2] = hL + rho * hR
	aux := cs.Mul(rho, hR) // rho * hR
	aux = cs.Add(hL, aux)  // hL + rho * hR
	prefoldedCopy.Table[2] = aux
	// cipher(a,b,c) = (1-a) * (1-b) * c
	// prefoldedCipher.Table[2] = (1 - hL) + rho * (1 - hR) = (1 + rho) - (hL + rho * hR)
	aux = cs.Sub(cs.Constant(1), aux) // 1 - (hL + rho * hR)
	aux = cs.Add(aux, rho)            // (1 + rho) - (hL + rho * hR)
	prefoldedCipher.Table[2] = aux
	// zeros everywhere else
	for i := 0; i < len(prefoldedCopy.Table); i++ {
		if i != 2 {
			prefoldedCopy.Table[i] = cs.Constant(0)
			prefoldedCipher.Table[i] = cs.Constant(0)
		}
	}

	return prefoldedCopy, prefoldedCipher

}
