#!/bin/sh
# -*- coding: utf-8 -*-

protoc --grpc-web_out=import_style=typescript,mode=grpcwebtext:frontend/grpc -I grpc grpc/topics.proto
