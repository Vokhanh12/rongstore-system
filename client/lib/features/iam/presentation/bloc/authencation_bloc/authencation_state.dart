part of 'authencation_bloc.dart';

@immutable
abstract class AuthencationState extends Equatable {
  const AuthencationState();

  @override
  List<Object?> get props => [];
}

class InitAuthencationState extends AuthencationState {
  const InitAuthencationState();

  @override
  List<Object?> get props => [];
}

class LoadingAuthencationState extends AuthencationState {
  const LoadingAuthencationState();

  @override
  List<Object?> get props => [];
}

class EmptyAuthencationState extends AuthencationState {
  const EmptyAuthencationState();

  @override
  List<Object?> get props => [];
}

class ErrorAuthencationState extends AuthencationState {
  final String message;
  const ErrorAuthencationState({required this.message});

  @override
  List<Object?> get props => [message];
}

/// ////////////////////////////////////////////////////////////////////////////////////////////////
class LoginState extends AuthencationState {
  const LoginState();

  @override
  List<Object?> get props => [];
}
