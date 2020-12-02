package gkr

// SetRoundConstants sets the round constants of MiMC.
func (gkr *CircuitGKR) setRoundConstants() {

	for i := range gkr.RoundConstants {
		gkr.RoundConstants[nLayers-1-i].Assign(arks[i])
	}
}

// SetVLAndVR sets the values of VL and VR for the final check of the intermediate sumchecks
func (gkr *CircuitGKR) setVLAndVR() {

	for i := range gkr.VLClaimed {
		gkr.VLClaimed[i].Assign(claimedVLs[nLayers-1-i])
	}

	for i := range gkr.VRClaimed {
		gkr.VRClaimed[i].Assign(claimedVRs[nLayers-1-i])
	}
}

// setQInitial sets QPrimeInitial for the first round of the first sumcheck
func (gkr *CircuitGKR) setQPrimeInitial() {

	for i := range gkr.QPrimeInitial {
		gkr.QPrimeInitial[i].Assign(initialQPrime[i])
	}
}

// setPolynomials sets the polynomials of the sumchecks
func (gkr *CircuitGKR) setPolynomials() {

	for layer := range sumcheckProofs {
		// filling in HLPoly of layer: degHL+1 coefficients
		for d := range gkr.HLPolynomials[layer].Coefficients {
			gkr.HLPolynomials[layer].Coefficients[d].Assign(sumcheckProofs[nLayers-1-layer][0][d])
		}
		// filling in HRPoly of layer: degHR+1 coefficients
		for d := range gkr.HRPolynomials[layer].Coefficients {
			gkr.HRPolynomials[layer].Coefficients[d].Assign(sumcheckProofs[nLayers-1-layer][1][d])
		}
		// filling in HPrimePolys of layer
		// bN iterations
		for varIndex := range gkr.HPrimePolynomials[layer] {
			// degHPrime + 1 coefficients
			for d := range gkr.HPrimePolynomials[layer][varIndex].Coefficients {
				gkr.HPrimePolynomials[layer][varIndex].Coefficients[d].Assign(sumcheckProofs[nLayers-1-layer][2+varIndex][d])
			}
		}
	}
}

// setInputs sets the inputs of the gkr
// we are currently hashing [0, 0, 0, 0, 0, 0, 0, 0]
func (gkr *CircuitGKR) setInputs() {

	for i := range gkr.VInput.Table {
		gkr.VInput.Table[i].Assign(0)
	}
}

// setOutputs sets the outputs of the gkr
func (gkr *CircuitGKR) setOutputs() {

	for i := range gkr.VOutput.Table {
		gkr.VOutput.Table[i].Assign(outputs[i])
	}
}

func (gkr *CircuitGKR) setPublicInputs() {
	gkr.setRoundConstants()
	gkr.setInputs()
	gkr.setOutputs()
	gkr.setVLAndVR()
	gkr.setPolynomials()
	gkr.setQPrimeInitial()
}
