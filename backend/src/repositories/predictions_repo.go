package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/f1tipping/backend/src/models"
)

// PredictionsRepository handles database operations for predictions
type PredictionsRepository struct {
	db *sql.DB
}

// NewPredictionsRepository creates a new PredictionsRepository
func NewPredictionsRepository(db *sql.DB) *PredictionsRepository {
	return &PredictionsRepository{db: db}
}

// Create inserts a new prediction
func (r *PredictionsRepository) Create(prediction *models.Prediction) error {
	// Escape driver and team IDs for array storage
	driverIDs := make([]string, len(prediction.DriverIDs))
	for i, id := range prediction.DriverIDs {
		driverIDs[i] = id
	}
	teamIDs := make([]string, len(prediction.TeamIDs))
	for i, id := range prediction.TeamIDs {
		teamIDs[i] = id
	}

	result, err := r.db.Exec(
		"INSERT INTO predictions (user_id, submit_time, sprint_points, race_points, total_points, driver_ids, team_ids) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, user_id, submit_time, sprint_points, race_points, total_points",
		prediction.UserID,
		prediction.SubmitTime,
		prediction.SprintPoints,
		prediction.RacePoints,
		prediction.TotalPoints,
		driverIDs,
		teamIDs,
	)
	if err != nil {
		return fmt.Errorf("failed to create prediction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("prediction insert failed: no rows affected")
	}

	return nil
}

// GetByID retrieves a prediction by ID
func (r *PredictionsRepository) GetByID(id string) (*models.Prediction, error) {
	var p models.Prediction
	var sprintPoints, racePoints, totalPoints sql.NullInt32
	var driverIDs, teamIDs sql.NullString

	err := r.db.QueryRow(
		"SELECT user_id, submit_time, sprint_points, race_points, total_points, driver_ids, team_ids FROM predictions WHERE id = $1",
		id,
	).Scan(&p.UserID, &p.SubmitTime, &sprintPoints, &racePoints, &totalPoints, &driverIDs, &teamIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction by id: %w", err)
	}

	p.ID = id

	if sprintPoints.Valid {
		p.SprintPoints = int(sprintPoints.Int32)
	}
	if racePoints.Valid {
		p.RacePoints = int(racePoints.Int32)
	}
	if totalPoints.Valid {
		p.TotalPoints = int(totalPoints.Int32)
	}

	// Parse arrays
	if driverIDs.Valid {
		p.DriverIDs, err = parseStringArray(driverIDs.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse driver_ids: %w", err)
		}
	}
	if teamIDs.Valid {
		p.TeamIDs, err = parseStringArray(teamIDs.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse team_ids: %w", err)
		}
	}

	return &p, nil
}

// GetByUser retrieves all predictions for a user
func (r *PredictionsRepository) GetByUser(userID string) ([]models.Prediction, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, submit_time, sprint_points, race_points, total_points, driver_ids, team_ids FROM predictions WHERE user_id = $1 ORDER BY submit_time DESC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query predictions for user: %w", err)
	}
	defer rows.Close()

	var predictions []models.Prediction
	for rows.Next() {
		var p models.Prediction
		var sprintPoints, racePoints, totalPoints sql.NullInt32
		var driverIDs, teamIDs sql.NullString

		if err := rows.Scan(
			&p.ID, &p.UserID, &p.SubmitTime, &sprintPoints, &racePoints, &totalPoints, &driverIDs, &teamIDs,
		); err != nil {
			return nil, fmt.Errorf("failed to scan prediction: %w", err)
		}

		if sprintPoints.Valid {
			p.SprintPoints = int(sprintPoints.Int32)
		}
		if racePoints.Valid {
			p.RacePoints = int(racePoints.Int32)
		}
		if totalPoints.Valid {
			p.TotalPoints = int(totalPoints.Int32)
		}

		if driverIDs.Valid {
			p.DriverIDs, err = parseStringArray(driverIDs.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse driver_ids: %w", err)
			}
		}
		if teamIDs.Valid {
			p.TeamIDs, err = parseStringArray(teamIDs.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse team_ids: %w", err)
			}
		}

		predictions = append(predictions, p)
	}

	return predictions, nil
}

// GetByUserAndTime retrieves prediction for a user at a specific time
func (r *PredictionsRepository) GetByUserAndTime(userID string, submitTime time.Time) (*models.Prediction, error) {
	var p models.Prediction
	var sprintPoints, racePoints, totalPoints sql.NullInt32
	var driverIDs, teamIDs sql.NullString

	err := r.db.QueryRow(
		"SELECT user_id, submit_time, sprint_points, race_points, total_points, driver_ids, team_ids FROM predictions WHERE user_id = $1 AND submit_time = $2",
		userID, submitTime,
	).Scan(&p.UserID, &p.SubmitTime, &sprintPoints, &racePoints, &totalPoints, &driverIDs, &teamIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get prediction by user and time: %w", err)
	}

	if sprintPoints.Valid {
		p.SprintPoints = int(sprintPoints.Int32)
	}
	if racePoints.Valid {
		p.RacePoints = int(racePoints.Int32)
	}
	if totalPoints.Valid {
		p.TotalPoints = int(totalPoints.Int32)
	}

	if driverIDs.Valid {
		p.DriverIDs, err = parseStringArray(driverIDs.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse driver_ids: %w", err)
		}
	}
	if teamIDs.Valid {
		p.TeamIDs, err = parseStringArray(teamIDs.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse team_ids: %w", err)
		}
	}

	return &p, nil
}

// UpdatePoints updates sprint, race, and total points for a prediction
func (r *PredictionsRepository) UpdatePoints(predID string, sprintPoints int, racePoints int, totalPoints int) error {
	result, err := r.db.Exec(
		"UPDATE predictions SET sprint_points = $1, race_points = $2, total_points = $3 WHERE id = $4",
		sprintPoints, racePoints, totalPoints, predID,
	)
	if err != nil {
		return fmt.Errorf("failed to update points: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("prediction update failed: no rows affected")
	}

	return nil
}

// GetAll retrieves all predictions
func (r *PredictionsRepository) GetAll() ([]models.Prediction, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, submit_time, sprint_points, race_points, total_points, driver_ids, team_ids FROM predictions ORDER BY submit_time DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query all predictions: %w", err)
	}
	defer rows.Close()

	var predictions []models.Prediction
	for rows.Next() {
		var p models.Prediction
		var sprintPoints, racePoints, totalPoints sql.NullInt32
		var driverIDs, teamIDs sql.NullString

		if err := rows.Scan(
			&p.ID, &p.UserID, &p.SubmitTime, &sprintPoints, &racePoints, &totalPoints, &driverIDs, &teamIDs,
		); err != nil {
			return nil, fmt.Errorf("failed to scan prediction: %w", err)
		}

		if sprintPoints.Valid {
			p.SprintPoints = int(sprintPoints.Int32)
		}
		if racePoints.Valid {
			p.RacePoints = int(racePoints.Int32)
		}
		if totalPoints.Valid {
			p.TotalPoints = int(totalPoints.Int32)
		}

		if driverIDs.Valid {
			p.DriverIDs, err = parseStringArray(driverIDs.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse driver_ids: %w", err)
			}
		}
		if teamIDs.Valid {
			p.TeamIDs, err = parseStringArray(teamIDs.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse team_ids: %w", err)
			}
		}

		predictions = append(predictions, p)
	}

	return predictions, nil
}

// parseStringArray parses a PostgreSQL text array string
func parseStringArray(s string) ([]string, error) {
	if s == "" {
		return []string{}, nil
	}
	// Handle PostgreSQL array format: "{"elem1", "elem2", "elem3"}"
	s = s[1 : len(s)-1] // Remove braces
	elements := make([]string, 0)
	current := ""
	inQuotes := false
	for _, c := range s {
		if c == '"' {
			inQuotes = !inQuotes
		}
		if c == ',' && !inQuotes {
			elements = append(elements, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		elements = append(elements, current)
	}
	return elements, nil
}
