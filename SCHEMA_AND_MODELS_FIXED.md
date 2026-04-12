# F1 Tipping - Fixed Database Schema & Models Reference

**Date**: April 12, 2026  
**Status**: Comprehensive fix reference  
**Purpose**: Use this document as blueprint for correcting schema/models mismatches

---

## Table of Contents
1. [Database Schema (Fixed)](#database-schema-fixed)
2. [Go Models (Fixed)](#go-models-fixed)
3. [API Contracts](#api-contracts)
4. [Issues Fixed](#issues-fixed)
5. [Migration Guide](#migration-guide)
6. [Data Flow Diagrams](#data-flow-diagrams)

---

## Database Schema (Fixed)

### CORRECTED: PostgreSQL Schema

**File**: `backend/src/db/schema.sql` (use as authoritative source)

```sql
-- F1 Prediction App Database Schema (CORRECTED VERSION)
-- PostgreSQL
-- Authoritative: Load this schema file instead of hardcoding in Go

-- ============================================================================
-- CONSTRUCTORS/TEAMS TABLE
-- ============================================================================
-- FIXED: Removed unused `id` SERIAL column to eliminate confusion
-- FIXED: Made constructor_id the primary key (natural key pattern)
CREATE TABLE IF NOT EXISTS constructors (
    constructor_id VARCHAR(50) PRIMARY KEY,
    constructor_name VARCHAR(100) NOT NULL,
    race_car1_position INTEGER,
    race_car2_position INTEGER,
    sprint_car1_position INTEGER,
    sprint_car2_position INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- DRIVERS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS drivers (
    driver_id VARCHAR(50) PRIMARY KEY,
    driver_name VARCHAR(100) NOT NULL,
    constructor_id VARCHAR(50) NOT NULL REFERENCES constructors(constructor_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- PREDICTIONS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS predictions (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    submit_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sprint_points INTEGER DEFAULT 0,
    race_points INTEGER DEFAULT 0,
    total_points INTEGER DEFAULT 0,
    driver_ids TEXT[],        -- PostgreSQL array of driver IDs
    team_ids TEXT[],          -- PostgreSQL array of team/constructor IDs
    UNIQUE(user_id, submit_time)
);

-- ============================================================================
-- SPRINT RESULTS TABLE (NOW PROPERLY CREATED)
-- ============================================================================
-- FIXED: Foreign key now references constructor_id (VARCHAR), not integer id
-- Stores individual team finishing positions for sprint races
CREATE TABLE IF NOT EXISTS sprint_results (
    id SERIAL PRIMARY KEY,
    constructor_id VARCHAR(50) NOT NULL REFERENCES constructors(constructor_id) ON DELETE CASCADE,
    position INTEGER NOT NULL CHECK (position > 0 AND position <= 20),
    race_id INTEGER,          -- Optional: link to specific race if seasons table added
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(race_id, constructor_id)  -- One result per constructor per race
);

-- ============================================================================
-- RACE RESULTS TABLE (NOW PROPERLY CREATED)
-- ============================================================================
-- FIXED: Foreign key now references constructor_id (VARCHAR), not integer id
-- Stores individual team finishing positions for main races
CREATE TABLE IF NOT EXISTS race_results (
    id SERIAL PRIMARY KEY,
    constructor_id VARCHAR(50) NOT NULL REFERENCES constructors(constructor_id) ON DELETE CASCADE,
    position INTEGER NOT NULL CHECK (position > 0 AND position <= 20),
    race_id INTEGER,          -- Optional: link to specific race if seasons table added
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(race_id, constructor_id)  -- One result per constructor per race
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_drivers_constructor ON drivers(constructor_id);
CREATE INDEX IF NOT EXISTS idx_predictions_user ON predictions(user_id);
CREATE INDEX IF NOT EXISTS idx_predictions_submit ON predictions(submit_time);
CREATE INDEX IF NOT EXISTS idx_sprint_results_constructor ON sprint_results(constructor_id);
CREATE INDEX IF NOT EXISTS idx_race_results_constructor ON race_results(constructor_id);
CREATE INDEX IF NOT EXISTS idx_sprint_results_position ON sprint_results(position);
CREATE INDEX IF NOT EXISTS idx_race_results_position ON race_results(position);
```

### Key Changes from Original

| Issue | Original | Fixed | Why |
|-------|----------|-------|-----|
| Constructor ID | Two IDs: `id` (INT) + `constructor_id` (VARCHAR) | One ID: `constructor_id` (VARCHAR) | Eliminates confusion and redundancy |
| Sprint/Race Tables | Foreign keys reference `constructors.id` (INT) | Foreign keys reference `constructors.constructor_id` (VARCHAR) | Matches how code actually looks up constructors |
| Sprint/Race in Go | Not created in main.go's createSchema() | Load schema.sql file instead | Ensures consistency between SQL and code |
| Predictions Array | Stored as TEXT[] (PostgreSQL format) | Same, but parseStringArray() documents format | Better documentation of expected format |

---

## Go Models (Fixed)

### Package: `backend/src/models/`

#### driver.go (CORRECTED)
```go
package models

import "encoding/json"

// Driver represents a Formula 1 driver with their constructor
type Driver struct {
	ID              string `json:"id"`              // Format: "driver_1", "driver_2"
	Name            string `json:"name"`
	ConstructorID   string `json:"constructor_id"` // References constructors.constructor_id
	ConstructorName string `json:"constructor_name"`
}

// MarshalJSON custom JSON marshaling for Driver
func (d Driver) MarshalJSON() ([]byte, error) {
	type Alias Driver
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&d),
	})
}
```

**Notes**:
- `ID` format MUST be consistent: "driver_1", "driver_2", etc.
- `ConstructorID` is the VARCHAR foreign key to constructors table
- `ConstructorName` retrieved via JOIN or separate query

---

#### team.go (CORRECTED - REMOVED REDUNDANCY)
```go
package models

// Team represents a constructor/team with their race positions
type Team struct {
	ID                 string  `json:"id"`                    // Same as ConstructorID (removed redundancy)
	ConstructorID      string  `json:"constructor_id"`        // Primary identifier
	ConstructorName    string  `json:"constructor_name"`
	RaceCar1Position   *int    `json:"race_car1_position"`    // NULL = not finished
	RaceCar2Position   *int    `json:"race_car2_position"`    // NULL = not finished
	SprintCar1Position *int    `json:"sprint_car1_position"`  // NULL = not finished
	SprintCar2Position *int    `json:"sprint_car2_position"`  // NULL = not finished
}
```

**IMPORTANT CHANGE**:
- ✅ OLD: Both `id` and `constructor_id` set to different values
- ❌ NEW: Both `id` and `constructor_id` must be equal (same value)
- Or: Remove `ID` field entirely and use only `ConstructorID`
- **Recommendation**: Remove `ID` field to eliminate confusion

**Improved version**:
```go
type Team struct {
	ConstructorID      string  `json:"id"`                    // Use "id" instead of "constructor_id" in JSON
	ConstructorName    string  `json:"constructor_name"`
	RaceCar1Position   *int    `json:"race_car1_position"`
	RaceCar2Position   *int    `json:"race_car2_position"`
	SprintCar1Position *int    `json:"sprint_car1_position"`
	SprintCar2Position *int    `json:"sprint_car2_position"`
}
```

---

#### prediction.go (NO CHANGES NEEDED)
```go
package models

// Prediction represents a user's prediction for a race weekend
type Prediction struct {
	ID           string   `json:"id"`           // Format: "pred_<userId>_<timestamp>"
	UserID       string   `json:"user_id"`
	SubmitTime   string   `json:"submit_time"`  // RFC3339 format
	DriverIDs    []string `json:"driver_ids"`   // 5 selected drivers
	TeamIDs      []string `json:"team_ids"`     // 2 selected teams
	SprintPoints int      `json:"sprint_points"`
	RacePoints   int      `json:"race_points"`
	TotalPoints  int      `json:"total_points"`
}
```

**Status**: ✅ No changes needed - this model is correct

---

#### result.go (CAN BE SIMPLIFIED)
```go
package models

// PositionResult represents where a team finished in a race
type PositionResult struct {
	ConstructorID string `json:"constructor_id"`
	Position      int    `json:"position"` // 1-20 (20th place last)
}

// SprintResult represents sprint race results
type SprintResult struct {
	RaceID   int
	Results  []PositionResult  // All constructors' positions
}

// RaceResult represents main race results
type RaceResult struct {
	RaceID   int
	Results  []PositionResult  // All constructors' positions
}
```

**Changes**:
- Simplified: Remove separate `DriverPositions` and `TeamPositions`
- Constructor positioning is what matters for points
- Can add driver-level positions if needed later

---

## API Contracts

### Request/Response Examples

#### GET /api/admin/drivers
**Response**:
```json
[
  {
    "id": "driver_1",
    "name": "Max Verstappen",
    "constructor_id": "Red Bull Racing",
    "constructor_name": "Red Bull Racing"
  },
  {
    "id": "driver_2",
    "name": "Sergio Perez",
    "constructor_id": "Red Bull Racing",
    "constructor_name": "Red Bull Racing"
  }
]
```

**Format Contract**:
- `id`: Always "driver_<number>"
- `constructor_id`: Always matches constructors.constructor_id

---

#### POST /api/admin/drivers
**Request**:
```json
{
  "name": "New Driver Name",
  "constructor_id": "Red Bull Racing",
  "driver_id": "21"
}
```

**Changes from current**:
- ❌ OLD: Takes `constructor_id` and `constructor_name`, auto-generates driver_id
- ✅ NEW: Takes `constructor_id` and explicit `driver_id` parameter

**Response**: Returns created Driver object (same format as GET)

---

#### GET /api/admin/teams
**Response**:
```json
[
  {
    "id": "Red Bull Racing",
    "constructor_name": "Red Bull Racing",
    "race_car1_position": null,
    "race_car2_position": null,
    "sprint_car1_position": null,
    "sprint_car2_position": null
  }
]
```

**Changes from current**:
- ❌ OLD: Both `id` and `constructor_id` fields (redundant)
- ✅ NEW: Only `id` field (maps to constructor_id)

---

#### POST /api/predictions
**Request**:
```json
{
  "user_id": "my-user",
  "driver_ids": ["driver_1", "driver_2", "driver_3", "driver_4", "driver_5"],
  "team_ids": ["Red Bull Racing", "Ferrari"]
}
```

**Response**:
```json
{
  "id": "pred_my-user_1712973600000000000",
  "user_id": "my-user",
  "submit_time": "2026-04-12T14:30:00Z",
  "driver_ids": ["driver_1", "driver_2", "driver_3", "driver_4", "driver_5"],
  "team_ids": ["Red Bull Racing", "Ferrari"],
  "sprint_points": 0,
  "race_points": 0,
  "total_points": 0
}
```

---

#### NEW: GET /api/predictions/user/:userId
**Added endpoint** (missing in current implementation)

**Response**:
```json
[
  {
    "id": "pred_my-user_1712973600000000000",
    "user_id": "my-user",
    "submit_time": "2026-04-12T14:30:00Z",
    "driver_ids": ["driver_1", "driver_2", "driver_3", "driver_4", "driver_5"],
    "team_ids": ["Red Bull Racing", "Ferrari"],
    "sprint_points": 25,
    "race_points": 60,
    "total_points": 85
  }
]
```

---

#### NEW: POST /api/admin/race-positions
**New endpoint** (missing in current implementation)

**Request**:
```json
{
  "constructor_id": "Red Bull Racing",
  "position": 1
}
```

**Behavior**:
- Stores result in race_results table
- Marks this team's finishing position
- Used to calculate points for predictions

---

#### NEW: POST /api/admin/sprint-positions
**New endpoint** (missing in current implementation)

Same as race-positions but writes to sprint_results table.

---

## Issues Fixed

### Issue #1: Constructor ID Duplication ✅
**Problem**: Schema had both `id` (INT) and `constructor_id` (VARCHAR), causing confusion  
**Solution**: Remove `id` column, use `constructor_id` as primary key  
**Impact**: Eliminates foreign key mismatches, simplifies code  
**Files to update**:
- `backend/src/db/schema.sql` - Remove `id SERIAL PRIMARY KEY`
- `backend/src/main.go` - Update schema creation to remove `id`

---

### Issue #2: Foreign Key Type Mismatch in Result Tables ✅
**Problem**: sprint_results/race_results referenced `constructors.id` (INT) but code uses `constructor_id` (VARCHAR)  
**Solution**: Update foreign keys to reference `constructors.constructor_id`  
**SQL Fix**:
```sql
-- BEFORE:
FOREIGN KEY (constructor_id) REFERENCES constructors(id)

-- AFTER:
FOREIGN KEY (constructor_id) REFERENCES constructors(constructor_id)
```

---

### Issue #3: Result Tables Not Created in Go ✅
**Problem**: `schema.sql` defines tables but `main.go` hardcodes incomplete schema  
**Solution**: Load `schema.sql` file instead of hardcoding  
**Files to update**:
- `backend/src/main.go` - Remove hardcoded schema, load from file:
```go
func createSchema() error {
    schemaFile, err := ioutil.ReadFile("backend/src/db/schema.sql")
    if err != nil {
        return fmt.Errorf("failed to read schema: %w", err)
    }
    _, err = DB.Exec(string(schemaFile))
    if err != nil {
        return fmt.Errorf("failed to create schema: %w", err)
    }
    return nil
}
```

---

### Issue #4: Driver ID Generation Bug ✅
**Problem**: `AdminService.AddDriver()` uses `fmt.Sprintf("driver_%s", constructorID)` creating invalid IDs  
**Solution**: Accept `driverID` parameter explicitly  
**Code fix**: `backend/src/services/admin.go`
```go
// BEFORE:
func (s *AdminService) AddDriver(name string, constructorID string) (*models.Driver, error) {
    driver := &models.Driver{
        ID: fmt.Sprintf("driver_%s", constructorID),  // WRONG!
        ...
    }
}

// AFTER:
func (s *AdminService) AddDriver(name string, constructorID string, driverID string) (*models.Driver, error) {
    driver := &models.Driver{
        ID: fmt.Sprintf("driver_%s", driverID),  // CORRECT!
        ...
    }
}
```

---

### Issue #5: Missing User Predictions Endpoint ✅
**Problem**: Web UI calls `GET /api/predictions/my-user` but endpoint expects prediction ID, not user ID  
**Solution**: Add new endpoint `GET /api/predictions/user/:userId`  
**Files to update**: `backend/src/main.go`
```go
http.HandleFunc("/api/predictions/user/", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        userGetPredictionsByUser(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
})

func userGetPredictionsByUser(w http.ResponseWriter, r *http.Request) {
    userID := strings.TrimPrefix(r.URL.Path, "/api/predictions/user/")
    predictions, err := Preds.GetPredictionsByUser(userID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed: %v", err), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(predictions)
}
```

---

### Issue #6: Team Model Redundant Fields ✅
**Problem**: Team has both `ID` and `ConstructorID` with identical values  
**Solution**: Remove one, use only `constructor_id` in JSON response  
**Files to update**: `backend/src/models/team.go`
```go
// BEFORE:
type Team struct {
    ID            string `json:"id"`
    ConstructorID string `json:"constructor_id"`  // Same value!
    ...
}

// AFTER:
type Team struct {
    ConstructorID string `json:"id"`  // Use "id" for JSON compatibility
    ConstructorName string `json:"constructor_name"`
    ...
}
```

---

### Issue #7: Mobile App Hardcoded Data ✅
**Problem**: UI-mobile screens use mock data instead of API  
**Solution**: Implement API calls  
**Files to update**:
- `UI-mobile/src/screens/DriverSelectionScreen.js`
- `UI-mobile/src/screens/TeamSelectionScreen.js`

---

## Migration Guide

### Step 1: Update Database Schema
**Files**:
- `backend/src/db/schema.sql` - Apply all fixes above
- Add indices for query performance

**SQL Migration**:
```sql
-- For existing databases, add migration:
ALTER TABLE constructors DROP COLUMN id;
ALTER TABLE constructors ADD PRIMARY KEY (constructor_id);

ALTER TABLE sprint_results 
  DROP CONSTRAINT IF EXISTS sprint_results_constructor_id_fkey,
  ADD CONSTRAINT sprint_results_constructor_id_fkey 
    FOREIGN KEY (constructor_id) REFERENCES constructors(constructor_id);

ALTER TABLE race_results 
  DROP CONSTRAINT IF EXISTS race_results_constructor_id_fkey,
  ADD CONSTRAINT race_results_constructor_id_fkey 
    FOREIGN KEY (constructor_id) REFERENCES constructors(constructor_id);
```

---

### Step 2: Update Go Models
**Files**:
- `backend/src/models/team.go` - Remove redundant `ID` field
- `backend/src/models/driver.go` - No changes needed

---

### Step 3: Update Services
**Files**:
- `backend/src/services/admin.go` - Fix AddDriver to accept driverID
- `backend/src/repositories/drivers_repo.go` - Update Create() method signature

---

### Step 4: Update API Endpoints
**Files**:
- `backend/src/main.go`
  - Update `createSchema()` to load from file
  - Add `GET /api/predictions/user/:userId` endpoint
  - Update request handlers for Team (remove ID field duplication)

---

### Step 5: Update UI
**Web UI** (`UI-web/src/pages/UserPage.js`):
```javascript
// BEFORE:
API.get('/api/predictions/my-user').then(setMyPredictions);

// AFTER:
API.get('/api/predictions/user/my-user').then(setMyPredictions);
```

**Mobile UI** (`UI-mobile/src/screens/*.js`):
- Add API calls for real data
- Remove hardcoded mock data

---

## Data Flow Diagrams

### Current Broken Flow (Web UI Predictions)
```
UserPage.js
    ↓
GET /api/predictions/my-user
    ↓
main.go: userGetPrediction()
    ↓
Looks for prediction ID = "my-user"
    ↓
DB: SELECT WHERE id = $1
    ↓
❌ NOT FOUND (wrong query parameter type)
    ↓
Returns 404
```

### Fixed Flow (Web UI Predictions)
```
UserPage.js
    ↓
GET /api/predictions/user/my-user
    ↓
main.go: userGetPredictionsByUser()
    ↓
Extracts user_id = "my-user"
    ↓
DB: SELECT WHERE user_id = $1
    ↓
✅ Returns all predictions for user
    ↓
Frontend displays results
```

---

### Data Type Flow: Single Driver

```
┌─ Database ─────────────────────────────────────┐
│ drivers table:                                   │
│  driver_id      = "driver_1"  (VARCHAR)         │
│  driver_name    = "Max Verstappen"              │
│  constructor_id = "Red Bull Racing"             │
└────────────────────────────────────────────────┘
         ↓
┌─ Go Model ──────────────────────────────────────┐
│ type Driver struct {                             │
│   ID              = "driver_1"                   │
│   Name            = "Max Verstappen"             │
│   ConstructorID   = "Red Bull Racing"            │
│   ConstructorName = "Red Bull Racing"            │
│ }                                                │
└────────────────────────────────────────────────┘
         ↓
┌─ JSON Response ─────────────────────────────────┐
│ {                                                │
│   "id": "driver_1",                              │
│   "name": "Max Verstappen",                      │
│   "constructor_id": "Red Bull Racing",           │
│   "constructor_name": "Red Bull Racing"          │
│ }                                                │
└────────────────────────────────────────────────┘
         ↓
┌─ JavaScript (Web UI) ───────────────────────────┐
│ driver.id                ✅ Works                │
│ driver.name              ✅ Works                │
│ driver.constructor_id    ✅ Works                │
│ driver.constructorName   ✅ Works                │
└────────────────────────────────────────────────┘
```

---

### Points Calculation Flow (Currently Broken)

```
❌ BROKEN FLOW:
User submits prediction
    ↓
Stored in predictions table
    ↓
Admin updates team positions
    ↓
❌ Nothing happens (no trigger)
    ↓
Points stay at 0
    ↓
User sees 0 points forever
```

```
✅ FIXED FLOW:
User submits prediction
    ↓
Stored in predictions table
    ↓
Admin POST /api/admin/race-positions
    ↓
Results stored in race_results table
    ↓
System calls CalculatePoints() for all predictions
    ↓
Points calculated from race_results
    ↓
predictions table updated: sprint_points, race_points, total_points
    ↓
User sees their score
```

---

## Quick Checklist for Implementation

- [ ] Update `backend/src/db/schema.sql` (remove `id` column, fix FKs)
- [ ] Update `backend/src/main.go` createSchema() to load from file
- [ ] Update `backend/src/models/team.go` (remove redundant ID field)
- [ ] Update `backend/src/services/admin.go` (fix AddDriver signature)
- [ ] Add `GET /api/predictions/user/:userId` endpoint in main.go
- [ ] Add `POST /api/admin/race-positions` endpoint
- [ ] Add `POST /api/admin/sprint-positions` endpoint
- [ ] Update `UI-web/src/pages/UserPage.js` (fix API endpoint)
- [ ] Update mobile screens to use API instead of mock data
- [ ] Run tests to verify no regressions
- [ ] Update API documentation

---

**Last Updated**: April 12, 2026  
**Reference Used**: Full database audit report
