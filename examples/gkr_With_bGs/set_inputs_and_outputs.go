package gkr

func (gkr *FullGKRWithBGsCircuit) setInputs() {

	for i := range gkr.VInput.Table {
		gkr.VInput.Table[i].Assign(0)
	}
}

func (gkr *FullGKRWithBGsCircuit) setOutputs() {

	for i := range gkr.VOutput.Table {
		gkr.VOutput.Table[i].Assign(outputs[i])
	}
}
