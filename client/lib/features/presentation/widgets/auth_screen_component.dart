
import 'package:flutter/material.dart';
import 'package:rongchoi_application/core/config/app_dimensions.dart';
import 'package:rongchoi_application/core/config/app_typography.dart';
import 'package:rongchoi_application/core/config/space.dart';
import 'package:rongchoi_application/core/constants/assets.dart';
import 'package:rongchoi_application/core/constants/corlos.dart';

Widget authTopColumn(bool isFromSignUp) {
  return Column(
    children: [
      Space.yf(1.5),
      const Padding(
        padding: EdgeInsets.only(left: 60, right: 60),
        child: Image(image: AssetImage(AppAssets.logoSplashScreen)),
      ),
      Space.yf(1.5),
    ],
  );
}

Widget authBottomButton(bool isFromSignUp, void Function()? onTap) {
  return GestureDetector(
    onTap: onTap,
    child: Container(
      height: AppDimensions.normalize(35),
      decoration:
          BoxDecoration(color: AppColors.antiqueRuby, borderRadius: BorderRadius.only(topLeft: Radius.circular(AppDimensions.normalize(7.5)), topRight: Radius.circular(AppDimensions.normalize(7.5)))),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            isFromSignUp ? "Already have an account?" : "Don’t have an account?",
            style: AppText.b2?.copyWith(color: Colors.white),
          ),
          Space.yf(.5),
          Text(
            isFromSignUp ? "Login".toUpperCase() : "SIGN UP",
            style: AppText.h3b?.copyWith(color: Colors.white),
          ),
        ],
      ),
    ),
  );
}