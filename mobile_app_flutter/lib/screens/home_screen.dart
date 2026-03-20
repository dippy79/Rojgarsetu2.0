import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../blocs/auth/auth_bloc.dart';
import '../blocs/jobs/jobs_bloc.dart';
import 'jobs_list_screen.dart';
import 'courses_screen.dart';
import 'gov_jobs_screen.dart';
import 'videos_screen.dart';
import '../components/filter_bar.dart';
import '../models/job.dart';
import '../core/di/service_locator.dart';

class HomeScreen extends StatefulWidget {
  final String role;

  const HomeScreen({super.key, required this.role});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedIndex = 0;

  late final List<Widget> _tabs;

  @override
  void initState() {
    super.initState();
    _tabs = _buildTabs();
    context.read<JobsBloc>().add(FetchJobs(page: 1, limit: 10));
  }

  List<Widget> _buildTabs() {
    if (widget.role == 'candidate') {
      return [
        const JobsListScreen(),
        const GovJobsScreen(),
        const CoursesScreen(),
        const VideosScreen(),
      ];
    } else {
      return [
        const JobsListScreen(showCompany: true),
        const Placeholder(), // Company jobs
        const Placeholder(), // Post job
        const Placeholder(), // Analytics
      ];
    }
  }

  final List<FilterOption> _jobFilters = const [
    FilterOption('All', ''),
    FilterOption('Full-time', 'full_time'),
    FilterOption('Part-time', 'part_time'),
    FilterOption('Remote', 'remote'),
    FilterOption('IT', 'it'),
    FilterOption('Sales', 'sales'),
  ];

  @override
  Widget build(BuildContext context) {
    final authBloc = context.watch<AuthBloc>();

    return Scaffold(
      appBar: AppBar(
        title: Text('${widget.role.toUpperCase()} Dashboard'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () => context.read<AuthBloc>().add(LogoutRequested()),
          ),
        ],
      ),
      body: IndexedStack(
        index: _selectedIndex,
        children: _tabs,
      ),
      bottomNavigationBar: widget.role == 'candidate'
        ? NavigationBar(
            selectedIndex: _selectedIndex,
            onDestinationSelected: (index) => setState(() => _selectedIndex = index),
            destinations: const [
              NavigationDestination(
                icon: Icon(Icons.work_outline),
                selectedIcon: Icon(Icons.work),
                label: 'Jobs',
              ),
              NavigationDestination(
                icon: Icon(Icons.account_balance),
                selectedIcon: Icon(Icons.account_balance),
                label: 'Gov Jobs',
              ),
              NavigationDestination(
                icon: Icon(Icons.school_outlined),
                selectedIcon: Icon(Icons.school),
                label: 'Courses',
              ),
              NavigationDestination(
                icon: Icon(Icons.play_circle_outline),
                selectedIcon: Icon(Icons.play_circle),
                label: 'Videos',
              ),
            ],
          )
        : null,
      floatingActionButton: widget.role == 'company' 
        ? FloatingActionButton(
            onPressed: () {}, // Navigate to post job
            child: const Icon(Icons.add),
          )
        : null,
    );
  }
}

