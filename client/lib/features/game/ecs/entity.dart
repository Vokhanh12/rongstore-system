import 'component.dart';

class Entity {
  final int id;
  final Map<Type, Component> components = {};

  Entity(this.id);

  void add(Component c) => components[c.runtimeType] = c;

  T? get<T extends Component>() => components[T] as T?;

  bool has<T extends Component>() => components.containsKey(T);

  void remove<T extends Component>() => components.remove(T);
}