class ServerException implements Exception {}

class DatabaseException implements Exception {}

class SharedPreferencesException implements Exception {}

class ConnectionException implements Exception {}

class CacheException implements Exception {}

//route
class RouteException implements Exception {
  final String message;
  const RouteException(this.message);
}