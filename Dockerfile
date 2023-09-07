# Build stage
FROM golang:1.21-alpine AS build

RUN apk update && \
    apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main

# Final image
FROM alpine

RUN apk update && \
    apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 3000

CMD ["./main"]