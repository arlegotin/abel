package svd

import "gonum.org/v1/gonum/mat"

// Factorize : performs singulare value decomposition.
// Will panic if decomposition failed
func factorize(matrix *mat.Dense) (*mat.Dense, *[]float64, *mat.Dense) {
	SVD := mat.SVD{}

	// Factorization itself
	ok := SVD.Factorize(matrix, mat.SVDThin)

	if !ok {
		panic("Couldn't factorize matrix")
	}

	// Extract U
	var U mat.Dense
	SVD.UTo(&U)

	// Extract V
	var V mat.Dense
	SVD.VTo(&V)

	// Extract sigma
	sigma := SVD.Values(nil)

	return &U, &sigma, &V
}
