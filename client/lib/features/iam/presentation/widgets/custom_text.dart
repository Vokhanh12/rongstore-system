

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:rongchoi_application/features/presentation/bloc/tranlation_bloc/tranlation_bloc.dart';
import 'package:rongchoi_application/features/presentation/utils/tranlation_util.dart';

class CustomText extends StatelessWidget {
  final String text;
  final TextStyle? style;
  final TextAlign? textAlign;
  final int? maxLines;
  final TextOverflow? overflow;

  const CustomText({
    super.key,
    required this.text,
    this.style,
    this.textAlign,
    this.maxLines,
    this.overflow,
  });

  @override
  Widget build(BuildContext context) {
    return BlocSelector<TranlationBloc, TranlationState, String>(
      selector: (state) {
        if (state is GetAllTranlationsLocalState) {
          return TranlationUtil.getTranlationsByCode(state.tranlationItems, text);
        }
        return text; // Trả về text gốc nếu chưa có dữ liệu dịch
      },
      builder: (context, translatedText) {
        return Text(
          translatedText,
          style: style,
          textAlign: textAlign,
          maxLines: maxLines,
          overflow: overflow,
        );
      },
    );
  }
}
