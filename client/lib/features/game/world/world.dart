import 'package:rongchoi_application/features/game/world/entiy_builder.dart';

class World {
  late EntiyBuilder entiyBuilder;

  World() {
    this.entiyBuilder = entiyBuilderSingleton;
  }
}
