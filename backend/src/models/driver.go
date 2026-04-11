// Package models defines the core data structures for the F1 prediction app
package models

import "encoding/json"

// Driver represents a Formula 1 driver with their constructor
type Driver struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ConstructorID   string `json:"constructor_id"`
	ConstructorName string `json:"constructor_name"`
}

// MarshalJSON custom JSON marshaling for Driver
func (d Driver) MarshalJSON() ([]byte, error) {
	type Alias Driver
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&d),
	})
}
