part of 'form_auto_bloc.dart';

@immutable
abstract class FormAutoState extends Equatable {
  const FormAutoState();

  @override
  List<Object?> get props => [];
}

class InitFormAutoState extends FormAutoState {
  const InitFormAutoState();

  @override
  List<Object?> get props => [];
}

class LoadingFormAutoState extends FormAutoState {
  const LoadingFormAutoState();

  @override
  List<Object?> get props => [];
}

class EmptyFormAutoState extends FormAutoState {
  const EmptyFormAutoState();

  @override
  List<Object?> get props => [];
}

class ErrorFormAutoState extends FormAutoState {
  final String message;
  const ErrorFormAutoState({required this.message});

  @override
  List<Object?> get props => [message];
}

/// ////////////////////////////////////////////////////////////////////////////////////////////////
class GetFormDataState extends FormAutoState {
  final dynamic dynamicData;

  const GetFormDataState({required this.dynamicData});

  @override
  List<Object?> get props => [];
}

class SetFormDataState extends FormAutoState {
  final dynamic dynamicData;

  const SetFormDataState({required this.dynamicData});

  @override
  List<Object?> get props => [dynamicData];
}

class SetValueInFormDataState extends FormAutoState {
  final Map<String, dynamic> valueKey;

  const SetValueInFormDataState({required this.valueKey});

  @override
  List<Object?> get props => [valueKey];
}

class submitFormDataState extends FormAutoState {
  final dynamic dynamicData;
  const submitFormDataState({required this.dynamicData});

  @override
  List<Object?> get props => [dynamicData];
}