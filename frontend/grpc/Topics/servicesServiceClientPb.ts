/**
 * @fileoverview gRPC-Web generated client stub for topics
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as topics_messages_pb from '../topics/messages_pb';

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
    topics_messages_pb.TopicInfo,
    (request: topics_messages_pb.TopicRequest) => {
      return request.serializeBinary();
    },
    topics_messages_pb.TopicInfo.deserializeBinary
  );

  topic(
    request: topics_messages_pb.TopicRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/topics.TopicService/Topic',
      request,
      metadata || {},
      this.methodInfoTopic);
  }

}

