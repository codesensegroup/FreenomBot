FROM alpine:latest

ADD config.toml      /root/FreenomBot/
COPY resources        /root/FreenomBot/resources/
ADD FreenomBot_amd64 /root/FreenomBot/
ADD run /root/

WORKDIR /root/

EXPOSE 8080

ENTRYPOINT ["/bin/sh", "/root/run"]