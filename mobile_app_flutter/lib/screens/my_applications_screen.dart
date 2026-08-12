import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import '../../core/di/service_locator.dart';
import '../../services/api_service.dart';
import '../../core/constants/api_constants.dart';
import '../../models/application.dart';
import '../../components/status_badge.dart';
import 'package:go_router/go_router.dart';

class MyApplicationsScreen extends StatefulWidget {
  const MyApplicationsScreen({super.key});

  @override
  State<MyApplicationsScreen> createState() => _MyApplicationsScreenState();
}

class _MyApplicationsScreenState extends State<MyApplicationsScreen> {
  final RefreshController _refreshController = RefreshController(initialRefresh: false);
  final List<Application> _applications = [];
  bool _isLoading = true;
  String? _error;
  int _page = 1;
  bool _hasMore = true;

  @override
  void initState() {
    super.initState();
    _loadApplications();
  }

  Future<void> _loadApplications({bool refresh = false}) async {
    if (refresh) {
      _page = 1;
      _applications.clear();
      _refreshController.refreshCompleted();
    }
    try {
      final apiService = sl<ApiService>();
      final params = {'page': _page, 'limit': 10};
      final response = await apiService.get(ApiConstants.myApplications, params: params);
      final data = response.data['data'] as List;
      final newApps = data.map((json) => Application.fromJson(json)).toList();
      setState(() {
        _applications.addAll(newApps);
        _hasMore = newApps.length == 10;
        _page++;
        _isLoading = false;
        _error = null;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
      _refreshController.refreshFailed();
    }
  }

  Future<void> _onRefresh() async {
    await _loadApplications(refresh: true);
  }

  Future<void> _onLoading() async {
    if (_hasMore) {
      await _loadApplications();
      _refreshController.loadComplete();
    } else {
      _refreshController.loadNoData();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('My Applications')),
      body: _isLoading
        ? const Center(child: CircularProgressIndicator())
        : _error != null
          ? Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error, size: 64, color: Colors.red),
                  const SizedBox(height: 16),
                  Text(_error!),
                  ElevatedButton(
                    onPressed: () => _loadApplications(refresh: true),
                    child: const Text('Retry'),
                  ),
                ],
              ),
            )
          : SmartRefresher(
              controller: _refreshController,
              enablePullDown: true,
              enablePullUp: _hasMore,
              header: const WaterDropHeader(),
              onRefresh: _onRefresh,
              onLoading: _onLoading,
              child: _applications.isEmpty
                ? const Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.inbox_outlined, size: 64),
                        SizedBox(height: 16),
                        Text('No applications yet'),
                        Text('Apply to jobs to see them here'),
                      ],
                    ),
                  )
                : ListView.builder(
                    itemCount: _applications.length,
                    itemBuilder: (context, index) {
                      final app = _applications[index];
                      return Card(
                        margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
                        child: ListTile(
                          leading: CircleAvatar(
                            backgroundColor: app.statusColor.withOpacity(0.2),
                            child: Icon(
                              switch (app.status) {
                                'hired' => Icons.check_circle,
                                'shortlisted' => Icons.star,
                                'reviewed' => Icons.visibility,
                                'applied' => Icons.send,
                                'rejected' => Icons.cancel,
                                _ => Icons.help_outline,
                              },
                              color: app.statusColor,
                            ),
                          ),
                          title: Text(app.job?.title ?? 'Job Title'),
                          subtitle: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(app.job?.companyName ?? 'Company'),
                              StatusBadge(status: app.status),
                              if (app.coverLetter != null && app.coverLetter!.isNotEmpty)
                                Padding(
                                  padding: const EdgeInsets.only(top: 4),
                                  child: Text(
                                    app.coverLetter!,
                                    maxLines: 2,
                                    overflow: TextOverflow.ellipsis,
                                    style: Theme.of(context).textTheme.bodySmall,
                                  ),
                                ),
                              Text('Applied: ${app.appliedAt?.toLocal().toString().split(' ')[0] ?? ''}'),
                            ],
                          ),
                          trailing: const Icon(Icons.arrow_forward_ios),
                          onTap: () {
                            context.go('/jobs/${app.job?.id}');
                          },
                        ),
                      );
                    },
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
