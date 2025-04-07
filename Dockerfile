FROM golang:1.24-alpine AS builder

ADD . /app
WORKDIR /app

ARG VERSION
# ARG BRANCH
# ARG COMMIT

RUN apk update
# RUN apk add --no-cache git gcc musl-dev
RUN go build -o /tmp/app ./cmd/app

FROM alpine:latest

# RUN apk add --no-cache ca-certificates tzdata
# RUN update-ca-certificates

ARG SERVER_CERT_PATH
ARG SERVER_KEY_PATH

COPY --from=builder /tmp/app /usr/bin/app
COPY --from=builder /app/db/migrations /db/migrations
COPY --from=builder /app/docs /docs

RUN chmod +x /usr/bin/app

RUN adduser -D -u 1000 app

USER app

ENTRYPOINT [ "app" ]
