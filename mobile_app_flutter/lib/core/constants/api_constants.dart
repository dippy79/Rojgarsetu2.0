class ApiConstants {
  static const String baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://10.0.2.2:8080',
  );

  static const String register   = '/api/v1/auth/register';
  static const String login      = '/api/v1/auth/login';
  static const String logout     = '/api/v1/auth/logout';
  static const String refresh    = '/api/v1/auth/refresh';

  static const String jobs       = '/api/v1/jobs';
  static const String govJobs    = '/api/v1/gov-jobs';
  static const String privJobs   = '/api/v1/priv-jobs';
  static const String courses    = '/api/v1/courses';
  static const String videos     = '/api/v1/videos';

  static const String candidates      = '/api/v1/candidates';
  static const String candidateMe     = '/api/v1/candidates/me';
  static const String myApplications  = '/api/v1/candidates/me/applications';

  static const String companies  = '/api/v1/companies';
  static const String companyMe  = '/api/v1/companies/me';
  static const String myJobs     = '/api/v1/companies/me/jobs';

  static const String applications = '/api/v1/applications';
}
