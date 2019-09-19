#!/bin/sh
# -*- coding: utf-8 -*-

protoc \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:frontend/grpc \
  --go_out=plugins=grpc:backend/grpc_front \
  -I grpc grpc/topics.proto
