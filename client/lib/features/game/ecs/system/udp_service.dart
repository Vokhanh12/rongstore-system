import 'dart:io';
import 'dart:convert';
import 'dart:async';

class UdpService {
  RawDatagramSocket? _socket;
  InternetAddress serverAddr;
  int serverPort;

  final StreamController<Map<String, dynamic>> _messagesController =
      StreamController.broadcast();

  Stream<Map<String, dynamic>> get messages => _messagesController.stream;

  UdpService({required this.serverAddr, required this.serverPort});

  Future<void> connect() async {
    _socket = await RawDatagramSocket.bind(InternetAddress.anyIPv4, 0);
    _socket!.listen((event) {
      if (event == RawSocketEvent.read) {
        final datagram = _socket!.receive();
        if (datagram != null) {
          final msg = utf8.decode(datagram.data);
          try {
            final decoded = json.decode(msg) as Map<String, dynamic>;
            _messagesController.add(decoded); // phát ra stream
          } catch (e) {
            print('UDP decode error: $e');
          }
        }
      }
    });
  }

  void send(Map<String, dynamic> msg) {
    if (_socket == null) return;
    final data = utf8.encode(json.encode(msg));
    _socket!.send(data, serverAddr, serverPort);
  }

  void close() {
    _socket?.close();
    _messagesController.close();
  }
}
