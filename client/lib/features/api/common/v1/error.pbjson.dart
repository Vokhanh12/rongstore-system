// This is a generated file - do not edit.
//
// Generated from common/v1/error.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use errorDetailDescriptor instead')
const ErrorDetail$json = {
  '1': 'ErrorDetail',
  '2': [
    {'1': 'field', '3': 1, '4': 1, '5': 9, '10': 'field'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
    {'1': 'code', '3': 3, '4': 1, '5': 9, '10': 'code'},
  ],
};

/// Descriptor for `ErrorDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List errorDetailDescriptor = $convert.base64Decode(
    'CgtFcnJvckRldGFpbBIUCgVmaWVsZBgBIAEoCVIFZmllbGQSGAoHbWVzc2FnZRgCIAEoCVIHbW'
    'Vzc2FnZRISCgRjb2RlGAMgASgJUgRjb2Rl');

@$core.Deprecated('Use errorDescriptor instead')
const Error$json = {
  '1': 'Error',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 9, '10': 'code'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'details',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.common.v1.ErrorDetail',
      '10': 'details'
    },
  ],
};

/// Descriptor for `Error`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List errorDescriptor = $convert.base64Decode(
    'CgVFcnJvchISCgRjb2RlGAEgASgJUgRjb2RlEhgKB21lc3NhZ2UYAiABKAlSB21lc3NhZ2USMA'
    'oHZGV0YWlscxgDIAMoCzIWLmNvbW1vbi52MS5FcnJvckRldGFpbFIHZGV0YWlscw==');
