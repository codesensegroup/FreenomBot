
FROM golang:alpine AS build-env
RUN mkdir /app
WORKDIR /app/
RUN apk add --update --no-cache ca-certificates git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o app

FROM alpine:latest
#ADD config.toml /root/app/
COPY --from=build-env /app/app /root/app/
ADD run /root/

WORKDIR /root/

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "/root/run"]