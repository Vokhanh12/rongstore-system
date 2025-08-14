import 'dart:async';

import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:meta/meta.dart';
import 'package:rongchoi_application/core/constants/constants.dart';
import 'package:rongchoi_application/features/domain/entities/tranlations_entity.dart';
import 'package:rongchoi_application/features/domain/usecases/tranlation_usecases/get_all_tranlations_local_usecase.dart';

part 'tranlation_state.dart';
part 'tranlation_event.dart';

class TranlationBloc extends Bloc<TranlationEvent, TranlationState> {
  final GetAllTranlationsLocalUsecase getAllTranlationsLocalUsecase;

  TranlationBloc({required this.getAllTranlationsLocalUsecase})
      : super(const InitTranlationState()) {
    on<GetAllTranlationsLocalEvent>(getAllTranlationsLocalEvent);
  }

  Future<void> getAllTranlationsLocalEvent(event, emit) async {
    emit(const LoadingTranlationState());
    await Future.delayed(Duration(seconds: 5));

    try {
      final function = await getAllTranlationsLocalUsecase
          .call(ParamsGetAllTranlationsLocalUsecase());
      function.fold((failure) {
        emit(const ErrorTranlationState(
            message: Constants.databaseFailureMessage));
      }, (data) {
        emit(GetAllTranlationsLocalState(tranlationItems: data));
      });
    } catch (e) {
      emit(ErrorTranlationState(message: e.toString()));
    }
  }
}
