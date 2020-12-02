package gkr

import (
	"github.com/consensys/gnark/frontend"
)

// Polynomial encodes a polynomial in terms of its coefficients:
// a0 + a1X + ... + ad X^d <--> {a0, a1, ... , ad}
type Polynomial struct {
	Coefficients []frontend.Variable
}

// eval returns p(x) and adds the associated computational constraints to cs
func (p *Polynomial) eval(cs *frontend.ConstraintSystem, x frontend.Variable) (res frontend.Variable) {

	res = cs.Constant(0)

	for i := len(p.Coefficients) - 1; i >= 0; i-- {
		if i != len(p.Coefficients)-1 {
			res = cs.Mul(res, x)
		}
		res = cs.Add(res, p.Coefficients[i])
	}

	return res
}

// zeroAndOne returns P(0) + P(1)
func (p *Polynomial) zeroAndOne(cs *frontend.ConstraintSystem) frontend.Variable {

	// coeffsInterface is required for cs.Add(a, b, coeffsInterface[1:]...) to be accepted.
	coeffsInterface := make([]interface{}, len(p.Coefficients))
	for i, coeff := range p.Coefficients {
		coeffsInterface[i] = coeff
	}

	res := cs.Add(p.Coefficients[0], p.Coefficients[0], coeffsInterface[1:]...)

	return res
}
