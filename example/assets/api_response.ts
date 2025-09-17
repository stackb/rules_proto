/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "example.assets";

export interface ApiResponse {
  code: number;
  status: string;
}

function createBaseApiResponse(): ApiResponse {
  return { code: 0, status: "" };
}

export const ApiResponse = {
  encode(message: ApiResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.code !== 0) {
      writer.uint32(8).int32(message.code);
    }
    if (message.status !== "") {
      writer.uint32(18).string(message.status);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApiResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApiResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.code = reader.int32();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.status = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ApiResponse {
    return {
      code: isSet(object.code) ? Number(object.code) : 0,
      status: isSet(object.status) ? String(object.status) : "",
    };
  },

  toJSON(message: ApiResponse): unknown {
    const obj: any = {};
    if (message.code !== 0) {
      obj.code = Math.round(message.code);
    }
    if (message.status !== "") {
      obj.status = message.status;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ApiResponse>, I>>(base?: I): ApiResponse {
    return ApiResponse.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ApiResponse>, I>>(object: I): ApiResponse {
    const message = createBaseApiResponse();
    message.code = object.code ?? 0;
    message.status = object.status ?? "";
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

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
