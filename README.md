# SaveSync

A modern, deduplicating backup solution with content-defined chunking (CDC) and support for multiple storage backends.

## Features

- **Content-Defined Chunking (CDC)**: Efficient deduplication using rolling hash algorithm
- **Multiple Storage Backends**: Local filesystem, S3-compatible storage, and SFTP
- **Web Interface**: Modern React-based UI for managing backups, snapshots, and targets
- **RESTful API**: Complete HTTP API with Swagger documentation
- **Scheduled Backups**: Configurable backup schedules (manual, hourly, daily, weekly, cron)
- **Snapshot Management**: Browse file trees, restore snapshots, and download manifests
- **Job Tracking**: Real-time backup job status and history

## Architecture

### Backend (Go)
- **HTTP Server**: Chi router with middleware support
- **Database**: SQLite for metadata storage
- **Storage Abstraction**: Pluggable backend interface
- **Logging**: Structured logging with Zap
- **Metrics**: Prometheus metrics for observability

### Frontend (React + TypeScript)
- **Framework**: Vite + React 18 + TypeScript
- **UI Components**: shadcn/ui with Tailwind CSS
- **State Management**: React Query (server state) + Zustand (client state)
- **Routing**: React Router v6
- **Styling**: Dark/light mode support with semantic theming

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/savesync.git
cd savesync

# Start services
docker compose up -d

# Access the application
# Frontend: http://localhost:5173
# Backend API: http://localhost:8080/api
# Swagger Docs: http://localhost:8080/swagger/index.html
```

### Manual Setup

#### Backend

```bash
cd backend

# Install dependencies
go mod download

# Generate Swagger documentation
swag init -g cmd/savesyncd/main.go

# Run the server
go run cmd/savesyncd/main.go
```

#### Frontend

```bash
cd frontend

# Install dependencies
pnpm install

# Start development server
pnpm dev
```

## Configuration

### Environment Variables

Backend configuration:
- `PORT`: HTTP server port (default: 8080)
- `DB_PATH`: SQLite database path (default: `./data/savesync.db`)
- `DATA_DIR`: Data directory for local storage (default: `./data`)

Frontend configuration:
- `VITE_API_URL`: Backend API URL (default: `http://localhost:8080/api`)

## API Documentation

Interactive API documentation is available via Swagger UI when running the backend:
```
http://localhost:8080/swagger/index.html
```

To regenerate Swagger documentation:
```bash
cd backend
swag init -g cmd/savesyncd/main.go
```

## Testing

### Backend Tests

```bash
cd backend

# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### CI/CD

The project includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that:
- Runs backend tests
- Generates Swagger documentation
- Builds Docker images for both backend and frontend

## Project Structure

```
savesync/
├── backend/
│   ├── cmd/savesyncd/        # Application entry point
│   ├── internal/
│   │   ├── app/              # Business logic layer
│   │   ├── domain/           # Domain models and interfaces
│   │   └── infra/            # Infrastructure (HTTP, DB, storage)
│   └── docs/                 # Generated Swagger docs
├── frontend/
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── hooks/            # React Query hooks
│   │   ├── pages/            # Route pages
│   │   └── lib/              # Utilities and API client
│   └── public/               # Static assets
└── docker-compose.yml        # Docker orchestration
```

## Storage Backends

### Local Filesystem
Store backups on the local filesystem.

### S3-Compatible
Compatible with Amazon S3 and S3-compatible services (MinIO, DigitalOcean Spaces, etc.).

### SFTP
Backup to remote servers via SFTP protocol.

## Development

### Prerequisites
- Go 1.24+
- Node.js 20+
- pnpm (for frontend)
- Docker & Docker Compose (optional)

### Building for Production

```bash
# Using Docker Compose
docker compose build

# Manual builds
# Backend
cd backend
go build -o savesyncd cmd/savesyncd/main.go

# Frontend
cd frontend
pnpm build
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with [Go](https://golang.org/)
- UI powered by [shadcn/ui](https://ui.shadcn.com/)
- Icons from [Lucide](https://lucide.dev/)
