FROM golang:1.23-alpine AS go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=go-builder /server .
COPY data/ ./data/

EXPOSE 8080
CMD ["./server"]
