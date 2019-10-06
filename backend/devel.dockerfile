FROM golang:alpine
ARG PKGNAME
RUN apk --no-cache --update upgrade && apk --no-cache add git gcc libc-dev ca-certificates
ENV GO111MODULE=on
ENV PKGNAME=${PKGNAME}
RUN mkdir -p /opt/code /etc/real
VOLUME [ "/opt/code", "/etc/real" ]
WORKDIR /opt/code
ENTRYPOINT [ "./run.sh" ]
