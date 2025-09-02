import 'component.dart';

class Entity {
  final int id;
  final Map<Type, Component> _components = {};

  Entity(this.id);

  void add(Component c) => _components[c.runtimeType] = c;

  T? get<T extends Component>() => _components[T] as T?;

  bool has<T extends Component>() => _components.containsKey(T);

  void remove<T extends Component>() => _components.remove(T);
}

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
        if (!e._components.containsKey(t)) {
          ok = false;
          break;
        }
      }
      if (ok) yield e;
    }
  }
}
