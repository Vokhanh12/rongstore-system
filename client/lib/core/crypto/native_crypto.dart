import 'package:flutter/services.dart';

class NativeCrypto {
  static const _ch = MethodChannel('native_crypto');

  static Future<String> getPublicKey() async {
    final v = await _ch.invokeMethod<String>('getPublicKey');
    if (v == null) throw 'null public key';
    return v;
  }

  static Future<String> deriveSharedSecret(String serverPubB64) async {
    final v = await _ch.invokeMethod<String>('deriveSharedSecret', {
      'serverPub': serverPubB64,
    });
    if (v == null) throw 'null shared secret';
    return v;
  }
}
