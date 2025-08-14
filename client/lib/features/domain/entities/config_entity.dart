import 'package:equatable/equatable.dart';
import 'package:objectbox/objectbox.dart';

@Entity()
class ConfigEntity extends Equatable {
  const ConfigEntity({required this.updateTranlation});

  @Id(assignable: true)
  final bool updateTranlation;

  @override
  List<Object?> get props => [updateTranlation];
}
