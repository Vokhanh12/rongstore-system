import 'dart:ui';
import 'package:flame/components.dart' as comp;
import 'package:flame/game.dart';
import 'package:flame/sprite.dart';
import 'package:flame_rive/flame_rive.dart';
import 'package:flutter/material.dart' as mtr;
import 'package:rongchoi_application/features/game/ecs/component.dart';
import 'package:rongchoi_application/features/game/ecs/entity.dart';

class FlameSystem {
  final FlameGame game;
  final Map<Entity, comp.SpriteAnimationComponent> cacheSAC = {};
  final Map<Entity, RiveComponent> cacheRC = {};

  FlameSystem({required this.game}) {
    game.images.prefix = '';
  }

  Future<void> onLoad(World world) async {
    for (final e in world.entities) {
      final animData = e.get<AnimationData>();
      final pos = e.get<Position>();
      final cus = e.get<CustomSprite>();
      final size2d = e.get<Size2D>();
      final riveAnimationData = e.get<RiveAnimationData>();
      final riveData = e.get<RiveData>();

      if (animData != null && pos != null) {
        final image = await game.images.load(animData.asset);

        final spriteSheet = SpriteSheet(
          image: image,
          srcSize: Vector2(
            image.width / animData.cols,
            image.height / animData.rows,
          ),
        );

        final anim = spriteSheet.createAnimation(
          row: 0,
          stepTime: animData.stepTime,
          to: animData.cols,
        );

        final flameComp = comp.SpriteAnimationComponent(
          animation: anim,
          position: Vector2(pos.x, pos.y),
          size: Vector2(128, 128),
        );

        flameComp.anchor = comp.Anchor.center;
        cacheSAC[e] = flameComp;
        game.add(flameComp);
      } else if (cus != null && pos != null && size2d != null) {
        final sprite = comp.SpriteComponent()
          ..sprite = cus.sprite
          ..size = Vector2(size2d.w, size2d.h)
          ..position = Vector2(pos.x, pos.y);
        game.add(sprite);
      } else if (riveData != null && pos != null && size2d != null) {
        

        final riveComponent = RiveComponent(artboard: riveData.artboard)
          ..x = pos.x
          ..y = pos.y
          ..width = size2d.w
          ..height = size2d.h;

        cacheRC[e] = riveComponent;
        game.add(riveComponent);
      }
    }
  }

  void render(World world, Canvas canvas) {
    for (final e in world.entities) {
      final pos = e.get<Position>();
      final size = e.get<Size2D>();
      final app = e.get<Appearance>();
      final cusSprite = e.get<CustomSprite>();
      final anim = e.get<AnimationData>();

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

  void update(World world) {
    for (final e in world.entities) {
      final pos = e.get<Position>();
      final anim = e.get<AnimationData>();
      final size2d = e.get<Size2D>();
      final riveAnimationData = e.get<RiveAnimationData>();
      final riveData = e.get<RiveData>();

     
    }
  }
}
