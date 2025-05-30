# Build stage for Go backend
FROM golang:1.21-alpine AS go-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Build stage for SvelteKit frontend
FROM node:18-alpine AS node-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./
RUN npm ci

# Copy frontend source
COPY frontend/ ./

# Build the frontend
RUN npm run build

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the Go binary
COPY --from=go-builder /app/main .

# Copy static files
COPY --from=node-builder /app/frontend/build ./static/
COPY static/ ./static/

# Create directory for SQLite database
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set environment variables
ENV DATABASE_URL=sqlite:///data/payments.db
ENV PORT=8080

CMD ["./main"]