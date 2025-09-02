import Flutter
import UIKit

@main
@objc class AppDelegate: FlutterAppDelegate {
  override func application(
    _ application: UIApplication,
    didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
  ) -> Bool {
    GeneratedPluginRegistrant.register(with: self)

    if let registrar = self.registrar(forPlugin: "NativeCryptoPlugin") {
      NativeCryptoPlugin.register(with: registrar)
    } else {
      print("⚠️ NativeCryptoPlugin registrar is nil — check target membership")
    }

    return super.application(application, didFinishLaunchingWithOptions: launchOptions)
  }
}
