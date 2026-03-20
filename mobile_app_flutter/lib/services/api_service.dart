import 'package:dio/dio.dart';
import '../core/constants/api_constants.dart';
import '../core/storage/token_storage.dart';

class ApiService {
  late final Dio _dio;
  final TokenStorage _storage;

  ApiService(this._storage) {
    _dio = Dio(BaseOptions(
      baseUrl:        ApiConstants.baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 10),
      headers: {'Content-Type': 'application/json'},
    ));

    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await _storage.getAccessToken();
        if (token != null && token.isNotEmpty) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        handler.next(options);
      },
      onError: (error, handler) async {
        if (error.response?.statusCode == 401) {
          final refreshed = await _tryRefresh();
          if (refreshed) {
            final token = await _storage.getAccessToken();
            error.requestOptions.headers['Authorization'] =
                'Bearer $token';
            final retry = await _dio.fetch(error.requestOptions);
            return handler.resolve(retry);
          }
        }
        handler.next(error);
      },
    ));
  }

  Future<bool> _tryRefresh() async {
    try {
      final rt = await _storage.getRefreshToken();
      if (rt == null || rt.isEmpty) return false;
      final resp = await _dio.post(
          ApiConstants.refresh, data: {'token': rt});
      final newToken = resp.data['token'] as String?;
      if (newToken == null) return false;
      await _storage.saveTokens(
        accessToken:  newToken,
        refreshToken: rt,
        userID:       await _storage.getUserID() ?? '',
        role:         await _storage.getRole()   ?? '',
      );
      return true;
    } catch (_) {
      await _storage.clearAll();
      return false;
    }
  }

  Future<Response> get(String path,
      {Map<String, dynamic>? params}) =>
      _dio.get(path, queryParameters: params);

  Future<Response> post(String path, {dynamic data}) =>
      _dio.post(path, data: data);

  Future<Response> put(String path, {dynamic data}) =>
      _dio.put(path, data: data);

  Future<Response> patch(String path, {dynamic data}) =>
      _dio.patch(path, data: data);

  Future<Response> delete(String path) =>
      _dio.delete(path);
}
      });
      return response.data['status'] == 'success';
    } catch (e) {
      return false;
    }
  }

  // Applications
  Future<bool> applyForJob(String jobId) async {
    try {
      final response = await _dio.post('/users/apply-job', data: {
        'jobId': jobId,
      });
      return response.data['status'] == 'success';
    } catch (e) {
      return false;
    }
  }

  Future<List<Application>> getApplications() async {
    final response = await _dio.get('/users/applications');
    if (response.data['status'] == 'success') {
      final apps = (response.data['data']['applications'] as List)
          .map((json) => Application.fromJson(json))
          .toList();
      return apps;
    }
    return [];
  }

  // Notifications
  Future<List<Notification>> getNotifications({bool unreadOnly = false}) async {
    final response = await _dio.get('/users/notifications', queryParameters: {
      'unreadOnly': unreadOnly.toString(),
    });
    if (response.data['status'] == 'success') {
      final notifications = (response.data['data']['notifications'] as List)
          .map((json) => Notification.fromJson(json))
          .toList();
      return notifications;
    }
    return [];
  }

  Future<bool> markNotificationRead(String notificationId) async {
    try {
      final response = await _dio.put('/users/notifications/$notificationId/read');
      return response.data['status'] == 'success';
    } catch (e) {
      return false;
    }
  }
}

