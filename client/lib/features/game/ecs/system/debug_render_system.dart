import 'dart:ui';
import 'package:flutter/material.dart' as mtr;
import 'package:rongchoi_application/features/game/ecs/component.dart';
import 'package:rongchoi_application/features/game/ecs/entity.dart';
import 'package:vector_math/vector_math.dart';

class DebugRenderSystem {
  void render(World world, Canvas canvas) {
    for (final e in world.entities) {
      final pos = e.get<Position>();
      final size = e.get<Size2D>();
      final app = e.get<Appearance>();
      final cusSprite = e.get<CustomSprite>();

      if (pos != null && size != null && app != null) {
        final rect = Rect.fromLTWH(
          pos.x,
          pos.y,
          size.w,
          size.h,
        );

        final fillPaint = Paint()..color = app.color.withOpacity(0.8);
        canvas.drawRect(rect, fillPaint);

        final strokePaint = Paint()
          ..color = mtr.Colors.black
          ..style = PaintingStyle.stroke
          ..strokeWidth = 1.5;
        canvas.drawRect(rect, strokePaint);
      } else if (pos != null && size != null && cusSprite != null) {
        cusSprite.sprite.render(
          canvas,
          position: Vector2(pos.x, pos.y),
          size: Vector2(size.w, size.h),
        );
      }
    }
  }
}
