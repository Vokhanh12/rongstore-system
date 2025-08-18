import 'package:dartz/dartz.dart';
import 'package:equatable/equatable.dart';
import 'package:rongchoi_application/core/error/failure.dart';
import 'package:rongchoi_application/core/usecase/usecase.dart';
import 'package:rongchoi_application/features/iam/domain/entities/tranlations_entity.dart';
import 'package:rongchoi_application/features/iam/domain/repositories/iam_repository.dart';
class LoginUsecase
    implements UseCase<List<TranlationsEntity>, ParamsLoginUsecase> {
  final AuthencationRepository authRepo;

  LoginUsecase(this.authRepo);

  @override
  Future<Either<Failure, List<TranlationsEntity>>> call(
      ParamsLoginUsecase params) {
    // TODO: implement call
    throw UnimplementedError();
  }
}

class ParamsLoginUsecase extends Equatable {
  const ParamsLoginUsecase();

  @override
  List<Object> get props => [];

  @override
  String toString() {
    return 'ParamsGetAllTranlationsLocalUsecase Params{}';
  }
}
