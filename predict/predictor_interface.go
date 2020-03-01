package svd

import (
	"database/sql"
)

// PredictorKey : represents key for predictor storing
// type PredictorKey [16]byte

// Predictor : intervase describing predictor
type Predictor interface {
	ReadSQLRows(*sql.Rows)
	GetValue(ID, ID) float64
	GetSubjectsForObject(ID) *IDValues
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

// Predictors : structure for predictors storing
// type Predictors map[PredictorKey]Predictor

// Key :
// func (predictors *Predictors) Key(str string) PredictorKey {
// 	return md5.Sum([]byte(str))
// }

// Get :
// func (predictors *Predictors) Get(str string) Predictor {
// 	key := predictors.Key(str)

// 	predictor, ok := (*predictors)[key]
// }
