-- F1 Prediction App Database Schema
-- PostgreSQL

-- Constructors/Teams table
CREATE TABLE IF NOT EXISTS constructors (
    id SERIAL PRIMARY KEY,
    constructor_id VARCHAR(50) UNIQUE NOT NULL,
    constructor_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Drivers table
CREATE TABLE IF NOT EXISTS drivers (
    id SERIAL PRIMARY KEY,
    driver_id VARCHAR(50) UNIQUE NOT NULL,
    driver_name VARCHAR(100) NOT NULL,
    constructor_id INTEGER REFERENCES constructors(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Predictions table
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

-- Sprint results table
CREATE TABLE IF NOT EXISTS sprint_results (
    id SERIAL PRIMARY KEY,
    constructor_id INTEGER REFERENCES constructors(id),
    position INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Race results table
CREATE TABLE IF NOT EXISTS race_results (
    id SERIAL PRIMARY KEY,
    constructor_id INTEGER REFERENCES constructors(id),
    position INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_drivers_constructor ON drivers(constructor_id);
CREATE INDEX idx_sprint_results_position ON sprint_results(position);
CREATE INDEX idx_race_results_position ON race_results(position);
