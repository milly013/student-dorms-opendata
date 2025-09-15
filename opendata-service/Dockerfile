# Build stage
FROM golang:1.25-alpine AS build

WORKDIR /app

# Kopiraj go.mod i go.sum
COPY go.mod go.sum ./

# Preuzmi dependencies
RUN go mod download

# Kopiraj ostatak koda
COPY . .

# Build binary
RUN CGO_ENABLED=0 go build -o opendata-service main.go

# Final stage
FROM alpine:3.20

WORKDIR /app

# Kopiraj binary iz build stage-a
COPY --from=build /app/opendata-service .

# Expose port
EXPOSE 8082

# Start service
CMD ["./opendata-service"]
