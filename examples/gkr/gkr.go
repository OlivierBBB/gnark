package gkr

const nLayers = 2
const bN = 1
const bG = 1

// FullGKR contains the circuit data for an nLayers deep GKR circuit
// Note: the input folding is not optimized.
type FullGKR struct {
	ouputFolding      SingleFold                          // to produce the initial claim
	sumcheckVerifiers [nLayers][bN + 2*bG]SumcheckVerfier // round verifications
	inputFoldingLeft  SingleFold                          // to finish off the proof
	inputFoldingRight SingleFold                          // to finish off the proof
}
