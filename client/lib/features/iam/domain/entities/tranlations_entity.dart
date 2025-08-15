import 'package:equatable/equatable.dart';
import 'package:objectbox/objectbox.dart';

@Entity()
class TranlationsEntity extends Equatable {
  const TranlationsEntity(
      {required this.id,
      required this.code,
      required this.tranlationVi,
      required this.tranlationEn});

  @Id(assignable: true)
  final int id;
  final String code;
  final String tranlationVi;
  final String tranlationEn;

  @override
  List<Object?> get props => [
        id,
        code,
        tranlationVi,
        tranlationEn,
      ];
}
