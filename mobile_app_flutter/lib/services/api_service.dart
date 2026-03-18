import 'dart:convert';
import 'package:dio/dio.dart';
import '../models/job.dart';
import '../models/course.dart';
import '../models/user.dart';

class ApiService {
  static const String baseUrl = 'http://localhost:3000/api';
  late final Dio _dio;
  String? _accessToken;

  ApiService() {
    _dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 30),
      receiveTimeout: const Duration(seconds: 30),
      headers: {
        'Content-Type': 'application/json',
      },
    ));

    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) {
        if (_accessToken != null) {
          options.headers['Authorization'] = 'Bearer $_accessToken';
        }
        return handler.next(options);
      },
      onError: (error, handler) {
        return handler.next(error);
      },
    ));
  }

  void setToken(String? token) {
    _accessToken = token;
  }

  // Auth
  Future<Map<String, dynamic>> register({
    required String name,
    required String email,
    required String password,
  }) async {
    final response = await _dio.post('/auth/register', data: {
      'name': name,
      'email': email,
      'password': password,
    });
    return response.data;
  }

  Future<Map<String, dynamic>> login({
    required String email,
    required String password,
  }) async {
    final response = await _dio.post('/auth/login', data: {
      'email': email,
      'password': password,
    });
    if (response.data['status'] == 'success') {
      _accessToken = response.data['data']['accessToken'];
    }
    return response.data;
  }

  Future<void> logout() async {
    _accessToken = null;
  }

  // Jobs
  Future<List<Job>> getJobs({
    int page = 1,
    int limit = 20,
    String? location,
    String? source,
    String? jobType,
    String? search,
  }) async {
    final response = await _dio.get('/jobs', queryParameters: {
      'page': page,
      'limit': limit,
      if (location != null) 'location': location,
      if (source != null) 'source': source,
      if (jobType != null) 'jobType': jobType,
      if (search != null) 'search': search,
    });

    if (response.data['status'] == 'success') {
      final jobs = (response.data['data']['jobs'] as List)
          .map((json) => Job.fromJson(json))
          .toList();
      return jobs;
    }
    return [];
  }

  Future<Job?> getJobDetails(String jobId) async {
    try {
      final response = await _dio.get('/jobs/$jobId');
      if (response.data['status'] == 'success') {
        return Job.fromJson(response.data['data']);
      }
    } catch (e) {
      return null;
    }
    return null;
  }

  Future<List<Job>> getRecommendedJobs() async {
    final response = await _dio.get('/recommendations/jobs');
    if (response.data['status'] == 'success') {
      final jobs = (response.data['data']['jobs'] as List)
          .map((json) => Job.fromJson(json))
          .toList();
      return jobs;
    }
    return [];
  }

  // Courses
  Future<List<Course>> getCourses({
    int page = 1,
    int limit = 20,
    String? provider,
    String? level,
    bool? free,
  }) async {
    final response = await _dio.get('/courses', queryParameters: {
      'page': page,
      'limit': limit,
      if (provider != null) 'provider': provider,
      if (level != null) 'level': level,
      if (free != null) 'free': free.toString(),
    });

    if (response.data['status'] == 'success') {
      final courses = (response.data['data']['courses'] as List)
          .map((json) => Course.fromJson(json))
          .toList();
      return courses;
    }
    return [];
  }

  Future<List<Course>> getRecommendedCourses() async {
    final response = await _dio.get('/recommendations/courses');
    if (response.data['status'] == 'success') {
      final courses = (response.data['data']['courses'] as List)
          .map((json) => Course.fromJson(json))
          .toList();
      return courses;
    }
    return [];
  }

  // User Profile
  Future<User?> getProfile() async {
    try {
      final response = await _dio.get('/users/profile');
      if (response.data['status'] == 'success') {
        return User.fromJson(response.data['data']);
      }
    } catch (e) {
      return null;
    }
    return null;
  }

  Future<bool> updateProfile({
    String? name,
    String? phone,
    List<String>? skills,
    String? resumeUrl,
  }) async {
    try {
      final response = await _dio.put('/users/profile', data: {
        if (name != null) 'name': name,
        if (phone != null) 'phone': phone,
        if (skills != null) 'skills': skills,
        if (resumeUrl != null) 'resumeUrl': resumeUrl,
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

