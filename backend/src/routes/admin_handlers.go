package routes

import (
	"encoding/json"
	"net/http"
	"sync"
)

// AdminService handles driver and team management
type AdminService struct {
	drivers map[string]*models.Driver
	teams   map[string]*models.Team
	mu      sync.RWMutex
}

// NewAdminService creates a new AdminService instance
func NewAdminService() *AdminService {
	return &AdminService{
		drivers: make(map[string]*models.Driver),
		teams:   make(map[string]*models.Team),
	}
}

// AddDriver adds a driver via admin API
func (s *AdminService) AddDriver(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		ConstructorID  string `json:"constructor_id"`
		ConstructorName string `json:"constructor_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	driverID := "driver_" + input.ConstructorID
	s.drivers[driverID] = &models.Driver{
		ID:           driverID,
		Name:         input.Name,
		ConstructorID: input.ConstructorID,
		ConstructorName: input.ConstructorName,
	}
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s.drivers[driverID])
}

// GetDrivers returns all drivers via admin API
func (s *AdminService) GetDrivers(w http.ResponseWriter) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	drivers := make([]*models.Driver, 0, len(s.drivers))
	for _, d := range s.drivers {
		drivers = append(drivers, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

// AddTeam adds a team via admin API
func (s *AdminService) AddTeam(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConstructorName string `json:"constructor_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	teamID := "team_" + input.ConstructorID
	s.teams[teamID] = &models.Team{
		ID:                teamID,
		ConstructorID:    teamID,
		ConstructorName:  input.ConstructorName,
		RaceCar1Position:  nil,
		RaceCar2Position:  nil,
		SprintCar1Position: nil,
		SprintCar2Position: nil,
	}
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s.teams[teamID])
}

// GetTeams returns all teams via admin API
func (s *AdminService) GetTeams(w http.ResponseWriter) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	teams := make([]*models.Team, 0, len(s.teams))
	for _, t := range s.teams {
		teams = append(teams, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

// UpdateRacePositions updates race positions for a team
func (s *AdminService) UpdateRacePositions(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConstructorID      string `json:"constructor_id"`
		RaceCar1Position   *int   `json:"race_car1_position"`
		RaceCar2Position   *int   `json:"race_car2_position"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	team, ok := s.teams[input.ConstructorID]
	if !ok {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	if input.RaceCar1Position != nil {
		team.RaceCar1Position = input.RaceCar1Position
	}
	if input.RaceCar2Position != nil {
		team.RaceCar2Position = input.RaceCar2Position
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(team)
}

// UpdateSprintPositions updates sprint positions for a team
func (s *AdminService) UpdateSprintPositions(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConstructorID      string `json:"constructor_id"`
		SprintCar1Position *int   `json:"sprint_car1_position"`
		SprintCar2Position *int   `json:"sprint_car2_position"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	team, ok := s.teams[input.ConstructorID]
	if !ok {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	if input.SprintCar1Position != nil {
		team.SprintCar1Position = input.SprintCar1Position
	}
	if input.SprintCar2Position != nil {
		team.SprintCar2Position = input.SprintCar2Position
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(team)
}
