import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import '../../blocs/jobs/jobs_bloc.dart';
import '../../blocs/jobs/jobs_event.dart';
import '../../components/job_card.dart';
import '../../components/filter_bar.dart';
import 'job_detail_screen.dart';

class JobsListScreen extends StatefulWidget {
  final bool showCompany;

  const JobsListScreen({super.key, this.showCompany = false});

  @override
  State<JobsListScreen> createState() => _JobsListScreenState();
}

class _JobsListScreenState extends State<JobsListScreen> {
  final RefreshController _refreshController = RefreshController(initialRefresh: false);
  final String _filterLocation = '';
  String _filterJobType = '';

  void _onRefresh() => context.read<JobsBloc>().add(FetchJobs(page: 1, limit: 10));

  void _onLoadMore() {
    final state = context.read<JobsBloc>().state;
    if (state is JobsLoaded && state.hasMore) {
      context.read<JobsBloc>().add(FetchJobs(
        page: state.jobs.length ~/ 10 + 1,
        limit: 10,
        location: _filterLocation,
        jobType: _filterJobType,
      ));
    }
  }

  void _onFilterSelected(String? value) {
    setState(() {
      _filterJobType = value ?? '';
    });
    context.read<JobsBloc>().add(FetchJobs(
      page: 1, 
      limit: 10,
      location: _filterLocation,
      jobType: value ?? '',
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SmartRefresher(
        controller: _refreshController,
        enablePullDown: true,
        enablePullUp: true,
        header: const WaterDropHeader(),
        onRefresh: _onRefresh,
        onLoading: _onLoadMore,
        child: CustomScrollView(
          slivers: [
            SliverAppBar(
              floating: true,
              title: Text('Jobs'),
              bottom: PreferredSize(
                preferredSize: const Size.fromHeight(56),
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: FilterBar(
                    options: const [
                      FilterOption('All', ''),
                      FilterOption('Full-time', 'full_time'),
                      FilterOption('Part-time', 'part_time'),
                      FilterOption('Remote', 'remote'),
                      FilterOption('IT', 'it'),
                      FilterOption('Sales', 'sales'),
                    ],
                    onSelected: _onFilterSelected,
                  ),
                ),
              ),
            ),
            BlocBuilder<JobsBloc, JobsState>(
              builder: (context, state) {
                return SliverList(
                  delegate: SliverChildBuilderDelegate(
                    (context, index) {
                      if (state is JobsInitial || state is JobsLoading) {
                        return ListTile(
                          leading: const CircularProgressIndicator(strokeWidth: 2),
                          title: const Text('Loading jobs...'),
                        );
                      }

                      if (state is JobsError) {
                        return ListTile(
                          leading: const Icon(Icons.error, color: Colors.red),
                          title: Text(state.message),
                          trailing: IconButton(
                            icon: const Icon(Icons.refresh),
                            onPressed: _onRefresh,
                          ),
                        );
                      }

                      if (state is JobsLoaded && state.jobs.isEmpty) {
                        return const ListTile(
                          leading: Icon(Icons.inbox_outlined),
                          title: Text('No jobs found'),
                          subtitle: Text('Try adjusting your filters'),
                        );
                      }

                      if (state is JobsLoaded) {
                        final job = state.jobs[index];
                        return JobCard(
                          key: ValueKey(job.id),
                          job: job,
                          showCompany: widget.showCompany,
                          onTap: () => Navigator.of(context).push(
                            MaterialPageRoute(
                              builder: (_) => JobDetailScreen(jobId: job.id),
                            ),
                          ),
                        );
                      }

                      return const SizedBox.shrink();
                    },
                    childCount: 1, // Dynamic handled by state
                  ),
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  @override
  void dispose() {
    _refreshController.dispose();
    super.dispose();
  }
}

