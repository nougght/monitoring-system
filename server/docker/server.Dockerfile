FROM golang:1.26-alpine AS build

WORKDIR /build

# Context is repo root so the go.mod replace ../shared/go resolves.
COPY shared/go ./shared/go
COPY server ./server

WORKDIR /build/server
RUN go build -o /build/app ./cmd/main.go


# light weight image for running
FROM alpine
WORKDIR /app

COPY --from=build /build/app ./app

ENTRYPOINT ["./app"]
