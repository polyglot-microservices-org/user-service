# --- Build Stage ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy application source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 go build -o /user-service

# --- Run Stage ---
FROM alpine:latest

WORKDIR /root/

# Copy only the compiled binary
COPY --from=builder /user-service .

EXPOSE 5000

# Single ENTRYPOINT (no CMD)
ENTRYPOINT ["./user-service"]

