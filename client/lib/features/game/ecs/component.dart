import 'dart:ui';

import 'package:flame/components.dart';

abstract class Component {}

class Position extends Component {
  double x, y;
  Position(this.x, this.y);
}

class Velocity extends Component {
  double dx, dy;
  Velocity({this.dx = 0, this.dy = 0});
}

class Size2D extends Component {
  double w, h;
  Size2D(this.w, this.h);
}

class Appearance extends Component {
  final Color color;
  Appearance(this.color);
}

class NetworkId extends Component {
  final String id;
  NetworkId(this.id);
}

class PlayerTag extends Component {}

class CollisionBox extends Component {
  final bool isStatic;
  CollisionBox({this.isStatic = false});
}

class Animation extends Component {
  final SpriteAnimation animation;
  SpriteAnimationComponent? flameComponent;

  Animation(this.animation) {
    flameComponent = SpriteAnimationComponent(
      animation: animation,
      size: Vector2(64, 64),
    );
  }
}

class AnimationData extends Component {
  final String asset;
  final int rows, cols;
  final double stepTime;
  AnimationData(this.asset, this.rows, this.cols, this.stepTime);
}

class Direction extends Component {
  bool facingLeft;
  Direction({this.facingLeft = false});
}

class CustomSprite extends Component{
  final Sprite sprite;

  CustomSprite(this.sprite); 
}