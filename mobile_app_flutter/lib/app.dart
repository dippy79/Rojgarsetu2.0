import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import 'blocs/auth/auth_bloc.dart';
import 'blocs/courses/courses_bloc.dart';
import 'blocs/gov_jobs/gov_jobs_bloc.dart';
import 'blocs/jobs/jobs_bloc.dart';
import 'blocs/videos/videos_bloc.dart';
import 'screens/candidate_profile_screen.dart';
import 'screens/company_profile_screen.dart';
import 'screens/home_screen.dart';
import 'screens/job_applications_screen.dart';
import 'screens/login_screen.dart';
import 'screens/my_applications_screen.dart';
import 'screens/register_screen.dart';
import 'screens/splash_screen.dart';
import 'theme.dart';

/// Helper class to convert BLoC Stream into a Listenable for GoRouter
class GoRouterRefreshStream extends ChangeNotifier {
  late final StreamSubscription<dynamic> _subscription;

  GoRouterRefreshStream(Stream<dynamic> stream) {
    notifyListeners();
    _subscription = stream.asBroadcastStream().listen(
          (_) => notifyListeners(),
        );
  }

  @override
  void dispose() {
    _subscription.cancel();
    super.dispose();
  }
}

class RojgarSetuApp extends StatefulWidget {
  const RojgarSetuApp({super.key});

  @override
  State<RojgarSetuApp> createState() => _RojgarSetuAppState();
}

class _RojgarSetuAppState extends State<RojgarSetuApp> {
  late final GoRouter _router;
  late final GoRouterRefreshStream _authRefreshStream;

  @override
  void initState() {
    super.initState();
    final authBloc = context.read<AuthBloc>();
    _authRefreshStream = GoRouterRefreshStream(authBloc.stream);

    _router = GoRouter(
      initialLocation: '/',
      refreshListenable: _authRefreshStream,
      routes: [
        GoRoute(
          path: '/',
          builder: (context, state) => const SplashScreen(),
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
            final role = state.pathParameters['role'] ?? 'candidate';
            return HomeScreen(role: role);
          },
        ),
        GoRoute(
          path: '/profile',
          builder: (context, state) {
            final authState = context.read<AuthBloc>().state;
            if (authState is AuthAuthenticated) {
              return authState.role == 'candidate'
                  ? const CandidateProfileScreen()
                  : const CompanyProfileScreen();
            }
            return const SizedBox.shrink();
          },
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
        final isAuthPath = state.matchedLocation == '/login' ||
            state.matchedLocation == '/register';

        // 1. App initialization / splash phase
        if (authState is AuthInitial || authState is AuthLoading) {
          return '/';
        }

        // 2. Unauthenticated user trying to access protected routes
        if (!isLoggedIn && !isAuthPath) {
          return '/login';
        }

        // 3. Authenticated user trying to access login/register/splash
        if (isLoggedIn && (isAuthPath || state.matchedLocation == '/')) {
          return '/home/${authState.role}';
        }

        return null;
      },
    );
  }

  @override
  void dispose() {
    _authRefreshStream.dispose();
    _router.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
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
        routerConfig: _router,
        debugShowCheckedModeBanner: false,
      ),
    );
  }
}