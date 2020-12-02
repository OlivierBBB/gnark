package gkr

func (gkr *FullGKRWithBGsCircuit) setPolynomials() {

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
