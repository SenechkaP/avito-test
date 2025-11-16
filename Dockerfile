FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/pr_service ./cmd

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache make bash ca-certificates curl

COPY --from=builder /app/pr_service /app/pr_service
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/Makefile ./Makefile
COPY --from=builder /app/.env .env

COPY docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ARG TARGETARCH
ENV GOOSE_VERSION v3.14.0
RUN echo "Downloading goose for arch=${TARGETARCH} (version ${GOOSE_VERSION})" \
 && curl -fL "https://github.com/pressly/goose/releases/download/${GOOSE_VERSION}/goose_linux_${TARGETARCH}" -o /usr/local/bin/goose \
 && chmod +x /usr/local/bin/goose

RUN /usr/local/bin/goose -version

ENTRYPOINT ["/app/entrypoint.sh"]
