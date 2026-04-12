package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/f1tipping/backend/src/repositories"
	"github.com/f1tipping/backend/src/services"
)

// Database connection configuration
var (
	DB     *sql.DB
	Admin  *services.AdminService
	Preds  *services.PredictionService
)

func main() {
	// Initialize database
	log.Println("Initializing database connection...")
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer DB.Close()

	// Initialize repositories
	log.Println("Initializing repositories...")
	driversRepo := repositories.NewDriverRepository(DB)
	teamsRepo := repositories.NewTeamRepository(DB)
	predictionsRepo := repositories.NewPredictionsRepository(DB)

	// Initialize services
	Admin = services.NewAdminService(driversRepo, teamsRepo)
	Preds = services.NewPredictionService(predictionsRepo)

	// Initialize schema
	log.Println("Creating database schema...")
	if err := createSchema(); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Setup CORS middleware
	setupCORS()

	// Initialize drivers and teams data
	log.Println("Loading initial data...")
	if err := loadInitialData(); err != nil {
		log.Printf("Warning: Failed to load initial data: %v", err)
	}

	// Initialize routes
	initRoutes()

	log.Println("F1 Prediction API server started successfully")
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(http.DefaultServeMux)))
}

func initDB() error {
	// Get environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "f1_prediction")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "postgres")

	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPass)

	log.Printf("Connecting to PostgreSQL at %s", dbHost)

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Set connection timeout
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	DB = db
	return nil
}

func createSchema() error {
	schema := `
-- Create constructors/teams table
CREATE TABLE IF NOT EXISTS constructors (
    id SERIAL PRIMARY KEY,
    constructor_id VARCHAR(50) UNIQUE NOT NULL,
    constructor_name VARCHAR(100) NOT NULL,
    race_car1_position INTEGER,
    race_car2_position INTEGER,
    sprint_car1_position INTEGER,
    sprint_car2_position INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create drivers table
CREATE TABLE IF NOT EXISTS drivers (
    id SERIAL PRIMARY KEY,
    driver_id VARCHAR(50) UNIQUE NOT NULL,
    driver_name VARCHAR(100) NOT NULL,
    constructor_id VARCHAR(50) REFERENCES constructors(constructor_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create predictions table
CREATE TABLE IF NOT EXISTS predictions (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    submit_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sprint_points INTEGER DEFAULT 0,
    race_points INTEGER DEFAULT 0,
    total_points INTEGER DEFAULT 0,
    driver_ids TEXT[],
    team_ids TEXT[],
    UNIQUE(user_id, submit_time)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_drivers_constructor ON drivers(constructor_id);
CREATE INDEX IF NOT EXISTS idx_predictions_user ON predictions(user_id);
CREATE INDEX IF NOT EXISTS idx_predictions_submit ON predictions(submit_time);
CREATE INDEX IF NOT EXISTS idx_constructors_id ON constructors(id);
`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("Database schema created successfully")
	return nil
}

func loadInitialData() error {
	// Check if constructors table has data
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM constructors").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count constructors: %w", err)
	}

	if count > 0 {
		log.Printf("Data already exists (count: %d), skipping initial data load", count)
		return nil
	}

	log.Println("No constructors found, loading initial data...")

	// Constructor and driver data
	constructorsAndDrivers := []struct {
		constructorName string
		constructorID   string
		drivers         []struct {
			name     string
			driverID string
		}
	}{
		{
			constructorName: "Red Bull Racing",
			constructorID:   "Red Bull Racing",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Max Verstappen", driverID: "1"},
				{name: "Sergio Perez", driverID: "2"},
			},
		},
		{
			constructorName: "McLaren",
			constructorID:   "McLaren",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Lando Norris", driverID: "3"},
				{name: "Oscar Piastri", driverID: "4"},
			},
		},
		{
			constructorName: "Mercedes",
			constructorID:   "Mercedes",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Lewis Hamilton", driverID: "5"},
				{name: "George Russell", driverID: "6"},
			},
		},
		{
			constructorName: "Ferrari",
			constructorID:   "Ferrari",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Charles Leclerc", driverID: "7"},
				{name: "Carlos Sainz", driverID: "8"},
			},
		},
		{
			constructorName: "Aston Martin",
			constructorID:   "Aston Martin",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Fernando Alonso", driverID: "9"},
				{name: "Lance Stroll", driverID: "10"},
			},
		},
		{
			constructorName: "Alpine",
			constructorID:   "Alpine",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Pierre Gasly", driverID: "11"},
				{name: "Esteban Ocon", driverID: "12"},
			},
		},
		{
			constructorName: "Williams",
			constructorID:   "Williams",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Alexander Albon", driverID: "13"},
				{name: "Logan Sargeant", driverID: "14"},
			},
		},
		{
			constructorName: "AlphaTauri",
			constructorID:   "AlphaTauri",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Yuki Tsunoda", driverID: "15"},
				{name: "Liam Lawson", driverID: "16"},
			},
		},
		{
			constructorName: "Haas",
			constructorID:   "Haas",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Nico Hulkenberg", driverID: "17"},
				{name: "Kevin Magnussen", driverID: "18"},
			},
		},
		{
			constructorName: "Kick Sauber",
			constructorID:   "Kick Sauber",
			drivers: []struct{
				name     string
				driverID string
			}{
				{name: "Valtteri Bottas", driverID: "19"},
				{name: "Zhou Guanyu", driverID: "20"},
			},
		},
	}

	for _, cd := range constructorsAndDrivers {
		constructor := fmt.Sprintf("team_%s", cd.constructorID)

		// Insert constructor
		result, err := DB.Exec(
			`INSERT INTO constructors (constructor_id, constructor_name, race_car1_position, race_car2_position)
			 VALUES ($1, $2, $3, $4)`,
			constructor, cd.constructorName, 0, 0,
		)
		if err != nil {
			log.Printf("Warning: Failed to insert constructor %s: %v", cd.constructorID, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			// Insert drivers for this constructor
			for _, driver := range cd.drivers {
				_, err := DB.Exec(
					`INSERT INTO drivers (driver_id, driver_name, constructor_id) VALUES ($1, $2, $3)`,
					fmt.Sprintf("driver_%s", driver.driverID), driver.name, constructor,
				)
				if err != nil {
					log.Printf("Warning: Failed to insert driver %s: %v", driver.name, err)
				}
			}
		}
	}

	log.Println("Initial data loaded successfully")
	return nil
}

func setupCORS() {
	// CORS is handled via corsHandler middleware
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func initRoutes() {
	// Admin handlers
	http.HandleFunc("/api/admin/drivers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			adminGetDrivers(w, r)
		case http.MethodPost:
			adminAddDriver(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/admin/teams", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			adminGetTeams(w, r)
		case http.MethodPost:
			adminAddTeam(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/admin/race-positions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			adminUpdateRacePositions(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/admin/sprint-positions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			adminUpdateSprintPositions(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// User handlers
	http.HandleFunc("/api/predictions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userGetPrediction(w, r)
		case http.MethodPost:
			userCreatePrediction(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// Admin routes
func adminGetDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := Admin.GetAllDrivers()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get drivers: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(drivers)
}

func adminAddDriver(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		ConstructorID  string `json:"constructor_id"`
		ConstructorName string `json:"constructor_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	driver, err := Admin.AddDriver(input.Name, input.ConstructorID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add driver: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(driver)
}

func adminGetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := Admin.GetAllTeams()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get teams: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(teams)
}

func adminAddTeam(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConstructorName string `json:"constructor_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	team, err := Admin.AddTeam(input.ConstructorName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add team: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(team)
}

func adminUpdateRacePositions(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConstructorID string `json:"constructor_id"`
		Car1          int    `json:"car1"`
		Car2          int    `json:"car2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	car1 := input.Car1
	car2 := input.Car2

	if err := Admin.UpdateRacePositions(input.ConstructorID, &car1, &car2); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update positions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(map[string]string{"status": "Race positions updated"})
}

func adminUpdateSprintPositions(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConstructorID string `json:"constructor_id"`
		Car1          int    `json:"car1"`
		Car2          int    `json:"car2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	car1 := input.Car1
	car2 := input.Car2

	if err := Admin.UpdateSprintPositions(input.ConstructorID, &car1, &car2); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update positions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(map[string]string{"status": "Sprint positions updated"})
}

// User routes
func userCreatePrediction(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID     string   `json:"user_id"`
		DriverIDs  []string `json:"driver_ids"`
		TeamIDs    []string `json:"team_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	prediction, err := Preds.CreatePrediction(input.UserID, input.DriverIDs, input.TeamIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create prediction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(prediction)
}

func userGetPrediction(w http.ResponseWriter, r *http.Request) {
	predID := strings.TrimPrefix(r.URL.Path, "/api/predictions/")
	if predID == "" {
		http.Error(w, "Prediction ID required", http.StatusBadRequest)
		return
	}

	prediction, err := Preds.GetPrediction(predID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get prediction: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(prediction)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
