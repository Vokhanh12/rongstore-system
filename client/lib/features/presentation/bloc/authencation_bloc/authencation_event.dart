part of 'authencation_bloc.dart';

@immutable
abstract class AuthencationEvent extends Equatable {
  const AuthencationEvent();

  @override
  List<Object?> get props => [];
}

class LoginEvent extends AuthencationEvent {
  const LoginEvent();

  @override
  List<Object> get props => [];
}
