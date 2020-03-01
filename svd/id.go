package svd

import (
	"fmt"
	"sort"
)

// ID : represents objects and subjects ID
type ID int

// IDs : represents slice of IDs
type IDs []ID

// IDMap : represents map {ID->int}
type IDMap map[ID]int

// MapIndexes : creates maps of indexes: {old->new} and {new->old}
func (ids *IDs) MapIndexes() *IDMap {
	idMap := IDMap{}
	index := 0

	for _, id := range *ids {
		if _, ok := idMap[id]; !ok {
			idMap[id] = index
			index++
		}
	}

	return &idMap
}

// Contains : check if element exists in the slice
func (ids *IDs) Contains(targetID ID) bool {
	for _, id := range *ids {
		if id == targetID {
			return true
		}
	}

	return false
}

// AppendIfNotExists : appends element to slice if it not exists
func (ids *IDs) AppendIfNotExists(newID ID) {
	if !ids.Contains(newID) {
		*ids = append(*ids, newID)
	}
}

// IDValue : stores pair (id, value)
type IDValue struct {
	ID    ID
	Value float64
}

// IDValues : stores slice of pairs
type IDValues []IDValue

// Stretch : sets capacity
func (list *IDValues) Stretch(capacity int) {
	*list = make(IDValues, 0, capacity)
}

// Set : sets value for id
func (list *IDValues) Set(id ID, value float64) {
	*list = append(*list, IDValue{id, value})
}

// Sort : sets value for id
func (list *IDValues) Sort() {
	sort.Slice(*list, func(i, j int) bool {
		return (*list)[i].Value > (*list)[j].Value
	})
}

// Filter : filters
func (list *IDValues) Filter(filter func(IDValue) bool) {
	n := 0
	for _, idValue := range *list {
		if filter(idValue) {
			(*list)[n] = idValue
			n++
		}
	}
	*list = (*list)[:n]
}

// Print : sets value for id
func (list *IDValues) Print() {
	for _, idValue := range *list {
		fmt.Printf("%d -> %f\n", idValue.ID, idValue.Value)
	}
}

// Len : sets value for id
func (list *IDValues) Len() int {
	return len(*list)
}

// Cut : sets value for id
func (list *IDValues) Cut(k int) {
	*list = (*list)[0:k]
}
