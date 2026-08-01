import 'package:equatable/equatable.dart';

class GovJob extends Equatable {
  final String id;
  final String title;
  final String? department;
  final String? location;
  final String? eligibility;
  final String? applicationDeadline;
  final String? source;
  final String? notificationUrl;
  final DateTime? postedAt;

  const GovJob({
    required this.id,
    required this.title,
    this.department,
    this.location,
    this.eligibility,
    this.applicationDeadline,
    this.source,
    this.notificationUrl,
    this.postedAt,
  });

  factory GovJob.fromJson(Map<String, dynamic> json) {
    return GovJob(
      id: json['id'] as String? ?? '',
      title: json['title'] as String? ?? '',
      department: json['department'] as String?,
      location: json['location'] as String?,
      eligibility: json['eligibility'] as String?,
      applicationDeadline: json['application_deadline'] as String?,
      source: json['source'] as String?,
      notificationUrl: json['notification_url'] as String?,
      postedAt: json['posted_at'] != null
          ? DateTime.tryParse(json['posted_at'].toString())
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'department': department,
      'location': location,
      'eligibility': eligibility,
      'application_deadline': applicationDeadline,
      'source': source,
      'notification_url': notificationUrl,
      'posted_at': postedAt?.toIso8601String(),
    };
  }

  /// Immutable State Updates (BLoC / Riverpod support)
  GovJob copyWith({
    String? id,
    String? title,
    String? department,
    String? location,
    String? eligibility,
    String? applicationDeadline,
    String? source,
    String? notificationUrl,
    DateTime? postedAt,
  }) {
    return GovJob(
      id: id ?? this.id,
      title: title ?? this.title,
      department: department ?? this.department,
      location: location ?? this.location,
      eligibility: eligibility ?? this.eligibility,
      applicationDeadline: applicationDeadline ?? this.applicationDeadline,
      source: source ?? this.source,
      notificationUrl: notificationUrl ?? this.notificationUrl,
      postedAt: postedAt ?? this.postedAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        title,
        department,
        location,
        eligibility,
        applicationDeadline,
        source,
        notificationUrl,
        postedAt,
      ];
}