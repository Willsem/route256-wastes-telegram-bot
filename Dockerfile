FROM golang:1.19.2-alpine3.16 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build ./cmd/bot

FROM alpine:3.16

WORKDIR /app

COPY --from=build /app/bot ./

CMD ["/app/bot", "-config", "/app/config.yaml"]