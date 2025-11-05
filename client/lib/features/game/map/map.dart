import 'package:rongchoi_application/features/game/tools/sprite_builder.dart';

abstract class IMap {
  late ISpriteBuilder spriteBuilder;

  IMap({required this.spriteBuilder});

  Future<void> renderLayout();
}

class MapHome extends IMap {
  MapHome({required super.spriteBuilder});

  @override
  Future<void> renderLayout() {
    throw UnimplementedError();
  }
}
