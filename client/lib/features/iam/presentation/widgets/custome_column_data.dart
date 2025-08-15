import 'package:flutter/material.dart';


class CustomeColumnData extends StatelessWidget {


  const CustomeColumnData({super.key, required this.children});

  final List<Widget> children;

  

  @override
  Widget build(BuildContext context) {


    return Column(
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