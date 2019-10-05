#!/bin/sh -e
# -*- coding: utf-8 -*-

run() {
  exec real_${PKGNAME}
}

depBuild() {
  go generate
  go build -o /usr/bin/real_${PKGNAME} ${PKGNAME}
  run
}

case ${MODE} in
devel)
  depBuild
;;
prod)
  run
;;
*)
  cat << EOD
Usage: ${0}
=========
Env Var:
  MODE=(devel|prod): Development Mode or Production Mode
  PKGNAME=[Package Name]: The name of the package
EOD
;;
esac
