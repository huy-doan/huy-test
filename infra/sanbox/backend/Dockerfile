FROM golang:1.24.1-alpine AS builder

ENV PATH /usr/local/go/bin:$PATH
ENV GOLANG_VERSION 1.24.1

RUN apk --no-cache add gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags "-w -s" -o /app/main ./src/main.go
RUN go build -trimpath -ldflags "-w -s" -o /app/command main.go

FROM alpine:3.21 AS runner

RUN apk update && \
    apk add --no-cache tzdata \
        ca-certificates \
        netcat-openbsd && \
    # Set up timezone data properly
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    rm -f /usr/bin/gpasswd /usr/bin/passwd /usr/bin/chfn /sbin/unix_chkpwd /usr/bin/expiry /usr/bin/chage /usr/bin/chsh && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.24.1/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

RUN mkdir -p /app/logs && \
    rm -f /app/logs/* && \
    touch /app/logs/db-backend.log && \
    chmod 755 /app/logs

USER root:root
# Set environment variables for timezone
ENV TZ=Asia/Tokyo

WORKDIR /app

# copy API
COPY --from=builder /app/main /app/main
# copy Batch
COPY --from=builder /app/command /app/command
# copy database
COPY --from=builder /app/database /app/database
# copy scripts
COPY --from=builder /app/ops/scripts /app/scripts
RUN chmod +x /app/scripts/*

# copy translation files
COPY --from=builder /app/src/lib/i18n/translations /app/src/lib/i18n/translations

CMD ["/app/main"]
