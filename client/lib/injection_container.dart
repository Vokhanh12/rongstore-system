import 'package:get_it/get_it.dart';
import 'package:rongchoi_application/features/data/datasources/db/database_helper.dart';
import 'package:rongchoi_application/features/data/datasources/local/tranlation_local_datasource.dart';
import 'package:rongchoi_application/features/data/repositories/tranlation_repository_impl.dart';
import 'package:rongchoi_application/features/domain/repositories/tranlation_repository.dart';
import 'package:rongchoi_application/features/domain/usecases/tranlation_usecases/get_all_tranlations_local_usecase.dart';
import 'package:rongchoi_application/features/presentation/bloc/tranlation_bloc/tranlation_bloc.dart';

final locator = GetIt.instance;

Future<void> initLocator() async {

  /// Database
  locator.registerLazySingleton<DatabaseHelper>(() => DatabaseHelper());

  /// DataSources ////////////////////////////////////////////////////////////////////////////////////////////////
  /// /// User DataSource
  locator.registerLazySingleton<TranlationLocalDataSource>(() => TrnaltionLocalDataSourceImpl(locator()));

  /// Reopsitory ////////////////////////////////////////////////////////////////////////////////////////////////
  /// /// Tranlations
  locator.registerLazySingleton<TranlationRepository>(() => TranlationRepositoryImpl(
      tranlationLocalDataSource: locator(),
    ),
  );

  /// Bloc ////////////////////////////////////////////////////////////////////////////////////////////////
  /// /// Group Bloc
  locator.registerFactory(
    () => TranlationBloc(
    getAllTranlationsLocalUsecase: locator()
    ),
  );

  /// Usecase ////////////////////////////////////////////////////////////////////////////////////////////////
  /// /// Group Bloc
  locator.registerLazySingleton<GetAllTranlationsLocalUsecase>(() => GetAllTranlationsLocalUsecase(locator()));


}