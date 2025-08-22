import Flutter
import UIKit

@main
@objc class AppDelegate: FlutterAppDelegate {
  override func application(
    _ application: UIApplication,
    didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
  ) -> Bool {
    GeneratedPluginRegistrant.register(with: self)

    // Đăng ký plugin thủ công — an toàn khi file plugin thuộc target Runner
    if let registrar = self.registrar(forPlugin: "NativeCryptoPlugin") {
      NativeCryptoPlugin.register(with: registrar)
    } else {
      // optional: log để debug nếu plugin không được compile vào target
      print("⚠️ NativeCryptoPlugin registrar is nil — check target membership")
    }

    return super.application(application, didFinishLaunchingWithOptions: launchOptions)
  }
}
