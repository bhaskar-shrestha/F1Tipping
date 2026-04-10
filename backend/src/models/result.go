package models

// PositionResult represents where a constructor finished in a race
type PositionResult struct {
	ConstructorID string `json:"constructor_id"`
	Position      int    `json:"position"` // 1st, 2nd, 3rd, etc.
}

// SprintResult represents sprint race results
type SprintResult struct {
	DriverPositions []PositionResult `json:"driver_positions"`
	TeamPositions   []PositionResult `json:"team_positions"`
}

// RaceResult represents main race results
type RaceResult struct {
	DriverPositions []PositionResult `json:"driver_positions"`
	TeamPositions   []PositionResult `json:"team_positions"`
}
