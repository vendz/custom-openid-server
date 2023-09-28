# # Build stage
# FROM golang:1.21-alpine AS build

# RUN apk update && \
#     apk add --no-cache git

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# RUN go build -o main

# # Final image
# FROM alpine

# RUN apk update && \
#     apk add --no-cache ca-certificates

# WORKDIR /app

# COPY --from=build /app/main .

# EXPOSE 3000

# CMD ["./main"]

FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN mkdir /app

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

EXPOSE 3000

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main