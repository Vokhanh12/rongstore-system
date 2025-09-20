import 'package:rongchoi_application/features/game/tools/grid_2d.dart';

abstract class Layout2d
{ 
  
  final Grid2d grid;

  Layout2d({required this.grid});

  Future <void> render();
}

class HouseLayout2d extends Layout2d{
  HouseLayout2d({required super.grid});

  @override
  Future<void> render() {
    // TODO: implement render
    throw UnimplementedError();
  }
} 

class ThousandRoadLayout2d extends Layout2d{
  ThousandRoadLayout2d({required super.grid});

  @override
  Future<void> render() {
    // TODO: implement render
    throw UnimplementedError();
  }

}

class WalkRoadLayout2d extends Layout2d{
  WalkRoadLayout2d({required super.grid});

  @override
  Future<void> render() {
    // TODO: implement render
    throw UnimplementedError();
  }

}

class SkyLayout2d extends Layout2d{
  SkyLayout2d({required super.grid});

  @override
  Future<void> render() {
    // TODO: implement render
    throw UnimplementedError();
  }

}




