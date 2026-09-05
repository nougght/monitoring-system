FROM golang:1.26-alpine AS build

WORKDIR /build

# Context is repo root so the go.mod replace ../shared/go resolves.
COPY shared/go ./shared/go
COPY server ./server

WORKDIR /build/server
RUN go build -o /build/monit-server ./cmd/main.go


# light weight image for running
FROM alpine
WORKDIR /app/server

COPY --from=build /build/monit-server .
COPY --from=build /build/server/config.yaml .
COPY bin/agent.exe /app/bin/agent.exe
COPY server/creds/root-ca.crt /app/server/creds/root-ca.crt
ENTRYPOINT ["./monit-server"]
