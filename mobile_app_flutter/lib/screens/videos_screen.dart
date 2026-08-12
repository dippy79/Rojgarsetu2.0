import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:url_launcher/url_launcher.dart';
import '../../blocs/videos/videos_bloc.dart';
import '../../models/video.dart';
import '../../components/filter_bar.dart';
import '../../core/di/service_locator.dart';

class VideosScreen extends StatefulWidget {
  const VideosScreen({super.key});

  @override
  State<VideosScreen> createState() => _VideosScreenState();
}

class _VideosScreenState extends State<VideosScreen> {
  final RefreshController _refreshController = RefreshController(initialRefresh: false);
  String _filterChannel = '';
  final String _filterCategory = '';

  @override
  void initState() {
    super.initState();
    context.read<VideosBloc>().add(const FetchVideos(page: 1, limit: 10));
  }

  void _onRefresh() {
    context.read<VideosBloc>().add(const FetchVideos(page: 1, limit: 10));
  }

  void _onFilterChannel(String? value) {
    setState(() => _filterChannel = value ?? '');
    context.read<VideosBloc>().add(FetchVideos(page: 1, limit: 10, channel: value ?? ''));
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
      appBar: AppBar(title: const Text('Videos')),
      body: BlocBuilder<VideosBloc, VideosState>(
        builder: (context, state) {
          if (state is VideosLoading) {
            return const Center(child: CircularProgressIndicator());
          }
          if (state is VideosError) {
            return ListView(
              children: [
                ListTile(
                  leading: const Icon(Icons.error, color: Colors.red),
                  title: Text(state.message),
                  trailing: IconButton(
                    icon: const Icon(Icons.refresh),
                    onPressed: () => context.read<VideosBloc>().add(const FetchVideos(page: 1, limit: 10)),
                  ),
                ),
              ],
            );
          }
          if (state is VideosLoaded && state.videos.isEmpty) {
            return const Center(child: Text('No videos found'));
          }
          if (state is VideosLoaded) {
            return SmartRefresher(
              controller: _refreshController,
              enablePullDown: true,
              header: const WaterDropHeader(),
              onRefresh: _onRefresh,
              child: ListView.builder(
                itemCount: state.videos.length,
                itemBuilder: (context, index) {
                  final video = state.videos[index];
                  return Card(
                    margin: const EdgeInsets.all(16),
                    child: ListTile(
                      leading: CachedNetworkImage(
                        imageUrl: video.thumbnail ?? '',
                        width: 60,
                        height: 60,
                        fit: BoxFit.cover,
                        placeholder: (context, url) => const CircularProgressIndicator(),
                        errorWidget: (context, url, error) => const Icon(Icons.video_library),
                      ),
                      title: Text(video.title),
                      subtitle: Text(video.channel),
                      trailing: ElevatedButton(
                        onPressed: () => _launchUrl(video.videoUrl),
                        child: const Text('Watch'),
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
