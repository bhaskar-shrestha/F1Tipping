package models

// Prediction represents a user's prediction for a race weekend
type Prediction struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	SubmitTime   string   `json:"submit_time"`
	DriverIDs    []string `json:"driver_ids"`  // 5 selected drivers
	TeamIDs      []string `json:"team_ids"`    // 2 selected teams
	SprintPoints int      `json:"sprint_points"` // Points earned from sprint
	RacePoints   int      `json:"race_points"`   // Points earned from race
	TotalPoints  int      `json:"total_points"`  // sprint + race points
}
