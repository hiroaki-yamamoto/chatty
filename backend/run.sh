#!/bin/sh -e
# -*- coding: utf-8 -*-

go mod download
go build -o /usr/bin/app ./${PKGNAME}
exec app
