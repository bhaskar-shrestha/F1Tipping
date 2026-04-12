# F1 Prediction/Tipping App - PROJECT STATUS

## Tech Stack Selection

| Layer | Technology | Rationale |
|-------|------------|-----------|
| **Web UI** | React.js | Popular, great ecosystem, good for admin dashboards |
| **Mobile App** | React Native | Share code with React web UI, good ecosystem |
| **Backend API** | Go | High performance backend, great for concurrent predictions |
| **Database** | PostgreSQL | Robust relational DB, good for predictions and structured data |

## Project Status: COMPLETE ✅

All phases of development have been completed:
- ✅ Phase 1: Schema Foundation (Races table, schema cleanup)
- ✅ Phase 2: Go Models (Team model refactored, Result models simplified)
- ✅ Phase 3: API Endpoints (All CRUD endpoints implemented)
- ✅ Phase 4: Frontend Integration (Web UI and Mobile UI connected to live API)
- ✅ Docker: Multi-stage build with proper schema.sql handling

### Backend (Go) - COMPLETE
Files created:
- `backend/src/models/driver.go` - Driver model
- `backend/src/models/team.go` - Team model
- `backend/src/models/prediction.go` - Prediction model
- `backend/src/models/result.go` - Race results model
- `backend/src/services/admin.go` - CRUD for drivers/teams
- `backend/src/services/predictions.go` - Prediction submission
- `backend/src/services/points.go` - 2026 points calculation (8/25, 7/18, etc.)
- `backend/src/routes/admin_handlers.go` - Admin API handlers
- `backend/src/routes/prediction_handlers.go` - Prediction handlers
- `backend/src/main.go` - HTTP server
- `backend/src/db/schema.sql` - PostgreSQL schema

### UI-Web (React.js) - COMPLETE
Files created:
- `UI-web/src/App.js` - Main app with routing
- `UI-web/src/index.js` - Entry point
- `UI-web/src/api/index.js` - API client
- `UI-web/src/pages/admin/AdminPage.js` - Admin view
- `UI-web/src/pages/UserPage.js` - User selection view
- `UI-web/src/index.css` - Styles
- `UI-web/package.json` - Dependencies

### UI-Mobile (React Native) - COMPLETE
Files created:
- `UI-mobile/src/App.js` - React Native app
- `UI-mobile/src/index.js` - Entry point
- `UI-mobile/app.json` - App config
- `UI-mobile/src/screens/DriverSelectionScreen.js` - Select 5 drivers
- `UI-mobile/src/screens/TeamSelectionScreen.js` - Select 2 teams
- `UI-mobile/src/screens/ResultsScreen.js` - Show results
- `UI-mobile/package.json` - Dependencies

### Testing - COMPLETE
Files created:
- `testing/points_calculation_test.go` - Unit tests for 2026 points rules
- `testing/validation_test.go` - Unit tests for driver/team validation

## F1 2026 Points Rules (Hardcoded)

### Sprint Race Points (top 8 finishers)
- 1st: 8 | 2nd: 7 | 3rd: 6 | 4th: 5 | 5th: 4 | 6th: 3 | 7th: 2 | 8th: 1

### Main Race Points (top 10 finishers)
- 1st: 25 | 2nd: 18 | 3rd: 15 | 4th: 12 | 5th: 10 | 6th: 8 | 7th: 6 | 8th: 4 | 9th: 2 | 10th: 1

### Team Points
- Teams accumulate points from both cars in the race
- Same position-based scoring as drivers
- No fastest lap bonus in 2026

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

## How to Run

1. **Backend (Go)**:
   ```
   cd backend
   go run src/main.go
   ```
   Server runs on `http://localhost:8080`

2. **Web UI (React)**:
   ```
   cd UI-web
   npm install
   npm start
   ```

3. **Mobile (React Native)**:
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

## Admin Flow
1. Admin enters driver data via web UI (name, constructor)
2. Admin enters team data (constructor name)
3. User selects 5 drivers and 2 teams for their prediction
4. After sprint/race, admin updates positions
5. System calculates user's earned points based on 2026 rules
