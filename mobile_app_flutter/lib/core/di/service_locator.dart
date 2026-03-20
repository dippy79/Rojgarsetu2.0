import 'package:get_it/get_it.dart';
import '../../services/api_service.dart';
import '../../services/auth_service.dart';
import '../storage/token_storage.dart';
import '../../blocs/auth/auth_bloc.dart';

final sl = GetIt.instance;

void setupServiceLocator() {
  sl.registerLazySingleton<TokenStorage>(() => TokenStorage());

  sl.registerLazySingleton<ApiService>(
      () => ApiService(sl<TokenStorage>()));

  sl.registerLazySingleton<AuthService>(
      () => AuthService(sl<ApiService>(), sl<TokenStorage>()));

  sl.registerFactory<AuthBloc>(
      () => AuthBloc(authService: sl<AuthService>()));
}
