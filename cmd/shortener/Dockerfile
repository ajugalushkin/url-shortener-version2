# syntax=docker/dockerfile:1
FROM golang:1.22.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /shortener-api ./cmd/shortener/main.go

FROM alpine AS server
WORKDIR /
COPY --from=builder /shortener-api /shortener-api
EXPOSE 8080
#USER nonroot:nonroot
ENTRYPOINT ["/shortener-api"]