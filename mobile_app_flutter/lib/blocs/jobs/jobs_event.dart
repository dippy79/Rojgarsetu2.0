part of 'jobs_bloc.dart';

abstract class JobsEvent extends Equatable {
  const JobsEvent();

  @override
  List<Object?> get props => [];
}

class FetchJobs extends JobsEvent {
  final int page;
  final int limit;
  final String location;
  final String jobType;

  const FetchJobs({
    this.page = 1,
    this.limit = 20,
    this.location = '',
    this.jobType = '',
  });

  @override
  List<Object?> get props => [page, limit, location, jobType];
}
