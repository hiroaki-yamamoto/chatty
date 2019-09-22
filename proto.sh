#!/bin/sh -ex
# -*- coding: utf-8 -*-

backendDir=backend/generated
frontendDir=frontend/generated

generateBackend() {
  mkdir -p ${backendDir}
  for fle in grpc/*; do
    protoc \
      --go_out=plugins=grpc:${backendDir} \
      -I grpc ${fle}/*.proto
  done
}

generateFrontend() {
  mkdir -p ${frontendDir}
  for fle in grpc/*; do
    protoc \
      --grpc-web_out=import_style=typescript,mode=grpcwebtext:${frontendDir} \
      -I grpc ${fle}/*.proto
  done
}

case $1 in
  "backend")
    generateBackend
    ;;
  "frontend")
    generateFrontend
    ;;
  "clean")
    rm -rf ${backendDir} ${frontendDir}
    ;;
  *)
    generateBackend
    generateFrontend
    ;;
esac
