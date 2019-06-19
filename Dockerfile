# Build stage
FROM golang:1.12.5-alpine3.9 AS build
RUN apk add --no-cache git
WORKDIR /app/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o kelsier main.go

# Final stage
FROM alpine:3.9.4
LABEL maintainer="danielkvist@protonmail.com"
RUN apk add --no-cache ca-certificates
RUN adduser -D -g '' daniel && \
    mkdir /app/ && chown daniel /app/
USER daniel
COPY --from=build /app/kelsier /app/kelsier
ENTRYPOINT ["./app/kelsier"]
