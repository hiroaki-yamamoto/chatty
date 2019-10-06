#!/bin/sh -e
# -*- coding: utf-8 -*-

go generate ./${PKGNAME}
go build -o /usr/bin/app ./${PKGNAME}
exec app
