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

3. Start the development environment:
```bash
make up
```

This will start:
- PostgreSQL database
- Redis cache
- Backend server
- Frontend development server

- `make up` - Start all services
- `make down` - Stop all services
- `make backend` - Start only the backend service
- `make frontend` - Start only the frontend service
- `make logs` - View logs from all services
- `make clean` - Remove all containers and volumes

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
