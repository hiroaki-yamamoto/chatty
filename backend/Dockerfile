FROM golang:alpine AS builder
ARG PKGNAME
RUN apk --no-cache --update upgrade && apk --no-cache add git gcc libc-dev ca-certificates
ENV GO111MODULE=on
ENV PKGNAME=${PKGNAME}
RUN mkdir -p /opt/code
ADD . /opt/code
WORKDIR /opt/code
RUN go generate ./${PKGNAME}
RUN go build -o /tmp/build ./${PKGNAME}

FROM alpine AS executer
RUN mkdir -p /etc/real
VOLUME [ "/etc/real" ]
RUN apk --no-cache --update upgrade && apk --no-cache add ca-certificates
COPY --from=builder /tmp/build /usr/bin/app
ENTRYPOINT [ "app" ]
