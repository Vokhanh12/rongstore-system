
import 'package:rongchoi_application/core/error/exception.dart';
import 'package:rongchoi_application/features/iam/data/datasources/db/database_helper.dart';
import 'package:rongchoi_application/features/iam/domain/entities/tranlations_entity.dart';

abstract class TranlationLocalDataSource {
  Future<List<TranlationsEntity>> getALlTranlationLocal();
}

class TrnaltionLocalDataSourceImpl implements TranlationLocalDataSource {

  final DatabaseHelper databaseHelper;
  TrnaltionLocalDataSourceImpl(this.databaseHelper);

  @override
  Future<List<TranlationsEntity>> getALlTranlationLocal() async{
    try{
      return await databaseHelper.getAllTranslationsLocal();
    } catch (e) {
      throw DatabaseException();
    }
  }

}