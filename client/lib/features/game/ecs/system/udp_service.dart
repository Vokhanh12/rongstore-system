import 'dart:io';
import 'dart:convert';

class UdpService {
  RawDatagramSocket? _socket;
  InternetAddress serverAddr;
  int serverPort;

  UdpService({required this.serverAddr, required this.serverPort});

  Future<void> connect() async {
    _socket = await RawDatagramSocket.bind(InternetAddress.anyIPv4, 0);
    _socket!.listen((event) {
      if (event == RawSocketEvent.read) {
        final datagram = _socket!.receive();
        if (datagram != null) {
          final msg = utf8.decode(datagram.data);
          print('Received: $msg');
        }
      }
    });
  }

  void send(Map<String, dynamic> msg) {
    print(_socket);

    if (_socket == null) return;
    final data = utf8.encode(json.encode(msg));
    print(data);
    _socket!.send(data, serverAddr, serverPort);
  }

  void close() {
    _socket?.close();
  }
}
