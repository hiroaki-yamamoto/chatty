#!/bin/sh -e
# -*- coding: utf-8 -*-

backendDir=backend/rpc
frontendDir=frontend/rpc

generateBackend() {
  mkdir -p ${backendDir}
  # Intermediate protocol between backend and frontend.
  for fle in grpc/*.proto; do
    protoc --go_out=plugins=grpc:${backendDir} -I grpc ${fle}
  done
  # **Internal** protocol between backend and **backend**
  for fle in backend/*/grpc/*.proto; do
    mkdir -p $(dirname $(dirname ${fle}))/rpc/
    protoc \
      --go_out=plugins=grpc:$(dirname $(dirname ${fle}))/rpc/ \
      -I $(dirname ${fle}) ${fle}
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
    rm -r ${backendDir} ${frontendDir} backend/*/rpc
    ;;
  *)
    generateBackend
    generateFrontend
    ;;
esac
