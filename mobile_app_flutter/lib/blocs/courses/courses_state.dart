part of 'courses_bloc.dart';

abstract class CoursesState extends Equatable {
  const CoursesState();

  @override
  List<Object?> get props => [];
}

class CoursesInitial extends CoursesState {}

class CoursesLoading extends CoursesState {}

class CoursesLoaded extends CoursesState {
  final List<Course> courses;
  final int count;
  final bool hasMore;

  const CoursesLoaded({
    required this.courses,
    required this.count,
    required this.hasMore,
  });

  @override
  List<Object?> get props => [courses, count, hasMore];
}

class CoursesError extends CoursesState {
  final String message;

  const CoursesError({required this.message});

  @override
  List<Object?> get props => [message];
}
