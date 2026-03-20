import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class TokenStorage {
  static const _accessKey  = 'access_token';
  static const _refreshKey = 'refresh_token';
  static const _userIDKey  = 'user_id';
  static const _roleKey    = 'role';

  final FlutterSecureStorage _storage;

  TokenStorage() : _storage = const FlutterSecureStorage(
aOptions: AndroidOptions(),
    iOptions: IOSOptions(
        accessibility: KeychainAccessibility.first_unlock),
  );

  Future<void> saveTokens({
    required String accessToken,
    required String refreshToken,
    required String userID,
    required String role,
  }) async {
    await Future.wait([
      _storage.write(key: _accessKey,  value: accessToken),
      _storage.write(key: _refreshKey, value: refreshToken),
      _storage.write(key: _userIDKey,  value: userID),
      _storage.write(key: _roleKey,    value: role),
    ]);
  }

  Future<String?> getAccessToken()  async =>
      _storage.read(key: _accessKey);
  Future<String?> getRefreshToken() async =>
      _storage.read(key: _refreshKey);
  Future<String?> getUserID()       async =>
      _storage.read(key: _userIDKey);
  Future<String?> getRole()         async =>
      _storage.read(key: _roleKey);

  Future<void> clearAll() async => _storage.deleteAll();
}
