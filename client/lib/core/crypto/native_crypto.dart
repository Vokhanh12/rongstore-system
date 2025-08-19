import 'dart:convert';
import 'package:flutter/services.dart';

class NativeCrypto {
  static const MethodChannel _ch = MethodChannel('native_crypto');

  /// Trả public key client (base64 of 0x04||X||Y)
  static Future<String> getPublicKey() async {
    final res = await _ch.invokeMethod<String>('getPublicKey');
    if (res == null) throw Exception('getPublicKey returned null');
    return res;
  }

  /// Truyền serverPub (base64) và encryptedSessionData (base64).
  /// Native sẽ:
  ///  - ECDH(private, serverPub)
  ///  - derive AES key (SHA256(shared))
  ///  - decrypt encryptedSessionData (AES-GCM, nonce prefixed)
  ///  - lưu AES key vào native memory cho session
  ///
  /// Trả { "ok": true, "sessionId": "<id>" } hoặc lỗi.
  static Future<Map<String, dynamic>> processServerHandshake({
    required String serverPublicKeyBase64,
    required String encryptedSessionDataBase64,
  }) async {
    final res = await _ch.invokeMethod<Map>('processServerHandshake', {
      'serverPub': serverPublicKeyBase64,
      'encryptedSessionData': encryptedSessionDataBase64,
    });
    if (res == null) throw Exception('processServerHandshake returned null');
    return Map<String, dynamic>.from(res);
  }

  /// Yêu cầu native mã hoá payload bytes (trả base64 ciphertext)
  static Future<String> encryptPayloadBytes(List<int> plainBytes) async {
    final b64 = base64Encode(plainBytes);
    final res = await _ch.invokeMethod<String>('encryptPayload', {'plain': b64});
    if (res == null) throw Exception('encryptPayload returned null');
    return res;
  }

  /// Giải mã payload (nếu cần)
  static Future<List<int>> decryptPayloadBase64(String cipherBase64) async {
    final res = await _ch.invokeMethod<String>('decryptPayload', {'cipher': cipherBase64});
    if (res == null) throw Exception('decryptPayload returned null');
    return base64Decode(res);
  }

  /// Xóa session key ở native
  static Future<void> clearSession() async {
    await _ch.invokeMethod('clearSession');
  }
}
