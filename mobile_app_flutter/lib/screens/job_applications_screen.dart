import 'package:flutter/material.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import '../../core/di/service_locator.dart';
import '../../services/api_service.dart';
import '../../core/constants/api_constants.dart';
import '../../components/status_badge.dart';
import 'package:go_router/go_router.dart';

class JobApplicationsScreen extends StatefulWidget {
  final String jobId;

  const JobApplicationsScreen({super.key, required this.jobId});

  @override
  State<JobApplicationsScreen> createState() => _JobApplicationsScreenState();
}

class _JobApplicationsScreenState extends State<JobApplicationsScreen> {
  final RefreshController _refreshController = RefreshController(initialRefresh: false);
  List<dynamic> _applications = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadApplications();
  }

  Future<void> _loadApplications() async {
    try {
      final apiService = sl<ApiService>();
      final response = await apiService.get('${ApiConstants.jobs}/${widget.jobId}/applications');
      setState(() {
        _applications = response.data['data'];
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  Future<void> _updateStatus(String appId, String newStatus) async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Update Status'),
        content: Text('Set status to \$newStatus?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context, false), child: const Text('Cancel')),
          TextButton(onPressed: () => Navigator.pop(context, true), child: const Text('Update')),
        ],
      ),
    );
    if (confirm != true || !mounted) return;

    try {
      final apiService = sl<ApiService>();
      await apiService.patch('${ApiConstants.applications}/\$appId/status', data: {'status': newStatus});
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Status updated')),
        );
        _loadApplications();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to update: \$e')),
        );
      }
    }
  }

  Future<void> _onRefresh() async {
    await _loadApplications();
    _refreshController.refreshCompleted();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Job Applications')),
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
                  ElevatedButton(onPressed: _loadApplications, child: const Text('Retry')),
                ],
              ),
            )
          : SmartRefresher(
              controller: _refreshController,
              enablePullDown: true,
              header: const WaterDropHeader(),
              onRefresh: _onRefresh,
              child: _applications.isEmpty
                ? const Center(child: Text('No applications'))
                : ListView.builder(
                    itemCount: _applications.length,
                    itemBuilder: (context, index) {
                      final app = _applications[index];
                      return Card(
                        margin: const EdgeInsets.all(16),
                        child: Padding(
                          padding: const EdgeInsets.all(16),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  Expanded(child: Text(app['candidate']['full_name'] ?? 'Candidate')),
                                  StatusBadge(status: app['status'] ?? 'applied'),
                                ],
                              ),
                              const SizedBox(height: 8),
                              Text(app['candidate']['email'] ?? '', style: const TextStyle(fontStyle: FontStyle.italic)),
                              const SizedBox(height: 16),
                              DropdownButton<String>(
                                value: app['status'],
                                isExpanded: true,
                                items: ['applied', 'reviewed', 'shortlisted', 'rejected', 'hired']
                                  .map((status) => DropdownMenuItem(value: status, child: StatusBadge(status: status)))
                                  .toList(),
                                onChanged: (value) => _updateStatus(app['id'], value!),
                              ),
                            ],
                          ),
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
