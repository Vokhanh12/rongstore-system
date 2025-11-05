import 'package:rongchoi_application/features/game/ecs/component/base_component.dart';
import 'dart:ui';

class SpriteLayoutComponent extends BaseComponent {
  final bool showFrame;       // có vẽ khung quanh entity không
  final bool showPosition;    // có vẽ điểm position không
  final Color frameColor;     // màu khung
  final Color positionColor;  // màu điểm position
  final double frameStroke;   // độ dày khung
  final double pointRadius;   // bán kính điểm position

  SpriteLayoutComponent({
    this.showFrame = true,
    this.showPosition = true,
    this.frameColor = const Color(0xFFFF0000),     // đỏ
    this.positionColor = const Color(0xFF00FF00),  // xanh lá
    this.frameStroke = 1.0,
    this.pointRadius = 4.0,
  });
}
