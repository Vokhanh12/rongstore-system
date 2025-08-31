// This is a generated file - do not edit.
//
// Generated from common/v1/base_request.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class BaseRequest extends $pb.GeneratedMessage {
  factory BaseRequest() => create();

  BaseRequest._();

  factory BaseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BaseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BaseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'common.v1'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BaseRequest clone() => BaseRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BaseRequest copyWith(void Function(BaseRequest) updates) =>
      super.copyWith((message) => updates(message as BaseRequest))
          as BaseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BaseRequest create() => BaseRequest._();
  @$core.override
  BaseRequest createEmptyInstance() => create();
  static $pb.PbList<BaseRequest> createRepeated() => $pb.PbList<BaseRequest>();
  @$core.pragma('dart2js:noInline')
  static BaseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BaseRequest>(create);
  static BaseRequest? _defaultInstance;
}

const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
