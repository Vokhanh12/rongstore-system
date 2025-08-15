part of 'tranlation_bloc.dart';

@immutable
abstract class TranlationState extends Equatable {
  const TranlationState();

  @override
  List<Object?> get props => [];
}

class InitTranlationState extends TranlationState {
  const InitTranlationState();

  @override
  List<Object?> get props => [];
}

class LoadingTranlationState extends TranlationState {
  const LoadingTranlationState();

  @override
  List<Object?> get props => [];
}

class EmptyTranlationState extends TranlationState {
  const EmptyTranlationState();

  @override
  List<Object?> get props => [];
}

class ErrorTranlationState extends TranlationState {
  final String message;
  const ErrorTranlationState({required this.message});

  @override
  List<Object?> get props => [message];
}

/// ////////////////////////////////////////////////////////////////////////////////////////////////
class GetAllTranlationsLocalState extends TranlationState {
  final List<TranlationsEntity> tranlationItems;

  const GetAllTranlationsLocalState({required this.tranlationItems});

  @override
  List<Object?> get props => [tranlationItems];
}