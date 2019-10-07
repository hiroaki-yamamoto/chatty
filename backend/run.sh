#!/bin/sh -e
# -*- coding: utf-8 -*-

go mod download
go generate ./${PKGNAME}
go build -o /usr/bin/app ./${PKGNAME}
exec app
