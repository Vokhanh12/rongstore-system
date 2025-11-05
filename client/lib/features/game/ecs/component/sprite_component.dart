import 'package:rongchoi_application/features/game/ecs/component/base_component.dart';
import 'dart:ui';

class Size2D extends BaseComponent {
  double width;
  double height;

  Size2D({required this.width, required this.height});
}

class SheetAnimation extends BaseComponent {
  String assetPath; 
  int frameCount;   
  double frameRate; 
  bool loop;     

  SheetAnimation({
    required this.assetPath,
    this.frameCount = 1,
    this.frameRate = 0.1,
    this.loop = true,
  });
}

class RiveAnimation extends BaseComponent {
  String assetPath;    
  String animationName; 
  bool autoplay;
  bool loop;

  RiveAnimation({
    required this.assetPath,
    required this.animationName,
    this.autoplay = true,
    this.loop = true,
  });
}

class Appearance extends BaseComponent {
  String? assetPath;    
  Color? tint;          
  double scale;
  double rotation;      
  bool flipX;
  bool flipY;
  double opacity;

  Appearance({
    this.assetPath,
    this.tint,
    this.scale = 1.0,
    this.rotation = 0.0,
    this.flipX = false,
    this.flipY = false,
    this.opacity = 1.0,
  });
}
