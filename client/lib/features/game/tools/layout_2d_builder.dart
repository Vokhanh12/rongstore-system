import 'package:flutter/widgets.dart';
import 'package:rongchoi_application/features/game/tools/grid_2d.dart';

abstract class ILayout2dBuilder
{ 
  
  final Grid2d grid;
  final Canvas canvas;

  ILayout2dBuilder({required this.grid, required this.canvas});

  Future <void> render();
}

class HouseLayout2dBuilder extends ILayout2dBuilder{
  HouseLayout2dBuilder({required super.grid, required super.canvas});

  @override
  Future<void> render() {
    // TODO: implement render
    throw UnimplementedError();
  }
} 

class ThousandRoadLayout2dBuilder extends ILayout2dBuilder{
  ThousandRoadLayout2dBuilder({required super.grid, required super.canvas});

  @override
  Future<void> render() {
    // TODO: implement render
    throw UnimplementedError();
  }

}

class WalkRoadLayout2dBuilder extends ILayout2dBuilder{
  WalkRoadLayout2dBuilder({required super.grid, required super.canvas});

  @override
  Future<void> render() {
    // TODO: implement render
    throw UnimplementedError();
  }

}

class SkyLayout2dBuilder extends ILayout2dBuilder{
  SkyLayout2dBuilder({required super.grid, required super.canvas});

  @override
  Future<void> render() {
    // TODO: implement render
    throw UnimplementedError();
  }
}




