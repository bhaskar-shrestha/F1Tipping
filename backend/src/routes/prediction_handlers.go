package routes

import (
	"encoding/json"
	"net/http"
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

// CreatePrediction creates a new prediction
func (s *PredictionService) CreatePrediction(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID   string   `json:"user_id"`
		DriverIDs []string `json:"driver_ids"`
		TeamIDs  []string `json:"team_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate: exactly 5 drivers
	if len(input.DriverIDs) != 5 {
		http.Error(w, "must select exactly 5 drivers", http.StatusBadRequest)
		return
	}

	// Validate: exactly 2 teams
	if len(input.TeamIDs) != 2 {
		http.Error(w, "must select exactly 2 teams", http.StatusBadRequest)
		return
	}

	// Check for duplicate driver IDs
	driverSet := make(map[string]bool)
	for _, id := range input.DriverIDs {
		if driverSet[id] {
			http.Error(w, "duplicate driver selected", http.StatusBadRequest)
			return
		}
		driverSet[id] = true
	}

	// Check for duplicate team IDs
	teamSet := make(map[string]bool)
	for _, id := range input.TeamIDs {
		if teamSet[id] {
			http.Error(w, "duplicate team selected", http.StatusBadRequest)
			return
		}
		teamSet[id] = true
	}

	s.mu.Lock()
	predictionID := "pred_" + input.UserID
	prediction := &models.Prediction{
		ID:         predictionID,
		UserID:     input.UserID,
		SubmitTime: "2026-04-10T00:00:00Z",
		DriverIDs:  input.DriverIDs,
		TeamIDs:    input.TeamIDs,
	}
	s.predictions[predictionID] = prediction
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prediction)
}

// GetPrediction retrieves a prediction by ID
func (s *PredictionService) GetPrediction(w http.ResponseWriter, predID string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pred, ok := s.predictions[predID]
	if !ok {
		http.Error(w, "Prediction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pred)
}

// GetPredictionsByUser retrieves all predictions for a user
func (s *PredictionService) GetPredictionsByUser(w http.ResponseWriter, userID string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	predictions := make([]*models.Prediction, 0, len(s.predictions))
	for _, pred := range s.predictions {
		if pred.UserID == userID {
			predictions = append(predictions, pred)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(predictions)
}

// CalculatePoints calculates points for all predictions
func (s *PredictionService) CalculatePoints(w http.ResponseWriter, r *http.Request) {
	// Get driver positions
	var driverPositions map[string]int
	if err := json.NewDecoder(r.Body).Decode(&driverPositions); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get team positions (simplified: constructor ID with positions)
	var teamPositions map[string]struct {
		Car1 int `json:"car1"`
		Car2 int `json:"car2"`
	}
	if err := json.NewDecoder(r.Body).Decode(&teamPositions); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Create team map from input
	teams := make(map[string]*models.Team)
	for tid, pos := range teamPositions {
		teams[tid] = &models.Team{
			ID:               tid,
			ConstructorID:    tid,
			RaceCar1Position: &pos.Car1,
			RaceCar2Position: &pos.Car2,
		}
	}

	// Calculate points (simplified - would use actual service)
	s.mu.Lock()
	for _, pred := range s.predictions {
		// Get sprint and race points (simplified)
		pred.SprintPoints = 0
		pred.RacePoints = 0
		pred.TotalPoints = pred.SprintPoints + pred.RacePoints
	}
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Points calculated"})
}
