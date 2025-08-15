part of 'form_auto_bloc.dart';

@immutable
abstract class FormAutoEvent extends Equatable {
  const FormAutoEvent();

  @override
  List<Object?> get props => [];
}

class SetDataFormDataEvent extends FormAutoEvent {
  const SetDataFormDataEvent();

  @override
  List<Object> get props => [];
}

class GetDataFormDataEvent extends FormAutoEvent {
  const GetDataFormDataEvent();

  @override
  List<Object> get props => [];
}
