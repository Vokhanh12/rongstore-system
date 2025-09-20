import 'package:flame/components.dart' as comp;
import 'package:flame/game.dart';
import 'package:flame/sprite.dart';
import 'package:rongchoi_application/features/game/ecs/component.dart';
import 'package:rongchoi_application/features/game/ecs/entity.dart';

import '../../world/world.dart';

class DebugFlameRenderSystem {
  final FlameGame game;
  final Map<Entity, comp.SpriteAnimationComponent> cache = {};

  DebugFlameRenderSystem({required this.game});

  Future<void> ensureLoaded(World world) async {
    for (final e in world.entities) {
      final animData = e.get<AnimationData>();
      final pos = e.get<Position>();

      if (animData != null && pos != null && !cache.containsKey(e)) {
        game.images.prefix = '';
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

        cache[e] = flameComp;
        game.add(flameComp);
      }
    }
  }

  void sync(World world) {
    for (final e in world.entities) {

      if (!cache.containsKey(e)) {
        final animData = e.get<AnimationData>();
        final pos = e.get<Position>();
        final size = e.get<Size2D>();
        final cusSprite = e.get<CustomSprite>();


        if (animData != null && pos != null) {
          final image = game.images.fromCache(animData.asset);
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
            to: animData.cols * animData.rows,
          );

          final flameComp = comp.SpriteAnimationComponent(
            animation: anim,
            position: Vector2(pos.x, pos.y),
            size: Vector2(128, 128),
            anchor: comp.Anchor.center,
          );

          cache[e] = flameComp;
          game.add(flameComp);
        } 
      }

      final pos = e.get<Position>();
      final dir = e.get<Direction>();
      final sprite = cache[e];

      if (sprite != null) {
        if (pos != null) {
          sprite.position = Vector2(pos.x, pos.y); 
        }
        if (dir != null) {
          sprite.scale.x = dir.facingLeft ? -1 : 1;
        }
      }
    }
  }
}
