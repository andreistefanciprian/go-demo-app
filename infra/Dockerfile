## Build
FROM golang:1.22-alpine AS build

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go ./

# Build with optimizations and proper architecture
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -a -ldflags="-w -s" -o go-demo-app .

## Deploy
FROM gcr.io/distroless/static-debian12:nonroot

# Copy ca-certificates from build stage
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=build /app/go-demo-app /go-demo-app

# Expose port (documentation)
EXPOSE 8080

ENTRYPOINT ["/go-demo-app"]