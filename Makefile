
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

build-agent:
	cd agent && GOOS=linux GOARCH=amd64 go build -o ../bin/agent ./cmd
	cd agent && GOOS=windows GOARCH=amd64 go build -o ../bin/agent.exe ./cmd

run-server:
	cd server && go run cmd/main.go

timescale-up:
	docker compose --env-file ./server/.env -f ./server/docker/docker-compose.yaml -p monitoring-system up -d --build timescale

server-build:
	docker compose -f ./server/docker/docker-compose.yaml -p monitoring-system --build monitoring-backend

server-up:
	docker compose --env-file ./server/.env -f ./server/docker/docker-compose.yaml -p monitoring-system up -d  --build monitoring-backend

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


gen-swag-agent:
	cd agent && swag init -g cmd/main.go -o api --parseDependency --parseInternal 

gen-swag-server:
	cd server && swag init -g cmd/main.go -o api --parseDependency --parseInternal 



# temp
gen-ca:
# root — самоподписанный
	openssl ecparam -genkey -name prime256v1 -noout -out ./rootCA/root-ca.key
	openssl req -x509 -new -key ./rootCA/root-ca.key -days 3650 \
		-subj "/CN=Monitoring Root CA" -out ./rootCA/root-ca.crt

# intermediate — CSR подписывается root
	openssl ecparam -genkey -name prime256v1 -noout -out ./server/creds/intermediate-ca.key
	openssl req -new -key ./server/creds/intermediate-ca.key -subj "/CN=Monitoring Intermediate CA" -out ./server/creds/intermediate-ca.csr
	openssl x509 -req -in ./server/creds/intermediate-ca.csr -CA ./rootCA/root-ca.crt -CAkey ./rootCA/root-ca.key \
		-CAcreateserial -days 1825 -extfile rootCA/int.cnf \
		-out ./server/creds/intermediate-ca.crt