import 'dart:io';
import 'dart:math' as math;
import 'dart:ui';
import 'package:flame/events.dart';
import 'package:flame/experimental.dart';
import 'package:flame/game.dart';
import 'package:flame/input.dart';
import 'package:flutter/material.dart';
import 'package:rongchoi_application/core/constants/game_assets.dart';
import 'package:rongchoi_application/features/game/ecs/system/debug_flame_render_system.dart';
import 'package:uuid/uuid.dart';

import 'package:rongchoi_application/features/game/ecs/system/collision_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/debug_render_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/movement_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/udp_service.dart';

import 'ecs/component.dart';
import 'ecs/entity.dart';
import 'package:flame/components.dart' hide World;
import 'ecs/entity.dart' as ecs;
import 'ecs/component.dart' as comp;

class Snapshot {
  final double timestamp;
  final double x;
  final double y;

  Snapshot({required this.timestamp, required this.x, required this.y});
}

class MyGame extends FlameGame with TapDetector, HasKeyboardHandlerComponents {
  final ecs.World ecsWorld = ecs.World();

  late final MovementSystem movementSystem;
  late final CollisionSystem collisionSystem;
  late final DebugRenderSystem renderSystem;
  late final DebugFlameRenderSystem flameRenderSystem;

  double worldW = 0;
  double worldH = 0;

  static const double speed = 200.0;
  static const double arrivalThreshold = 6.0;

  double shopTop = 0.0;
  double shopBottom = 0.0;
  double walkTop = 0.0;
  double walkBottom = 0.0;
  double streetTop = 0.0;
  double streetBottom = 0.0;

  Entity? localPlayer;
  double? _targetX;
  double? _targetY;

  UdpService? _udp;
  final sendInterval = 0.05;
  double _sendAcc = 0;

  final String playerId = Uuid().v4();
  final Map<String, Entity> remotePlayers = {};
  final Map<String, List<Snapshot>> remoteSnapshots = {};

  @override
  Future<void> onLoad() async {
    await super.onLoad();

    worldW = size.x;
    worldH = size.y;

    final sectionHeight = worldH / 3;
    shopTop = 0;
    shopBottom = sectionHeight;
    walkTop = shopBottom;
    walkBottom = walkTop + sectionHeight;
    streetTop = walkBottom;
    streetBottom = worldH;

    camera.setBounds(Rectangle.fromLTRB(0, 0, worldW, worldH),
        considerViewport: true);

    movementSystem = MovementSystem(worldWidth: worldW, worldHeight: worldH);
    collisionSystem = CollisionSystem();
    renderSystem = DebugRenderSystem();
    flameRenderSystem = DebugFlameRenderSystem(game: this);

    _udp = UdpService(
        serverAddr: InternetAddress("100.114.31.30"), serverPort: 8080);

    try {
      await _udp!.connect();
      print('UDP connect success');
    } catch (e) {
      print('UDP connect error: $e');
    }

    _spawnLocalPlayer();

    await flameRenderSystem.ensureLoaded(ecsWorld);

    _udp!.messages.listen((msg) {
      _handleIncoming(msg);
    });
  }

  void _spawnLocalPlayer() {
    final me = ecsWorld.create()
      ..add(PlayerTag())
      ..add(NetworkId(playerId))
      ..add(Position(300.0, walkTop + 20.0))
      ..add(comp.Velocity())
      ..add(Size2D(28.0, 28.0))
      ..add(CollisionBox())
      ..add(Direction(facingLeft: true))
      ..add(AnimationData(AppGameAssets.catRun, 2, 3, 0.1));
    localPlayer = me;
  }

  void _handleIncoming(Map<String, dynamic> msg) {
    final type = msg['type'];
    final timestamp = DateTime.now().millisecondsSinceEpoch / 1000.0;

    if (type == 'snapshot' && msg['players'] is List) {
      final players = (msg['players'] as List).cast<Map<String, dynamic>>();
      for (var p in players) {

        final id = p['id'] as String;

        if (id == playerId) continue;

        final targetX = (p['x'] as num?)?.toDouble() ?? 0.0;
        final targetY = (p['y'] as num?)?.toDouble() ?? 0.0;

        if (!remotePlayers.containsKey(id)) {
            final e = ecsWorld.create()
            ..add(PlayerTag())
            ..add(NetworkId(id))
            ..add(Position(targetX, targetY))
            ..add(Size2D(28.0, 28.0))
            ..add(CollisionBox())
            ..add(AnimationData(AppGameAssets.catRun, 2, 3, 0.1));
          remotePlayers[id] = e;
          remoteSnapshots[id] = [Snapshot(timestamp: timestamp, x: targetX, y: targetY)];

        } else {
          remoteSnapshots[id]!
              .add(Snapshot(timestamp: timestamp, x: targetX, y: targetY));
          if (remoteSnapshots[id]!.length > 5)
            remoteSnapshots[id]!.removeAt(0);
        }
      }
    }
  }

  void _updateRemotePlayers(double dt) {
    final now = DateTime.now().millisecondsSinceEpoch / 1000.0;
    remotePlayers.forEach((id, entity) {
      final pos = entity.get<Position>();
      final snapshots = remoteSnapshots[id];
      if (pos != null && snapshots != null && snapshots.length >= 2) {
        final renderTime = now - 0.1;
        Snapshot? prev, next;
        for (int i = 0; i < snapshots.length - 1; i++) {
          if (snapshots[i].timestamp <= renderTime &&
              snapshots[i + 1].timestamp >= renderTime) {
            prev = snapshots[i];
            next = snapshots[i + 1];
            break;
          }
        }
        if (prev != null && next != null) {
          final t =
              (renderTime - prev.timestamp) / (next.timestamp - prev.timestamp);
          pos.x = prev.x + (next.x - prev.x) * t;
          pos.y = prev.y + (next.y - prev.y) * t;
        } else if (snapshots.isNotEmpty) {
          pos.x += (snapshots.last.x - pos.x) * 0.15;
          pos.y += (snapshots.last.y - pos.y) * 0.15;
        }
      }
    });
  }

  @override
  void update(double dt) {
    super.update(dt);

    if (_targetX != null && _targetY != null) {
      _setVelocityTowardsTarget();
    }

    movementSystem.update(ecsWorld, dt);
    collisionSystem.update(ecsWorld);

    final p = localPlayer?.get<Position>();
    final s = localPlayer?.get<Size2D>();
    if (p != null && s != null) {
      if (p.y < walkTop) p.y = walkTop;
      if (p.y + s.h > walkBottom) p.y = walkBottom - s.h;
      camera.moveTo(Vector2(p.x, p.y));
    }

    _sendAcc += dt;
    if (_sendAcc >= sendInterval) {
      _sendAcc = 0;
      if (localPlayer != null) {
        final pos = localPlayer!.get<Position>()!;
        _udp?.send({
          'type': 'player_update',
          'id': playerId,
          'x': pos.x,
          'y': pos.y,
        });
      }
    }

    _updateRemotePlayers(dt);

    flameRenderSystem.sync(ecsWorld);
  }

  void _setVelocityTowardsTarget() {
    if (localPlayer == null || _targetX == null || _targetY == null) return;
    final pos = localPlayer!.get<Position>();
    final vel = localPlayer!.get<comp.Velocity>();
    if (pos == null || vel == null) return;

    final dx = _targetX! - pos.x;
    final dy = _targetY! - pos.y;
    final dist = math.sqrt(dx * dx + dy * dy);
    if (dist <= arrivalThreshold) {
      vel.dx = 0;
      vel.dy = 0;
      _targetX = null;
      _targetY = null;
    } else {
      final nx = dx / dist;
      final ny = dy / dist;
      vel.dx = nx * speed;
      vel.dy = ny * speed;
    }
  }

  @override
  void render(Canvas canvas) {
    super.render(canvas);
    
    renderSystem.render(ecsWorld, canvas);

    canvas.drawRect(Rect.fromLTWH(0, shopTop, worldW, shopBottom - shopTop),
        Paint()..color = const Color(0x55FF0000));

    canvas.drawRect(Rect.fromLTWH(0, walkTop, worldW, walkBottom - walkTop),
        Paint()..color = const Color.fromARGB(84, 55, 255, 0));

    canvas.drawRect(
        Rect.fromLTWH(0, streetTop, worldW, streetBottom - streetTop),
        Paint()..color = const Color.fromARGB(83, 3, 24, 255));
  }

  @override
  void onTapDown(TapDownInfo info) {
    final tapPos = info.eventPosition.global;
    if (tapPos.y < walkTop || tapPos.y > walkBottom) return;

    _targetX = tapPos.x;
    _targetY = tapPos.y;
  }

  @override
  void onRemove() {
    _udp?.close();
    super.onRemove();
  }
}
