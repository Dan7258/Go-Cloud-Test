FROM golang:1.23.8 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloud .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/cloud .
RUN chmod +x ./cloud
CMD ["./cloud"]