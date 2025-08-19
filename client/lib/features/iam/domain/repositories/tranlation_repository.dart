import 'package:dartz/dartz.dart';
import 'package:rongchoi_application/core/error/failure.dart';
import 'package:rongchoi_application/features/iam/domain/entities/tranlations_entity.dart';

abstract class TranlationRepository{
  /// Local
  Future<Either<Failure, List<TranlationsEntity>>> getAllTranlationsLocal();

  /// Remote
  Future<Either<Failure, List<TranlationsEntity>>> getAllTranlationsRemote();
}
