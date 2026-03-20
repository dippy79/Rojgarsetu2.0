import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'core/di/service_locator.dart';
import 'blocs/auth/auth_bloc.dart';
import 'blocs/jobs/jobs_bloc.dart';
import 'blocs/courses/courses_bloc.dart';
import 'app.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  setupServiceLocator();

  runApp(
    MultiBlocProvider(
      providers: [
        BlocProvider(
          create: (context) => sl<AuthBloc>()..add(CheckAuthStatus()),
        ),
        BlocProvider(create: (context) => JobsBloc()),
        BlocProvider(create: (context) => CoursesBloc()),
      ],
      child: const RojgarSetuApp(),
    ),
  );
}

