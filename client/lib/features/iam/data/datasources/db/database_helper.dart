import 'package:objectbox/objectbox.dart';
import 'package:path/path.dart' as path;
import 'package:path_provider/path_provider.dart';
import 'package:rongchoi_application/features/iam/domain/entities/tranlations_entity.dart';
import 'package:rongchoi_application/objectbox.g.dart';


class DatabaseHelper {
    static DatabaseHelper? _databaseHelper;
  DatabaseHelper._instance() {
    _databaseHelper = this;
  }

  factory DatabaseHelper() => _databaseHelper ?? DatabaseHelper._instance();

  static Store? _store;

  Future<Store?> get store async {
    _store ??= await _create();
    return _store;
  }

  Future<Store> _create() async {
    final appDir = await getApplicationDocumentsDirectory();
    final dbPath = path.join(appDir.path, "MessengerObjectBox");
    final store = openStore(directory: dbPath);
    return store;
  }

  void close() async {
    try {
      _store?.close();
    } catch (e) { return; }
  }

  Future<bool> saveTranlationLocal(TranlationsEntity tranlationsItem) async {
    final db = await store;
    await db!.box<TranlationsEntity>().putAsync(tranlationsItem, mode: PutMode.put);
    return true;
  }

 Future<List<TranlationsEntity>> getAllTranslationsLocal() async {
    final db = await store;
    final query = db!.box<TranlationsEntity>().query().build();
    final results = query.find();
    query.close();
    return results;
  }
}