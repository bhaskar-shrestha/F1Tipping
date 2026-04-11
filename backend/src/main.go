package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

// Service for managing drivers and teams
type Driver struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ConstructorID string `json:"constructor_id"`
	ConstructorName string `json:"constructor_name"`
}

type Team struct {
	ID            string `json:"id"`
	ConstructorID string `json:"constructor_id"`
	ConstructorName string `json:"constructor_name"`
}

type AdminService struct {
	mu     sync.RWMutex
	drivers []Driver
	teams   []Team
}

func NewAdminService() *AdminService {
	return &AdminService{
		drivers: []Driver{},
		teams:   []Team{},
	}
}

func (s *AdminService) AddDriver(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		ConstructorID  string `json:"constructor_id"`
		ConstructorName string `json:"constructor_name"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	s.mu.Lock()
	s.drivers = append(s.drivers, Driver{
		ID:           input.ConstructorID,
		Name:         input.Name,
		ConstructorID: input.ConstructorID,
		ConstructorName: input.ConstructorName,
	})
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func (s *AdminService) GetDrivers(w http.ResponseWriter) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.drivers)
}

func (s *AdminService) AddTeam(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConstructorName string `json:"constructor_name"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	s.mu.Lock()
	s.teams = append(s.teams, Team{
		ID:           "team_" + input.ConstructorName,
		ConstructorID: "team_" + input.ConstructorName,
		ConstructorName: input.ConstructorName,
	})
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func (s *AdminService) GetTeams(w http.ResponseWriter) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.teams)
}

// Service for managing predictions
type Prediction struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	SubmitTime string `json:"submit_time"`
	DriverIDs  []string `json:"driver_ids"`
	TeamIDs    []string `json:"team_ids"`
}

type PredictionService struct {
	mu        sync.RWMutex
	predictions []Prediction
}

func NewPredictionService() *PredictionService {
	return &PredictionService{
		predictions: []Prediction{},
	}
}

func (s *PredictionService) CreatePrediction(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID    string   `json:"user_id"`
		DriverIDs []string `json:"driver_ids"`
		TeamIDs   []string `json:"team_ids"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	if len(input.DriverIDs) != 5 {
		http.Error(w, "must select exactly 5 drivers", http.StatusBadRequest)
		return
	}

	if len(input.TeamIDs) != 2 {
		http.Error(w, "must select exactly 2 teams", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.predictions = append(s.predictions, Prediction{
		ID:         "pred_" + input.UserID,
		UserID:     input.UserID,
		SubmitTime: "2026-04-10T00:00:00Z",
		DriverIDs:  input.DriverIDs,
		TeamIDs:    input.TeamIDs,
	})
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func (s *PredictionService) GetPrediction(w http.ResponseWriter, predID string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.predictions {
		if p.ID == predID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

// Global instances
var (
	adminService  = NewAdminService()
	predService   = NewPredictionService()
)

func main() {
	log.Println("F1 Prediction API")
	log.Println("Admin: POST /api/admin/drivers - {\"name\":\"Driver Name\", \"constructor_id\":\"c1\", \"constructor_name\":\"Team Name\"}")
	log.Println("Admin: POST /api/admin/teams - {\"constructor_name\":\"Team Name\"}")
	log.Println("User: POST /api/predictions - {\"user_id\":\"u1\", \"driver_ids\":[\"d1\",\"d2\",\"d3\",\"d4\",\"d5\"], \"team_ids\":[\"t1\",\"t2\"]}")

	http.HandleFunc("/api/admin/drivers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			adminService.GetDrivers(w)
		case "POST":
			adminService.AddDriver(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/admin/teams", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			adminService.GetTeams(w)
		case "POST":
			adminService.AddTeam(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/predictions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			predService.GetPrediction(w, r.URL.Path[1:])
		case "POST":
			predService.CreatePrediction(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	log.Printf("Server listening on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
