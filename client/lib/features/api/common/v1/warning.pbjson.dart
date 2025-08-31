// This is a generated file - do not edit.
//
// Generated from common/v1/warning.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use warningDescriptor instead')
const Warning$json = {
  '1': 'Warning',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 9, '10': 'code'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'metadata',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.common.v1.Warning.MetadataEntry',
      '10': 'metadata'
    },
  ],
  '3': [Warning_MetadataEntry$json],
};

@$core.Deprecated('Use warningDescriptor instead')
const Warning_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Warning`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List warningDescriptor = $convert.base64Decode(
    'CgdXYXJuaW5nEhIKBGNvZGUYASABKAlSBGNvZGUSGAoHbWVzc2FnZRgCIAEoCVIHbWVzc2FnZR'
    'I8CghtZXRhZGF0YRgDIAMoCzIgLmNvbW1vbi52MS5XYXJuaW5nLk1ldGFkYXRhRW50cnlSCG1l'
    'dGFkYXRhGjsKDU1ldGFkYXRhRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKA'
    'lSBXZhbHVlOgI4AQ==');
