import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../blocs/courses/courses_bloc.dart';
import '../../models/course.dart';
import '../../components/filter_bar.dart';
import '../../core/di/service_locator.dart';

class CoursesScreen extends StatefulWidget {
  const CoursesScreen({super.key});

  @override
  State<CoursesScreen> createState() => _CoursesScreenState();
}

class _CoursesScreenState extends State<CoursesScreen> {
  final RefreshController _refreshController = RefreshController(initialRefresh: false);
  String _filterProvider = '';
  String _filterMode = '';
  String _filterLevel = '';

  @override
  void initState() {
    super.initState();
    context.read<CoursesBloc>().add(const FetchCourses(page: 1, limit: 10));
  }

  void _onRefresh() {
    context.read<CoursesBloc>().add(const FetchCourses(page: 1, limit: 10));
  }

  void _onFilterProvider(String? value) {
    setState(() => _filterProvider = value ?? '');
    context.read<CoursesBloc>().add(FetchCourses(page: 1, limit: 10, provider: value ?? ''));
  }

  void _launchUrl(String? url) async {
    if (url == null || url.isEmpty) return;
    final uri = Uri.parse(url);
    if (await canLaunchUrl(uri)) {
      await launchUrl(uri, mode: LaunchMode.externalApplication);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Courses')),
      body: BlocBuilder<CoursesBloc, CoursesState>(
        builder: (context, state) {
          if (state is CoursesLoading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (state is CoursesError) {
            return ListView(
              children: [
                ListTile(
                  leading: const Icon(Icons.error, color: Colors.red),
                  title: Text(state.message),
                  trailing: IconButton(
                    icon: const Icon(Icons.refresh),
                    onPressed: () => context.read<CoursesBloc>().add(const FetchCourses(page: 1, limit: 10)),
                  ),
                ),
              ],
            );
          }
          if (state is CoursesLoaded && state.courses.isEmpty) {
            return const Center(child: Text('No courses found'));
          }
          if (state is CoursesLoaded) {
            return SmartRefresher(
              controller: _refreshController,
              enablePullDown: true,
              header: const WaterDropHeader(),
              onRefresh: _onRefresh,
              child: ListView.builder(
                itemCount: state.courses.length,
                itemBuilder: (context, index) {
                  final course = state.courses[index];
                  return Card(
                    margin: const EdgeInsets.all(16),
                    child: ListTile(
                      title: Text(course.title),
                      subtitle: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('Provider: \${course.provider}'),
                          Text('Mode: \${course.mode}'),
                          Text('Level: \${course.level}'),
                        ],
                      ),
                      trailing: ElevatedButton(
                        onPressed: () => _launchUrl(course.url),
                        child: const Text('Enroll'),
                      ),
                    ),
                  );
                },
              ),
            );
          }
          return const SizedBox();
        },
      ),
    );
  }
}
