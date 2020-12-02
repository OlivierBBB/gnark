package gkr

import (
	"github.com/consensys/gnark/frontend"
)

// PrefoldedCopyAndCipher generates two prefolded gate BKTs
//		prefoldedCopy(?,?)		:= Copy(Q,?,?)
// 		prefoldedCipher(?,?)	:= Cipher(Q,?,?)
// these BKTs contain a nonzero value only for (?,?) == (0,1) i.e. in BKT.Table[1]
func PrefoldedCopyAndCipher(cs *frontend.ConstraintSystem, Q frontend.Variable) (prefoldedCopy, prefoldedCipher PrefoldedGateBKT) {

	// set prefoldedCopy.Table[1]
	prefoldedCopy.Table[2] = Q
	// put zeros everywhere else
	for i := 0; i < len(prefoldedCopy.Table); i++ {
		if i != 2 {
			prefoldedCopy.Table[i] = cs.Constant(0)
		}
	}

	// set prefoldedCipher.Table[1]
	prefoldedCipher.Table[1] = cs.Sub(1, Q)
	// put zeros everywhere else
	for i := 0; i < len(prefoldedCipher.Table); i++ {
		if i != 2 {
			prefoldedCipher.Table[i] = cs.Constant(0)
		}
	}

	return prefoldedCopy, prefoldedCipher
}

// PrefoldedCopyAndCipherLinComb generates two linear combinations of prefolded gate BKTs
//		prefoldedCopy(?,?) := lambda * Copy(hL,?,?) + rho * Copy(hR,?,?)
//		prefoldedCipher(?,?) := lambda * Cipher(hL,?,?) + rho * Cipher(hR,?,?)
// these BKTs contain a nonzero value only for (?,?) == (0,1) i.e. in BKT.Table[1]
func PrefoldedCopyAndCipherLinComb(cs *frontend.ConstraintSystem, lambda, rho, hL, hR frontend.Variable) (prefoldedCopy, prefoldedCipher PrefoldedGateBKT) {

	// populate prefoldedCopy
	// copy(a,b,c) = a * (1-b) * c
	// prefoldedCopy.Table[01] = lambda * hL + rho * hR

	// cs.Println("hL:", hL)
	// cs.Println("hR:", hR)
	aux1 := cs.Mul(lambda, hL)                  // lambda * hL
	aux2 := cs.Mul(rho, hR)                     // rho * hR
	prefoldedCopy.Table[2] = cs.Add(aux1, aux2) // lambda * hL + rho * hR
	// force zeros everywhere else --- is this necessary ?
	for i := 0; i < len(prefoldedCopy.Table); i++ {
		if i != 2 {
			prefoldedCopy.Table[i] = cs.Constant(0)
		}
	}

	// populate prefoldedCipher
	// cipher(a,b,c) = (1-a) * (1-b) * c
	// prefoldedCipher.Table[01] = lambda * (1 - hL) + rho * (1 - hR)
	aux1 = cs.Sub(lambda, aux1) // lambda - lambda * hL
	aux2 = cs.Sub(rho, aux2)    // rho - rho * hR
	prefoldedCipher.Table[2] = cs.Add(aux1, aux2)
	// force zeros everywhere else --- is this necessary ?
	for i := 0; i < len(prefoldedCipher.Table); i++ {
		if i != 2 {
			prefoldedCipher.Table[i] = cs.Constant(0)
		}
	}

	// cs.Println("prefoldedCopy[2]:", prefoldedCopy.Table[2])
	// cs.Println("prefoldedCipher[2]:", prefoldedCipher.Table[2])

	return prefoldedCopy, prefoldedCipher

}
