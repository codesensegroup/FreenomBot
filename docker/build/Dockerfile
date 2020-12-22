
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go mod vender;
RUN go build -o app

FROM alpine:latest
ADD config.toml      /root/app/
COPY resources        /root/app/resources/
COPY --from=build-env /src/app /root/app/
ADD run /root/

WORKDIR /root/

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "/root/run"]