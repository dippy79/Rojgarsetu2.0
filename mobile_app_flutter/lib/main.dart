import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'app.dart';
import 'services/auth_service.dart';
import 'services/api_service.dart';
import 'services/notification_service.dart';
import 'blocs/auth/auth_bloc.dart';
import 'blocs/jobs/jobs_bloc.dart';
import 'blocs/courses/courses_bloc.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize services
  final authService = AuthService();
  final apiService = ApiService();
  final notificationService = NotificationService();
  
  await notificationService.initialize();
  
  runApp(
    MultiProvider(
      providers: [
        Provider<AuthService>.value(value: authService),
        Provider<ApiService>.value(value: apiService),
        Provider<NotificationService>.value(value: notificationService),
        BlocProvider<AuthBloc>(
          create: (context) => AuthBloc(authService: authService)..add(CheckAuthStatus()),
        ),
        BlocProvider<JobsBloc>(
          create: (context) => JobsBloc(apiService: apiService),
        ),
        BlocProvider<CoursesBloc>(
          create: (context) => CoursesBloc(apiService: apiService),
        ),
      ],
      child: const RojgarSetuApp(),
    ),
  );
}

