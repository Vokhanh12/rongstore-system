
// layout_component.dart

import 'dart:ui';

import 'package:rongchoi_application/features/game/ecs/component/base_component.dart';

class GridTileLayout extends BaseComponent {
  final int rows;        
  final int cols;        
  final double tileSize; 
  final double width;    
  final double height;   
  final Color lineColor; 
  final double strokeWidth; 

  GridTileLayout({
    this.rows = 10,
    this.cols = 10,
    this.tileSize = 32.0,
    this.width = 320.0,
    this.height = 320.0,
    this.lineColor = const Color.fromARGB(255, 238, 255, 0),
    this.strokeWidth = 1.0,
  });
}

class ParallelogramGridLayout extends BaseComponent {
  final double tileWidth;   
  final double tileHeight;  
  final int rows;
  final int cols;
  final double startX;      
  final double startY;      
  final Color color;        
  final double strokeWidth; 

  ParallelogramGridLayout({
    this.tileWidth = 44.14,     
    this.tileHeight = 22.07,    
    this.rows = 3,
    this.cols = 25,
    this.startX = 18,
    this.startY = 23 * 32,      
    this.color = const Color.fromARGB(255, 238, 0, 255),
    this.strokeWidth = 1.0,
  });
}

