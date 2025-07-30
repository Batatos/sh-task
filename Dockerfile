FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod ./

RUN go mod download && go mod tidy

COPY . .

RUN go mod download && go mod tidy && CGO_ENABLED=0 go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/main .



RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./main"] 