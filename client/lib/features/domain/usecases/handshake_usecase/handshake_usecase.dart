import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import 'package:rongchoi_application/core/error/failure.dart';
import 'package:rongchoi_application/core/usecase/usecase.dart';
import 'package:rongchoi_application/features/domain/entities/tranlations_entity.dart';
import 'package:rongchoi_application/features/domain/repositories/authencation_repository.dart';
import 'package:rongchoi_application/features/domain/repositories/tranlation_repository.dart';

class HandshakeUsecase
    implements UseCase<List<TranlationsEntity>, ParamsHandshakeUsecase> {
  final AuthencationRepository authRepo;

  HandshakeUsecase(this.authRepo);

  @override
  Future<Either<Failure, List<TranlationsEntity>>> call(
      ParamsHandshakeUsecase params) {
    // TODO: implement call
    throw UnimplementedError();
  }
}

class ParamsHandshakeUsecase extends Equatable {
  const ParamsHandshakeUsecase();

  @override
  List<Object> get props => [];

  @override
  String toString() {
    return 'ParamsGetAllTranlationsLocalUsecase Params{}';
  }
}
