package svd

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// USVM2 : struct for storing results of SVD-factorization and index-maps
type USVM2 struct {
	U          *mat.Dense
	Sigma      *[]float64
	V          *mat.Dense
	ObjectMap  *IDMap
	SubjectMap *IDMap
}

// ReadSQLRows : processes and stores SQL-rows
func (usvm2 *USVM2) ReadSQLRows(rows *sql.Rows) {
	records := make(Records, 0)

	records.ReadSQLRows(rows)

	matrix, objectMap, subjectMap := records.ToMatrixAndMaps()

	U, sigma, V := factorize(matrix)

	usvm2.U = U
	usvm2.Sigma = sigma
	usvm2.V = V
	usvm2.ObjectMap = objectMap
	usvm2.SubjectMap = subjectMap
}

// Cut : decreases dims of data
func (usvm2 *USVM2) Cut(k int) (bool, string) {

	if usvm2.isEmpty() {
		return false, "Is empty"
	}

	const minK = 1

	if k < minK {
		return false, fmt.Sprintf("k=%d must be at least %d", k, minK)
	}

	if k < len(*usvm2.Sigma) {
		newSigma := (*usvm2.Sigma)[0:k]
		usvm2.Sigma = &newSigma

		m, _ := (*usvm2.U).Dims()
		newU := (*usvm2.U).Slice(0, m, 0, k).(*mat.Dense)
		usvm2.U = newU

		n, _ := (*usvm2.V).Dims()
		newV := (*usvm2.V).Slice(0, n, 0, k).(*mat.Dense)
		usvm2.V = newV
	}

	return true, ""
}

// Compact : compression without loss
func (usvm2 *USVM2) Compact() (bool, string) {

	if usvm2.isEmpty() {
		return false, "Is empty"
	}

	for k, value := range *usvm2.Sigma {
		if value <= 0 {
			return usvm2.Cut(k)
		}
	}

	return true, ""
}

// Compress : compression with loss
func (usvm2 *USVM2) Compress(savedPart float64) (bool, string) {

	if usvm2.isEmpty() {
		return false, "Is empty"
	}

	const minPart = 0.0
	const maxPart = 1.0

	if savedPart <= minPart || maxPart < savedPart {
		return false, fmt.Sprintf("savedPart=%f must be in interval (%f, %f]", savedPart, minPart, maxPart)
	}

	fullSigmaSum := 0.0

	for _, value := range *usvm2.Sigma {
		fullSigmaSum += value
	}

	targetSum := fullSigmaSum * savedPart
	sigmaSum := 0.0
	cutK := 0

	for k, value := range *usvm2.Sigma {
		sigmaSum += value
		cutK = k + 1

		if sigmaSum >= targetSum {
			break
		}
	}

	return usvm2.Cut(cutK)
}

// GetSize : returns sizes of data
func (usvm2 *USVM2) GetSize() (int, int, int, int, int) {

	if usvm2.isEmpty() {
		return 0, 0, 0, 0, 0
	}

	m, k := (*usvm2.U).Dims()
	n, _ := (*usvm2.V).Dims()
	objectsLength := len(*usvm2.ObjectMap)
	subjectsLength := len(*usvm2.SubjectMap)

	return m, k, n, objectsLength, subjectsLength
}

// Marshal : exports struct
func (usvm2 *USVM2) Marshal() ([]byte, error) {
	return json.Marshal(usvm2)
}

// Unmarshal : imports struct
func (usvm2 *USVM2) Unmarshal(marshaled []byte) error {
	return json.Unmarshal(marshaled, usvm2)
}

// ObjectIndex : returns object's index
func (usvm2 *USVM2) ObjectIndex(object ID) (int, string) {

	if usvm2.isEmpty() {
		return -1, "Is empty"
	}

	return (*usvm2.ObjectMap)[object], ""
}

// SubjectIndex : returns subjects's index
func (usvm2 *USVM2) SubjectIndex(subject ID) (int, string) {

	if usvm2.isEmpty() {
		return -1, "Is empty"
	}

	return (*usvm2.SubjectMap)[subject], ""
}

// GetValue : return value for object and subject
func (usvm2 *USVM2) GetValue(object, subject ID) float64 {

	if usvm2.isEmpty() {
		return 0
	}

	uIndex, _ := usvm2.ObjectIndex(object)
	u := (*usvm2.U).RawRowView(uIndex)

	vIndex, _ := usvm2.SubjectIndex(subject)
	v := (*usvm2.V).RawRowView(vIndex)

	value := 0.0

	for k := range u {
		value += u[k] * (*usvm2.Sigma)[k] * v[k]
	}

	return value
}

// GetSubjectsForObject : returns subjects sorted by best value
func (usvm2 *USVM2) GetSubjectsForObject(object ID) *IDValues {

	_, _, _, l, _ := usvm2.GetSize()

	list := IDValues{}

	list.Stretch(l)

	if l == 0 {
		return &list
	}

	for subject := range *usvm2.SubjectMap {
		value := usvm2.GetValue(object, subject)
		list.Set(subject, value)
	}

	list.Sort()

	return &list

}

// PrintSigma : prints sigma
func (usvm2 *USVM2) PrintSigma() {
	fmt.Println(usvm2.Sigma)
}

// isEmpty : returns true if no data
func (usvm2 *USVM2) isEmpty() bool {

	if usvm2.U == nil {
		return true
	}

	if (*usvm2.U).IsEmpty() {
		return true
	}

	if usvm2.Sigma == nil {
		return true
	}

	if len(*usvm2.Sigma) == 0 {
		return true
	}

	if usvm2.V == nil {
		return true
	}

	if (*usvm2.V).IsEmpty() {
		return true
	}

	if usvm2.ObjectMap == nil {
		return true
	}

	if len(*usvm2.ObjectMap) == 0 {
		return true
	}

	if usvm2.SubjectMap == nil {
		return true
	}

	if len(*usvm2.SubjectMap) == 0 {
		return true
	}

	return false
}
