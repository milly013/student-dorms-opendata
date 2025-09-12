# Build stage
FROM golang:1.25-alpine AS build
WORKDIR /app

# Instaliraj git i certifikate
RUN apk add --no-cache git ca-certificates

# Kopiraj Go mod fajlove i preuzmi dependencies
COPY go.mod go.sum ./
RUN go mod download

# Kopiraj ceo kod
COPY . .

# Build aplikaciju
RUN CGO_ENABLED=0 go build -o api-gateway main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
RUN adduser -D app
USER app

# Kopiraj build binary iz build stage
COPY --from=build /app/api-gateway .

EXPOSE 8000

ENTRYPOINT ["./api-gateway"]
