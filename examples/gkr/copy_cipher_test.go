package gkr

import (
	"testing"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
)

func (ccc CopyCipherCircuit) FoldCopyCipherAndEval(cs *frontend.ConstraintSystem) (fldCo, fldCi, fldCoLC, fldCiLC frontend.Variable) {

	copy, cipher := PrefoldedCopyAndCipher(cs, ccc.Q)
	copyLC, cipherLC := PrefoldedCopyAndCipherLinComb(cs, ccc.Rho, ccc.HL, ccc.HR)

	fldCo = copy.Fold(cs, ccc.A, ccc.B)
	fldCi = cipher.Fold(cs, ccc.C, ccc.D)

	fldCoLC = copyLC.Fold(cs, ccc.E, ccc.F)
	fldCiLC = cipherLC.Fold(cs, ccc.G, ccc.H)

	return
}

type CopyCipherCircuit struct {
	Q                      frontend.Variable
	HL, HR                 frontend.Variable
	Lambda, Rho            frontend.Variable
	A, B, C, D, E, F, G, H frontend.Variable
	ExpectedFoldCo         frontend.Variable
	ExpectedFoldCi         frontend.Variable
	ExpectedFoldCoLC       frontend.Variable
	ExpectedFoldCiLC       frontend.Variable
}

func (ccc *CopyCipherCircuit) Define(curveID gurvy.ID, cs *frontend.ConstraintSystem) error {

	FldCo, FldCi, FldCoLC, FldCiLC := ccc.FoldCopyCipherAndEval(cs)

	cs.AssertIsEqual(ccc.ExpectedFoldCo, FldCo)
	cs.AssertIsEqual(ccc.ExpectedFoldCi, FldCi)
	cs.AssertIsEqual(ccc.ExpectedFoldCiLC, FldCiLC)
	cs.AssertIsEqual(ccc.ExpectedFoldCoLC, FldCoLC)

	return nil
}

func TestCopyCipher(t *testing.T) {

	var ccc CopyCipherCircuit

	r1cs, err := frontend.Compile(gurvy.BN256, &ccc)

	assert := groth16.NewAssert(t)

	assert.NoError(err)

	{
		var witness CopyCipherCircuit

		A, B, C, D := 3124545, 1234123456, 21356378, 23436787654
		E, F, G, H := 986542325, 873673457639, 7623752483, 876234512
		Q, HR, HL, Lambda, Rho := 2346234, 123413, 4645686, 2314321, 8734215

		witness.A.Assign(A)
		witness.B.Assign(B)
		witness.C.Assign(C)
		witness.D.Assign(D)

		witness.E.Assign(E)
		witness.F.Assign(F)
		witness.G.Assign(G)
		witness.H.Assign(H)

		witness.Q.Assign(Q)
		witness.HR.Assign(HR)
		witness.HL.Assign(HL)
		witness.Lambda.Assign(Lambda)
		witness.Rho.Assign(Rho)

		// EFCo := Q
		// EFCi := 0
		// EFCoLC := 0
		// EFCiLC := 0

		witness.ExpectedFoldCo.Assign(0)
		witness.ExpectedFoldCi.Assign(0)
		witness.ExpectedFoldCoLC.Assign(0)
		witness.ExpectedFoldCiLC.Assign(0)

		assert.ProverSucceeded(r1cs, &witness)
	}
}
