class Course {
  final String id;
  final String title;
  final String provider;
  final String? url;
  final List<String>? skills;
  final String? duration;
  final String? level;
  final String? thumbnailUrl;
  final bool isFree;
  final double? price;

  Course({
    required this.id,
    required this.title,
    required this.provider,
    this.url,
    this.skills,
    this.duration,
    this.level,
    this.thumbnailUrl,
    this.isFree = true,
    this.price,
  });

  factory Course.fromJson(Map<String, dynamic> json) {
    return Course(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      provider: json['provider'] ?? '',
      url: json['url'],
      skills: json['skills'] != null ? List<String>.from(json['skills']) : null,
      duration: json['duration'],
      level: json['level'],
      thumbnailUrl: json['thumbnailUrl'],
      isFree: json['isFree'] ?? true,
      price: json['price']?.toDouble(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'provider': provider,
      'url': url,
      'skills': skills,
      'duration': duration,
      'level': level,
      'thumbnailUrl': thumbnailUrl,
      'isFree': isFree,
      'price': price,
    };
  }
}

