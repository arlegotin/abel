package svd

import (
	"database/sql"

	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

// Record : structure to store rows from DB
type Record struct {
	object  ID
	subject ID
	value   float64
}

// Records : structure to store slice or records
type Records []Record

// ReadSQLRows : reads and saves SQL-rows
func (records *Records) ReadSQLRows(rows *sql.Rows) {
	*records = *(rowsToRecords(rows))
}

// ToMatrixAndMaps : converst stored rows into matrix and index-maps
func (records *Records) ToMatrixAndMaps() (*mat.Dense, *IDMap, *IDMap) {
	return recordsToMatrixAndMaps(records)
}

// rowsToRecords : converts SQL-rows to records
func rowsToRecords(rows *sql.Rows) *Records {
	records := make(Records, 0)

	for rows.Next() {
		var (
			object  int
			subject int
			value   float64
		)

		if err := rows.Scan(
			&object,
			&subject,
			&value,
		); err != nil {
			continue
		}

		record := Record{ID(object), ID(subject), value}

		records = append(records, record)
	}

	return &records
}

// recordsToObjectsAndSubjects : extracts unique objects and subjects from records
func recordsToObjectsAndSubjects(records *Records) (*IDs, *IDs) {
	objects := make(IDs, 0)
	subjects := make(IDs, 0)

	for _, record := range *records {
		objects.AppendIfNotExists(record.object)
		subjects.AppendIfNotExists(record.subject)
	}

	return &objects, &subjects
}

// recordsToMatrixAndMaps : converts records to matrix (using sparse matrix first) and index maps
func recordsToMatrixAndMaps(records *Records) (*mat.Dense, *IDMap, *IDMap) {
	objects, subjects := recordsToObjectsAndSubjects(records)

	objectMap := objects.MapIndexes()
	subjectMap := subjects.MapIndexes()

	m := len(*objects)
	n := len(*subjects)

	// Create sparse matrix
	sparseMatrix := sparse.NewDOK(m, n)

	// Fill sparse matrix with records
	for _, record := range *records {
		sparseMatrix.Set(
			(*objectMap)[record.object],
			(*subjectMap)[record.subject],
			record.value,
		)
	}

	// Convert sparse matrix to dense (to use in mat.SVD)
	matrix := sparseMatrix.ToDense()

	return matrix, objectMap, subjectMap
}
