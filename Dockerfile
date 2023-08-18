# Build stage
FROM golang:alpine AS build

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

RUN go build -o main

# Final image
FROM alpine

RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 3000

CMD ["./main"]
