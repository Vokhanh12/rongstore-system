import 'package:flame/components.dart';
import 'package:rongchoi_application/core/constants/game_assets.dart';
import 'package:rongchoi_application/features/game/ecs/component.dart';
import 'package:rongchoi_application/features/game/ecs/entity.dart';
import 'package:rongchoi_application/features/game/tools/grid_2d.dart';
import 'package:rongchoi_application/features/game/world/entiy_builder.dart';

abstract class ISpriteBuilder {
  final EntiyBuilder entiyBuilder;

  ISpriteBuilder({required this.entiyBuilder});

  Future<Entity> thousandRoad();
  Future<Entity> sideWalk();
  Future<Entity> sideHouse();
  Future<Entity> cyclingInRoad();
  Future<Entity> blackSky();
  Future<Entity> player();
  Future<Entity> house();
}

class SpriteBuilder extends ISpriteBuilder {
  SpriteBuilder({required super.entiyBuilder});

  @override
  Future<Entity> blackSky() {
    throw UnimplementedError();
  }

  @override
  Future<Entity> cyclingInRoad() {
    throw UnimplementedError();
  }

  @override
  Future<Entity> house() async{
    this.images.prefix = "assets/";
    final image = await this.images.load('game/png/house.png');
    final sprite = Sprite(image);

    double cellSize = gridConfig.size;

    return entiyBuilder.create()
      ..add(Position(2 * cellSize, 5 * cellSize))
      ..add(CustomSprite(sprite))
      ..add(Size2D(13 * cellSize, 18 * cellSize));
  }

  @override
  Future<Entity> player() async{
     return entiyBuilder.create()
      ..add(PlayerTag())
      ..add(Position(300.0, 20.0))
      ..add(Velocity())
      ..add(Size2D(28.0, 28.0))
      ..add(CollisionBox())
      ..add(Direction(facingLeft: true))
      ..add(Transform(anchor: Anchor.center))
      ..add(AnimationData(asset: AppGameAssets.catRun, rows: 2, cols: 3, stepTime: 0.1));
  }

  @override
  Future<Entity> sideHouse() {
    // TODO: implement sideHouse
    throw UnimplementedError();
  }

  @override
  Future<Entity> sideWalk() {
    // TODO: implement sideWalk
    throw UnimplementedError();
  }

  @override
  Future<Entity> thousandRoad() {
    // TODO: implement thousandRoad
    throw UnimplementedError();
  }
}
