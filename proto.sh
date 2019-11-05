#!/bin/sh -e
# -*- coding: utf-8 -*-

backend() {
  ninja -C backend ${@}
}

frontend() {
  ninja -C frontend ${@}
}

case $1 in
  "backend")
    backend
    ;;
  "frontend")
    frontend
    ;;
  "clean")
    backend -t clean
    frontend -t clean
    ;;
  *)
    backend
    frontend
    ;;
esac
