// This is a generated file - do not edit.
//
// Generated from test/v1/test.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use testRequestDescriptor instead')
const TestRequest$json = {
  '1': 'TestRequest',
  '2': [
    {'1': 'test', '3': 1, '4': 1, '5': 9, '10': 'test'},
  ],
};

/// Descriptor for `TestRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testRequestDescriptor =
    $convert.base64Decode('CgtUZXN0UmVxdWVzdBISCgR0ZXN0GAEgASgJUgR0ZXN0');

@$core.Deprecated('Use testResponseDescriptor instead')
const TestResponse$json = {
  '1': 'TestResponse',
};

/// Descriptor for `TestResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testResponseDescriptor =
    $convert.base64Decode('CgxUZXN0UmVzcG9uc2U=');
