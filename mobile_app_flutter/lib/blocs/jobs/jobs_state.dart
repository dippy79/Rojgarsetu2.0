part of 'jobs_bloc.dart';

abstract class JobsState extends Equatable {
  const JobsState();

  @override
  List<Object?> get props => [];
}

class JobsInitial extends JobsState {}

class JobsLoading extends JobsState {}

class JobsLoaded extends JobsState {
  final List<Job> jobs;
  final int count;
  final bool hasMore;

  const JobsLoaded({
    required this.jobs,
    required this.count,
    required this.hasMore,
  });

  @override
  List<Object?> get props => [jobs, count, hasMore];
}

class JobsError extends JobsState {
  final String message;

  const JobsError({required this.message});

  @override
  List<Object?> get props => [message];
}
