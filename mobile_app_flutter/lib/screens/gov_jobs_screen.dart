import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import '../../blocs/gov_jobs/gov_jobs_bloc.dart';
import '../../models/gov_job.dart';
import '../../components/filter_bar.dart';
import '../../components/job_card.dart'; // reuse style
import 'package:url_launcher/url_launcher.dart';
import '../../core/di/service_locator.dart';

class GovJobsScreen extends StatefulWidget {
  const GovJobsScreen({super.key});

  @override
  State<GovJobsScreen> createState() => _GovJobsScreenState();
}

class _GovJobsScreenState extends State<GovJobsScreen> {
  final RefreshController _refreshController = RefreshController(initialRefresh: false);
  String _filterDepartment = '';
  String _filterLocation = '';

  @override
  void initState() {
    super.initState();
    context.read<GovJobsBloc>().add(const FetchGovJobs(page: 1, limit: 10));
  }

  void _onRefresh() {
    context.read<GovJobsBloc>().add(FetchGovJobs(page: 1, limit: 10));
  }

  void _onLoadMore() {
    final state = context.read<GovJobsBloc>().state;
    if (state is GovJobsLoaded && state.hasMore) {
      context.read<GovJobsBloc>().add(FetchGovJobs(
        page: (state.govJobs.length ~/ 10) + 1,
        limit: 10,
        department: _filterDepartment,
        location: _filterLocation,
      ));
    }
  }

  void _onFilterDepartment(String? value) {
    setState(() => _filterDepartment = value ?? '');
    context.read<GovJobsBloc>().add(FetchGovJobs(page: 1, limit: 10, department: value ?? ''));
  }

  void _launchUrl(String url) async {
    final uri = Uri.parse(url);
    if (await canLaunchUrl(uri)) {
      await launchUrl(uri, mode: LaunchMode.externalApplication);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            floating: true,
            title: const Text('Government Jobs'),
            bottom: PreferredSize(
              preferredSize: const Size.fromHeight(56),
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: FilterBar(
                  options: const [
                    FilterOption('All Departments', ''),
                    FilterOption('Railways', 'railways'),
                    FilterOption('Banking', 'banking'),
                    FilterOption('SSC', 'ssc'),
                    FilterOption('UPSC', 'upsc'),
                  ],
                  onSelected: _onFilterDepartment,
                ),
              ),
            ),
          ),
          BlocBuilder<GovJobsBloc, GovJobsState>(
            builder: (context, state) {
              if (state is GovJobsLoading) {
                return const SliverToBoxAdapter(
                  child: Center(child: CircularProgressIndicator()),
                );
              }
              if (state is GovJobsError) {
                return SliverToBoxAdapter(
                  child: ListTile(
                    leading: const Icon(Icons.error, color: Colors.red),
                    title: Text(state.message),
                    trailing: IconButton(
                      icon: const Icon(Icons.refresh),
                      onPressed: () => context.read<GovJobsBloc>().add(const FetchGovJobs(page: 1, limit: 10)),
                    ),
                  ),
                );
              }
              if (state is GovJobsLoaded && state.govJobs.isEmpty) {
                return SliverFillRemaining(
                  child: Center(child: Text('No gov jobs found')),
                );
              }
              if (state is GovJobsLoaded) {
                return SliverList(
                  delegate: SliverChildBuilderDelegate(
                    childCount: state.govJobs.length,
                    (context, index) {
                      final govJob = state.govJobs[index];
                      return Card(
                        margin: const EdgeInsets.all(16),
                        child: ListTile(
                          title: Text(govJob.title),
                          subtitle: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text('Department: \${govJob.department}'),
                              Text('Location: \${govJob.location}'),
                              Text('Deadline: \${govJob.deadline ?? 'N/A'}'),
                              Text('Source: \${govJob.source}'),
                            ],
                          ),
                          trailing: ElevatedButton(
                            onPressed: () => _launchUrl(govJob.notificationUrl),
                            child: const Text('Apply'),
                          ),
                        ),
                      );
                    },
                  ),
                );
              }
              return const SliverFillRemaining(child: SizedBox());
            },
          ),
        ],
      ),
      bottomNavigationBar: BlocBuilder<GovJobsBloc, GovJobsState>(
        builder: (context, state) {
          if (state is GovJobsLoaded && state.hasMore) {
            return BottomAppBar(
              child: SmartRefresher(
                controller: _refreshController,
                enablePullUp: true,
                onLoading: _onLoadMore,
                child: const SizedBox(height: 50),
              ),
            );
          }
          return const SizedBox();
        },
      ),
    );
  }
}
