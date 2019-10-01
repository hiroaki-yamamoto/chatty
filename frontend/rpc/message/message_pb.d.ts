import * as jspb from "google-protobuf"

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';

export class Message extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getSendername(): string;
  setSendername(value: string): void;

  getPosttime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setPosttime(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasPosttime(): boolean;
  clearPosttime(): void;

  getProfile(): string;
  setProfile(value: string): void;

  getMessage(): string;
  setMessage(value: string): void;

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
    sendername: string,
    posttime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    profile: string,
    message: string,
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

