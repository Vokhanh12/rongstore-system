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
import 'package:rongchoi_application/features/game/ecs/system/flame_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/movement_system.dart';
import 'package:rongchoi_application/features/game/ecs/system/udp_service.dart';

import 'ecs/component.dart';
import 'ecs/component.dart' as comp;
import 'ecs/entity.dart';
import 'world/world.dart' as ecs;

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
  Entity? localBlackSkyInRoad;

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
  Color backgroundColor() => const Color(0xFF87CEEB); // xanh trời

  @override
  Future<void> onLoad() async {
    await super.onLoad();

    // Kích thước màn hình thực tế (pixel logic)
    final screenSize = Vector2(
      window.physicalSize.width / window.devicePixelRatio,
      window.physicalSize.height / window.devicePixelRatio,
    );

    // Ví dụ bạn muốn map 30x17 tile (giống chuẩn 1920x1080)
    gridCols = 19; // padding -1
    gridRows = 33; // padding -1

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
    await _spawnLocalSideWalk();
    await _spawnLocalSideHouse();
    await _spawnCyclingInRoad();
    await _spawnBlackSky();
    await _spawnLocalPlayer();
    await _spawnLocalHouse();

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
      ..add(comp.Transform(anchor: Anchor.center))
      ..add(AnimationData(
          asset: AppGameAssets.catRun, rows: 2, cols: 3, stepTime: 0.1));
    localPlayer = me;
  }

  Future<void> _spawnLocalThoudsandRoad() async {
    this.images.prefix = "assets/";
    final image = await this.images.load('game/png/thousand-road.png');
    final sprite = Sprite(image);

    final roadEntity = ecsWorld.create()
      ..add(Position(0, 26 * tileSize))
      ..add(CustomSprite(sprite))
      ..add(Size2D(43 * tileSize, 3 * tileSize));

    localThoudsandRoad = roadEntity;
  }

  Future<void> _spawnLocalSideWalk() async {
    this.images.prefix = "assets/";
    final image = await this.images.load('game/png/side-walk.png');
    final sprite = Sprite(image);

    final roadEntity = ecsWorld.create()
      ..add(Position(0, 23 * tileSize))
      ..add(CustomSprite(sprite))
      ..add(Size2D(43 * tileSize, 3 * tileSize));

    localThoudsandRoad = roadEntity;
  }

   Future<void> _spawnLocalHouse() async {
    this.images.prefix = "assets/";
    final image = await this.images.load('game/png/house.png');
    final sprite = Sprite(image);

    final roadEntity = ecsWorld.create()
      ..add(Position(2 * tileSize, 5 * tileSize))
      ..add(CustomSprite(sprite))
      ..add(Size2D(13 * tileSize, 18 * tileSize));

    localThoudsandRoad = roadEntity;
  }

  Future<void> _spawnLocalSideHouse() async {
    this.images.prefix = "assets/";
    final image = await this.images.load('game/png/side-walk.png');
    final sprite = Sprite(image);

    final roadEntity = ecsWorld.create()
      ..add(Position(0, 20 * tileSize))
      ..add(CustomSprite(sprite))
      ..add(Size2D(43 * tileSize, 3 * tileSize));

    localThoudsandRoad = roadEntity;
  }

  Future<void> _spawnCyclingInRoad() async {
    final artboard = await loadArtboard(RiveFile.asset(
      'assets/game/rive/cycling-in-the-road.riv',
    ));

    final controller = StateMachineController.fromArtboard(
      artboard,
      "State Machine 1",
    );

    artboard.addController(controller!);

    final cyclingEntity = ecsWorld.create()
      ..add(RiveData(artboard: artboard))
      ..add(RiveAnimationData(
          x1: 0, y1: 4 * tileSize, x2: 16 * tileSize, y2: 4 * tileSize))
      ..add(Position(0, 14 * tileSize))
      ..add(Size2D(300, 300));

    localCyclingInRoad = cyclingEntity;
  }

  Future<void> _spawnBlackSky() async {
    final artboard = await loadArtboard(RiveFile.asset(
      'assets/game/rive/black-sky-with-shooting-stars.riv',
    ));

    final controller = SimpleAnimation("Animation 1");
    artboard.addController(controller);

    final blackSkyEntity = ecsWorld.create()
      ..add(RiveData(artboard: artboard))
      ..add(Position(0, 0))
      ..add(comp.Transform(anchor: Anchor.topLeft))
      ..add(Size2D(1113.7 - (1113.7 * 34 / 100), 259 - (259 * 34 / 100)));

    localBlackSkyInRoad = blackSkyEntity;
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

    if (loadedAsset) {
      flameSystem.update(ecsWorld, dt);
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
  void render(Canvas canvas) async {
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

    final Paint paint = Paint()
      ..color = Colors.red
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1;

    const double w = 64 / 1.45;
    const double h = 32 / 1.45;
    const int rows = 3;
    const int cols = 25;

    final double x0 = 18;
    final double y0 = 23 * tileSize;

    for (int r = 0; r < rows; r++) {
      for (int c = 0; c < cols; c++) {
        final double x = x0 + c * w - r * (w / 2);
        final double y = y0 + r * h;

        final Path parallelogram = Path()
          ..moveTo(x, y)
          ..lineTo(x + w, y)
          ..lineTo(x + w - w / 2, y + h)
          ..lineTo(x - w / 2, y + h)
          ..close();

        canvas.drawPath(parallelogram, paint);
      }
    }

    final Paint paintg = Paint()
      ..color = const Color.fromARGB(255, 238, 0, 255)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1;

    const double wg = 64 / 1.45;
    const double hg = 32 / 1.45;
    const int rowsg = 3;
    const int colsg = 10;

    final double x0g = 40;
    final double y0g = 20 * tileSize;

    for (int r = 0; r < rowsg; r++) {
      for (int c = 0; c < colsg; c++) {
        final double x = x0g + c * wg - r * (wg / 2);
        final double y = y0g + r * hg;

        final Path parallelogram = Path()
          ..moveTo(x, y)
          ..lineTo(x + wg, y)
          ..lineTo(x + wg - wg / 2, y + hg)
          ..lineTo(x - wg / 2, y + hg)
          ..close();

        canvas.drawPath(parallelogram, paintg);
      }
    }

    final posCIR = localCyclingInRoad?.get<Position>();
    if (posCIR != null) {
      final dotPaint = Paint()..color = Colors.yellow;
      canvas.drawCircle(Offset(posCIR.x, posCIR.y), 5, dotPaint);

      final textPainter = TextPainter(
        text: TextSpan(
          text:
              "(${posCIR.x.toStringAsFixed(1)}, ${posCIR.y.toStringAsFixed(1)})",
          style: const TextStyle(color: Colors.yellow, fontSize: 14),
        ),
        textDirection: TextDirection.ltr,
      );
      textPainter.layout();
      textPainter.paint(canvas, Offset(posCIR.x, posCIR.y - 20));
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

    // Hiển thị dot localPlayer
    final pos1 = localBlackSkyInRoad?.get<Position>();
    if (pos1 != null) {
      final dotPaint = Paint()..color = Colors.red;
      canvas.drawCircle(Offset(pos1.x, pos1.y), 5, dotPaint);

      final textPainter = TextPainter(
        text: TextSpan(
          text: "(${pos1.x.toStringAsFixed(1)}, ${pos1.y.toStringAsFixed(1)})",
          style: const TextStyle(color: Colors.red, fontSize: 14),
        ),
        textDirection: TextDirection.ltr,
      );
      textPainter.layout();
      textPainter.paint(canvas, Offset(pos1.x, pos1.y - 20));
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

    // Vẽ grid xanh bằng tileSize
    final double startX = 0;
    final double startY = 5 * tileSize;
    final double gridWidth = tileSize * gridCols;
    final double gridHeight = 8 * tileSize;

    final Paint gridPaintBlue = Paint()
      ..color = Colors.blue
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1;

// Vẽ các đường ngang
    for (double y = startY; y <= startY + gridHeight; y += tileSize) {
      canvas.drawLine(
        Offset(startX, y),
        Offset(startX + gridWidth, y),
        gridPaintBlue,
      );
    }

// Vẽ các đường dọc
    for (double x = startX; x <= startX + gridWidth; x += tileSize) {
      canvas.drawLine(
        Offset(x, startY),
        Offset(x, startY + gridHeight),
        gridPaintBlue,
      );
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
