FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o ./bin/app main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /app/bin/app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./app"]