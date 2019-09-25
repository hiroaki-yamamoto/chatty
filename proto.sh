#!/bin/sh -e
# -*- coding: utf-8 -*-

backendDir=backend
frontendDir=frontend/grpc

generateBackend() {
  for fle in grpc/*.proto; do
    out=${backendDir}/$(basename ${fle/.proto/})/server
    mkdir -p ${out}
    protoc --go_out=plugins=grpc:${out} -I grpc ${fle}
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
    rm -rf ${backendDir}/**/server/*.pb.go ${frontendDir}
    ;;
  *)
    generateBackend
    generateFrontend
    ;;
esac
