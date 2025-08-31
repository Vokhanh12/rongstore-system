// This is a generated file - do not edit.
//
// Generated from common/v1/base_response.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use baseResponseDescriptor instead')
const BaseResponse$json = {
  '1': 'BaseResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {
      '1': 'data',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Any',
      '10': 'data'
    },
    {
      '1': 'metadata',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.common.v1.Metadata',
      '10': 'metadata'
    },
    {
      '1': 'error',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.common.v1.Error',
      '10': 'error'
    },
    {
      '1': 'pagination',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.common.v1.Pagination',
      '10': 'pagination'
    },
    {
      '1': 'warnings',
      '3': 6,
      '4': 3,
      '5': 11,
      '6': '.common.v1.Warning',
      '10': 'warnings'
    },
    {
      '1': 'details',
      '3': 7,
      '4': 3,
      '5': 11,
      '6': '.common.v1.BaseResponse.DetailsEntry',
      '10': 'details'
    },
  ],
  '3': [BaseResponse_DetailsEntry$json],
};

@$core.Deprecated('Use baseResponseDescriptor instead')
const BaseResponse_DetailsEntry$json = {
  '1': 'DetailsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `BaseResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List baseResponseDescriptor = $convert.base64Decode(
    'CgxCYXNlUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2VzcxIoCgRkYXRhGAIgASgLMh'
    'QuZ29vZ2xlLnByb3RvYnVmLkFueVIEZGF0YRIvCghtZXRhZGF0YRgDIAEoCzITLmNvbW1vbi52'
    'MS5NZXRhZGF0YVIIbWV0YWRhdGESJgoFZXJyb3IYBCABKAsyEC5jb21tb24udjEuRXJyb3JSBW'
    'Vycm9yEjUKCnBhZ2luYXRpb24YBSABKAsyFS5jb21tb24udjEuUGFnaW5hdGlvblIKcGFnaW5h'
    'dGlvbhIuCgh3YXJuaW5ncxgGIAMoCzISLmNvbW1vbi52MS5XYXJuaW5nUgh3YXJuaW5ncxI+Cg'
    'dkZXRhaWxzGAcgAygLMiQuY29tbW9uLnYxLkJhc2VSZXNwb25zZS5EZXRhaWxzRW50cnlSB2Rl'
    'dGFpbHMaOgoMRGV0YWlsc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUg'
    'V2YWx1ZToCOAE=');
