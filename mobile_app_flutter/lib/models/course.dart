import 'package:equatable/equatable.dart';

class Course extends Equatable {
  final String id;
  final String title;
  final String? provider;
  final String? mode;
  final String? level;
  final String? description;
  final String? url;
  final DateTime createdAt;

  const Course({
    required this.id,
    required this.title,
    this.provider,
    this.mode,
    this.level,
    this.description,
    this.url,
    required this.createdAt,
  });

  factory Course.fromJson(Map<String, dynamic> json) {
    return Course(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      provider: json['provider'],
      mode: json['mode'],
      level: json['level'],
      description: json['description'],
      url: json['url'],
      createdAt: json['created_at'] != null ? DateTime.parse(json['created_at']) : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'provider': provider,
      'mode': mode,
      'level': level,
      'description': description,
      'url': url,
      'created_at': createdAt.toIso8601String(),
    };
  }

  @override
  List<Object?> get props => [id, title, createdAt];
}

