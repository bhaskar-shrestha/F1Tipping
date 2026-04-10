package services

import (
	"errors"
	"fmt"
	"sync"
)

// AdminService handles CRUD operations for drivers and teams (constructor)
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

// AddDriver adds a new driver
func (s *AdminService) AddDriver(name string, constructorID string) (*models.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	driverID := fmt.Sprintf("driver_%d", len(s.drivers)+1)
	driver := &models.Driver{
		ID:           driverID,
		Name:         name,
		ConstructorID: constructorID,
	}

	s.drivers[driverID] = driver
	return driver, nil
}

// GetDrivers returns all drivers
func (s *AdminService) GetDrivers() []*models.Driver {
	s.mu.RLock()
	defer s.mu.RUnlock()

	drivers := make([]*models.Driver, 0, len(s.drivers))
	for _, driver := range s.drivers {
		drivers = append(drivers, driver)
	}
	return drivers
}

// AddTeam adds a new team (constructor)
func (s *AdminService) AddTeam(name string) (*models.Team, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	teamID := fmt.Sprintf("team_%d", len(s.teams)+1)
	team := &models.Team{
		ID:                 teamID,
		ConstructorID:      teamID,
		ConstructorName:    name,
		RaceCar1Position:   nil,
		RaceCar2Position:   nil,
		SprintCar1Position: nil,
		SprintCar2Position: nil,
	}

	s.teams[teamID] = team
	return team, nil
}

// GetTeams returns all teams
func (s *AdminService) GetTeams() []*models.Team {
	s.mu.RLock()
	defer s.mu.RUnlock()

	teams := make([]*models.Team, 0, len(s.teams))
	for _, team := range s.teams {
		teams = append(teams, team)
	}
	return teams
}

// UpdateRacePositions updates race positions for a team (constructor)
// Positions are stored for both cars (Car1 and Car2 from same constructor)
func (s *AdminService) UpdateRacePositions(constructorID string, car1Pos *int, car2Pos *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	team, ok := s.teams[constructorID]
	if !ok {
		return errors.New("constructor not found")
	}

	if car1Pos != nil {
		team.RaceCar1Position = car1Pos
	}
	if car2Pos != nil {
		team.RaceCar2Position = car2Pos
	}

	return nil
}

// UpdateSprintPositions updates sprint positions for a team (constructor)
func (s *AdminService) UpdateSprintPositions(constructorID string, car1Pos *int, car2Pos *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	team, ok := s.teams[constructorID]
	if !ok {
		return errors.New("constructor not found")
	}

	if car1Pos != nil {
		team.SprintCar1Position = car1Pos
	}
	if car2Pos != nil {
		team.SprintCar2Position = car2Pos
	}

	return nil
}
