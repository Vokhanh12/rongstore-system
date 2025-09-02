import 'package:flame/game.dart' hide Route; 
import 'package:flutter/material.dart'; 
import 'package:rongchoi_application/core/error/exception.dart';
import 'package:rongchoi_application/features/game/mygame.dart';
import 'package:rongchoi_application/features/iam/presentation/screen/login.dart';
import 'package:rongchoi_application/features/iam/presentation/screen/splash.dart';

sealed class AppRouter {
  static const String splash = '/';
  static const String login = '/login';
  static const String root = '/root';
  static const String game = '/game';

  static Route<dynamic> onGenerateRoute(RouteSettings routeSettings) {
    switch (routeSettings.name) {
      case splash:
        return MaterialPageRoute(builder: (_) => const SplashScreen());

      case login:
        return MaterialPageRoute(builder: (_) => const LoginScreen());

      case game:
        return MaterialPageRoute(
          builder: (_) => GameWidget(
            game: MyGame(), 
          ),
        );

      default:
        throw const RouteException('Route not found!');
    }
  }
}
