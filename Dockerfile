FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/pr_service ./cmd

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

FROM alpine:3.18 AS stage-1

WORKDIR /app

RUN apk add --no-cache make bash ca-certificates curl

COPY --from=builder /app/pr_service /app/pr_service
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/Makefile ./Makefile
COPY --from=builder /app/.env .env

COPY docker/entrypoint.sh /app/entrypoint.sh

RUN sed -i 's/\r$//' /app/entrypoint.sh && chmod +x /app/entrypoint.sh

COPY --from=builder /go/bin/goose /usr/local/bin/goose

ENTRYPOINT ["/app/entrypoint.sh"]