import * as jspb from "google-protobuf"

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as status_pb from './status_pb';

export class Message extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getTopicid(): string;
  setTopicid(value: string): void;

  getSendername(): string;
  setSendername(value: string): void;

  getPosttime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setPosttime(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasPosttime(): boolean;
  clearPosttime(): void;

  getMessage(): string;
  setMessage(value: string): void;

  getBump(): boolean;
  setBump(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Message.AsObject;
  static toObject(includeInstance: boolean, msg: Message): Message.AsObject;
  static serializeBinaryToWriter(message: Message, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Message;
  static deserializeBinaryFromReader(message: Message, reader: jspb.BinaryReader): Message;
}

export namespace Message {
  export type AsObject = {
    id: string,
    topicid: string,
    sendername: string,
    posttime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    message: string,
    bump: boolean,
  }
}

export class MessageRequest extends jspb.Message {
  getTopicid(): string;
  setTopicid(value: string): void;

  getStartfrom(): number;
  setStartfrom(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MessageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MessageRequest): MessageRequest.AsObject;
  static serializeBinaryToWriter(message: MessageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MessageRequest;
  static deserializeBinaryFromReader(message: MessageRequest, reader: jspb.BinaryReader): MessageRequest;
}

export namespace MessageRequest {
  export type AsObject = {
    topicid: string,
    startfrom: number,
  }
}

export class PostRequest extends jspb.Message {
  getTopicid(): string;
  setTopicid(value: string): void;

  getName(): string;
  setName(value: string): void;

  getBump(): boolean;
  setBump(value: boolean): void;

  getMessage(): string;
  setMessage(value: string): void;

  getRecaptcha(): string;
  setRecaptcha(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PostRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PostRequest): PostRequest.AsObject;
  static serializeBinaryToWriter(message: PostRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PostRequest;
  static deserializeBinaryFromReader(message: PostRequest, reader: jspb.BinaryReader): PostRequest;
}

export namespace PostRequest {
  export type AsObject = {
    topicid: string,
    name: string,
    bump: boolean,
    message: string,
    recaptcha: string,
  }
}

