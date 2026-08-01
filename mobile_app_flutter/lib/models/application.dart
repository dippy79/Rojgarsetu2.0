import 'package:equatable/equatable.dart';
import 'package:flutter/material.dart';

// Assuming Job model is in another file, import it accordingly
import 'job.dart';

class JobApplication extends Equatable {
  final String id;
  final String status; // applied, reviewed, shortlisted, rejected, hired
  final DateTime? appliedAt;
  final DateTime? updatedAt;
  final Job? job;
  final String? coverLetter;

  const JobApplication({
    required this.id,
    required this.status,
    this.appliedAt,
    this.updatedAt,
    this.job,
    this.coverLetter,
  });

  factory JobApplication.fromJson(Map<String, dynamic> json) {
    return JobApplication(
      id: json['id'] as String? ?? '',
      status: json['status'] as String? ?? 'applied',
      appliedAt: json['created_at'] != null 
          ? DateTime.tryParse(json['created_at'].toString()) 
          : null,
      updatedAt: json['updated_at'] != null 
          ? DateTime.tryParse(json['updated_at'].toString()) 
          : null,
      job: json['job'] != null ? Job.fromJson(json['job'] as Map<String, dynamic>) : null,
      coverLetter: json['cover_letter'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'status': status,
      'created_at': appliedAt?.toIso8601String(),
      'updated_at': updatedAt?.toIso8601String(),
      'job': job?.toJson(),
      'cover_letter': coverLetter,
    };
  }

  /// UI Helper for Status Colors
  Color get statusColor {
    return switch (status.toLowerCase()) {
      'hired' => Colors.green,
      'shortlisted' => Colors.orange,
      'reviewed' => Colors.blue,
      'applied' => Colors.grey,
      'rejected' => Colors.red,
      _ => Colors.grey,
    };
  }

  /// Immutable State Updates (BLoC / Riverpod ke liye)
  JobApplication copyWith({
    String? id,
    String? status,
    DateTime? appliedAt,
    DateTime? updatedAt,
    Job? job,
    String? coverLetter,
  }) {
    return JobApplication(
      id: id ?? this.id,
      status: status ?? this.status,
      appliedAt: appliedAt ?? this.appliedAt,
      updatedAt: updatedAt ?? this.updatedAt,
      job: job ?? this.job,
      coverLetter: coverLetter ?? this.coverLetter,
    );
  }

  @override
  List<Object?> get props => [
        id,
        status,
        appliedAt,
        updatedAt,
        job,
        coverLetter,
      ];
}