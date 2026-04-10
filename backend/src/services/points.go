package services

import "fmt"

// F1 2026 Points Rules

// SprintPoints returns points for sprint positions (top 8 finishers)
// Position 1-8: 8, 7, 6, 5, 4, 3, 2, 1 points
var SprintPoints = []int{8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// RacePoints returns points for race positions (top 10 finishers)
// Position 1-10: 25, 18, 15, 12, 10, 8, 6, 4, 2, 1 points
var RacePoints = []int{25, 18, 15, 12, 10, 8, 6, 4, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// Constants
const (
	SPRINT_POINTS = SprintPoints
	RACE_POINTS   = RacePoints
)

// calculatePoints calculates points for a list of driver IDs based on their positions
func calculatePoints(driverIDs []string, positions map[string]int, points []int) int {
	total := 0
	for _, driverID := range driverIDs {
		pos := positions[driverID]
		total += points[pos]
	}
	return total
}

// calculateTeamPoints calculates points for a list of team IDs based on both cars' positions
func calculateTeamPoints(teamIDs []string, teams map[string]*models.Team, points []int) int {
	total := 0
	for _, teamID := range teamIDs {
		team, ok := teams[teamID]
		if !ok {
			continue
		}

		// Get positions for both cars (or single position if available)
		var car1Pos int = 0
		if team.RaceCar1Position != nil {
			car1Pos = *team.RaceCar1Position
		} else if team.RaceCar2Position != nil {
			car1Pos = *team.RaceCar2Position
		}

		var car2Pos int = 0
		if team.RaceCar2Position != nil {
			car2Pos = *team.RaceCar2Position
		} else if team.RaceCar1Position != nil {
			car2Pos = *team.RaceCar1Position
		}

		total += points[car1Pos]
		total += points[car2Pos]
	}
	return total
}

// predictionID generates a prediction ID from user ID
func predictionID(userID string) string {
	return fmt.Sprintf("pred_%s", userID)
}
