# ---------- Build Stage ----------
# FROM golang:1.22-alpine AS builder

# WORKDIR /app

# # Install git (required for go mod in alpine sometimes)
# RUN apk add --no-cache git

# # Copy go mod files first (better layer caching)
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy rest of source
# COPY . .

# # Build static binary
# RUN CGO_ENABLED=0 go build -o app ./server

# # ---------- Runtime Stage ----------
# FROM gcr.io/distroless/base-debian12

# WORKDIR /app

# # Copy binary from builder
# COPY --from=builder /app/app .

# # Expose gRPC port
# EXPOSE 50051

# # Run as non-root user
# USER nonroot:nonroot

# ENTRYPOINT ["/app/app"]

FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app ./server

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 50051
USER nonroot:nonroot
CMD ["/app/app"]