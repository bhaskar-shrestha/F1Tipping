-- F1 Prediction App Database Schema (CORRECTED VERSION)
-- PostgreSQL
-- Authoritative schema file - load from application instead of hardcoding

-- ============================================================================
-- RACES TABLE (Sports Management)
-- ============================================================================
-- Tracks race calendar for seasons
-- Allows tracking multiple seasons and linking results to specific races
CREATE TABLE IF NOT EXISTS races (
    race_id SERIAL PRIMARY KEY,
    race_name VARCHAR(100) NOT NULL,          -- "Chinese Grand Prix", "Monaco GP"
    race_date DATE NOT NULL,
    race_type VARCHAR(20) NOT NULL,           -- 'sprint' or 'main'
    season INT NOT NULL,                      -- 2026, 2027, etc.
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(season, race_name, race_type)
);

CREATE INDEX IF NOT EXISTS idx_races_season_date ON races(season, race_date);
CREATE INDEX IF NOT EXISTS idx_races_type ON races(race_type);

-- ============================================================================
-- CONSTRUCTORS/TEAMS TABLE (FIXED)
-- ============================================================================
-- FIXED: Removed redundant `id` SERIAL - use constructor_id as primary key
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
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    submit_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sprint_points INTEGER DEFAULT 0,
    race_points INTEGER DEFAULT 0,
    total_points INTEGER DEFAULT 0,
    driver_ids TEXT[],        -- PostgreSQL array format
    team_ids TEXT[],          -- PostgreSQL array format
    UNIQUE(user_id, submit_time)
);

-- ============================================================================
-- SPRINT RESULTS TABLE (FIXED)
-- ============================================================================
-- FIXED: Foreign key now references constructor_id (VARCHAR), not int id
-- NEW: Linked to races table for full race context
CREATE TABLE IF NOT EXISTS sprint_results (
    id SERIAL PRIMARY KEY,
    race_id INTEGER NOT NULL REFERENCES races(race_id) ON DELETE CASCADE,
    constructor_id VARCHAR(50) NOT NULL REFERENCES constructors(constructor_id) ON DELETE CASCADE,
    position INTEGER NOT NULL CHECK (position > 0 AND position <= 20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(race_id, constructor_id)
);

-- ============================================================================
-- RACE RESULTS TABLE (FIXED)
-- ============================================================================
-- FIXED: Foreign key now references constructor_id (VARCHAR), not int id
-- NEW: Linked to races table for full race context
CREATE TABLE IF NOT EXISTS race_results (
    id SERIAL PRIMARY KEY,
    race_id INTEGER NOT NULL REFERENCES races(race_id) ON DELETE CASCADE,
    constructor_id VARCHAR(50) NOT NULL REFERENCES constructors(constructor_id) ON DELETE CASCADE,
    position INTEGER NOT NULL CHECK (position > 0 AND position <= 20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(race_id, constructor_id)
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_drivers_constructor ON drivers(constructor_id);
CREATE INDEX IF NOT EXISTS idx_predictions_user ON predictions(user_id);
CREATE INDEX IF NOT EXISTS idx_predictions_submit ON predictions(submit_time);
CREATE INDEX IF NOT EXISTS idx_sprint_results_constructor ON sprint_results(constructor_id);
CREATE INDEX IF NOT EXISTS idx_sprint_results_race ON sprint_results(race_id);
CREATE INDEX IF NOT EXISTS idx_race_results_constructor ON race_results(constructor_id);
CREATE INDEX IF NOT EXISTS idx_race_results_race ON race_results(race_id);
