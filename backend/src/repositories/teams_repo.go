package repositories

import (
	"database/sql"
	"fmt"

	"github.com/f1tipping/backend/src/models"
)

// TeamRepository handles database operations for teams (constructors)
type TeamRepository struct {
	db *sql.DB
}

// NewTeamRepository creates a new TeamRepository
func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

// GetAll retrieves all teams
func (r *TeamRepository) GetAll() ([]models.Team, error) {
	rows, err := r.db.Query("SELECT constructor_id, constructor_name FROM constructors ORDER BY constructor_id")
	if err != nil {
		return nil, fmt.Errorf("failed to query constructors: %w", err)
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		var constructorID, name string
		if err := rows.Scan(&constructorID, &name); err != nil {
			return nil, fmt.Errorf("failed to scan constructor: %w", err)
		}
		t.ID = constructorID
		t.ConstructorName = name
		t.RaceCar1Position = nil
		t.RaceCar2Position = nil
		t.SprintCar1Position = nil
		t.SprintCar2Position = nil
		teams = append(teams, t)
	}

	return teams, nil
}

// GetByID retrieves a team by ID
func (r *TeamRepository) GetByID(id string) (*models.Team, error) {
	var t models.Team
	var constructorID, name string
	err := r.db.QueryRow("SELECT constructor_id, constructor_name FROM constructors WHERE constructor_id = $1", id).
		Scan(&constructorID, &name)
	if err != nil {
		return nil, fmt.Errorf("failed to get constructor by id: %w", err)
	}
	t.ID = constructorID
	t.ConstructorName = name
	t.RaceCar1Position = nil
	t.RaceCar2Position = nil
	t.SprintCar1Position = nil
	t.SprintCar2Position = nil

	return &t, nil
}

// Create inserts a new team
func (r *TeamRepository) Create(team *models.Team) error {
	err := r.db.QueryRow(
		"INSERT INTO constructors (constructor_id, constructor_name) VALUES ($1, $2) RETURNING constructor_id, constructor_name",
		team.ID, team.ConstructorName,
	).Scan(&team.ID, &team.ConstructorName)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

// UpdateRacePositions updates race positions for a team
func (r *TeamRepository) UpdateRacePositions(constructorID string, car1Pos *int, car2Pos *int) error {
	var car1PosVal, car2PosVal sql.NullInt32
	if car1Pos != nil {
		car1PosVal = sql.NullInt32{Int32: int32(*car1Pos), Valid: true}
	}
	if car2Pos != nil {
		car2PosVal = sql.NullInt32{Int32: int32(*car2Pos), Valid: true}
	}

	result, err := r.db.Exec(
		"UPDATE constructors SET race_car1_position = $1, race_car2_position = $2, updated_at = CURRENT_TIMESTAMP WHERE constructor_id = $3",
		car1PosVal, car2PosVal, constructorID,
	)
	if err != nil {
		return fmt.Errorf("failed to update race positions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("team update failed: no rows affected")
	}

	return nil
}

// UpdateSprintPositions updates sprint positions for a team
func (r *TeamRepository) UpdateSprintPositions(constructorID string, car1Pos *int, car2Pos *int) error {
	var car1PosVal, car2PosVal sql.NullInt32
	if car1Pos != nil {
		car1PosVal = sql.NullInt32{Int32: int32(*car1Pos), Valid: true}
	}
	if car2Pos != nil {
		car2PosVal = sql.NullInt32{Int32: int32(*car2Pos), Valid: true}
	}

	result, err := r.db.Exec(
		"UPDATE constructors SET sprint_car1_position = $1, sprint_car2_position = $2, updated_at = CURRENT_TIMESTAMP WHERE constructor_id = $3",
		car1PosVal, car2PosVal, constructorID,
	)
	if err != nil {
		return fmt.Errorf("failed to update sprint positions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("team update failed: no rows affected")
	}

	return nil
}

// GetTeamWithPositions retrieves a team with position info
func (r *TeamRepository) GetTeamWithPositions(constructorID string) (*models.Team, error) {
	var t models.Team
	var constructorIDStr, name string
	var car1Pos, car2Pos, sprintCar1Pos, sprintCar2Pos sql.NullInt32

	err := r.db.QueryRow(
		`SELECT constructor_id, constructor_name,
			COALESCE(race_car1_position, 0) as race_c1,
			COALESCE(race_car2_position, 0) as race_c2,
			COALESCE(sprint_car1_position, 0) as sprint_c1,
			COALESCE(sprint_car2_position, 0) as sprint_c2
		FROM constructors WHERE constructor_id = $1`,
		constructorID,
	).Scan(&constructorIDStr, &name, &car1Pos, &car2Pos, &sprintCar1Pos, &sprintCar2Pos)
	if err != nil {
		return nil, fmt.Errorf("failed to get team with positions: %w", err)
	}

	t.ID = constructorIDStr
	t.ConstructorName = name

	t.RaceCar1Position = nil
	t.RaceCar2Position = nil
	t.SprintCar1Position = nil
	t.SprintCar2Position = nil

	if car1Pos.Valid {
		pos := int(car1Pos.Int32)
		t.RaceCar1Position = &pos
	}
	if car2Pos.Valid {
		pos := int(car2Pos.Int32)
		t.RaceCar2Position = &pos
	}
	if sprintCar1Pos.Valid {
		pos := int(sprintCar1Pos.Int32)
		t.SprintCar1Position = &pos
	}
	if sprintCar2Pos.Valid {
		pos := int(sprintCar2Pos.Int32)
		t.SprintCar2Position = &pos
	}

	return &t, nil
}
