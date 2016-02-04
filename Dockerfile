FROM        sdurrheimer/alpine-glibc
MAINTAINER  Thibault NORMAND <me@zenithar.org>

WORKDIR /gopath/src/github.com/Zenithar/goproject
COPY    . /gopath/src/github.com/Zenithar/goproject

RUN apk add --update -t build-deps tar openssl git make bash \
    && source ./scripts/goenv.sh /go /gopath \
    && make build \
    && cp goproject_server /bin/ \
    && mkdir /app \
    && apk del --purge build-deps \
    && rm -rf /go /gopath /var/cache/apk/*

EXPOSE     3000 5555
WORKDIR    /app
VOLUME     ["/app"]
ENTRYPOINT [ "/bin/goproject_server" ]
CMD        [ "" ]
