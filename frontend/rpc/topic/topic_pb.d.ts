import * as jspb from "google-protobuf"

export class TopicInfo extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getTitle(): string;
  setTitle(value: string): void;

  getNumreplies(): number;
  setNumreplies(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TopicInfo.AsObject;
  static toObject(includeInstance: boolean, msg: TopicInfo): TopicInfo.AsObject;
  static serializeBinaryToWriter(message: TopicInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TopicInfo;
  static deserializeBinaryFromReader(message: TopicInfo, reader: jspb.BinaryReader): TopicInfo;
}

export namespace TopicInfo {
  export type AsObject = {
    id: string,
    title: string,
    numreplies: number,
  }
}

export class TopicRequest extends jspb.Message {
  getBoardid(): string;
  setBoardid(value: string): void;

  getStartfrom(): number;
  setStartfrom(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TopicRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TopicRequest): TopicRequest.AsObject;
  static serializeBinaryToWriter(message: TopicRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TopicRequest;
  static deserializeBinaryFromReader(message: TopicRequest, reader: jspb.BinaryReader): TopicRequest;
}

export namespace TopicRequest {
  export type AsObject = {
    boardid: string,
    startfrom: number,
  }
}

