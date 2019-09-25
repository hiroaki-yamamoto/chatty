#!/bin/sh -e
# -*- coding: utf-8 -*-

backendDir=backend/rpc
frontendDir=frontend/grpc

generateBackend() {
  for fle in grpc/*.proto; do
    mkdir -p ${backendDir}
    protoc --go_out=plugins=grpc:${backendDir} -I grpc ${fle}
  done
}

generateFrontend() {
  mkdir -p ${frontendDir}
  for fle in grpc/*.proto; do
    out=${frontendDir}/$(basename ${fle/.proto/})
    mkdir -p ${out}
    protoc \
      --grpc-web_out=import_style=typescript,mode=grpcwebtext:${out} \
      -I grpc ${fle}
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
    rm -rf ${backendDir}/rpc ${frontendDir}
    ;;
  *)
    generateBackend
    generateFrontend
    ;;
esac
