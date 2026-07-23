import 'api_service.dart';
import '../core/constants/api_constants.dart';
import '../core/storage/token_storage.dart';

class AuthService {
  final ApiService   _api;
  final TokenStorage _storage;

  AuthService(this._api, this._storage);

  Future<Map<String, dynamic>> register({
    required String name,
    required String email,
    required String password,
    required String role,
    String? phone,
  }) async {
    final resp = await _api.post(ApiConstants.register, data: {
      'name':     name,
      'email':    email,
      'password': password,
      'role':     role,
      if (phone != null) 'phone': phone,
    });
    return resp.data as Map<String, dynamic>;
  }

  Future<void> login({
    required String email,
    required String password,
  }) async {
    final resp = await _api.post(ApiConstants.login, data: {
      'email':    email,
      'password': password,
    });
    final data = resp.data as Map<String, dynamic>;
    await _storage.saveTokens(
      accessToken:  data['token']         as String? ?? '',
      refreshToken: data['refresh_token'] as String? ?? '',
      userID:       data['id']            as String? ?? '',
      role:         data['role']          as String? ?? '',
    );
  }

  Future<void> logout() async {
    try { await _api.post(ApiConstants.logout); } catch (_) {}
    await _storage.clearAll();
  }

  Future<bool> isLoggedIn() async {
    final token = await _storage.getAccessToken();
    return token != null && token.isNotEmpty;
  }

  Future<String?> getCurrentUserID() => _storage.getUserID();
  Future<String?> getCurrentRole()   => _storage.getRole();
}

