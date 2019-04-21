FROM golang:1.12.4-alpine3.9 AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN apk add --no-cache git
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kelsier

FROM alpine:3.9
LABEL maintainer="danielkvist@protonmail.com"
RUN apk add --no-cache ca-certificates 
COPY --from=build /app/kelsier /app/kelsier
ENTRYPOINT ["/app/kelsier"]