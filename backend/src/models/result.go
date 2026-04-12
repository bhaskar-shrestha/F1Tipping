package models

// ConstructorPosition represents where a team/constructor finished in a race
type ConstructorPosition struct {
	ConstructorID string `json:"constructor_id"`
	Position      int    `json:"position"` // 1-20 (1st through 20th)
}

// SprintResult represents sprint race results
// Contains all constructor positions for a specific sprint race
type SprintResult struct {
	RaceID    int                   `json:"race_id"`
	Positions []ConstructorPosition `json:"positions"` // All constructors' final positions
}

// RaceResult represents main race results
// Contains all constructor positions for a specific main race
type RaceResult struct {
	RaceID    int                   `json:"race_id"`
	Positions []ConstructorPosition `json:"positions"` // All constructors' final positions
}
