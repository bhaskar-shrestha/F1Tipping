# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Required Skills - F1 Racing 2026

This app focuses on Formula 1 **2026 racing rules only**. Key racing concepts:

- **22 Grands Prix** per season (March to December)
- **6 Sprint weekends**: Chinese, Miami, Canadian, British, Dutch, and Singapore Grands Prix
- **Qualifying**: Knock-out format (Q1: 18 min, Q2: 15 min, Q3: 12 min)
- **Overtaking**: DRS replaced by new overtake mode; active aero in front/rear wings
- **Pit Stops**: At least one required; both tyre compounds must be used; 4 tyres under 2 seconds
- **Race Finish**: Ends when leader completes predetermined laps

### Points System (2026)

| Sprint Race (top 8) | Race (top 10) |
|---------------------|---------------|
| 1st: 8 | 1st: 25 |
| 2nd: 7 | 2nd: 18 |
| 3rd: 6 | 3rd: 15 |
| 4th: 5 | 4th: 12 |
| 5th: 4 | 5th: 10 |
| 6th: 3 | 6th: 8 |
| 7th: 2 | 7th: 6 |
| 8th: 1 | 8th: 4 |
| 9th: 0 | 9th: 2 |
| 10th: 0 | 10th: 1 |

Teams accumulate points from both cars. No fastest lap bonus in 2026.

## Tech Stack

| Layer | Technology |
|-------|------------|
| **Web UI** | React.js |
| **Mobile App** | React Native |
| **Backend API** | Go |
| **Database** | PostgreSQL (schema in `backend/src/db/schema.sql`) |

## Implementation Status

**Status: COMPLETE** - All core features implemented and refined

### Recent Refinements (April 2026)

**Phase 1 - Schema Foundation**
- Updated database schema with `races` table for 2026 F1 calendar management
- Removed redundant `constructors.id` column (kept `constructor_id` as natural key)
- Fixed foreign key references across all result tables
- Seeded 2026 race calendar (28 races: 14 weekends × 2 types)
- Schema now loaded from file (`schema.sql`) instead of hardcoded

**Phase 2 - Go Models**
- Removed redundant `ConstructorID` field from Team model (kept only `ID`)
- Simplified result models: renamed `PositionResult` → `ConstructorPosition`
- Unified race/sprint results into single `Positions[]` array structure
- Added `race_id` field to result models for proper race context

**Phase 3 - API Endpoints**
- Added `GET /api/predictions/user/:userId` for querying user predictions
- Implemented `POST /api/admin/race-positions` for updating race results
- Implemented `POST /api/admin/sprint-positions` for updating sprint results
- All endpoints return consistent snake_case JSON fields

**Phase 4 - Frontend Updates**
- Web UI: Fixed API endpoint paths and field name mappings
- Mobile UI: Connected DriverSelectionScreen, TeamSelectionScreen, ResultsScreen to live API
- Both platforms now use correct snake_case field names from backend
- Mobile fallback to mock data when API unavailable

See `PROJECT_PLAN.md` for detailed project documentation.

### Backend (Go) - COMPLETE

## API Endpoints

### Admin Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/drivers` | Get all drivers |
| POST | `/api/admin/drivers` | Add driver (body: `{name, constructor_id, constructor_name}`) |
| GET | `/api/admin/teams` | Get all teams |
| POST | `/api/admin/teams` | Add team (body: `{constructor_name}`) |
| POST | `/api/admin/race-positions` | Update race positions |
| POST | `/api/admin/sprint-positions` | Update sprint positions |

### User Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/predictions` | Submit prediction (body: `{user_id, driver_ids[x5], team_ids[x2]}`) |
| GET | `/api/predictions/:id` | Get prediction by ID |
| GET | `/api/predictions/user/:userId` | Get all predictions for a user |

### UI-Web (React.js) - COMPLETE
- **Admin Page**: Add/edit drivers and teams
- **User Page**: Select 5 drivers and 2 teams for predictions
- **Results Display**: Show earned points

### UI-Mobile (React Native) - COMPLETE
- **DriverSelectionScreen**: Mobile-optimized driver selection
- **TeamSelectionScreen**: Mobile-optimized team selection
- **ResultsScreen**: Display results on mobile

### Testing - COMPLETE
- `points_calculation_test.go`: Unit tests for 2026 points rules
- `validation_test.go`: Unit tests for input validation

## How to Run

1. **Backend (Go)**: Requires Go 1.20+
   ```
   cd backend
   go run src/main.go
   ```
   Server runs on `http://localhost:8080`

2. **Web UI (React)**: Requires Node.js 16+
   ```
   cd UI-web
   npm install
   npm start
   ```

3. **Mobile (React Native)**: Requires Expo CLI
   ```
   cd UI-mobile
   npm install
   npx expo start
   ```

## Testing

```bash
cd testing
go test -v .
```

## Development Guidelines

- Keep UI implementations modular and reusable
- Backend provides consistent API contract for all UI clients
- Each skill/UI implementation is in its own folder
- Admin flow: Enter driver/team data → User selects → After sprint/race, update positions → System calculates points
