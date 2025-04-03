FROM golang:1.24-alpine3.21 AS dev

ENV CGO_ENABLED 0
ENV GO111MODULE auto
ENV GOLANG_VERSION 1.24.1
# don't auto-upgrade the gotoolchain
# https://github.com/docker-library/golang/issues/472
ENV GOTOOLCHAIN=local

RUN apk --no-cache add git

WORKDIR /app
COPY . /app

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.7.0 && \
    go install github.com/cosmtrek/air@v1.27.3 && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.3

CMD ["sh", "-c", "air -c /app/.air.toml"]
