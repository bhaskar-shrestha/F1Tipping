# F1 Tipping App - Docker Setup

## Status: Working

All services are now containerized and running successfully.

## Quick Start

### Build and Run All Services

```bash
# Start all services (builds and runs)
docker-compose up

# Or run in background
docker-compose up -d

# View logs
docker-compose logs -f
```

### Access the Services

| Service | URL | Description |
|---------|-----|-------------|
| Web UI | http://localhost | React frontend |
| Backend API | http://localhost:8080 | Go REST API |

### API Endpoints

```bash
# Get all drivers
curl http://localhost:8080/api/admin/drivers

# Get all teams
curl http://localhost:8080/api/admin/teams

# Add a driver
curl -X POST http://localhost:8080/api/admin/drivers \
  -H "Content-Type: application/json" \
  -d '{"name":"Max Verstappen","constructor_id":"ver","constructor_name":"Red Bull"}'

# Submit a prediction
curl -X POST http://localhost:8080/api/predictions \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user1","driver_ids":["ver","nor","pi","alp","sar"],"team_ids":["redbull","mcl"]}'

# Get predictions for a user
curl http://localhost:8080/api/predictions/user/user1

# Update race positions
curl -X POST http://localhost:8080/api/admin/race-positions \
  -H "Content-Type: application/json" \
  -d '{"constructor_id":"team_Red Bull","car1":1,"car2":2}'

# Update sprint positions
curl -X POST http://localhost:8080/api/admin/sprint-positions \
  -H "Content-Type: application/json" \
  -d '{"constructor_id":"team_Red Bull","car1":1,"car2":2}'
```

### Stop Services

```bash
# Stop all services
docker-compose down

# Remove containers, networks, and volumes
docker-compose down -v

# Remove images too
docker-compose down -v --rmi all
```

## Building Individual Services

### Build Backend Only

```bash
cd backend
docker build -t f1-tipping-backend .
```

### Build Web UI Only

```bash
cd UI-web
docker build -t f1-tipping-web .
```

## Docker Architecture

```
┌─────────────────────────────────────────┐
│           Docker Host                   │
├─────────────────────────────────────────┤
│  Container: f1-tipping-web          │   Port 80 ──┐
│    (React UI with nginx)              │           │
│  Container: f1-tipping-backend     │   Port 8080 ─┤
│    (Go API server)                    │           │
└─────────────────────────────────────────┘
```

Both containers are connected via an internal Docker network.

## Files Created/Modified

### Dockerfiles
- `backend/Dockerfile` - Multi-stage build for Go backend
- `UI-web/Dockerfile` - Multi-stage build for React app with nginx

### Configuration
- `docker-compose.yml` - Service orchestration
- `UI-web/nginx.conf` - Nginx configuration for serving static files
- `backend/go.mod` - Go module definition
- `backend/go.sum` - Go module checksums

### Application Files
- `backend/src/main.go` - Updated to use PORT environment variable
- `UI-web/public/index.html` - React app entry point

## Development

### Run Backend Directly

```bash
cd backend
PORT=8080 go run src/main.go
```

### Run Web UI Directly

```bash
cd UI-web
npm start
```

## Troubleshooting

### Backend not starting

```bash
# Check logs
docker-compose logs backend

# Rebuild
docker-compose build --no-cache
```

### Web UI not loading

```bash
# Check logs
docker-compose logs web

# Rebuild
docker-compose build --no-cache

# Try accessing directly
curl http://localhost
```

### Port already in use

```bash
# Change ports in docker-compose.yml
# backend: Change "- 8080:8080" to "- 8081:8080"
# web: Change "- 80:80" to "- 8082:80"
```

### Backend returns empty array

This is expected behavior as the backend uses in-memory storage. To add data, use the API endpoints shown above.

## Cleanup

```bash
# Remove all containers and volumes
docker-compose down -v

# Remove all images
docker system prune -f
```
