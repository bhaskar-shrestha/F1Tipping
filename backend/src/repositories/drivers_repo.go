package repositories

import (
	"database/sql"
	"fmt"

	"github.com/f1tipping/backend/src/models"
)

// DriverRepository handles database operations for drivers
type DriverRepository struct {
	db *sql.DB
}

// NewDriverRepository creates a new DriverRepository
func NewDriverRepository(db *sql.DB) *DriverRepository {
	return &DriverRepository{db: db}
}

// GetAll retrieves all drivers
func (r *DriverRepository) GetAll() ([]models.Driver, error) {
	rows, err := r.db.Query("SELECT driver_id, driver_name, constructor_id FROM drivers ORDER BY driver_id")
	if err != nil {
		return nil, fmt.Errorf("failed to query drivers: %w", err)
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		var driverID, name, constructorID string
		if err := rows.Scan(&driverID, &name, &constructorID); err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}
		d.ID = driverID
		d.Name = name
		d.ConstructorID = constructorID
		d.ConstructorName = "" // Will be set by service when fetching constructor
		drivers = append(drivers, d)
	}

	return drivers, nil
}

// GetByID retrieves a driver by ID
func (r *DriverRepository) GetByID(id string) (*models.Driver, error) {
	var d models.Driver
	var name, constructorID string
	err := r.db.QueryRow("SELECT driver_id, driver_name, constructor_id FROM drivers WHERE driver_id = $1", id).
		Scan(&d.ID, &name, &constructorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get driver by id: %w", err)
	}
	d.Name = name
	d.ConstructorID = constructorID
	d.ConstructorName = "" // Will be set by service

	return &d, nil
}

// Create inserts a new driver
func (r *DriverRepository) Create(driver *models.Driver) error {
	result, err := r.db.Exec(
		"INSERT INTO drivers (driver_id, driver_name, constructor_id) VALUES ($1, $2, $3) RETURNING driver_id, driver_name, constructor_id",
		driver.ID, driver.Name, driver.ConstructorID,
	)
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("driver insert failed: no rows affected")
	}

	return nil
}

// GetConstructor returns constructor info for a driver
func (r *DriverRepository) GetConstructor(constructorID string) (string, error) {
	var constructorName string
	err := r.db.QueryRow(
		"SELECT constructor_name FROM constructors WHERE constructor_id = $1",
		constructorID,
	).Scan(&constructorName)
	if err != nil {
		return "", fmt.Errorf("failed to get constructor: %w", err)
	}
	return constructorName, nil
}

// GetByConstructor returns all drivers for a constructor
func (r *DriverRepository) GetByConstructor(constructorID string) ([]models.Driver, error) {
	rows, err := r.db.Query("SELECT driver_id, driver_name, constructor_id FROM drivers WHERE constructor_id = $1 ORDER BY driver_id", constructorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query drivers for constructor: %w", err)
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		var driverID, name string
		var dbConstructorID string
		if err := rows.Scan(&driverID, &name, &dbConstructorID); err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}
		d.ID = driverID
		d.Name = name
		d.ConstructorID = dbConstructorID
		d.ConstructorName = constructorID
		drivers = append(drivers, d)
	}

	return drivers, nil
}
