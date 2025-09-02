import 'dart:ui';
import '../component.dart';
import '../entity.dart';

class RenderSystem {
  void render(World world, Canvas canvas) {
    for (final e in world.query([Position, Size2D, Appearance])) {
      final pos = e.get<Position>()!;
      final size = e.get<Size2D>()!;
      final app = e.get<Appearance>()!;
      final paint = Paint()..color = app.color;

      final rect = Rect.fromCenter(
        center: Offset(pos.x, pos.y),
        width: size.w,
        height: size.h,
      );

      canvas.drawRect(rect, paint);
    }
  }
}
