import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dio/dio.dart';
import '../../services/auth_service.dart';

// ── Events ──
abstract class AuthEvent {}

class CheckAuthStatus extends AuthEvent {}

class LoginRequested extends AuthEvent {
  final String email;
  final String password;
  LoginRequested({required this.email, required this.password});
}

class RegisterRequested extends AuthEvent {
  final String name;
  final String email;
  final String password;
  final String role;
  final String? phone;
  RegisterRequested({
    required this.name,
    required this.email,
    required this.password,
    required this.role,
    this.phone,
  });
}

class LogoutRequested extends AuthEvent {}

// ── States ──
abstract class AuthState {}
class AuthInitial         extends AuthState {}
class AuthLoading         extends AuthState {}
class AuthUnauthenticated extends AuthState {}

class AuthAuthenticated extends AuthState {
  final String userID;
  final String role;
  AuthAuthenticated({required this.userID, required this.role});
}

class AuthError extends AuthState {
  final String message;
  AuthError(this.message);
}

// ── Bloc ──
class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthService _authService;

  AuthBloc({required AuthService authService})
      : _authService = authService,
        super(AuthInitial()) {
    on<CheckAuthStatus>(_onCheckAuthStatus);
    on<LoginRequested>(_onLoginRequested);
    on<RegisterRequested>(_onRegisterRequested);
    on<LogoutRequested>(_onLogoutRequested);
  }

  Future<void> _onCheckAuthStatus(
    CheckAuthStatus event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());
    try {
      final loggedIn = await _authService.isLoggedIn();
      if (loggedIn) {
        final userID = await _authService.getCurrentUserID() ?? '';
        final role   = await _authService.getCurrentRole()   ?? '';
        emit(AuthAuthenticated(userID: userID, role: role));
      } else {
        emit(AuthUnauthenticated());
      }
    } catch (_) {
      emit(AuthUnauthenticated());
    }
  }

  Future<void> _onLoginRequested(
    LoginRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());
    try {
      await _authService.login(
          email: event.email, password: event.password);
      final userID = await _authService.getCurrentUserID() ?? '';
      final role   = await _authService.getCurrentRole()   ?? '';
      emit(AuthAuthenticated(userID: userID, role: role));
    } on DioException catch (e) {
      final msg = (e.response?.data as Map?)?['error']
          as String? ?? 'Login failed';
      emit(AuthError(msg));
    } catch (_) {
      emit(AuthError('Something went wrong'));
    }
  }

  Future<void> _onRegisterRequested(
    RegisterRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(AuthLoading());
    try {
      await _authService.register(
        name:     event.name,
        email:    event.email,
        password: event.password,
        role:     event.role,
        phone:    event.phone,
      );
      emit(AuthUnauthenticated());
    } on DioException catch (e) {
      final msg = (e.response?.data as Map?)?['error']
          as String? ?? 'Registration failed';
      emit(AuthError(msg));
    } catch (_) {
      emit(AuthError('Something went wrong'));
    }
  }

  Future<void> _onLogoutRequested(
    LogoutRequested event,
    Emitter<AuthState> emit,
  ) async {
    await _authService.logout();
    emit(AuthUnauthenticated());
  }
}
