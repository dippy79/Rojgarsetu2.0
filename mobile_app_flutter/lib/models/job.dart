import 'package:equatable/equatable.dart';

class Job extends Equatable {
  final String id;
  final String title;
  final String description;
  final String? location;
  final String? jobType;
  final int? salaryMin;
  final int? salaryMax;
  final List<String> skills;
  final bool isActive;
  final int views;
  final DateTime createdAt;

  const Job({
    required this.id,
    required this.title,
    required this.description,
    this.location,
    this.jobType,
    this.salaryMin,
    this.salaryMax,
    required this.skills,
    required this.isActive,
    required this.views,
    required this.createdAt,
  });

  factory Job.fromJson(Map<String, dynamic> json) {
    return Job(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      description: json['description'] ?? '',
      location: json['location'],
      jobType: json['job_type'],
      salaryMin: json['salary_min'],
      salaryMax: json['salary_max'],
      skills: List<String>.from(json['skills'] ?? []),
      isActive: json['is_active'] ?? false,
      views: json['views'] ?? 0,
      createdAt: json['created_at'] != null ? DateTime.parse(json['created_at']) : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'location': location,
      'job_type': jobType,
      'salary_min': salaryMin,
      'salary_max': salaryMax,
      'skills': skills,
      'is_active': isActive,
      'views': views,
      'created_at': createdAt.toIso8601String(),
    };
  }

  @override
  List<Object?> get props => [id, title, createdAt];
}

