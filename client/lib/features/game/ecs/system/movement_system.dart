import '../../world/world.dart';
import '../component.dart';

class MovementSystem {
  final double worldWidth;
  final double worldHeight;

  MovementSystem({required this.worldWidth, required this.worldHeight});

  void update(World world, double dt) {
    for (final e in world.query([Position, Velocity])) {
      final pos = e.get<Position>()!;
      final vel = e.get<Velocity>()!;
      
      pos.x += vel.dx * dt;
      pos.y += vel.dy * dt;

      pos.x = pos.x.clamp(0, worldWidth);
      pos.y = pos.y.clamp(0, worldHeight);

      final dir = e.get<Direction>();
      if (dir != null) {
        if (vel.dx < 0) {
          dir.facingLeft = true;
        } else if (vel.dx > 0) {
          dir.facingLeft = false;
        }
      }
    }
  }
}

