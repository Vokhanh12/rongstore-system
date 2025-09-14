import 'dart:io';
import 'dart:math' as math;
import 'dart:ui';
import 'package:flame/experimental.dart';

import 'package:flame/camera.dart';
import 'package:flame/events.dart';
import 'package:flame/game.dart';
import 'package:flame/input.dart';
import 'package:flame/components.dart' hide World;
import 'package:flame_rive/flame_rive.dart';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:uuid/uuid.dart';

import 'package:rongchoi_application/core/constants/game_assets.dart';
import 'package:rongchoi_application/features/game/ecs/system/collision_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/debug_flame_render_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/flame_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/movement_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/udp_service.dart';

import 'ecs/component.dart';
import 'ecs/component.dart' as comp;
import 'ecs/entity.dart';
import 'ecs/entity.dart' as ecs;

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
  late final FlameSystem flameSystem;

  static const double tileSize = 22.0;
  late int gridCols;
  late int gridRows;
  late double worldW;
  late double worldH;

  static const double speed = 200.0;
  static const double arrivalThreshold = 6.0;

  Entity? localPlayer;
  Entity? localThoudsandRoad;
  Entity? localCyclingInRoad;

  double? _targetX;
  double? _targetY;

  UdpService? _udp;
  final sendInterval = 0.05;
  double _sendAcc = 0;

  final String playerId = const Uuid().v4();
  final Map<String, Entity> remotePlayers = {};
  final Map<String, List<Snapshot>> remoteSnapshots = {};

  late Sprite roadTile;

  bool loadedAsset = false;

  @override
  Future<void> onLoad() async {
    await super.onLoad();

    // Kích thước màn hình thực tế (pixel logic)
    final screenSize = Vector2(
      window.physicalSize.width / window.devicePixelRatio,
      window.physicalSize.height / window.devicePixelRatio,
    );

    // Ví dụ bạn muốn map 30x17 tile (giống chuẩn 1920x1080)
    gridCols = 43; // padding -1
    gridRows = 19; // padding -1

    worldW = gridCols * tileSize;
    worldH = gridRows * tileSize;

    // Camera viewport theo kích thước màn hình
    camera.viewport = FixedResolutionViewport(resolution: screenSize);

    // Tính scale để map vừa khít màn hình
    final scaleX = screenSize.x / worldW;
    final scaleY = screenSize.y / worldH;
    camera.viewfinder.zoom = math.min(scaleX, scaleY);

    // Giới hạn camera không ra ngoài map
    camera.setBounds(
      Rectangle.fromLTWH(0, 0, worldW, worldH),
      considerViewport: true,
    );

    // Systems
    movementSystem = MovementSystem(worldWidth: worldW, worldHeight: worldH);
    collisionSystem = CollisionSystem();

    // Spawn local player
    await _spawnLocalThoudsandRoad();
    await _spawnCyclingInRoad();
    await _spawnLocalPlayer();
    loadedAsset = true;


    flameSystem = FlameSystem(game: this);
    flameSystem.onLoad(ecsWorld);

    // UDP
    _udp = UdpService(
      serverAddr: InternetAddress("100.114.31.30"),
      serverPort: 8080,
    );

    try {
      await _udp!.connect();
      print('UDP connect success');
    } catch (e) {
      print('UDP connect error: $e');
    }

    _udp!.messages.listen((msg) {
      _handleIncoming(msg);
    });

  }

  Future<void> _spawnLocalPlayer() async {
    final me = ecsWorld.create()
      ..add(PlayerTag())
      ..add(NetworkId(playerId))
      ..add(Position(300.0, 20.0))
      ..add(comp.Velocity())
      ..add(Size2D(28.0, 28.0))
      ..add(CollisionBox())
      ..add(Direction(facingLeft: true))
      ..add(AnimationData(
          asset: AppGameAssets.catRun, rows: 2, cols: 3, stepTime: 0.1));
    localPlayer = me;
  }

  Future<void> _spawnLocalThoudsandRoad() async {
    this.images.prefix = "assets/";
    final image = await this.images.load('game/png/thousand-road.png');
    final sprite = Sprite(image);

    final roadEntity = ecsWorld.create()
      ..add(Position(0, 16 * tileSize))
      ..add(CustomSprite(sprite))
      ..add(Size2D(43 * tileSize, 3 * tileSize));

    localThoudsandRoad = roadEntity;
  }

  Future<void> _spawnCyclingInRoad() async {

    final data = await rootBundle.load(
      'assets/game/rive/cycling-in-the-road.riv',
    );

    final file = RiveFile.import(data);

    final artboard = file.mainArtboard;

    artboard.addController(SimpleAnimation('GMKGSEP'));

    final cyclingEntity = ecsWorld.create()
      ..add(RiveData(artboard: artboard))
      ..add(RiveAnimationData(x1: 0, y1: 4 * tileSize, x2: 16 * tileSize, y2: 4 * tileSize))
      ..add(Position(0, 16 * tileSize))
      ..add(Size2D(43 * tileSize, 3 * tileSize));

    localCyclingInRoad = cyclingEntity;
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
            ..add(AnimationData(
          asset: AppGameAssets.catRun, rows: 2, cols: 3, stepTime: 0.1));
          remotePlayers[id] = e;
          remoteSnapshots[id] = [
            Snapshot(timestamp: timestamp, x: targetX, y: targetY)
          ];
        } else {
          remoteSnapshots[id]!
              .add(Snapshot(timestamp: timestamp, x: targetX, y: targetY));
          if (remoteSnapshots[id]!.length > 5) {
            remoteSnapshots[id]!.removeAt(0);
          }
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
    if (p != null) {
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

    if(loadedAsset){
      flameSystem.update(ecsWorld);
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

    // if(loadedAsset)
    // {
    //   flameSystem.render(ecsWorld, canvas);
    // }

    // Vẽ grid tile
    final gridPaint = Paint()
      ..color = const Color.fromARGB(255, 238, 255, 0)
      ..strokeWidth = 1;

    for (int r = 0; r <= gridRows; r++) {
      final y = r * tileSize;
      canvas.drawLine(Offset(0, y), Offset(worldW, y), gridPaint);
    }
    for (int c = 0; c <= gridCols; c++) {
      final x = c * tileSize;
      canvas.drawLine(Offset(x, 0), Offset(x, worldH), gridPaint);
    }

    // Hiển thị dot localPlayer
    final pos = localPlayer?.get<Position>();
    if (pos != null) {
      final dotPaint = Paint()..color = Colors.red;
      canvas.drawCircle(Offset(pos.x, pos.y), 5, dotPaint);

      final textPainter = TextPainter(
        text: TextSpan(
          text: "(${pos.x.toStringAsFixed(1)}, ${pos.y.toStringAsFixed(1)})",
          style: const TextStyle(color: Colors.red, fontSize: 14),
        ),
        textDirection: TextDirection.ltr,
      );
      textPainter.layout();
      textPainter.paint(canvas, Offset(pos.x, pos.y - 20));
    }

    // Hiển thị remotePlayers
    remotePlayers.forEach((id, entity) {
      final rPos = entity.get<Position>();
      if (rPos != null) {
        final dotPaint = Paint()..color = Colors.green;
        canvas.drawCircle(Offset(rPos.x, rPos.y), 4, dotPaint);

        final textPainter = TextPainter(
          text: TextSpan(
            text:
                "($id)\n(${rPos.x.toStringAsFixed(1)}, ${rPos.y.toStringAsFixed(1)})",
            style: const TextStyle(color: Colors.green, fontSize: 12),
          ),
          textDirection: TextDirection.ltr,
        );
        textPainter.layout();
        textPainter.paint(canvas, Offset(rPos.x, rPos.y - 20));
      }
    });

    // Chấm tại target
    if (_targetX != null && _targetY != null) {
      final targetPaint = Paint()..color = Colors.blue;
      canvas.drawCircle(Offset(_targetX!, _targetY!), 6, targetPaint);
    }
  }

  @override
  void onTapDown(TapDownInfo info) {
    final tapPos = info.eventPosition.global;
    _targetX = tapPos.x;
    _targetY = tapPos.y;
  }

  @override
  void onRemove() {
    _udp?.close();
    super.onRemove();
  }
}
