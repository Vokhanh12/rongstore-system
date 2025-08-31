import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:rongchoi_application/core/config/space.dart';
import 'package:rongchoi_application/core/constants/assets.dart';
import 'package:rongchoi_application/core/validator/validator.dart';
import 'package:extension_type_unions/extension_type_unions.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/auto_form.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/custom_text.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/custom_textformfield.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/custome_checkbox.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/custome_column_data.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/custome_elevated_button.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/custome_row_data.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();

  final _formKey = GlobalKey<FormState>();

  final Validators _validators = Validators();

  late FocusNode _usernameFocusNode;
  late FocusNode _passwordForcusNode;

  @override
  void initState() {
    // TODO: implement initState
    _usernameFocusNode = FocusNode();
    _passwordForcusNode = FocusNode();

    super.initState();
  }

  @override
  void dispose() {
    // TODO: implement dispose
    super.dispose();

    _usernameController.dispose();
    _passwordController.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SizedBox(
        height: double.infinity,
        child: Stack(
          children: [
            decorLeft01,
            decorRight02,
            decorRight03,
            //decorBottomLeft04,
            Padding(
              padding: Space.hf(1.3),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  AutoForm(
                      child:
                          Union3<CustomeColumnData, CustomeRowData, Widget>.in1(
                              CustomeColumnData(
                    children: [
                      CustomTextFormField(label: 'RC.Username'),
                      SizedBox(
                        height: 20,
                      ),
                      CustomTextFormField(label: 'RC.Password'),
                    ],
                  ))),
                  SizedBox(
                    height: 20,
                  ),
                  Row(
                    children: [
                      CustomeCheckbox(
                          text: "RC.RememberMe",
                          onChanged: (value) {},
                          value: false),
                      Spacer(),
                      CustomText(
                          text: 'RC.ForgotPassword',
                          textAlign: TextAlign.right),
                    ],
                  ),
                  SizedBox(
                    height: 20,
                  ),
                  CustomeElevatedButton(text: 'RC.Login'),
                  SizedBox(
                    height: 30,
                  ),
                  signUpPrompt,
                  SizedBox(
                    height: 20,
                  ),
                  CustomText(text: 'RC.Or'),
                ],
              ),

              // child: Form(
              //   key: _formKey,
              //   child: Column(
              //     children: [
              //       CustomTextFormField(
              //         label: "RC.Username",
              //         controller: _usernameController,
              //       ),
              //       Space.yf(1.3),
              //       CustomTextFormField(
              //           label: "RC.Password",
              //           controller: _passwordController,
              //           validator: _validators.validatePassword),
              //       Space.yf(.3),
              //       Row(
              //         children: [
              //           Row(
              //             children: [
              //               CustomeCheckbox(onChanged: (isChecked){}, value: true),
              //               CustomText(text: "RC.Remember")
              //             ],
              //           ),
              //           Spacer(),
              //           CustomText(
              //             text: "RC.ForgotPassword",
              //             style: AppText.b2,
              //           ),
              //         ],
              //       ),
              //       Space.yf(2.5),
              //       // BlocConsumer<SignInBloc, SignInState>(
              //       //   listener: (context, state) {
              //       //     if (state.status == SignInStatus.error) {
              //       //       showErrorAuthBottomSheet(context);
              //       //     }
              //       //     if (state.status == SignInStatus.success) {
              //       //       showSuccessfulAuthBottomSheet(context, false);
              //       //     }
              //       //   },
              //       //   builder: (context, state) {
              //       //     return customElevatedButton(
              //       //       onTap: () {
              //       //         if (_formKey.currentState!.validate()) {
              //       //           context.read<SignInBloc>().add(
              //       //                 SignInWithCredential(
              //       //                   email: _emailController.text.trim(),
              //       //                   password: _passwordController.text.trim(),
              //       //                 ),
              //       //               );
              //       //         }
              //       //       },
              //       //       text: (state.status == SignInStatus.submitting)
              //       //           ? AppStrings.wait
              //       //           : "Login".toUpperCase(),
              //       //       heightFraction: 20,
              //       //       width: double.infinity,
              //       //       color: AppColors.commonAmber,
              //       //     );
              //       //   },
              //       // ),
              //       Space.yf(2.5),
              //       CustomeElevatedButton(text: "RC.Login")
              //     ],
              //   ),
              // ),
            ),
          ],
        ),
      ),
    );
  }

  Widget get background => Positioned(
      top: 0.0,
      left: 0.0,
      right: 0.0,
      child: Container(
        width: MediaQuery.of(context).size.width,
        height: MediaQuery.of(context).size.height,
        color: Colors.white,
      ));

  // Decor 01
  Widget get decorLeft01 => Positioned(
        child: Align(
          alignment: Alignment.topLeft,
          child: SvgPicture.asset(
              width: MediaQuery.of(context).size.width * 0.2,
              height: MediaQuery.of(context).size.height * 0.2,
              AppAssets.loginDecore01),
        ),
      );

  // Decor 02
  Widget get decorRight02 => Positioned(
        top: MediaQuery.of(context).size.height / 7,
        right: 0,
        child: SvgPicture.asset(
            width: MediaQuery.of(context).size.width * 0.2,
            height: MediaQuery.of(context).size.height * 0.2,
            AppAssets.loginDecore02),
      );

  // Decor 03
  Widget get decorRight03 => Align(
        alignment: Alignment.topRight,
        child: SvgPicture.asset(
            width: MediaQuery.of(context).size.width * 0.2,
            height: MediaQuery.of(context).size.height * 0.2,
            AppAssets.loginDecore03),
      );

  // Decor 04
  Widget get decorBottomLeft04 => Positioned(
        bottom: 0,
        left: 0,
        child: SvgPicture.asset(
          width: MediaQuery.of(context).size.width * 0.2,
          height: MediaQuery.of(context).size.height * 0.2,
          AppAssets.loginDecore04,
        ),
      );

  Widget get signUpPrompt => SizedBox(
        width: MediaQuery.of(context).size.width,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CustomText(text: 'RC.HaveAccount', textAlign: TextAlign.right),
            SizedBox(
              width: 7,
            ),
            CustomText(text: 'RC.SignUpNow', textAlign: TextAlign.right),
            //tranlate
          ],
        ),
      );

  Widget get tranlate =>
      SvgPicture.asset(width: 30, height: 30, AppAssets.flagVi);

  Widget get tranlateFlag => SizedBox(
        width: MediaQuery.of(context).size.width,
        child: SvgPicture.asset(
          width: MediaQuery.of(context).size.width * 0.2,
          height: MediaQuery.of(context).size.height * 0.2,
          AppAssets.loginDecore04,
        ),
      );
}
