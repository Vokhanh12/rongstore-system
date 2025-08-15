import 'package:objectbox/objectbox.dart';

@Entity()
class TranslationsModel {
  @Id()
  int id;
  String code;
  String translationVi;
  String translationEn;

  TranslationsModel({
    this.id = 0, 
    required this.code,
    required this.translationVi,
    required this.translationEn,
  });
}