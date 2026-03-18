class User {
  final String id;
  final String name;
  final String email;
  final String? phone;
  final String? resumeUrl;
  final List<String>? skills;
  final DateTime? createdAt;

  User({
    required this.id,
    required this.name,
    required this.email,
    this.phone,
    this.resumeUrl,
    this.skills,
    this.createdAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      email: json['email'] ?? '',
      phone: json['phone'],
      resumeUrl: json['resumeUrl'],
      skills: json['skills'] != null ? List<String>.from(json['skills']) : null,
      createdAt: json['createdAt'] != null ? DateTime.parse(json['createdAt']) : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'email': email,
      'phone': phone,
      'resumeUrl': resumeUrl,
      'skills': skills,
      'createdAt': createdAt?.toIso8601String(),
    };
  }
}

class Application {
  final String id;
  final String status;
  final DateTime? appliedAt;
  final DateTime? updatedAt;
  final ApplicationJob? job;

  Application({
    required this.id,
    required this.status,
    this.appliedAt,
    this.updatedAt,
    this.job,
  });

  factory Application.fromJson(Map<String, dynamic> json) {
    return Application(
      id: json['id'] ?? '',
      status: json['status'] ?? 'pending',
      appliedAt: json['appliedAt'] != null ? DateTime.parse(json['appliedAt']) : null,
      updatedAt: json['updatedAt'] != null ? DateTime.parse(json['updatedAt']) : null,
      job: json['job'] != null ? ApplicationJob.fromJson(json['job']) : null,
    );
  }
}

class ApplicationJob {
  final String? id;
  final String? title;
  final String? location;
  final String? company;

  ApplicationJob({
    this.id,
    this.title,
    this.location,
    this.company,
  });

  factory ApplicationJob.fromJson(Map<String, dynamic> json) {
    return ApplicationJob(
      id: json['id'],
      title: json['title'],
      location: json['location'],
      company: json['company'],
    );
  }
}

class AppNotification {
  final String id;
  final String title;
  final String? message;
  final String type;
  final bool isRead;
  final Map<String, dynamic>? data;
  final DateTime? createdAt;

  AppNotification({
    required this.id,
    required this.title,
    this.message,
    required this.type,
    this.isRead = false,
    this.data,
    this.createdAt,
  });

  factory AppNotification.fromJson(Map<String, dynamic> json) {
    return AppNotification(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      message: json['message'],
      type: json['type'] ?? 'system',
      isRead: json['isRead'] ?? false,
      data: json['data'],
      createdAt: json['createdAt'] != null ? DateTime.parse(json['createdAt']) : null,
    );
  }
}

