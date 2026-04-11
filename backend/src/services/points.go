package services

import (
	"github.com/f1tipping/backend/src/models"
)

// F1 2026 Points Rules

// SprintPoints returns points for sprint positions (top 8 finishers)
// Position 1-8: 8, 7, 6, 5, 4, 3, 2, 1 points
var SprintPoints = []int{8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// RacePoints returns points for race positions (top 10 finishers)
// Position 1-10: 25, 18, 15, 12, 10, 8, 6, 4, 2, 1 points
var RacePoints = []int{25, 18, 15, 12, 10, 8, 6, 4, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// CalculateDriverPoints calculates points for a list of driver IDs based on their positions
func CalculateDriverPoints(driverIDs []string, positions map[string]int, points []int) int {
	total := 0
	for _, driverID := range driverIDs {
		pos, ok := positions[driverID]
		if !ok {
			continue
		}
		total += points[pos]
	}
	return total
}

// CalculateTeamPoints calculates points for a list of team IDs based on both cars' positions
func CalculateTeamPoints(teamIDs []string, teamMap map[string]struct{
	Car1 int `json:"car1"`
	Car2 int `json:"car2"`
}, points []int) int {
	total := 0
	for _, teamID := range teamIDs {
		team, ok := teamMap[teamID]
		if !ok {
			continue
		}

		// Get positions for both cars
		car1Pos := team.Car1
		car2Pos := team.Car2

		total += points[car1Pos]
		total += points[car2Pos]
	}
	return total
}

// CalculatePointsForPrediction calculates total points for a prediction
func CalculatePointsForPrediction(pred *models.Prediction, driverPositions map[string]int, teamPositions map[string]struct{
	Car1 int `json:"car1"`
	Car2 int `json:"car2"`
}) (sprintPoints, racePoints, totalPoints int) {
	// Calculate driver points
	driveSprint := CalculateDriverPoints(pred.DriverIDs, driverPositions, SprintPoints)
	driveRace := CalculateDriverPoints(pred.DriverIDs, driverPositions, RacePoints)

	// Calculate team points
	teamSprint := CalculateTeamPoints(pred.TeamIDs, teamPositions, SprintPoints)
	teamRace := CalculateTeamPoints(pred.TeamIDs, teamPositions, RacePoints)

	return driveSprint + teamSprint, driveRace + teamRace, (driveSprint + teamSprint) + (driveRace + teamRace)
}
