import 'dart:async';
import 'dart:ui';
import 'package:rongchoi_application/features/game/ecs/system/socket_service.dart';

import '../component.dart';
import '../entity.dart';

class NetworkSyncSystem {
  final SocketService socket;
  final Duration sendInterval;
  StreamSubscription<Map<String, dynamic>>? _sub;
  double _acc = 0;

  NetworkSyncSystem(
      {required this.socket,
      this.sendInterval = const Duration(milliseconds: 100)});

  void start(World world) {
    _sub = socket.messages.listen((msg) {
      if (msg['type'] == 'snapshot') {
        _applySnapshot(world, msg);
      } else if (msg['type'] == 'player_update') {
        _applyPlayerUpdate(world, msg);
      }
    });
  }

  void dispose() {
    _sub?.cancel();
  }

  void update(World world, double dt) {
    _acc += dt;
    if (_acc >= sendInterval.inMilliseconds / 1000.0) {
      _acc = 0;
// send local players
      for (final e in world.query([PlayerTag, NetworkId, Position])) {
        final id = e.get<NetworkId>()!.id;
        final pos = e.get<Position>()!;
        socket.send({
          'type': 'player_update',
          'id': id,
          'x': pos.x,
          'y': pos.y,
        });
      }
    }
  }

  void _applySnapshot(World world, Map<String, dynamic> msg) {
    final players = (msg['players'] as List?) ?? [];
    for (final p in players) {
      final id = p['id'] as String;
      final x = (p['x'] as num).toDouble();
      final y = (p['y'] as num).toDouble();
      final color = (p['color'] as int?) ?? 0xFF42A5F5;

      Entity? found;
      for (final e in world.query([NetworkId])) {
        if (e.get<NetworkId>()!.id == id) {
          found = e;
          break;
        }
      }
    }
  }

  void _applyPlayerUpdate(World world, Map<String, dynamic> msg) {
  final id = msg['id'] as String;
  final x = (msg['x'] as num).toDouble();
  final y = (msg['y'] as num).toDouble();

  // tìm entity đã tồn tại
  Entity? found;
  for (final e in world.query([NetworkId])) {
    if (e.get<NetworkId>()!.id == id) {
      found = e;
      break;
    }
  }

  if (found != null) {
    // update position
    final pos = found.get<Position>();
    if (pos != null) {
      pos.x = x;
      pos.y = y;
    }
  } else {
    // nếu chưa có thì tạo mới
    final e = world.create()
      ..add(NetworkId(id))
      ..add(Position(x, y))
      ..add(Size2D(28.0, 28.0))
      ..add(CollisionBox())
      ..add(Appearance(const Color(0xFFAB47BC))); // ví dụ màu khác cho người chơi khác
  }
}

}
