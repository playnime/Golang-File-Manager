# контейнер
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./server/cmd/server/main.go

# образ
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/server ./server

RUN mkdir -p /app/out

EXPOSE 9091

CMD ["./server"]