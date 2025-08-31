// This is a generated file - do not edit.
//
// Generated from common/v1/metadata.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use metadataDescriptor instead')
const Metadata$json = {
  '1': 'Metadata',
  '2': [
    {'1': 'timestamp', '3': 1, '4': 1, '5': 9, '10': 'timestamp'},
    {'1': 'request_id', '3': 2, '4': 1, '5': 9, '10': 'requestId'},
  ],
};

/// Descriptor for `Metadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataDescriptor = $convert.base64Decode(
    'CghNZXRhZGF0YRIcCgl0aW1lc3RhbXAYASABKAlSCXRpbWVzdGFtcBIdCgpyZXF1ZXN0X2lkGA'
    'IgASgJUglyZXF1ZXN0SWQ=');
