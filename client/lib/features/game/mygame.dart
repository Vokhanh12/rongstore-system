import 'dart:io';
import 'dart:math' as math;
import 'dart:ui';

import 'package:flame/events.dart';
import 'package:flame/experimental.dart';
import 'package:flame/game.dart';
import 'package:flame/input.dart';
import 'package:flutter/material.dart';

import 'package:rongchoi_application/features/game/ecs/system/collision_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/debug_render_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/movement_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/network_sync_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/socket_service.dart';
import 'package:rongchoi_application/features/game/ecs/system/udp_service.dart';

import 'ecs/component.dart';
import 'ecs/entity.dart';

import 'package:flame/components.dart' hide World;
import 'ecs/entity.dart' as ecs;
import 'ecs/component.dart' as comp;

class MyGame extends FlameGame with TapDetector, HasKeyboardHandlerComponents {
  final ecs.World ecsWorld = ecs.World();

  late final MovementSystem movementSystem;
  late final CollisionSystem collisionSystem;
  late final DebugRenderSystem renderSystem;
  late final NetworkSyncSystem netSystem;

  double worldW = 0;
  double worldH = 0;

  static const double speed = 200.0;
  static const double arrivalThreshold = 6.0;

  // Các vùng
  double shopTop = 0.0;
  double shopBottom = 0.0;
  double walkTop = 0.0;
  double walkBottom = 0.0;
  double streetTop = 0.0;
  double streetBottom = 0.0;

  Entity? localPlayer;

  double? _targetX;
  double? _targetY;

  SocketService? _socket;
  UdpService? _udp;

  final sendInterval = 0.1;
  double _sendAcc = 0;

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

    camera.setBounds(
      Rectangle.fromLTRB(0, 0, worldW, worldH),
      considerViewport: true,
    );

    movementSystem = MovementSystem(worldWidth: worldW, worldHeight: worldH);
    collisionSystem = CollisionSystem();
    renderSystem = DebugRenderSystem();

    //_socket = SocketService();
    //netSystem = NetworkSyncSystem(socket: _socket!);

    _udp =
        UdpService(serverAddr: InternetAddress("192.168.1.4"), serverPort: 8080);

    try {
      //await _socket!.connect('ws://6a3104eb7b7a.ngrok-free.app/game');
      await _udp!.connect();
      print('Udp connect success');
    } catch (e) {
      //print('Socket connect error: $e');
      print('Udp connect error: $e');
    }
    //netSystem.start(ecsWorld);
    //_udp.send(ecsWorld)
    _spawnDemoWorld();
  }

  void _spawnDemoWorld() {
    final me = ecsWorld.create()
      ..add(PlayerTag())
      ..add(NetworkId('local'))
      ..add(Position(300.0, walkTop + 20.0))
      ..add(comp.Velocity())
      ..add(Size2D(28.0, 28.0))
      ..add(CollisionBox())
      ..add(Appearance(const Color(0xFF42A5F5))); // xanh dương
    localPlayer = me;

    // obstacles
    for (int i = 0; i < 6; i++) {
      ecsWorld.create()
        ..add(Position(200.0 + i * 140.0, 520.0))
        ..add(Size2D(80.0, 40.0))
        ..add(CollisionBox(isStatic: true))
        ..add(Appearance(const Color(0xFF8D6E63)));
    }
  }

  @override
  void update(double dt) {
    super.update(dt);

    if (_targetX != null && _targetY != null) {
      _setVelocityTowardsTarget();
    }

    movementSystem.update(ecsWorld, dt);
    collisionSystem.update(ecsWorld);
    //netSystem.update(ecsWorld, dt);

    final p = localPlayer?.get<Position>();
    final s = localPlayer?.get<Size2D>();
    if (p != null && s != null) {
      if (p.y < walkTop) p.y = walkTop;
      if (p.y + s.h > walkBottom) p.y = walkBottom - s.h;
      camera.moveTo(Vector2(p.x, p.y));
    }

    if (localPlayer != null && _targetX != null && _targetY != null) {
      final pos = localPlayer!.get<Position>()!;
      final dx = _targetX! - pos.x;
      final dy = _targetY! - pos.y;
      final dist = math.sqrt(dx * dx + dy * dy);
      if (dist <= arrivalThreshold) {
        final vel = localPlayer!.get<comp.Velocity>();
        if (vel != null) {
          vel.dx = 0;
          vel.dy = 0;
        }
        _targetX = null;
        _targetY = null;
      }
    }

    _sendAcc += dt;
    if (_sendAcc >= sendInterval) {
      _sendAcc = 0;

      if (localPlayer != null) {
        final pos = localPlayer!.get<Position>()!;
        _udp?.send({
          'type': 'player_update',
          'id': 'local',
          'x': pos.x,
          'y': pos.y,
        });
        print('test3');
      }
    }
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

    // Vùng 1 - đỏ
    canvas.drawRect(
      Rect.fromLTWH(0, shopTop, worldW, shopBottom - shopTop),
      Paint()..color = const Color(0x55FF0000),
    );

    // Vùng 2 - xanh lá
    canvas.drawRect(
      Rect.fromLTWH(0, walkTop, worldW, walkBottom - walkTop),
      Paint()..color = const Color.fromARGB(84, 55, 255, 0),
    );

    // Vùng 3 - xanh dương
    canvas.drawRect(
      Rect.fromLTWH(0, streetTop, worldW, streetBottom - streetTop),
      Paint()..color = const Color.fromARGB(83, 3, 24, 255),
    );
  }

  @override
  void onTapDown(TapDownInfo info) {
    final tapPos = info.eventPosition.global;

    if (tapPos.y < walkTop || tapPos.y > walkBottom) {
      return;
    }

    _targetX = tapPos.x;
    _targetY = tapPos.y;
  }

  @override
  void onRemove() {
    //netSystem.dispose();
    _socket?.close();
    super.onRemove();
  }
}
