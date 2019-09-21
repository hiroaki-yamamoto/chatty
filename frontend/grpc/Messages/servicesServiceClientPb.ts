/**
 * @fileoverview gRPC-Web generated client stub for messages
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as messages_messages_pb from '../messages/messages_pb';

export class MessageServiceClient {
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

  methodInfoMessage = new grpcWeb.AbstractClientBase.MethodInfo(
    messages_messages_pb.Message,
    (request: messages_messages_pb.MessageRequest) => {
      return request.serializeBinary();
    },
    messages_messages_pb.Message.deserializeBinary
  );

  message(
    request: messages_messages_pb.MessageRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/messages.MessageService/Message',
      request,
      metadata || {},
      this.methodInfoMessage);
  }

}

