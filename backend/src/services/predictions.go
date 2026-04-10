package services

import (
	"errors"
	"sync"
)

// PredictionService handles prediction submissions
type PredictionService struct {
	predictions map[string]*models.Prediction
	mu          sync.RWMutex
}

// NewPredictionService creates a new PredictionService instance
func NewPredictionService() *PredictionService {
	return &PredictionService{
		predictions: make(map[string]*models.Prediction),
	}
}

// CreatePrediction creates a new prediction with validation
// Validates: exactly 5 driver IDs and 2 team IDs
func (s *PredictionService) CreatePrediction(userID string, driverIDs []string, teamIDs []string) (*models.Prediction, error) {
	if len(driverIDs) != 5 {
		return nil, errors.New("must select exactly 5 drivers")
	}

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

	s.mu.Lock()
	defer s.mu.Unlock()

	// Create prediction with 0 points (will be calculated after results)
	prediction := &models.Prediction{
		ID:         predictionID(userID),
		UserID:     userID,
		SubmitTime: "2026-04-10T00:00:00Z",
		DriverIDs:  driverIDs,
		TeamIDs:    teamIDs,
	}

	s.predictions[prediction.ID] = prediction
	return prediction, nil
}

// GetPrediction retrieves a prediction by ID
func (s *PredictionService) GetPrediction(predID string) (*models.Prediction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pred, ok := s.predictions[predID]
	if !ok {
		return nil, errors.New("prediction not found")
	}
	return pred, nil
}

// GetPredictionsByUser retrieves all predictions for a user
func (s *PredictionService) GetPredictionsByUser(userID string) []*models.Prediction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	predictions := make([]*models.Prediction, 0, len(s.predictions))
	for _, pred := range s.predictions {
		if pred.UserID == userID {
			predictions = append(predictions, pred)
		}
	}
	return predictions
}

// CalculatePoints calculates sprint and race points for a prediction
func (s *PredictionService) CalculatePoints(driverResults map[string]int, teams map[string]*models.Team) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, pred := range s.predictions {
		pred.SprintPoints = calculatePoints(pred.DriverIDs, driverResults, SprintPoints) +
			calculateTeamPoints(pred.TeamIDs, teams, SprintPoints)
		pred.RacePoints = calculatePoints(pred.DriverIDs, driverResults, RacePoints) +
			calculateTeamPoints(pred.TeamIDs, teams, RacePoints)
		pred.TotalPoints = pred.SprintPoints + pred.RacePoints
	}
	return nil
}
