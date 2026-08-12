import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../blocs/gov_jobs/gov_jobs_bloc.dart';
import '../../components/filter_bar.dart';
import '../../models/gov_job.dart';

class GovJobsScreen extends StatefulWidget {
  const GovJobsScreen({super.key});

  @override
  State<GovJobsScreen> createState() => _GovJobsScreenState();
}

class _GovJobsScreenState extends State<GovJobsScreen> {
  final RefreshController _refreshController =
      RefreshController(initialRefresh: false);
  String _filterDepartment = '';
  final String _filterLocation = '';

  @override
  void initState() {
    super.initState();
    _fetchJobs(page: 1);
  }

  @override
  void dispose() {
    _refreshController.dispose();
    super.dispose();
  }

  void _fetchJobs({required int page}) {
    context.read<GovJobsBloc>().add(
          FetchGovJobs(
            page: page,
            limit: 10,
            department: _filterDepartment,
            location: _filterLocation,
          ),
        );
  }

  void _onRefresh() {
    _fetchJobs(page: 1);
  }

  void _onLoadMore() {
    final state = context.read<GovJobsBloc>().state;
    if (state is GovJobsLoaded && state.hasMore) {
      final nextPage = (state.govJobs.length ~/ 10) + 1;
      _fetchJobs(page: nextPage);
    } else {
      _refreshController.loadNoData();
    }
  }

  void _onFilterDepartment(String? value) {
    setState(() => _filterDepartment = value ?? '');
    _fetchJobs(page: 1);
  }

  Future<void> _launchUrl(String? url) async {
    if (url == null || url.isEmpty) return;
    final uri = Uri.tryParse(url);
    if (uri != null && await canLaunchUrl(uri)) {
      await launchUrl(uri, mode: LaunchMode.externalApplication);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: BlocConsumer<GovJobsBloc, GovJobsState>(
        listener: (context, state) {
          // Reset smart refresher indicators when bloc completes
          if (state is GovJobsLoaded) {
            _refreshController.refreshCompleted();
            if (state.hasMore) {
              _refreshController.loadComplete();
            } else {
              _refreshController.loadNoData();
            }
          } else if (state is GovJobsError) {
            _refreshController.refreshFailed();
            _refreshController.loadFailed();
          }
        },
        builder: (context, state) {
          return SmartRefresher(
            controller: _refreshController,
            enablePullDown: true,
            enablePullUp: state is GovJobsLoaded && state.hasMore,
            onRefresh: _onRefresh,
            onLoading: _onLoadMore,
            child: CustomScrollView(
              slivers: [
                SliverAppBar(
                  floating: true,
                  pinned: true,
                  title: const Text('Government Jobs'),
                  bottom: PreferredSize(
                    preferredSize: const Size.fromHeight(60),
                    child: Padding(
                      padding: const EdgeInsets.all(8.0),
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
                _buildContentSliver(state),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildContentSliver(GovJobsState state) {
    if (state is GovJobsLoading) {
      return const SliverFillRemaining(
        child: Center(child: CircularProgressIndicator()),
      );
    }

    if (state is GovJobsError) {
      return SliverFillRemaining(
        hasScrollBody: false,
        child: Center(
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(Icons.error_outline, color: Colors.red, size: 48),
                const SizedBox(height: 12),
                Text(
                  state.message,
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.bodyLarge,
                ),
                const SizedBox(height: 16),
                ElevatedButton.icon(
                  onPressed: () => _fetchJobs(page: 1),
                  icon: const Icon(Icons.refresh),
                  label: const Text('Retry'),
                ),
              ],
            ),
          ),
        ),
      );
    }

    if (state is GovJobsLoaded) {
      if (state.govJobs.isEmpty) {
        return const SliverFillRemaining(
          hasScrollBody: false,
          child: Center(
            child: Text(
              'No government jobs found.',
              style: TextStyle(fontSize: 16, color: Colors.grey),
            ),
          ),
        );
      }

      return SliverPadding(
        padding: const EdgeInsets.all(12),
        sliver: SliverList(
          delegate: SliverChildBuilderDelegate(
            (context, index) {
              final govJob = state.govJobs[index];
              return Card(
                elevation: 2,
                margin: const EdgeInsets.symmetric(vertical: 6),
                child: Padding(
                  padding: const EdgeInsets.all(12),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        govJob.title,
                        style: Theme.of(context)
                            .textTheme
                            .titleMedium
                            ?.copyWith(fontWeight: FontWeight.bold),
                      ),
                      const SizedBox(height: 8),
                      Text('Department: ${govJob.department ?? 'N/A'}'),
                      Text('Location: ${govJob.location ?? 'N/A'}'),
                      Text(
                          'Deadline: ${govJob.applicationDeadline ?? 'N/A'}'),
                      Text('Source: ${govJob.source ?? 'N/A'}'),
                      const SizedBox(height: 8),
                      Align(
                        alignment: Alignment.centerRight,
                        child: ElevatedButton(
                          onPressed: govJob.notificationUrl != null
                              ? () => _launchUrl(govJob.notificationUrl)
                              : null,
                          child: const Text('Apply'),
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
            childCount: state.govJobs.length,
          ),
        ),
      );
    }

    return const SliverFillRemaining(child: SizedBox());
  }
}