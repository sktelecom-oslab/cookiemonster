FROM alpine:3.5

MAINTAINER Open System Lab

RUN apk --no-cache update

COPY bin/cookiemonster-linux-amd64 /usr/local/bin/cookiemonster

EXPOSE 8080

CMD ["/usr/local/bin/cookiemonster"]
