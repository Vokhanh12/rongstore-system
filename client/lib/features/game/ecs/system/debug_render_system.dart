import 'dart:ui';
import 'package:flame/game.dart';
import 'package:flutter/material.dart' as mtr;
import 'package:rongchoi_application/features/game/ecs/component.dart';
import 'package:rongchoi_application/features/game/ecs/entity.dart';

class DebugRenderSystem {
  void render(World world, Canvas canvas) {
    for (final e in world.entities) {
      final pos = e.get<Position>();
      final size = e.get<Size2D>();
      final app = e.get<Appearance>();
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
      }
    }
  }
}
