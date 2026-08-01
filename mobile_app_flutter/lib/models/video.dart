
import 'package:equatable/equatable.dart';

class Video extends Equatable {
  final String id;
  final String title;
  final String channel;
  final String videoUrl;
  final String? thumbnailUrl;
  final String? category;
  final DateTime? publishedAt;

  const Video({
    required this.id,
    required this.title,
    required this.channel,
    required this.videoUrl,
    this.thumbnailUrl,
    this.category,
    this.publishedAt,
  });

  factory Video.fromJson(Map<String, dynamic> json) {
    return Video(
      id: json['id'] as String? ?? '',
      title: json['title'] as String? ?? '',
      channel: json['channel'] as String? ?? '',
      videoUrl: json['video_url'] as String? ?? '',
      thumbnailUrl: json['thumbnail_url'] as String?,
      category: json['category'] as String?,
      publishedAt: json['published_at'] != null
          ? DateTime.tryParse(json['published_at'].toString())
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'channel': channel,
      'video_url': videoUrl,
      'thumbnail_url': thumbnailUrl,
      'category': category,
      'published_at': publishedAt?.toIso8601String(),
    };
  }

  /// Immutable State Updates (BLoC / Riverpod ke liye)
  Video copyWith({
    String? id,
    String? title,
    String? channel,
    String? videoUrl,
    String? thumbnailUrl,
    String? category,
    DateTime? publishedAt,
  }) {
    return Video(
      id: id ?? this.id,
      title: title ?? this.title,
      channel: channel ?? this.channel,
      videoUrl: videoUrl ?? this.videoUrl,
      thumbnailUrl: thumbnailUrl ?? this.thumbnailUrl,
      category: category ?? this.category,
      publishedAt: publishedAt ?? this.publishedAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        title,
        channel,
        videoUrl,
        thumbnailUrl,
        category,
        publishedAt,
      ];
}