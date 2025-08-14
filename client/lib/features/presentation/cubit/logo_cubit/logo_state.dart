part of 'logo_cubit.dart';

abstract class LogoState extends Equatable {
  const LogoState();

  @override
  List<Object?> get props => [];
}

class LogoStateInitial extends LogoState {}

class SwapLogoInitial extends LogoState {}

class SwapLogoLoading extends LogoState {}

class SwapLogoLoaded extends LogoState {
  const SwapLogoLoaded();

  @override
  List<Object?> get props => [];
}

class SwapLogoError extends LogoState {
  final String errorMessage;

  const SwapLogoError({required this.errorMessage});

  @override
  List<Object?> get props => [errorMessage];
}
