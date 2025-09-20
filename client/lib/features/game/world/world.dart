import 'package:rongchoi_application/features/game/ecs/entity.dart';

class World {
  final List<Entity> entities = [];
  int _nextId = 1;

  Entity create() {
    final e = Entity(_nextId++);
    entities.add(e);
    return e;
  }

  void destroy(Entity e) => entities.remove(e);

  Iterable<Entity> query(List<Type> allOf) sync* {
    for (final e in entities) {
      var ok = true;
      for (final t in allOf) {
        if (!e.components.containsKey(t)) {
          ok = false;
          break;
        }
      }
      if (ok) yield e;
    }
  }
}