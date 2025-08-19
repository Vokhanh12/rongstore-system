import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_svg/svg.dart';
import 'package:rongchoi_application/core/config/app_dimensions.dart';
import 'package:rongchoi_application/core/config/app_typography.dart';
import 'package:rongchoi_application/core/constants/corlos.dart';
import 'package:rongchoi_application/features/iam/presentation/bloc/tranlation_bloc/tranlation_bloc.dart';
import 'package:rongchoi_application/features/iam/presentation/utils/tranlation_util.dart';

class CustomTextFormField extends StatefulWidget {
  final String label;
  final String? title;
  final String? svgUrl;
  final String? Function(String?)? validator;

  const CustomTextFormField({
    Key? key,
    required this.label,
    this.title,
    this.svgUrl,
    this.validator,
  }) : super(key: key);

  @override
  _CustomTextFormFieldState createState() => _CustomTextFormFieldState();
}

class _CustomTextFormFieldState extends State<CustomTextFormField> {
  @override
  void initState() {
    super.initState();
    context.read<TranlationBloc>().add(GetAllTranlationsLocalEvent());
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<TranlationBloc, TranlationState>(
      builder: (context, state) {
        if (state is LoadingTranlationState) {
          return Center(child: CircularProgressIndicator());
        } else if (state is GetAllTranlationsLocalState) {
          return Column(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if(widget.title != null)
              Padding(
                padding: const EdgeInsets.only(left: 8, bottom: 5),
                child: Text(
                  widget.title!,
                  style: AppText.b1?.copyWith(
                      color: AppColors.TF_TEXT_COLOR,
                      fontWeight: FontWeight.w400),
                ),
              ),
              Container(
                  alignment: Alignment.center,
                  padding: const EdgeInsets.fromLTRB(30, 3, 20, 0),
                  height: 60,
                  decoration: BoxDecoration(
                    color: AppColors.TF_BOXDECORATION_COLOR,
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: TextFormField(
                    cursorColor: AppColors.TF_CURSOR_COLOR,
                    autovalidateMode: AutovalidateMode.onUserInteraction,
                    validator: widget.validator,
                    controller: TextEditingController(),
                    style: AppText.b2,
                    decoration: InputDecoration(
                      prefixIcon: widget.svgUrl == null
                          ? null
                          : Padding(
                              padding: EdgeInsets.only(
                                right: AppDimensions.normalize(10),
                                top: AppDimensions.normalize(1),
                              ),
                              child: SvgPicture.asset(
                                widget.svgUrl!,
                                colorFilter: const ColorFilter.mode(
                                  AppColors.deepTeal,
                                  BlendMode.srcIn,
                                ),
                              ),
                            ),
                      border: InputBorder.none,
                      focusedBorder: InputBorder.none,
                      enabledBorder: InputBorder.none,
                      errorBorder: InputBorder.none,
                      disabledBorder: InputBorder.none,
                      errorStyle: AppText.l1b?.copyWith(color: Colors.red),
                      errorMaxLines: 3,
                      hintText: TranlationUtil.getTranlationsByCode(
                          state.tranlationItems, widget.label),
                      labelStyle:
                          AppText.b1?.copyWith(color: AppColors.TF_TEXT_COLOR),
                    ),
                  )),
            ],
          );
        }

        return Container();
      },
    );
  }


}
