ROOT := $(shell pwd)
SERVICES := api-gateway auth-service user-service upload-service streaming-service search-service comment-service analytics-service notification-service scheduler-service

.DEFAULT_GOAL := up

.PHONY: up down build test lint dev-frontend dev-service logs ps clean

up:
	docker compose up --build -d

down:
	docker compose down

build:
	@for svc in $(SERVICES); do \
		echo "==> go build $$svc"; \
		cd services/$$svc && go build ./... 2>&1 | sed 's/^/  /'; \
		cd $(ROOT); \
	done
	@echo "==> npm run build (frontend)"
	cd frontend && npm run build 2>&1 | sed 's/^/  /'
	@echo "Build complete."

test:
	@for svc in $(SERVICES); do \
		echo "==> go test $$svc"; \
		cd services/$$svc && go test ./... 2>&1 | sed 's/^/  /'; \
		cd $(ROOT); \
	done

lint:
	@for svc in $(SERVICES); do \
		echo "==> go vet $$svc"; \
		cd services/$$svc && go vet ./... 2>&1 | sed 's/^/  /'; \
		cd $(ROOT); \
	done

dev-frontend:
	cd frontend && npm run dev

dev-service:
	@if [ -z "$(SVC)" ]; then \
		echo "Usage: make dev-service SVC=<service-name>"; \
		echo "Available: $(SERVICES)"; \
		exit 1; \
	fi
	cd services/$(SVC) && go run ./cmd/server

logs:
	docker compose logs -f

ps:
	docker compose ps

clean:
	docker compose down -v
	@echo "Volumes removed."
