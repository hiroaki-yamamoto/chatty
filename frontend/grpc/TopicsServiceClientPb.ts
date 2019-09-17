/**
 * @fileoverview gRPC-Web generated client stub for topics
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import {
  Message,
  MessageRequest,
  TopicInfo,
  TopicRequest} from './topics_pb';

export class TopicServiceClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: string; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodInfoTopic = new grpcWeb.AbstractClientBase.MethodInfo(
    TopicInfo,
    (request: TopicRequest) => {
      return request.serializeBinary();
    },
    TopicInfo.deserializeBinary
  );

  topic(
    request: TopicRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/topics.TopicService/Topic',
      request,
      metadata || {},
      this.methodInfoTopic);
  }

  methodInfoMessage = new grpcWeb.AbstractClientBase.MethodInfo(
    Message,
    (request: MessageRequest) => {
      return request.serializeBinary();
    },
    Message.deserializeBinary
  );

  message(
    request: MessageRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/topics.TopicService/Message',
      request,
      metadata || {},
      this.methodInfoMessage);
  }

}

