# FROM golang:1.25-alpine AS builder

# WORKDIR /app

# COPY go.mod ./
# RUN go mod download

# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o booking-service .

# # -------- runtime image --------
# FROM gcr.io/distroless/base-debian12

# WORKDIR /app
# COPY --from=builder /app/booking-service .

# EXPOSE 8002
# CMD ["./booking-service"]
FROM golang:1.25-alpine AS builder

WORKDIR /app

#installing air for hot reload
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]
