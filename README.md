# F1 Prediction App

A Formula 1 prediction/tipping web and mobile application where users select drivers and teams to predict performance before qualifying, sprint, and race events.

## Features

- **Driver Selection**: Select 5 drivers to predict
- **Team Selection**: Select 2 teams to predict
- **2026 F1 Rules**: Implements current F1 2026 points system
- **Cross-Platform**: Web interface (React.js) and Mobile app (React Native)
- **Admin Dashboard**: Manage drivers, teams, and race results
- **Docker Support**: Full containerization for easy deployment

## Tech Stack

| Component | Technology |
|-----------|------------|
| Web UI | React.js |
| Mobile App | React Native |
| Backend API | Go |
| Containerization | Docker |

## Project Structure

```
F1Tipping/
├── backend/              # Go backend API
│   └── src/
│       ├── models/      # Data models (driver, team, prediction, result)
│       ├── services/    # Business logic (admin, predictions, points)
│       ├── routes/      # API route handlers
│       ├── db/          # PostgreSQL schema
│       └── main.go      # HTTP server
├── UI-web/              # React.js web interface
│   └── src/
│       ├── pages/       # Admin and User views
│       ├── api/         # API client
│       └── App.js       # Main app component
├── UI-mobile/           # React Native mobile app
│   └── src/
│       ├── screens/     # Driver/Team selection screens
│       └── App.js       # Mobile app entry
└── testing/             # Go unit tests
```

## Getting Started

### Quick Start with Docker (Recommended)

The easiest way to run the application is using Docker:

```bash
# Start all services
docker-compose up

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

- Web UI: http://localhost
- Backend API: http://localhost:8080

See [DOCKER_README.md](DOCKER_README.md) for detailed Docker instructions.

### Prerequisites

- Docker and Docker Compose
- OR Go 1.20+, Node.js 16+, npm 8+ (for non-docker setup)

### Backend Setup

```bash
cd backend
go run src/main.go
```

The API server will start on `http://localhost:8080`

### Web UI Setup

```bash
cd UI-web
npm install
npm start
```

### Mobile App Setup

```bash
cd UI-mobile
npm install
npx expo start
```

## F1 2026 Points System

### Sprint Race (6 events)
Positions 1-8 receive points: 8, 7, 6, 5, 4, 3, 2, 1

### Main Race (22 events)
Positions 1-10 receive points: 25, 18, 15, 12, 10, 8, 6, 4, 2, 1

### Team Points
Both cars from the same constructor accumulate points separately.

## API Documentation

### Admin Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/drivers` | List all drivers |
| POST | `/api/admin/drivers` | Add a driver |
| GET | `/api/admin/teams` | List all teams |
| POST | `/api/admin/teams` | Add a team |

### User Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/predictions` | Submit prediction (5 drivers + 2 teams) |
| GET | `/api/predictions/:id` | Get prediction by ID |
| GET | `/api/predictions/user/:userId` | Get all predictions for a user |

## Database Schema

The application uses PostgreSQL with the following key tables:

- **drivers** - Driver information with constructor assignments
- **constructors** - Team/constructor data (using `constructor_id` as primary key)
- **predictions** - User predictions for races (stores driver_ids and team_ids as arrays)
- **sprint_results** - Sprint race results by constructor position
- **race_results** - Race results by constructor position
- **races** - F1 2026 race calendar (sprint and main races)

2026 Race Calendar:
- **22 Total Races** (14 race weekends)
- **6 Sprint Weekends**: China, Miami, Canada, Britain, Netherlands, Singapore
- **Automatic Seeding**: 28 race records seeded on startup (14 races × 2 types: sprint + main)

For detailed schema, see [schema.sql](backend/src/db/schema.sql).

## Development

- See [PROJECT_PLAN.md](PROJECT_PLAN.md) for detailed implementation
- See [CLAUDE.md](CLAUDE.md) for development guidelines
- Tests: `cd testing && go test -v .`

## License

Apache-2.0 license
