class Grid2d {
  final double col;
  final double row;
  final double size;

  const Grid2d({
    this.col = 19,
    this.row = 33,
    this.size = 22,
  });

  double get width => col * size;
  double get height => row * size;
  double get totalCells => col * row;
}

const gridConfig = Grid2d();
