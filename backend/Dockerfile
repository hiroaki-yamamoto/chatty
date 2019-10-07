FROM golang:alpine AS build
ARG PKGNAME
RUN apk --no-cache --update upgrade && apk --no-cache add git gcc libc-dev ca-certificates
ENV GO111MODULE=on
RUN mkdir -p /opt/code
COPY . /opt/code
WORKDIR /opt/code
RUN go mod download
RUN go generate ./${PKGNAME}
RUN go build -o /tmp/build ./${PKGNAME}

FROM alpine AS exec
RUN mkdir -p /etc/real
VOLUME [ "/etc/real" ]
RUN apk --no-cache --update upgrade && apk --no-cache add ca-certificates
COPY --from=build /tmp/build /usr/bin/app
ENTRYPOINT [ "app" ]
