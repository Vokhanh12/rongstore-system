
import 'package:equatable/equatable.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

part 'form_auto_state.dart';
part 'form_auto_event.dart';

class FormAutoBloc extends Bloc<FormAutoEvent, FormAutoState> {

  dynamic _data;

  FormAutoBloc(super.initialState);

  Future<void> getFormDataEvent(event, emit) async {
    emit(GetFormDataState(dynamicData: event));
  }

  Future<void> setFormDataEvent(event, emit) async {
    emit(SetFormDataState(dynamicData: event));
  }

  Future<void> setValueInFormDataEvent(event, emit) async {
    emit(SetValueInFormDataState(valueKey: event));
  }

   Future<void> SubmitFormDataEvent(event, emit) async {
    emit(submitFormDataState(dynamicData: event));
  }

}
