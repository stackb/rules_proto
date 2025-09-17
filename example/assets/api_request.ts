/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "example.assets";

export interface ApiRequest {
  header: string[];
}

function createBaseApiRequest(): ApiRequest {
  return { header: [] };
}

export const ApiRequest = {
  encode(message: ApiRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.header) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApiRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApiRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.header.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ApiRequest {
    return { header: Array.isArray(object?.header) ? object.header.map((e: any) => String(e)) : [] };
  },

  toJSON(message: ApiRequest): unknown {
    const obj: any = {};
    if (message.header?.length) {
      obj.header = message.header;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ApiRequest>, I>>(base?: I): ApiRequest {
    return ApiRequest.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ApiRequest>, I>>(object: I): ApiRequest {
    const message = createBaseApiRequest();
    message.header = object.header?.map((e) => e) || [];
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };
