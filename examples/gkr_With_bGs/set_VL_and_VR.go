package gkr

// SetRoundConstants sets the round constants of MiMC.
func (gkr *FullGKRWithBGsCircuit) setVLAndVR() {

	for i := range gkr.VLClaimed {
		gkr.VLClaimed[i].Assign(claimedVLs[nLayers-1-i])
	}

	for i := range gkr.VRClaimed {
		gkr.VRClaimed[i].Assign(claimedVRs[nLayers-1-i])
	}
}
