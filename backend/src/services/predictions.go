package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/f1tipping/backend/src/models"
	"github.com/f1tipping/backend/src/repositories"
)

// predictionID generates a unique prediction ID
func predictionID(userID string) string {
	return fmt.Sprintf("pred_%s_%d", userID, time.Now().UnixNano())
}

// PredictionService handles prediction submissions using repositories
type PredictionService struct {
	predictionsRepo *repositories.PredictionsRepository
}

// NewPredictionService creates a new PredictionService with repository
func NewPredictionService(predictionsRepo *repositories.PredictionsRepository) *PredictionService {
	return &PredictionService{
		predictionsRepo: predictionsRepo,
	}
}

// CreatePrediction creates a new prediction with validation
// Validates: exactly 5 driver IDs and 2 team IDs
func (s *PredictionService) CreatePrediction(userID string, driverIDs []string, teamIDs []string) (*models.Prediction, error) {
	// Validate: exactly 5 drivers
	if len(driverIDs) != 5 {
		return nil, errors.New("must select exactly 5 drivers")
	}

	// Validate: exactly 2 teams
	if len(teamIDs) != 2 {
		return nil, errors.New("must select exactly 2 teams")
	}

	// Check for duplicate driver IDs
	driverSet := make(map[string]bool)
	for _, id := range driverIDs {
		if driverSet[id] {
			return nil, errors.New("duplicate driver selected")
		}
		driverSet[id] = true
	}

	// Check for duplicate team IDs
	teamSet := make(map[string]bool)
	for _, id := range teamIDs {
		if teamSet[id] {
			return nil, errors.New("duplicate team selected")
		}
		teamSet[id] = true
	}

	// Create prediction with current timestamp
	prediction := &models.Prediction{
		ID:         predictionID(userID),
		UserID:     userID,
		SubmitTime: time.Now().UTC().Format(time.RFC3339),
		DriverIDs:  driverIDs,
		TeamIDs:    teamIDs,
	}

	// Create prediction in database
	if err := s.predictionsRepo.Create(prediction); err != nil {
		return nil, fmt.Errorf("failed to create prediction: %w", err)
	}

	return prediction, nil
}

// GetPrediction retrieves a prediction by ID
func (s *PredictionService) GetPrediction(predID string) (*models.Prediction, error) {
	return s.predictionsRepo.GetByID(predID)
}

// GetPredictionsByUser retrieves all predictions for a user
func (s *PredictionService) GetPredictionsByUser(userID string) ([]models.Prediction, error) {
	return s.predictionsRepo.GetByUser(userID)
}

// CalculatePoints calculates sprint and race points for a prediction
func (s *PredictionService) CalculatePoints(userID string, driverResults map[string]int, teamPositions map[string]struct {
	Car1 int `json:"car1"`
	Car2 int `json:"car2"`
}) error {
	// Find the most recent prediction for the user
	predictions, err := s.predictionsRepo.GetByUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user predictions: %w", err)
	}

	if len(predictions) == 0 {
		return errors.New("no predictions found for user")
	}

	// Get the most recent prediction
	prediction := predictions[0]

	// Calculate points
	totalSprintPoints := CalculateDriverPoints(prediction.DriverIDs, driverResults, SprintPoints)
	totalRacePoints := CalculateDriverPoints(prediction.DriverIDs, driverResults, RacePoints)
	totalTeamSprintPoints := CalculateTeamPoints(prediction.TeamIDs, teamPositions, SprintPoints)
	totalTeamRacePoints := CalculateTeamPoints(prediction.TeamIDs, teamPositions, RacePoints)

	sprintPoints := totalSprintPoints + totalTeamSprintPoints
	racePoints := totalRacePoints + totalTeamRacePoints
	totalPoints := sprintPoints + racePoints

	// Update points in database
	if err := s.predictionsRepo.UpdatePoints(prediction.ID, sprintPoints, racePoints, totalPoints); err != nil {
		return fmt.Errorf("failed to update points: %w", err)
	}

	// Update the local prediction struct
	prediction.SprintPoints = sprintPoints
	prediction.RacePoints = racePoints
	prediction.TotalPoints = totalPoints

	return nil
}
