#!/bin/sh -e
# -*- coding: utf-8 -*-

go mod download
if [ -n "${PKGNAME}" ]; then
  exec go test "./${PKGNAME}/..."
else
  exec go test ${1}
fi
