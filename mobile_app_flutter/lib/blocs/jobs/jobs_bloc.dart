import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:dio/dio.dart';
import '../../core/di/service_locator.dart';
import '../../models/job.dart';
import '../../core/constants/api_constants.dart';
import '../../services/api_service.dart';

part 'jobs_event.dart';
part 'jobs_state.dart';

class JobsBloc extends Bloc<JobsEvent, JobsState> {
  final ApiService _apiService = sl<ApiService>();

  JobsBloc() : super(JobsInitial()) {
    on<FetchJobs>(_onFetchJobs);
  }

  Future<void> _onFetchJobs(
    FetchJobs event,
    Emitter<JobsState> emit,
  ) async {
    emit(JobsLoading());
    try {
      final params = {
        'page': event.page,
        'limit': event.limit,
        if (event.location.isNotEmpty) 'location': event.location,
        if (event.jobType.isNotEmpty) 'job_type': event.jobType,
      };
      final response = await _apiService.get(ApiConstants.jobs, params: params);
      final data = response.data as Map<String, dynamic>;
      final jobsList = (data['data'] as List)
          .map((json) => Job.fromJson(json))
          .toList();
      final count = data['count'] as int? ?? 0;
      emit(JobsLoaded(jobs: jobsList, count: count, hasMore: jobsList.length == event.limit));
    } on DioException catch (e) {
      emit(JobsError(message: e.response?.data['error'] ?? 'Failed to fetch jobs'));
    } catch (e) {
      emit(JobsError(message: 'Something went wrong'));
    }
  }
}
