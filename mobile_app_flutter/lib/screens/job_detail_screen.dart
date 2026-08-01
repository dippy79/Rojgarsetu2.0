import 'package:flutter/material.dart';
import '../../models/job.dart';
import '../../components/status_badge.dart';
import '../../core/di/service_locator.dart';
import '../../core/constants/api_constants.dart';

class JobDetailScreen extends StatefulWidget {
  final String jobId;

  const JobDetailScreen({super.key, required this.jobId});

  @override
  State<JobDetailScreen> createState() => _JobDetailScreenState();
}

class _JobDetailScreenState extends State<JobDetailScreen> {
  Job? _job;
  bool _isLoading = true;
  String? _error;
  bool _isApplying = false;

  @override
  void initState() {
    super.initState();
    _loadJob();
  }

  Future<void> _loadJob() async {
    try {
      final apiService = sl<ApiService>();
      final response = await apiService.get('${ApiConstants.jobs}/${widget.jobId}');
      setState(() {
        _job = Job.fromJson(response.data['data']);
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  Future<void> _applyForJob() async {
    if (_job == null) return;

    setState(() => _isApplying = true);

    try {
      final apiService = sl<ApiService>();
      final response = await apiService.post(
        '${ApiConstants.jobs}/${widget.jobId}/apply',
        data: {'coverLetter': _coverLetterController.text},
      );

      if (response.data['status'] == 'success') {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Application submitted successfully!')),
          );
          _coverLetterController.clear();
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to apply: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isApplying = false);
    }
  }

  final TextEditingController _coverLetterController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Job Details')),
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
                    onPressed: _loadJob,
                    child: const Text('Retry'),
                  ),
                ],
              ),
            )
          : _job == null
            ? const Center(child: Text('Job not found'))
            : SingleChildScrollView(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            _job!.title,
                            style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                        StatusBadge(status: _job!.isActive ? 'Active' : 'Inactive'),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Text(
                      '₹${_job!.salaryMin ?? 0}/- - ₹${_job!.salaryMax ?? 0}/-',
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        color: Theme.of(context).colorScheme.primary,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text('Description', style: Theme.of(context).textTheme.titleMedium),
                            const SizedBox(height: 8),
                            Text(_job!.description),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text('Details', style: Theme.of(context).textTheme.titleMedium),
                            const SizedBox(height: 8),
                            Text('Location: ${_job!.location ?? "Remote"}'),
                            Text('Type: ${_job!.jobType ?? "Full-time"}'),
                            Text('Views: ${_job!.views}'),
                            if (_job!.skills.isNotEmpty) ...[
                              const SizedBox(height: 8),
                              Text('Skills:'),
                              const SizedBox(height: 4),
                              Wrap(
                                spacing: 6,
                                children: _job!.skills.map((skill) => Chip(
                                  label: Text(skill),
                                  backgroundColor: Colors.grey[100],
                                )).toList(),
                              ),
                            ],
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(height: 24),
                    Text(
                      'Cover Letter',
                      style: Theme.of(context).textTheme.titleMedium,
                    ),
                    const SizedBox(height: 8),
                    TextField(
                      controller: _coverLetterController,
                      maxLines: 4,
                      decoration: InputDecoration(
                        hintText: 'Tell us why you are perfect for this job...',
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),
                    SizedBox(
                      width: double.infinity,
                      height: 50,
                      child: ElevatedButton(
                        onPressed: _isApplying ? null : _applyForJob,
                        child: _isApplying
                          ? const SizedBox(
                              height: 20,
                              width: 20,
                              child: CircularProgressIndicator(strokeWidth: 2),
                            )
                          : const Text('Apply Now'),
                      ),
                    ),
                  ],
                ),
              ),
    );
  }

  @override
  void dispose() {
    _coverLetterController.dispose();
    super.dispose();
  }
}

