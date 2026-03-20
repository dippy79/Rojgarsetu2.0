import 'package:flutter/material.dart';
import '../../core/di/service_locator.dart';
import '../../services/api_service.dart';
import '../../core/constants/api_constants.dart';
import '../../components/status_badge.dart'; // for reuse if needed

class CandidateProfileScreen extends StatefulWidget {
  const CandidateProfileScreen({super.key});

  @override
  State<CandidateProfileScreen> createState() => _CandidateProfileScreenState();
}

class _CandidateProfileScreenState extends State<CandidateProfileScreen> {
  bool _isLoading = true;
  bool _isSaving = false;
  String? _error;
  Map<String, dynamic> _profile = {};
  final _formKey = GlobalKey<FormState>();
  final _fullNameController = TextEditingController();
  final _phoneController = TextEditingController();
  final _locationController = TextEditingController();
  final _bioController = TextEditingController();
  final _skillsController = TextEditingController();
  int? _experience;

  @override
  void initState() {
    super.initState();
    _loadProfile();
  }

  Future<void> _loadProfile() async {
    try {
      final apiService = sl<ApiService>();
      final response = await apiService.get(ApiConstants.candidatesMe);
      setState(() {
        _profile = response.data['data'];
        _fullNameController.text = _profile['full_name'] ?? '';
        _phoneController.text = _profile['phone'] ?? '';
        _locationController.text = _profile['location'] ?? '';
        _bioController.text = _profile['bio'] ?? '';
        _skillsController.text = (_profile['skills'] as List? ?? []).join(', ');
        _experience = _profile['experience'];
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load profile: \$_error')),
        );
      }
    }
  }

  Future<void> _saveProfile() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isSaving = true);

    try {
      final apiService = sl<ApiService>();
      final data = {
        'full_name': _fullNameController.text,
        'phone': _phoneController.text,
        'location': _locationController.text,
        'bio': _bioController.text,
        'skills': _skillsController.text.split(',').map((s) => s.trim()).where((s) => s.isNotEmpty).toList(),
        'experience': _experience,
      };
      await apiService.put(ApiConstants.candidateMe, data: data);


      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Profile updated successfully!')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to save profile: \$e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('My Profile'),
        actions: [
          IconButton(
            icon: const Icon(Icons.save),
            onPressed: _isSaving ? null : _saveProfile,
          ),
        ],
      ),
      body: _isLoading
        ? const Center(child: CircularProgressIndicator())
        : _error != null
          ? Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error, size: 64, color: Colors.red),
                  const SizedBox(height: 16),
                  const Text('Failed to load profile'),
                  ElevatedButton(
                    onPressed: _loadProfile,
                    child: const Text('Retry'),
                  ),
                ],
              ),
            )
          : Form(
              key: _formKey,
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    // Avatar placeholder
                    Center(
                      child: CircleAvatar(
                        radius: 50,
                        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
                        child: const Icon(Icons.person, size: 50),
                      ),
                    ),
                    const SizedBox(height: 24),
                    TextFormField(
                      controller: _fullNameController,
                      decoration: const InputDecoration(
                        labelText: 'Full Name',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.person),
                      ),
                      validator: (value) => value!.isEmpty ? 'Full name is required' : null,
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _phoneController,
                      decoration: const InputDecoration(
                        labelText: 'Phone',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.phone),
                      ),
                      validator: (value) => value!.length < 10 ? 'Valid phone required' : null,
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _locationController,
                      decoration: const InputDecoration(
                        labelText: 'Location',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.location_on),
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _bioController,
                      maxLines: 4,
                      decoration: const InputDecoration(
                        labelText: 'Bio',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.description),
                        alignLabelWithHint: true,
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _skillsController,
                      decoration: const InputDecoration(
                        labelText: 'Skills (comma separated)',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.psychology),
                        hintText: 'JavaScript, React, Node.js',
                      ),
                    ),
                    const SizedBox(height: 16),
                    DropdownButtonFormField<int>(
                      value: _experience,
                      decoration: const InputDecoration(
                        labelText: 'Experience (years)',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.work_history),
                      ),
                      items: List.generate(31, (index) => DropdownMenuItem(
                        value: index,
                        child: Text('$index years'),
                      )),
                      onChanged: (value) => setState(() => _experience = value),
                    ),
                    const SizedBox(height: 24),
                    ElevatedButton.icon(
                      onPressed: _isSaving ? null : _saveProfile,
                      icon: _isSaving 
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                          )
                        : const Icon(Icons.save),
                      label: Text(_isSaving ? 'Saving...' : 'Save Profile'),
                      style: ElevatedButton.styleFrom(padding: const EdgeInsets.symmetric(vertical: 16)),
                    ),
                    // Skills chips preview
                    if (_skillsController.text.isNotEmpty) ...[
                      const SizedBox(height: 16),
                      const Text('Skills Preview:', style: TextStyle(fontWeight: FontWeight.bold)),
                      const SizedBox(height: 8),
                      Wrap(
                        spacing: 6,
                        runSpacing: 6,
                        children: _skillsController.text.split(',').map((skill) {
                          final trimmed = skill.trim();
                          if (trimmed.isEmpty) return const SizedBox.shrink();
                          return Chip(
                            label: Text(trimmed),
                            backgroundColor: Theme.of(context).colorScheme.primaryContainer,
                          );
                        }).toList(),
                      ),
                    ],
                  ],
                ),
              ),
            ),
    );
  }

  @override
  void dispose() {
    _fullNameController.dispose();
    _phoneController.dispose();
    _locationController.dispose();
    _bioController.dispose();
    _skillsController.dispose();
    super.dispose();
  }
}
