import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:dio/dio.dart';
import '../../core/di/service_locator.dart';
import '../../models/course.dart';
import '../../core/constants/api_constants.dart';
import '../../services/api_service.dart';

part 'courses_event.dart';
part 'courses_state.dart';

class CoursesBloc extends Bloc<CoursesEvent, CoursesState> {
  final ApiService _apiService = sl<ApiService>();

  CoursesBloc() : super(CoursesInitial()) {
    on<FetchCourses>(_onFetchCourses);
  }

  Future<void> _onFetchCourses(
    FetchCourses event,
    Emitter<CoursesState> emit,
  ) async {
    emit(CoursesLoading());
    try {
      final params = {
        'page': event.page,
        'limit': event.limit,
        if (event.provider.isNotEmpty) 'provider': event.provider,
        if (event.level.isNotEmpty) 'level': event.level,
      };
      final response = await _apiService.get(ApiConstants.courses, params: params);
      final data = response.data as Map<String, dynamic>;
      final coursesList = (data['data'] as List)
          .map((json) => Course.fromJson(json))
          .toList();
      final count = data['count'] as int? ?? 0;
      emit(CoursesLoaded(courses: coursesList, count: count, hasMore: coursesList.length == event.limit));
    } on DioException catch (e) {
      emit(CoursesError(message: e.response?.data['error'] ?? 'Failed to fetch courses'));
    } catch (e) {
      emit(CoursesError(message: 'Something went wrong'));
    }
  }
}
