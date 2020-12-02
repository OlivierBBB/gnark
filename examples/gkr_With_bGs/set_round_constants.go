package gkr

// SetRoundConstants sets the round constants of MiMC.
func (gkr *FullGKRWithBGsCircuit) SetRoundConstants() {

	for i := range gkr.RoundConstants {
		gkr.RoundConstants[nLayers-1-i].Assign(arks[i])
	}
}
