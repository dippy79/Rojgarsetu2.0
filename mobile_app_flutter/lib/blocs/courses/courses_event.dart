part of 'courses_bloc.dart';

abstract class CoursesEvent extends Equatable {
  const CoursesEvent();

  @override
  List<Object?> get props => [];
}

class FetchCourses extends CoursesEvent {
  final int page;
  final int limit;
  final String provider;
  final String level;

  const FetchCourses({
    this.page = 1,
    this.limit = 20,
    this.provider = '',
    this.level = '',
  });

  @override
  List<Object?> get props => [page, limit, provider, level];
}
