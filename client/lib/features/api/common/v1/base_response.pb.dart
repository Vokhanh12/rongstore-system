// This is a generated file - do not edit.
//
// Generated from common/v1/base_response.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import '../../google/protobuf/any.pb.dart' as $0;
import 'error.pb.dart' as $2;
import 'metadata.pb.dart' as $1;
import 'pagination.pb.dart' as $3;
import 'warning.pb.dart' as $4;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class BaseResponse extends $pb.GeneratedMessage {
  factory BaseResponse({
    $core.bool? success,
    $0.Any? data,
    $1.Metadata? metadata,
    $2.Error? error,
    $3.Pagination? pagination,
    $core.Iterable<$4.Warning>? warnings,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? details,
  }) {
    final result = create();
    if (success != null) result.success = success;
    if (data != null) result.data = data;
    if (metadata != null) result.metadata = metadata;
    if (error != null) result.error = error;
    if (pagination != null) result.pagination = pagination;
    if (warnings != null) result.warnings.addAll(warnings);
    if (details != null) result.details.addEntries(details);
    return result;
  }

  BaseResponse._();

  factory BaseResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BaseResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BaseResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'common.v1'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'success')
    ..aOM<$0.Any>(2, _omitFieldNames ? '' : 'data', subBuilder: $0.Any.create)
    ..aOM<$1.Metadata>(3, _omitFieldNames ? '' : 'metadata',
        subBuilder: $1.Metadata.create)
    ..aOM<$2.Error>(4, _omitFieldNames ? '' : 'error',
        subBuilder: $2.Error.create)
    ..aOM<$3.Pagination>(5, _omitFieldNames ? '' : 'pagination',
        subBuilder: $3.Pagination.create)
    ..pc<$4.Warning>(6, _omitFieldNames ? '' : 'warnings', $pb.PbFieldType.PM,
        subBuilder: $4.Warning.create)
    ..m<$core.String, $core.String>(7, _omitFieldNames ? '' : 'details',
        entryClassName: 'BaseResponse.DetailsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('common.v1'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BaseResponse clone() => BaseResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BaseResponse copyWith(void Function(BaseResponse) updates) =>
      super.copyWith((message) => updates(message as BaseResponse))
          as BaseResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BaseResponse create() => BaseResponse._();
  @$core.override
  BaseResponse createEmptyInstance() => create();
  static $pb.PbList<BaseResponse> createRepeated() =>
      $pb.PbList<BaseResponse>();
  @$core.pragma('dart2js:noInline')
  static BaseResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BaseResponse>(create);
  static BaseResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get success => $_getBF(0);
  @$pb.TagNumber(1)
  set success($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSuccess() => $_has(0);
  @$pb.TagNumber(1)
  void clearSuccess() => $_clearField(1);

  @$pb.TagNumber(2)
  $0.Any get data => $_getN(1);
  @$pb.TagNumber(2)
  set data($0.Any value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(2)
  void clearData() => $_clearField(2);
  @$pb.TagNumber(2)
  $0.Any ensureData() => $_ensure(1);

  @$pb.TagNumber(3)
  $1.Metadata get metadata => $_getN(2);
  @$pb.TagNumber(3)
  set metadata($1.Metadata value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasMetadata() => $_has(2);
  @$pb.TagNumber(3)
  void clearMetadata() => $_clearField(3);
  @$pb.TagNumber(3)
  $1.Metadata ensureMetadata() => $_ensure(2);

  @$pb.TagNumber(4)
  $2.Error get error => $_getN(3);
  @$pb.TagNumber(4)
  set error($2.Error value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasError() => $_has(3);
  @$pb.TagNumber(4)
  void clearError() => $_clearField(4);
  @$pb.TagNumber(4)
  $2.Error ensureError() => $_ensure(3);

  @$pb.TagNumber(5)
  $3.Pagination get pagination => $_getN(4);
  @$pb.TagNumber(5)
  set pagination($3.Pagination value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasPagination() => $_has(4);
  @$pb.TagNumber(5)
  void clearPagination() => $_clearField(5);
  @$pb.TagNumber(5)
  $3.Pagination ensurePagination() => $_ensure(4);

  @$pb.TagNumber(6)
  $pb.PbList<$4.Warning> get warnings => $_getList(5);

  @$pb.TagNumber(7)
  $pb.PbMap<$core.String, $core.String> get details => $_getMap(6);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
