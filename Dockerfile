# Build stage
FROM golang:1.21-alpine AS builder

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Set Go environment to use proxy and avoid VCS
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org
ENV GONOPROXY=""
ENV GONOSUMDB=""
ENV GOPRIVATE=""

# Copy go mod and sum files (if they exist)
COPY go.mod ./
COPY go.su[m] ./

# Download dependencies and ensure go.sum is generated
RUN go mod download
RUN go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/config.env .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
