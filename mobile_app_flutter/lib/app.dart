import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'blocs/auth/auth_bloc.dart';
import 'blocs/jobs/jobs_bloc.dart';
import 'blocs/courses/courses_bloc.dart';
import 'theme.dart';
import 'screens/login_screen.dart';
import 'screens/register_screen.dart';
import 'screens/home_screen.dart';
import 'screens/candidate_profile_screen.dart';
import 'screens/company_profile_screen.dart';
import 'screens/my_applications_screen.dart';
import 'screens/job_applications_screen.dart';
import 'blocs/gov_jobs/gov_jobs_bloc.dart';
import 'blocs/videos/videos_bloc.dart';

class RojgarSetuApp extends StatelessWidget {
  const RojgarSetuApp({super.key});

  @override
  Widget build(BuildContext context) {
    final GoRouter router = GoRouter(
      initialLocation: '/',
      routes: [
        GoRoute(
          path: '/',
          builder: (context, state) => BlocBuilder<AuthBloc, AuthState>(
            builder: (context, state) {
              if (state is AuthLoading) {
                return const Scaffold(body: Center(child: CircularProgressIndicator()));
              }
              if (state is AuthAuthenticated) {
                return HomeScreen(role: state.role);
              }
              return const LoginScreen();
            },
          ),
        ),
        GoRoute(
          path: '/login',
          builder: (context, state) => const LoginScreen(),
        ),
        GoRoute(
          path: '/register',
          builder: (context, state) => const RegisterScreen(),
        ),
        GoRoute(
          path: '/home/:role',
          builder: (context, state) {
            final role = state.pathParameters['role']!;
            return HomeScreen(role: role);
          },
        ),
        GoRoute(
          path: '/profile',
          builder: (context, state) => BlocBuilder<AuthBloc, AuthState>(
            builder: (context, authState) {
              if (authState is! AuthAuthenticated) return const SizedBox();
              return authState.role == 'candidate' 
                ? const CandidateProfileScreen()
                : const CompanyProfileScreen();
            },
          ),
        ),
        GoRoute(
          path: '/my-applications',
          builder: (context, state) => const MyApplicationsScreen(),
        ),
        GoRoute(
          path: '/jobs/:jobId/applications',
          builder: (context, state) {
            final jobId = state.pathParameters['jobId']!;
            return JobApplicationsScreen(jobId: jobId);
          },
        ),
      ],
      redirect: (context, state) {
        final authState = context.read<AuthBloc>().state;
        final isLoggedIn = authState is AuthAuthenticated;
        final isAuthPath = state.uri.toString() == '/login' || state.uri.toString() == '/register';
        if (!isLoggedIn && !isAuthPath) {
          return '/login';
        }
        if (isLoggedIn && isAuthPath) {
          return '/';
        }
        return null;
      },
    );

    return MultiBlocProvider(
      providers: [
        BlocProvider.value(value: context.read<AuthBloc>()),
        BlocProvider.value(value: context.read<JobsBloc>()),
        BlocProvider.value(value: context.read<CoursesBloc>()),
        BlocProvider.value(value: context.read<GovJobsBloc>()),
        BlocProvider.value(value: context.read<VideosBloc>()),
      ],
      child: MaterialApp.router(
        title: 'RojgarSetu 2.0',
        theme: AppTheme.lightTheme,
        darkTheme: AppTheme.darkTheme,
        themeMode: ThemeMode.system,
        routerConfig: router,
        debugShowCheckedModeBanner: false,
      ),
    );
  }
}
