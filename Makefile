
.PHONY: run-agent, 

# include .env


run-agent:
	cd agent && go run cmd/main.go


run-server:
	cd server && go run cmd/main.go