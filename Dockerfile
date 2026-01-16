FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go-as-your-backend .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

ENV ENV=dev

COPY --from=builder /app/go-as-your-backend .

EXPOSE 8080

CMD ["./go-as-your-backend"]

