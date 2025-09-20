import 'package:rongchoi_application/features/game/map/spawn.dart';

abstract class IMap {
  final Spawn spawn;

  IMap({required this.spawn});

  Future<void> init();
}

class MapHome extends IMap {
  MapHome({required super.spawn});

  @override
  Future<void> init() {
    // TODO: implement init
    throw UnimplementedError();
  }
}
