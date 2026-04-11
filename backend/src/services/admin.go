package services

import (
	"fmt"

	"github.com/f1tipping/backend/src/models"
	"github.com/f1tipping/backend/src/repositories"
)

// AdminService handles CRUD operations for drivers and teams using repositories
type AdminService struct {
	driversRepo  *repositories.DriverRepository
	teamsRepo    *repositories.TeamRepository
}

// NewAdminService creates a new AdminService with repositories
func NewAdminService(driversRepo *repositories.DriverRepository, teamsRepo *repositories.TeamRepository) *AdminService {
	return &AdminService{
		driversRepo:  driversRepo,
		teamsRepo:    teamsRepo,
	}
}

// AddDriver adds a new driver
func (s *AdminService) AddDriver(name string, constructorID string) (*models.Driver, error) {
	driver := &models.Driver{
		ID:           fmt.Sprintf("driver_%s", constructorID),
		Name:         name,
		ConstructorID: constructorID,
	}

	// Create driver in database
	if err := s.driversRepo.Create(driver); err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	return driver, nil
}

// GetAllDrivers retrieves all drivers with their constructor names
func (s *AdminService) GetAllDrivers() ([]models.Driver, error) {
	drivers, err := s.driversRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all drivers: %w", err)
	}

	// Fetch constructor names for each driver
	constructorCache := make(map[string]string)
	for i := range drivers {
		constructorName, err := s.driversRepo.GetConstructor(drivers[i].ConstructorID)
		if err == nil {
			drivers[i].ConstructorName = constructorName
			constructorCache[drivers[i].ConstructorID] = constructorName
		}
	}

	return drivers, nil
}

// AddTeam adds a new team (constructor)
func (s *AdminService) AddTeam(constructorName string) (*models.Team, error) {
	team := &models.Team{
		ConstructorID:      fmt.Sprintf("team_%s", constructorName),
		ConstructorName:    constructorName,
		RaceCar1Position:   nil,
		RaceCar2Position:   nil,
		SprintCar1Position: nil,
		SprintCar2Position: nil,
	}

	// Create team in database
	if err := s.teamsRepo.Create(team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	// Update team with info from database
	updatedTeam, err := s.teamsRepo.GetByID(team.ConstructorID)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh team: %w", err)
	}
	updatedTeam.ID = team.ID

	return updatedTeam, nil
}

// GetAllTeams retrieves all teams
func (s *AdminService) GetAllTeams() ([]models.Team, error) {
	return s.teamsRepo.GetAll()
}

// UpdateRacePositions updates race positions for a team
func (s *AdminService) UpdateRacePositions(constructorID string, car1Pos *int, car2Pos *int) error {
	if err := s.teamsRepo.UpdateRacePositions(constructorID, car1Pos, car2Pos); err != nil {
		return fmt.Errorf("failed to update race positions: %w", err)
	}
	return nil
}

// UpdateSprintPositions updates sprint positions for a team
func (s *AdminService) UpdateSprintPositions(constructorID string, car1Pos *int, car2Pos *int) error {
	if err := s.teamsRepo.UpdateSprintPositions(constructorID, car1Pos, car2Pos); err != nil {
		return fmt.Errorf("failed to update sprint positions: %w", err)
	}
	return nil
}

// GetTeamWithPositions retrieves a team with position information
func (s *AdminService) GetTeamWithPositions(constructorID string) (*models.Team, error) {
	return s.teamsRepo.GetTeamWithPositions(constructorID)
}

// GetDriversByConstructor retrieves all drivers for a constructor
func (s *AdminService) GetDriversByConstructor(constructorID string) ([]models.Driver, error) {
	return s.driversRepo.GetByConstructor(constructorID)
}
