import '../component.dart';
import '../entity.dart';

class CollisionSystem {
  void update(World world) {
    final colliders = world.query([Position, Size2D, CollisionBox]).toList();

    for (var i = 0; i < colliders.length; i++) {
      for (var j = i + 1; j < colliders.length; j++) {
        final a = colliders[i];
        final b = colliders[j];

        final pa = a.get<Position>()!;
        final pb = b.get<Position>()!;
        final sa = a.get<Size2D>()!;
        final sb = b.get<Size2D>()!;
        final ca = a.get<CollisionBox>()!;
        final cb = b.get<CollisionBox>()!;

        final ax1 = pa.x - sa.w / 2, ax2 = pa.x + sa.w / 2;
        final ay1 = pa.y - sa.h / 2, ay2 = pa.y + sa.h / 2;

        final bx1 = pb.x - sb.w / 2, bx2 = pb.x + sb.w / 2;
        final by1 = pb.y - sb.h / 2, by2 = pb.y + sb.h / 2;

        final overlapX = (ax1 < bx2) && (ax2 > bx1);
        final overlapY = (ay1 < by2) && (ay2 > by1);

        if (overlapX && overlapY) {
          final dx1 = ax2 - bx1; 
          final dx2 = bx2 - ax1; 
          final dy1 = ay2 - by1; 
          final dy2 = by2 - ay1; 
          final minX = (dx1.abs() < dx2.abs()) ? dx1 : -dx2;
          final minY = (dy1.abs() < dy2.abs()) ? dy1 : -dy2;

          if (minX.abs() < minY.abs()) {
            _separate(pa, pb, ca, cb, minX, 0);
          } else {
            _separate(pa, pb, ca, cb, 0, minY);
          }
        }
      }
    }
  }

  void _separate(Position pa, Position pb, CollisionBox ca, CollisionBox cb,
      double sx, double sy) {
    if (ca.isStatic && !cb.isStatic) {
      pb.x += sx;
      pb.y += sy;
    } else if (!ca.isStatic && cb.isStatic) {
      pa.x -= sx;
      pa.y -= sy;
    } else if (!ca.isStatic && !cb.isStatic) {
      pa.x -= sx / 2;
      pa.y -= sy / 2;
    }
  }
}
