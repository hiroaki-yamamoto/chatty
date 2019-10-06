ARG PKGNAME

FROM golang:alpine AS builder
RUN apk --no-cache --update upgrade && apk --no-cache add git gcc libc-dev ca-certificates
ENV GO111MODULE=on
ENV PKGNAME=${PKGNAME}
RUN mkdir -p /opt/code /etc/real
COPY . /opt/code
WORKDIR /opt/code
RUN go generate ./${PKGNAME}
RUN go build -o /tmp/build ./${PKGNAME}

FROM alpine AS executer
ENV PKGNAME=${PKGNAME}
RUN apk --no-cache --update upgrade && apk --no-cache add ca-certificates
COPY --from=builder /tmp/build /usr/bin/app
ENTRYPOINT [ "app" ]
