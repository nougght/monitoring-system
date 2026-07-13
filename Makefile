
.PHONY: lint lint-agent lint-server lint-shared run-agent run-server

# include .env


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