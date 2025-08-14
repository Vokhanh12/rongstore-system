sealed class AppConfig {
  static const String version = "v1.2412.23";

  static const String username = "FUMA";
  static const String hostname = "https://api-demo.softz.vn";
  static const String timeRange = '1';
}

sealed class AppStrings {
  static const String fontFamily = 'Montserrat';
  static const String loading = "LOADING";
  static const String wait = "Đợi...";
}

sealed class Permission {
  static const String POM_PHIEU_XUATHANG = 'POM_PHIEU_XUATHANG';
  static const String POM_PHIEU_NHAPHANG = 'POM_PHIEU_NHAPHANG';
}

sealed class GoodDocumentType {
  static const String PXKHO = "PXKHO";
  static const String PNKHO = "PNKHO";
}

sealed class NumberFormats {
  static const String positiveNegativeWithDecimal = "#,##0.0;-#,##0.0";
}