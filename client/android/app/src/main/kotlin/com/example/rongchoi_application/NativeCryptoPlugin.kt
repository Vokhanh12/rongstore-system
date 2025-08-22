package com.example.rongchoi_application

import android.content.Context
import android.os.Build
import android.util.Base64
import io.flutter.embedding.engine.plugins.FlutterPlugin
import io.flutter.plugin.common.MethodCall
import io.flutter.plugin.common.MethodChannel
import io.flutter.plugin.common.BinaryMessenger
import java.math.BigInteger
import java.security.*
import java.security.interfaces.ECPublicKey
import java.security.spec.ECGenParameterSpec
import java.security.spec.ECPublicKeySpec
import java.security.KeyFactory
import javax.crypto.Cipher
import javax.crypto.KeyAgreement
import javax.crypto.spec.GCMParameterSpec
import javax.crypto.spec.SecretKeySpec
import java.security.KeyStore
import java.security.spec.ECParameterSpec
import kotlin.experimental.and
import java.security.spec.ECPoint
import java.security.spec.X509EncodedKeySpec
import java.security.MessageDigest

class NativeCryptoPlugin : FlutterPlugin, MethodChannel.MethodCallHandler {
    companion object {
        private const val CHANNEL = "native_crypto"
        private const val KEY_ALIAS = "ecdh_identity"
        private const val AES_KEY_LEN = 32 // bytes (256 bits)
        private const val GCM_TAG_LEN = 128 // bits
        private const val GCM_NONCE_LEN = 12 // bytes

        /** Manual register function for embedding v2 (call from MainActivity if you want). */
        @JvmStatic
        fun register(messenger: BinaryMessenger, context: Context) {
            val plugin = NativeCryptoPlugin()
            plugin.appContext = context
            plugin.channel = MethodChannel(messenger, CHANNEL).also {
                it.setMethodCallHandler(plugin)
            }
        }
    }

    private var channel: MethodChannel? = null
    internal lateinit var appContext: Context

    // in-memory session key -- keep ephemeral
    @Volatile
    private var sessionAesKey: ByteArray? = null

    override fun onAttachedToEngine(binding: FlutterPlugin.FlutterPluginBinding) {
        appContext = binding.applicationContext
        channel = MethodChannel(binding.binaryMessenger, CHANNEL).also {
            it.setMethodCallHandler(this)
        }
    }

    override fun onDetachedFromEngine(binding: FlutterPlugin.FlutterPluginBinding) {
        channel?.setMethodCallHandler(null)
        channel = null
        clearSession()
    }

    override fun onMethodCall(call: MethodCall, result: MethodChannel.Result) {
        when (call.method) {
            "getPublicKey" -> {
                try {
                    val pub = ensureAndGetPublicKey()
                    val uncompressed = ecPublicKeyToUncompressedPoint(pub)
                    result.success(Base64.encodeToString(uncompressed, Base64.NO_WRAP))
                } catch (e: Exception) {
                    result.error("KEY_ERROR", e.message, null)
                }
            }

            "processServerHandshake" -> {
                val serverPub = call.argument<String>("serverPub")
                val encSessionB64 = call.argument<String>("encryptedSessionData")
                if (serverPub.isNullOrBlank() || encSessionB64.isNullOrBlank()) {
                    result.error("BAD_ARGS", "serverPub and encryptedSessionData are required", null)
                    return
                }
                try {
                    val shared = deriveSharedSecretFromServerPub(serverPub)
                    val aesKey = sha256(shared) // 32 bytes
                    val sessionPlain = aesGcmOpen(aesKey, Base64.decode(encSessionB64, Base64.NO_WRAP))
                    sessionAesKey = aesKey
                    result.success(mapOf("ok" to true))
                } catch (e: Exception) {
                    result.error("HANDSHAKE_FAIL", e.message, null)
                }
            }

            "encryptPayload" -> {
                val plainB64 = call.argument<String>("plain")
                if (plainB64.isNullOrBlank()) {
                    result.error("BAD_ARGS", "plain required", null)
                    return
                }
                val key = sessionAesKey
                if (key == null) {
                    result.error("NO_SESSION", "no session key available; run handshake", null)
                    return
                }
                try {
                    val cipher = aesGcmSeal(key, Base64.decode(plainB64, Base64.NO_WRAP))
                    result.success(Base64.encodeToString(cipher, Base64.NO_WRAP))
                } catch (e: Exception) {
                    result.error("ENC_FAIL", e.message, null)
                }
            }

            "decryptPayload" -> {
                val cipherB64 = call.argument<String>("cipher")
                if (cipherB64.isNullOrBlank()) {
                    result.error("BAD_ARGS", "cipher required", null)
                    return
                }
                val key = sessionAesKey
                if (key == null) {
                    result.error("NO_SESSION", "no session key available; run handshake", null)
                    return
                }
                try {
                    val plain = aesGcmOpen(key, Base64.decode(cipherB64, Base64.NO_WRAP))
                    result.success(Base64.encodeToString(plain, Base64.NO_WRAP))
                } catch (e: Exception) {
                    result.error("DEC_FAIL", e.message, null)
                }
            }

            "clearSession" -> {
                clearSession()
                result.success(null)
            }

            else -> result.notImplemented()
        }
    }

    // ------------------- Key generation / retrieval -------------------

    @Throws(Exception::class)
    private fun ensureAndGetPublicKey(): PublicKey {
        val ks = KeyStore.getInstance("AndroidKeyStore")
        ks.load(null)

        if (ks.containsAlias(KEY_ALIAS)) {
            val entry = ks.getEntry(KEY_ALIAS, null) as KeyStore.PrivateKeyEntry
            return entry.certificate.publicKey
        }

        // create keypair
        createEcKeyPairInKeystore("secp521r1")
        val entry = ks.getEntry(KEY_ALIAS, null) as KeyStore.PrivateKeyEntry
        return entry.certificate.publicKey
    }

    @Throws(Exception::class)
    private fun createEcKeyPairInKeystore(curveName: String) {
        val kpg = KeyPairGenerator.getInstance("EC", "AndroidKeyStore")
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
            val gSpec = ECGenParameterSpec(curveName)
            val keySpec = android.security.keystore.KeyGenParameterSpec.Builder(
                KEY_ALIAS,
                (android.security.keystore.KeyProperties.PURPOSE_AGREE_KEY or android.security.keystore.KeyProperties.PURPOSE_SIGN)
            )
                .setAlgorithmParameterSpec(gSpec)
                .setDigests(android.security.keystore.KeyProperties.DIGEST_SHA256)
                .setUserAuthenticationRequired(false)
                .build()
            kpg.initialize(keySpec)
            kpg.generateKeyPair()
        } else {
            throw RuntimeException("Android < M not supported for keystore EC generation")
        }
    }

    // Convert EC PublicKey to 0x04||X||Y
    @Throws(Exception::class)
    private fun ecPublicKeyToUncompressedPoint(pub: PublicKey): ByteArray {
        val kf = KeyFactory.getInstance("EC")
        val pk: ECPublicKey = when (pub) {
            is ECPublicKey -> pub
            else -> {
                val kpub = kf.generatePublic(X509EncodedKeySpec(pub.encoded)) as ECPublicKey
                kpub
            }
        }
        val point = pk.w
        val fieldSize = ((pk.params.curve.field.fieldSize + 7) / 8)
        val xb = toFixedLength(point.affineX.toByteArray(), fieldSize)
        val yb = toFixedLength(point.affineY.toByteArray(), fieldSize)
        return byteArrayOf(0x04) + xb + yb
    }

    private fun toFixedLength(input: ByteArray, size: Int): ByteArray {
        val bi = if (input[0] == 0.toByte()) input.copyOfRange(1, input.size) else input
        if (bi.size == size) return bi
        if (bi.size > size) return bi.copyOfRange(bi.size - size, bi.size)
        val out = ByteArray(size)
        System.arraycopy(bi, 0, out, size - bi.size, bi.size)
        return out
    }

    // ------------------- ECDH / Derive / AES-GCM -------------------
    @Throws(Exception::class)
    private fun deriveSharedSecretFromServerPub(serverPubB64: String): ByteArray {
        val serverPubBytes = Base64.decode(serverPubB64, Base64.NO_WRAP)
        if (serverPubBytes.isEmpty() || serverPubBytes[0] != 0x04.toByte()) {
            throw IllegalArgumentException("server public key must be uncompressed point (0x04||X||Y)")
        }

        val ks = KeyStore.getInstance("AndroidKeyStore")
        ks.load(null)
        val privEntry = ks.getEntry(KEY_ALIAS, null) as? KeyStore.PrivateKeyEntry
            ?: throw IllegalStateException("Private key $KEY_ALIAS not found in AndroidKeyStore")
        val privKey = privEntry.privateKey

        // Lấy params từ certificate public key (an toàn)
        val certPub = privEntry.certificate.publicKey
        val ecPubFromCert = certPub as? ECPublicKey
            ?: throw IllegalStateException("Stored certificate public key is not EC")
        val ecParams: ECParameterSpec = ecPubFromCert.params

        // kích thước field (bytes)
        val fieldSize = ((ecParams.curve.field.fieldSize + 7) / 8)

        val expectedLen = 1 + 2 * fieldSize
        if (serverPubBytes.size != expectedLen) {
            throw IllegalArgumentException("unexpected server public key length ${serverPubBytes.size}, expected $expectedLen")
        }

        val x = serverPubBytes.copyOfRange(1, 1 + fieldSize)
        val y = serverPubBytes.copyOfRange(1 + fieldSize, 1 + 2 * fieldSize)
        val bx = BigInteger(1, x)
        val by = BigInteger(1, y)
        val w = ECPoint(bx, by)
        val pubSpec = ECPublicKeySpec(w, ecParams)

        val kf = KeyFactory.getInstance("EC")
        val serverPubKey = kf.generatePublic(pubSpec)

        val ka = KeyAgreement.getInstance("ECDH")
        ka.init(privKey) // dùng PrivateKey wrapper trực tiếp
        ka.doPhase(serverPubKey, true)
        return ka.generateSecret()
    }


    private fun sha256(input: ByteArray): ByteArray {
        val md = MessageDigest.getInstance("SHA-256")
        return md.digest(input)
    }

    private fun aesGcmOpen(aesKey: ByteArray, combined: ByteArray): ByteArray {
        if (combined.size < GCM_NONCE_LEN + 1) throw IllegalArgumentException("ciphertext too short")
        val nonce = combined.copyOfRange(0, GCM_NONCE_LEN)
        val cipherBytes = combined.copyOfRange(GCM_NONCE_LEN, combined.size)
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        val spec = GCMParameterSpec(GCM_TAG_LEN, nonce)
        val key = SecretKeySpec(aesKey, "AES")
        cipher.init(Cipher.DECRYPT_MODE, key, spec)
        return cipher.doFinal(cipherBytes)
    }

    private fun aesGcmSeal(aesKey: ByteArray, plain: ByteArray): ByteArray {
        val nonce = ByteArray(GCM_NONCE_LEN)
        SecureRandom().nextBytes(nonce)
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        val spec = GCMParameterSpec(GCM_TAG_LEN, nonce)
        val key = SecretKeySpec(aesKey, "AES")
        cipher.init(Cipher.ENCRYPT_MODE, key, spec)
        val ct = cipher.doFinal(plain) // includes tag
        return nonce + ct
    }

    private fun clearSession() {
        sessionAesKey?.fill(0)
        sessionAesKey = null
    }
}
