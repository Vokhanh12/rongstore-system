import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

part 'logo_state.dart';

class LogoCubit extends Cubit<LogoState> {
  LogoCubit() : super(SwapLogoInitial());
}
