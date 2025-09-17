/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { ApiRequest } from "./api_request";
import { ApiResponse } from "./api_response";

export const protobufPackage = "example.assets";

export interface ApiExchange {
  request: ApiRequest | undefined;
  response: ApiResponse | undefined;
}

function createBaseApiExchange(): ApiExchange {
  return { request: undefined, response: undefined };
}

export const ApiExchange = {
  encode(message: ApiExchange, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.request !== undefined) {
      ApiRequest.encode(message.request, writer.uint32(10).fork()).ldelim();
    }
    if (message.response !== undefined) {
      ApiResponse.encode(message.response, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ApiExchange {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseApiExchange();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.request = ApiRequest.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.response = ApiResponse.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ApiExchange {
    return {
      request: isSet(object.request) ? ApiRequest.fromJSON(object.request) : undefined,
      response: isSet(object.response) ? ApiResponse.fromJSON(object.response) : undefined,
    };
  },

  toJSON(message: ApiExchange): unknown {
    const obj: any = {};
    if (message.request !== undefined) {
      obj.request = ApiRequest.toJSON(message.request);
    }
    if (message.response !== undefined) {
      obj.response = ApiResponse.toJSON(message.response);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ApiExchange>, I>>(base?: I): ApiExchange {
    return ApiExchange.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ApiExchange>, I>>(object: I): ApiExchange {
    const message = createBaseApiExchange();
    message.request = (object.request !== undefined && object.request !== null)
      ? ApiRequest.fromPartial(object.request)
      : undefined;
    message.response = (object.response !== undefined && object.response !== null)
      ? ApiResponse.fromPartial(object.response)
      : undefined;
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
