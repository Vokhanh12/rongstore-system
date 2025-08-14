// This is a generated file - do not edit.
//
// Generated from test/v1/test.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'test.pb.dart' as $0;

export 'test.pb.dart';

@$pb.GrpcServiceName('test.TestService')
class TestServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  TestServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.TestResponse> test(
    $0.TestRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$test, request, options: options);
  }

  // method descriptors

  static final _$test = $grpc.ClientMethod<$0.TestRequest, $0.TestResponse>(
      '/test.TestService/Test',
      ($0.TestRequest value) => value.writeToBuffer(),
      $0.TestResponse.fromBuffer);
}

@$pb.GrpcServiceName('test.TestService')
abstract class TestServiceBase extends $grpc.Service {
  $core.String get $name => 'test.TestService';

  TestServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.TestRequest, $0.TestResponse>(
        'Test',
        test_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TestRequest.fromBuffer(value),
        ($0.TestResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.TestResponse> test_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.TestRequest> $request) async {
    return test($call, await $request);
  }

  $async.Future<$0.TestResponse> test(
      $grpc.ServiceCall call, $0.TestRequest request);
}
