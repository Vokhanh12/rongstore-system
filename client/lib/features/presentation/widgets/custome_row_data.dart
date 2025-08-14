import 'package:flutter/material.dart';


class CustomeRowData extends StatelessWidget {


  const CustomeRowData({super.key, required this.children});

  final List<Widget> children;

  

  @override
  Widget build(BuildContext context) {


    return Row(
      children: children
    );
  }

}


// for (var row in controllers)
//           for (var entry in row.entries)
//             Padding(
//               padding: const EdgeInsets.symmetric(vertical: 8.0),
//               child:  CustomTextFormField(
//                       label: "RC." + entry.key,
//                       controller: entry.value,
//                       validator: validators.validatePassword),
//             ),
//         const SizedBox(width: 10),