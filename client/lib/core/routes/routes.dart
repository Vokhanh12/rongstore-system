import 'package:flutter/material.dart';
import 'package:rongchoi_application/core/error/exception.dart';
import 'package:rongchoi_application/features/presentation/screen/login.dart';
import 'package:rongchoi_application/features/presentation/screen/splash.dart';


sealed class AppRouter {
  static const String splash = '/';
  static const String login = '/login';
  static const String root = '/root';

  static Route<dynamic> onGenerateRoute(RouteSettings routeSettings) {
    switch (routeSettings.name) {
      case splash:
        return MaterialPageRoute(builder: (_) => const SplashScreen());
      case login:
        return MaterialPageRoute(builder: (_) => const LoginScreen());


      default:
        throw const RouteException('Route not found!');
    }
  }
}