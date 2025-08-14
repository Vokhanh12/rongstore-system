import 'dart:async';

import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:meta/meta.dart';
import 'package:rongchoi_application/core/constants/constants.dart';
import 'package:rongchoi_application/features/domain/usecases/authencation_usecase/login_usecase.dart';
import 'package:rongchoi_application/features/presentation/bloc/tranlation_bloc/tranlation_bloc.dart';

part 'authencation_state.dart';
part 'authencation_event.dart';

class AuthencationBloc extends Bloc<AuthencationEvent, AuthencationState> {
  final LoginUsecase loginUsecase;

  AuthencationBloc({required this.loginUsecase})
      : super(const InitAuthencationState()) {
    on<LoginEvent>(loginEvent);
  }

  Future<void> loginEvent(event, emit) async {
    emit(const LoadingTranlationState());
    await Future.delayed(Duration(seconds: 2));

    try {
      final function = await loginUsecase.call(ParamsLoginUsecase());
      function.fold((failure) {
        emit(const ErrorAuthencationState(
            message: Constants.databaseFailureMessage));
      }, (data) {
        emit(LoginState());
      });
    } catch (e) {
      emit(ErrorAuthencationState(message: e.toString()));
    }
  }
}
