
.PHONY: lint lint-agent lint-server lint-shared run-agent run-server

include ./server/.env


lint: lint-agent lint-server lint-shared

lint-agent:
	docker run -t --rm -v .:/app -w /app/agent golangci/golangci-lint:v2.12.2 \
		golangci-lint run --no-config -E govet,staticcheck

lint-server:
	docker run -t --rm -v .:/app -w /app/server golangci/golangci-lint:v2.12.2 \
		golangci-lint run --no-config -E govet,staticcheck

lint-shared:
	docker run -t --rm -v .:/app -w /app/shared/go golangci/golangci-lint:v2.12.2 \
		golangci-lint run --no-config -E govet,staticcheck
	
run-agent:
	cd agent && go run cmd/main.go


run-server:
	cd server && go run cmd/main.go


timescale-up:
	docker compose -f ./server/docker/docker-compose.yaml -p monitoring-system up -d --build timescale

server-build:
	docker compose -f ./server/docker/docker-compose.yaml -p monitoring-system build monitoring-backend

server-up:
	docker compose -f ./server/docker/docker-compose.yaml -p monitoring-system up -d monitoring-backend

migrate-create: 
ifndef MIGRATE_NAME
	echo $(error migrate name is required, use `make migrate-create MIGRATE_NAME=your_migration_name`)
else
	migrate create -ext sql -dir server/internal/storage/timescale/migrations -seq ${MIGRATE_NAME}
endif

migrate-force:
ifndef MIGRATE_VERSION
	echo $(error migrate version is required, use `make migrate-force MIGRATE_VERSION=your_migration_version`)
else
	migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable -path server/internal/storage/timescale/migrations force ${MIGRATE_VERSION}
endif