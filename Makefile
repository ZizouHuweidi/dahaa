
up:
	docker-compose up --build

down:
	docker-compose down

logs:
	docker-compose logs -f

backend:
	cd backend && go run cmd/server/main.go

frontend:
	cd frontend && npm start

restart:
	docker-compose down && docker-compose up --build
