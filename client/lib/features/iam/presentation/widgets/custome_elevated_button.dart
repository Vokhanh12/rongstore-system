import 'package:flutter/material.dart';
import 'package:rongchoi_application/features/presentation/widgets/custom_text.dart';

class CustomeElevatedButton extends StatelessWidget {
  final String text;
  final TextStyle? style;
  final TextAlign? textAlign;
  final int? maxLines;
  final TextOverflow? overflow;
  final double? width;
  final double? height;

  const CustomeElevatedButton(
      {super.key,
      required this.text,
      this.style,
      this.textAlign,
      this.maxLines,
      this.overflow,
      this.width = double.infinity,
      this.height = 55});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
        width: width,
        height: height,
        child: ElevatedButton(
          onPressed: () {},
          child: CustomText(text: text),
          style: ButtonStyle(
            shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(10.0),
            )),
          ),
        ));
  }
}
