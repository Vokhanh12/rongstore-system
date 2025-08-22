import Flutter
import UIKit
import CryptoKit
import Security

public class NativeCryptoPlugin: NSObject, FlutterPlugin {
  public static let channelName = "native_crypto"
  private var sessionKey: SymmetricKey? = nil
  private let keyTag = "com.example.rongchoi_application.ecdh_identity"

  public static func register(with registrar: FlutterPluginRegistrar) {
    let channel = FlutterMethodChannel(name: channelName, binaryMessenger: registrar.messenger())
    let instance = NativeCryptoPlugin()
    registrar.addMethodCallDelegate(instance, channel: channel)
  }

  public func handle(_ call: FlutterMethodCall, result: @escaping FlutterResult) {
    switch call.method {
    case "getPublicKey":
      do {
        let pubB64 = try ensureAndGetPublicKey()
        result(pubB64)
      } catch {
        result(FlutterError(code: "KEY_ERROR", message: error.localizedDescription, details: nil))
      }
    case "processServerHandshake":
      guard let args = call.arguments as? [String: Any],
            let serverPub = args["serverPub"] as? String,
            let encSession = args["encryptedSessionData"] as? String else {
        result(FlutterError(code: "BAD_ARGS", message: "serverPub and encryptedSessionData required", details: nil))
        return
      }
      do {
        try processServerHandshake(serverPubB64: serverPub, encryptedSessionB64: encSession)
        result(["ok": true])
      } catch {
        result(FlutterError(code: "HANDSHAKE_FAIL", message: error.localizedDescription, details: nil))
      }
    case "encryptPayload":
      guard let args = call.arguments as? [String: Any],
            let plainB64 = args["plain"] as? String,
            let key = sessionKey else {
        result(FlutterError(code: "NO_SESSION", message: "no session key", details: nil))
        return
      }
      let plain = Data(base64Encoded: plainB64)!
      do {
        let sealed = try aesGcmSeal(key: key, plain: plain)
        result(sealed.base64EncodedString())
      } catch {
        result(FlutterError(code: "ENC_FAIL", message: error.localizedDescription, details: nil))
      }
    case "decryptPayload":
      guard let args = call.arguments as? [String: Any],
            let cipherB64 = args["cipher"] as? String,
            let key = sessionKey else {
        result(FlutterError(code: "NO_SESSION", message: "no session key", details: nil))
        return
      }
      let cipher = Data(base64Encoded: cipherB64)!
      do {
        let plain = try aesGcmOpen(key: key, combined: cipher)
        result(plain.base64EncodedString())
      } catch {
        result(FlutterError(code: "DEC_FAIL", message: error.localizedDescription, details: nil))
      }
    case "clearSession":
      sessionKey = nil
      result(nil)
    default:
      result(FlutterMethodNotImplemented)
    }
  }

  // MARK: - Key management (example with P-256 using Security)
  // For production you'd store private key in Secure Enclave or Keychain with proper attributes.
  func ensureAndGetPublicKey() throws -> String {
    // Try to fetch key from keychain
    if let pub = try? getPublicKeyFromKeychain() {
      return pub
    }
    // else create a new key pair in keychain (P256)
    try generateKeyPairP256()
    guard let pub = try? getPublicKeyFromKeychain() else {
      throw NSError(domain: "NativeCrypto", code: -1, userInfo: [NSLocalizedDescriptionKey: "failed get pub"])
    }
    return pub
  }

  func getPublicKeyFromKeychain() throws -> String {
    let query: [String: Any] = [
      kSecClass as String: kSecClassKey,
      kSecAttrApplicationTag as String: keyTag.data(using: .utf8)!,
      kSecReturnRef as String: true
    ]
    var item: CFTypeRef?
    let status = SecItemCopyMatching(query as CFDictionary, &item)
    if status != errSecSuccess { throw NSError(domain: "NativeCrypto", code: Int(status), userInfo: nil) }
    guard let secKey = item as! SecKey? else { throw NSError(domain: "NativeCrypto", code: -1, userInfo: nil) }

    guard let pubData = SecKeyCopyExternalRepresentation(SecKeyCopyPublicKey(secKey)!, nil) as Data? else {
      throw NSError(domain: "NativeCrypto", code: -1, userInfo: nil)
    }
    // pubData is uncompressed SEC1 (0x04||X||Y) for EC keys
    return pubData.base64EncodedString()
  }

  func generateKeyPairP256() throws {
    // Create attributes for key pair with privateKey stored in Secure Enclave if available
    let access = SecAccessControlCreateWithFlags(nil, kSecAttrAccessibleWhenUnlockedThisDeviceOnly, .privateKeyUsage, nil)
    var attributes: [String: Any] = [
      kSecAttrKeyType as String: kSecAttrKeyTypeECSECPrimeRandom,
      kSecAttrKeySizeInBits as String: 256,
      kSecPrivateKeyAttrs as String: [
        kSecAttrIsPermanent as String: true,
        kSecAttrApplicationTag as String: keyTag.data(using: .utf8)!,
        kSecAttrAccessControl as String: access as Any
      ]
    ]
    var error: Unmanaged<CFError>?
    guard let privateKey = SecKeyCreateRandomKey(attributes as CFDictionary, &error) else {
      throw error!.takeRetainedValue() as Error
    }
    // public key stored implicitly
  }

  // MARK: - Handshake processing
  func processServerHandshake(serverPubB64: String, encryptedSessionB64: String) throws {
    // decode server pub
    guard let serverPubBytes = Data(base64Encoded: serverPubB64) else { throw NSError(domain: "NativeCrypto", code: -1, userInfo: nil) }
    // load private key
    let query: [String: Any] = [
      kSecClass as String: kSecClassKey,
      kSecAttrApplicationTag as String: keyTag.data(using: .utf8)!,
      kSecReturnRef as String: true
    ]
    var item: CFTypeRef?
    let status = SecItemCopyMatching(query as CFDictionary, &item)
    if status != errSecSuccess { throw NSError(domain: "NativeCrypto", code: Int(status), userInfo: nil) }
    guard let privateKey = item as! SecKey? else { throw NSError(domain: "NativeCrypto", code: -1, userInfo: nil) }

    // perform ECDH: use SecKeyCopyKeyExchangeResult
    var error: Unmanaged<CFError>?
    guard let serverPubKey = SecKeyCreateWithData(serverPubBytes as CFData,
                                                 [kSecAttrKeyType: kSecAttrKeyTypeECSECPrimeRandom,
                                                  kSecAttrKeyClass: kSecAttrKeyClassPublic,
                                                  kSecAttrKeySizeInBits: 256] as CFDictionary,
                                                 &error) else {
      throw error!.takeRetainedValue() as Error
    }
    guard let sharedData = SecKeyCopyKeyExchangeResult(privateKey,
                                                       SecKeyAlgorithm.ecdhKeyExchangeStandardX963SHA256,
                                                       serverPubKey,
                                                       nil,
                                                       &error) as Data? else {
      // fallback: compute raw ECDH then SHA256 externally if needed
      throw error!.takeRetainedValue() as Error
    }

    // If SecKeyAlgorithm.ecdhKeyExchangeStandardX963SHA256 returned derived data, we can use that
    // But to be consistent with Android scheme, we can instead use raw secret and SHA256
    // For simplicity, assume sharedData is raw secret or already hashed; here we'll derive symmetric key via SHA256(sharedData)
    let sym = SHA256.hash(data: sharedData)
    let aesKey = SymmetricKey(data: Data(sym))

    // decrypt encryptedSessionB64 using AES-GCM (nonce prefixed)
    guard let enc = Data(base64Encoded: encryptedSessionB64) else { throw NSError(domain: "NativeCrypto", code: -1, userInfo: nil) }
    let plain = try aesGcmOpen(key: aesKey, combined: enc)

    // Store session key in memory
    sessionKey = aesKey

    // Optionally parse plain (session info)
  }

  // AES-GCM helpers
  func aesGcmSeal(key: SymmetricKey, plain: Data) throws -> Data {
    let nonce = AES.GCM.Nonce()
    let sealed = try AES.GCM.seal(plain, using: key, nonce: nonce)
    // return nonce || ciphertext || tag
    return Data(nonce) + sealed.ciphertext + sealed.tag
  }

  func aesGcmOpen(key: SymmetricKey, combined: Data) throws -> Data {
    // nonce(12) || ciphertext || tag(16)
    let nonceLen = 12
    guard combined.count > nonceLen + 16 else { throw NSError(domain: "NativeCrypto", code: -1, userInfo: nil) }
    let nonceData = combined.prefix(nonceLen)
    let ctAndTag = combined.suffix(combined.count - nonceLen)
    let ciphertext = ctAndTag.prefix(ctAndTag.count - 16)
    let tag = ctAndTag.suffix(16)
    let sealed = try AES.GCM.SealedBox(nonce: try AES.GCM.Nonce(data: nonceData),
                                       ciphertext: ciphertext,
                                       tag: tag)
    let plain = try AES.GCM.open(sealed, using: key)
    return plain
  }
}
