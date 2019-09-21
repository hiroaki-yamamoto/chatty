#!/bin/sh
# -*- coding: utf-8 -*-

mkdir -p backend/$(echo grpc/*)
mkdir -p frontend/grpc
for fle in grpc/frontend/**/*.proto; do
  protoc \
    --grpc-web_out=import_style=typescript,mode=grpcwebtext:frontend/grpc \
    --go_out=plugins=grpc:backend/grpc/frontend \
    -I grpc/frontend \
    ${fle}
done

# for fle in grpc/backend/**/services.proto; do
#   protoc \
#     --go_out=plugins=grpc:backend/grpc/backend \
#     -I grpc ${fle}
# done
