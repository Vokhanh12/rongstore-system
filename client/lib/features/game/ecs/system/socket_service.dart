import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';

class SocketService {
  WebSocketChannel? _channel;
  final _controller = StreamController<Map<String, dynamic>>.broadcast();

  Stream<Map<String, dynamic>> get messages => _controller.stream;

  Future<void> connect(String url) async {
    _channel = WebSocketChannel.connect(Uri.parse(url));
    _channel!.stream.listen((data) {
      try {
        final jsonMap = json.decode(data) as Map<String, dynamic>;
        _controller.add(jsonMap);
      } catch (_) {}
    }, onError: (e) {
    }, onDone: () {
    });
  }

  void send(Map<String, dynamic> message) {
    final ch = _channel;
    if (ch == null) return;
    ch.sink.add(json.encode(message));
  }

  void close() {
    _channel?.sink.close();
    _channel = null;
  }
}
