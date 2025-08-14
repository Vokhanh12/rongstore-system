import 'package:flutter/material.dart';
import 'package:rongchoi_application/core/config/app_dimensions.dart';
import 'package:rongchoi_application/core/constants/corlos.dart';
import 'package:rongchoi_application/core/constants/strings.dart';

class AppText {
  // Button
  static TextStyle? btn;

  // Headings (Level 1, 2, 3)
  static TextStyle? h1;
  static TextStyle? h1b;
  static TextStyle? h1d;
  static TextStyle? h1a;
  static TextStyle? h2;
  static TextStyle? h2b;
  static TextStyle? h2d;
  static TextStyle? h2a;
  static TextStyle? h3;
  static TextStyle? h3b;
  static TextStyle? h3d;
  static TextStyle? h3a;

  // Body
  static TextStyle? b1;
  static TextStyle? b1b;
  static TextStyle? b1d;
  static TextStyle? b1a;
  static TextStyle? b2;
  static TextStyle? b2b;
  static TextStyle? b2d;
  static TextStyle? b2a;
  // Level 3 Body (ví dụ)
  static TextStyle? b3;
  static TextStyle? b3b;
  static TextStyle? b3d;
  static TextStyle? b3a;

  // Label
  static TextStyle? l1;
  static TextStyle? l1b;
  static TextStyle? l1d;
  static TextStyle? l1a;
  static TextStyle? l2;
  static TextStyle? l2b;
  static TextStyle? l2d;
  static TextStyle? l2a;
  // Level 3 Label (ví dụ)
  static TextStyle? l3;
  static TextStyle? l3b;
  static TextStyle? l3d;
  static TextStyle? l3a;

  static init() {
    const b = FontWeight.bold;
    const baseStyle = TextStyle(
      fontFamily: AppStrings.fontFamily,
      letterSpacing: 1,
    );

    // Headings
    h1 = baseStyle.copyWith(fontSize: AppDimensions.font(12.0));
    h1b = h1!.copyWith(fontWeight: b);
    h1d = h1!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorDefault);
    h1a = h1!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorAtive);

    h2 = baseStyle.copyWith(
        fontSize: AppDimensions.font(10), fontWeight: FontWeight.w600);
    h2b = h2!.copyWith(fontWeight: b);
    h2d = h2!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorDefault);
    h2a = h2!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorAtive);

    h3 = baseStyle.copyWith(fontSize: AppDimensions.font(8));
    h3b = h3!.copyWith(fontWeight: b);
    h3d = h3!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorDefault);
    h3a = h3!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorAtive);

    // Body Text
    b1 = baseStyle.copyWith(fontSize: AppDimensions.font(7));
    b1b = b1!.copyWith(fontWeight: b);
    b1a = baseStyle.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorDefault);
    b1d = b1!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorAtive);

    b2 = baseStyle.copyWith(fontSize: AppDimensions.font(6.25));
    b2b = b2!.copyWith(fontWeight: b);
    b2a = baseStyle.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorDefault);
    b2d = b2!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorAtive);

    // Level 3 Body Text (ví dụ: nhỏ hơn b2)
    b3 = baseStyle.copyWith(fontSize: AppDimensions.font(5.5));
    b3b = b3!.copyWith(fontWeight: b);
    b3a = baseStyle.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorDefault);
    b3d = b3!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorAtive);

    // Label Text
    l1 = baseStyle.copyWith(fontSize: AppDimensions.font(5));
    l1b = l1!.copyWith(fontWeight: b);
    l1a = baseStyle.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorDefault);
    l1d = l1!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorAtive);

    l2 = baseStyle.copyWith(fontSize: AppDimensions.font(4));
    l2b = l2!.copyWith(fontWeight: b);
    l2a = baseStyle.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorDefault);
    l2d = l2!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorAtive);

    // Level 3 Label Text (ví dụ: nhỏ hơn l2)
    l3 = baseStyle.copyWith(fontSize: AppDimensions.font(3));
    l3b = l3!.copyWith(fontWeight: b);
    l3a = baseStyle.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorDefault);
    l3d = l3!.copyWith(
        fontWeight: FontWeight.w600, color: AppColors.colorAtive);

    // Button Style
    btn = baseStyle.copyWith(
        fontSize: AppDimensions.font(6), fontWeight: FontWeight.w600);
  }
}
