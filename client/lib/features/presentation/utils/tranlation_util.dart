import 'package:rongchoi_application/features/domain/entities/tranlations_entity.dart';

class TranlationUtil {
  static String getTranlationsByCode(
      List<TranlationsEntity> filterData, String tranlationCode) {
    var a = filterData
        .firstWhere(
          (item) => item.code == tranlationCode,
          orElse: () => TranlationsEntity(
              id: 0, code: 'none', tranlationEn: 'none', tranlationVi: 'none'),
        )
        .tranlationEn;
    return a;
  }
}
