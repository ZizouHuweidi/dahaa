# Dahaa - Multiplayer Bluffing Game

Dahaa is a multiplayer bluffing and deduction game where players attempt to fool each other with fake answers while identifying the correct one. The name "Dahaa" is an Arabic word referring to clever trickery or deception, aligning with the game's theme.

## Features

- Real-time multiplayer gameplay
- Multiple categories (History, Pop Culture, Science, etc.)
- Cross-platform support (iOS, Android, Web)
- Smart answer management to prevent duplicates
- Real-time scoring and game state updates

## Tech Stack

### Frontend
- React Native with Expo
- NativeWind for styling
- WebSocket for real-time communication
- React Native Web for web support

### Backend
- Go with Echo framework
- Gorilla WebSocket
- PostgreSQL for persistent storage
- Redis for game state management

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later
- Node.js 18 or later
- Expo CLI
- Make (for using Makefile commands)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/zizouhuweidi/dahaa.git
cd dahaa
```

2. Set up environment variables:
```bash
# Backend
cp backend/.env.example backend/.env
# Frontend
cp frontend/.env.example frontend/.env
```

3. Install dependencies:
```bash
make install-deps
```

4. Start the backend services (database, Redis, and backend server):
```bash
make up
```

5. In a separate terminal, start the frontend development server:
```bash
make frontend
```

## Development Commands

### Backend (Docker)
- `make up` - Start all services (PostgreSQL, Redis, backend)
- `make down` - Stop all services
- `make backend` - Start backend service without Docker
- `make test-backend` - Run backend tests
- `make db-migrate` - Run database migrations
- `make db-rollback` - Rollback database migrations
- `make lint-backend` - Run backend linter
- `make format-backend` - Format backend code
- `make generate` - Generate backend code

### Frontend (Local)
- `make frontend` - Start frontend development server
- `make test-frontend` - Run frontend tests
- `make lint-frontend` - Run frontend linter
- `make format-frontend` - Format frontend code

### Maintenance
- `make clean` - Remove all containers and volumes
- `make install-deps` - Install all dependencies

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
