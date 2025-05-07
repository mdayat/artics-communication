# syntax=docker/dockerfile:1
FROM alpine:3.21 AS base-alpine
WORKDIR /app

FROM golang:1.23.8-alpine3.21 AS base-go
WORKDIR /app

FROM base-go AS build
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY configs/ ./configs/
COPY internal/ ./internal/
COPY repository/ ./repository/

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/web
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed

FROM base-alpine AS final
COPY --from=build /app/app .
COPY --from=build /app/seed .
COPY .env .
ENTRYPOINT ["/app/app"]