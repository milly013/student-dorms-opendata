# Build stage
FROM golang:1.25-alpine AS build
WORKDIR /app

# Kopiraj Go mod fajlove i preuzmi dependencies
COPY go.mod go.sum ./
RUN go mod download

# Kopiraj ceo kod
COPY . .

# Build aplikaciju
RUN CGO_ENABLED=0 go build -o request-service main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
RUN adduser -D app
USER app

# Kopiraj build binary iz build stage
COPY --from=build /app/request-service .

EXPOSE 8083

ENTRYPOINT ["./request-service"]
