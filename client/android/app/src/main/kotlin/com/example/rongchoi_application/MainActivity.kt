package com.example.rongchoi_application

import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine

class MainActivity : FlutterActivity() {
    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)
        // Manual register to guarantee the channel is available early
        NativeCryptoPlugin.register(flutterEngine.dartExecutor.binaryMessenger, applicationContext)
    }
}
