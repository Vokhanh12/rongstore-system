import 'package:flutter/material.dart';

sealed class AppColors {
  /// Background color for scaffold
  static const Color scaffoldBackground = Color(0xffeeeeee);
  static const Color scaffoldBackground1 = Color(0xFFF2F6FC);

  /// A deep reddish color
  static const Color antiqueRuby = Color(0xffb93273);

  /// A deep teal color
  static const Color deepTeal = Color(0xff2e7b79);

  /// Color for grey text
  static const Color greyText = Color(0xffA8A8A8);

  /// Another grey color
  static const Color anotherGrey = Color(0xff666666);

  /// A light grey color
  static const Color lightGrey = Color(0xffC7C7C7);

  /// Color for tabs
  static const Color tabColor = Color(0xffF5732A);

  /// Common color used in the app
  static const Color commonBlue = Color(0xFF0050B4);
  static const Color commonBlue1 = Color(0xFF133776);

  static const Color commonLightBlue = Colors.lightBlue;
  static const Color commonWhite = Color.fromARGB(255, 255, 255, 255);
  static const Color commonRed = Color(0xFFF44245);
  static const Color commonGrey = Color(0xFFF1F4F6);
  static const Color commonBorderCard = Color(0xFFEFF2F4);

  static const Color commonBorderUnderListView = Color(0xFFEFF2F4);

  static const Color commonButtonColorDefault = Color(0xFFF2F8FD);
  static const Color commonButtonColorAtive = Color(0xFF0050B4);
  static const Color commonTextColor = Color(0xFFA8B4C6);
  static const Color commonTextColor1 = Color(0xFF9FA7B3);
  static const Color commonAppbarColor = Color(0xFFF2F8FD);
  static const Color commonHeaderColor = Color(0xFFF2F8FD);

  /// Color for success messages in snack bars
  static const Color successColor = Color(0xFF0D7BD4);

  /// Color for error messages in snack bars
  static const Color errorColor = Color(0xFFF44245);

  // Color for click change
  static const Color colorAtive = Color(0xFF0D7BD4);
  static const Color colorHover = Color(0xFFF5F8FF);
  static const Color colorDefault = Color(0xffA8A8A8);

  /// Color for unselected buttons
  static const Color unselectedButtonColor = Color(0xffEDEDED);

  /// A transparent color with 24% opacity
  static const Color transparentColor = Color(0xff484848); //with opacity 24%

  /////////////////////////////////////////////////////////////////
  // Color for custome text form field componet
  static const Color iceBlueBackgroupTf = Color(0xFFF0F5FA);
  static const Color blueHazeTextTf = Color(0xFFB4B9CA);
  // Color for custome text check box componet
  static const Color orangeShadeCkb = Color(0xFFFF7622);

  static const Color aliceBlueColor = Color(0xFFF0F5FA);

  static const Color _PRIMARY_COLOR = Color(0xFF454750);
  static const Color _ALICEBLUE_COLOR = Color(0xFFF0F5FA);
  static const Color _WHITE_COLOR = Color.fromARGB(255, 255, 255, 255);
  static const Color _ORANGE_SHADE = Color(0xFFFF7622);

  /// THEME APP ////////////////////////////////////////////////////////////////////////////////////////////////
  static Color MA_SFBACKGROUP_COLOR = _WHITE_COLOR;

  /// Custom Text Field ////////////////////////////////////////////////////////////////////////////////////////////////
  static Color TF_TEXT_COLOR = _PRIMARY_COLOR;
  static Color TF_CURSOR_COLOR = _PRIMARY_COLOR;
  static Color TF_BOXDECORATION_COLOR = _ALICEBLUE_COLOR;

  /// Custom Check Box ////////////////////////////////////////////////////////////////////////////////////////////////
}
