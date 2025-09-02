import 'dart:ui';


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
// Use a color for demo; replace by sprite handle if needed
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