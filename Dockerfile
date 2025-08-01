
    FROM golang:1.24.5-alpine AS builder
    WORKDIR /app
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux go build -o fechar_chats ./cmd
    
    FROM alpine:latest
    WORKDIR /app
    COPY --from=builder /app/fechar_chats ./fechar_chats
    COPY .env ./
    RUN apk add --no-cache ca-certificates
        
    ENTRYPOINT ["./fechar_chats"] 