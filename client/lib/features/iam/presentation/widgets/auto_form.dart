import 'package:extension_type_unions/extension_type_unions.dart';
import 'package:flutter/material.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/custome_column_data.dart';
import 'package:rongchoi_application/features/iam/presentation/widgets/custome_row_data.dart';

class AutoForm extends StatefulWidget {
  const AutoForm({
    super.key,
    required this.child,
  });

  final Union3<CustomeColumnData, CustomeRowData, Widget> child;
  @override
  State<AutoForm> createState() => _AutoFormState();
}

class _AutoFormState extends State<AutoForm> {
  final _formKey = GlobalKey<FormState>();
  final _controllerKeys = List<Map<String, TextEditingController>>;

  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
  }

//   void scanWidgets(BuildContext context) {
//   context.visitChildElements((element) {
//     if (element.widget is TextFormField) {
//       TextFormField field = element.widget as TextFormField;
//       print("Found TextFormField with label: ${field.decoration?.labelText}");
//     }
//     scanWidgets(element);
//   });
// }

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: widget.child.split(
        (customeColumnData) => customeColumnData,
        (customeRowData) => customeRowData,
        (widget) => widget,
      ),
    );
  }
}
