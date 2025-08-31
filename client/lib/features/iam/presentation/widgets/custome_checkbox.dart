import 'package:flutter/material.dart';
import 'package:rongchoi_application/core/constants/corlos.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/custom_text.dart';
class CustomeCheckbox extends StatelessWidget {
  const CustomeCheckbox(
      {super.key,
      required this.onChanged,
      required this.value,
      required this.text});

  final Function(bool?)? onChanged;
  final bool value;
  final text;

  @override
  Widget build(BuildContext context) {
    Color getColor(Set<WidgetState> states) {
      const Set<WidgetState> interactiveStates = <WidgetState>{
        WidgetState.pressed,
        WidgetState.hovered,
        WidgetState.focused,
      };
      if (states.any(interactiveStates.contains)) {
        return Colors.blue;
      }
      return Colors.white;
    }

    return Row(
      children: [
        SizedBox(
          width: 20,
          height: 20,
          child: Checkbox(
              checkColor: Colors.white,
              fillColor: WidgetStateProperty.resolveWith(getColor),
              value: value,
              side: BorderSide(color: Colors.blueGrey.shade100, width: 2),
              shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(4)),
              onChanged: onChanged),
        ),
        SizedBox(
          width: 5,
        ),
        CustomText(text: text)
      ],
    );
  }
}
