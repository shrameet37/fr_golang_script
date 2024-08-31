# FROM golang:1.21

FROM golang:alpine AS build-env
WORKDIR /app
RUN apk update && apk add gcc libc-dev librdkafka-dev pkgconf
ADD . /app
RUN cd /app && go build -tags musl -o goserver

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /app/goserver /app
COPY --from=build-env /app/us-east-1-bundle.pem /app
COPY --from=build-env /app/ap-south-1-bundle.pem /app

EXPOSE 8001

RUN adduser -s /bin/bash --disabled-password -u 1000 go-user

RUN chown -R go-user:go-user /app

USER 1000

ENTRYPOINT ./goserver