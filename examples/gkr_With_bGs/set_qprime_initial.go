package gkr

func (gkr *FullGKRWithBGsCircuit) setQInitial() {

	for i := range gkr.QPrimeInitial {
		gkr.QPrimeInitial[i].Assign(initialQPrime[i])
	}
}
