class Job {
  final String id;
  final String title;
  final String? location;
  final String? jobType;
  final int? salaryMin;
  final int? salaryMax;
  final String? eligibility;
  final String? description;
  final String? applicationUrl;
  final DateTime? postedAt;
  final String? source;
  final Company? company;

  Job({
    required this.id,
    required this.title,
    this.location,
    this.jobType,
    this.salaryMin,
    this.salaryMax,
    this.eligibility,
    this.description,
    this.applicationUrl,
    this.postedAt,
    this.source,
    this.company,
  });

  factory Job.fromJson(Map<String, dynamic> json) {
    return Job(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      location: json['location'],
      jobType: json['jobType'],
      salaryMin: json['salaryMin'],
      salaryMax: json['salaryMax'],
      eligibility: json['eligibility'],
      description: json['description'],
      applicationUrl: json['applicationUrl'],
      postedAt: json['postedAt'] != null ? DateTime.parse(json['postedAt']) : null,
      source: json['source'],
      company: json['company'] != null ? Company.fromJson(json['company']) : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'location': location,
      'jobType': jobType,
      'salaryMin': salaryMin,
      'salaryMax': salaryMax,
      'eligibility': eligibility,
      'description': description,
      'applicationUrl': applicationUrl,
      'postedAt': postedAt?.toIso8601String(),
      'source': source,
      'company': company?.toJson(),
    };
  }

  String get salaryRange {
    if (salaryMin == null && salaryMax == null) return 'Not disclosed';
    if (salaryMin != null && salaryMax != null) {
      return '₹${(salaryMin! / 100000).toStringAsFixed(1)}L - ₹${(salaryMax! / 100000).toStringAsFixed(1)}L';
    }
    if (salaryMin != null) return '₹${(salaryMin! / 100000).toStringAsFixed(1)}L+';
    return 'Up to ₹${(salaryMax! / 100000).toStringAsFixed(1)}L';
  }
}

class Company {
  final String? id;
  final String? name;
  final String? logo;
  final String? website;

  Company({
    this.id,
    this.name,
    this.logo,
    this.website,
  });

  factory Company.fromJson(Map<String, dynamic> json) {
    return Company(
      id: json['id'],
      name: json['name'],
      logo: json['logo'],
      website: json['website'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'logo': logo,
      'website': website,
    };
  }
}

