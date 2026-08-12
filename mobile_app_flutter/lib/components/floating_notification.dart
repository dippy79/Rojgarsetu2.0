import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../services/api_service.dart';
import '../services/auth_service.dart';
import '../core/storage/token_storage.dart';

/// Notification model matching the backend API response
class NotificationItem {
  final String id;
  final String userId;
  final String? enrollmentId;
  final String notificationType;
  final String channel;
  final String title;
  final String message;
  final Map<String, dynamic> payload;
  final DateTime sentAt;
  final DateTime? readAt;
  final DateTime? clickedAt;

  NotificationItem({
    required this.id,
    required this.userId,
    this.enrollmentId,
    required this.notificationType,
    required this.channel,
    required this.title,
    required this.message,
    required this.payload,
    required this.sentAt,
    this.readAt,
    this.clickedAt,
  });

  factory NotificationItem.fromJson(Map<String, dynamic> json) {
    return NotificationItem(
      id: json['id'] as String? ?? '',
      userId: json['user_id'] as String? ?? '',
      enrollmentId: json['enrollment_id'] as String?,
      notificationType: json['notification_type'] as String? ?? '',
      channel: json['channel'] as String? ?? 'in_app',
      title: json['title'] as String? ?? '',
      message: json['message'] as String? ?? '',
      payload: json['payload'] is String 
          ? jsonDecode(json['payload'] as String) as Map<String, dynamic>
          : (json['payload'] as Map<String, dynamic>? ?? {}),
      sentAt: DateTime.tryParse(json['sent_at'] as String? ?? '') ?? DateTime.now(),
      readAt: json['read_at'] != null 
          ? DateTime.tryParse(json['read_at'] as String)
          : null,
      clickedAt: json['clicked_at'] != null
          ? DateTime.tryParse(json['clicked_at'] as String)
          : null,
    );
  }

  bool get isUnread => readAt == null;

  int get daysUntilExpiry => payload['days_until_expiry'] as int? ?? 0;

  String get tradeName => payload['trade_name'] as String? ?? '';

  String get enrollmentIdFromPayload => payload['enrollment_id'] as String? ?? '';
}

/// Floating Notification Widget for Flutter
/// 
/// A floating action button with a dropdown notification panel
/// that follows US Enterprise design standards with glassmorphism effects.
class FloatingNotification extends StatefulWidget {
  final AuthService authService;
  final ApiService apiService;
  final TokenStorage tokenStorage;
  final VoidCallback? onNavigateToEnrollment;
  final VoidCallback? onNavigateToTrade;
  final Duration pollInterval;

  const FloatingNotification({
    super.key,
    required this.authService,
    required this.apiService,
    required this.tokenStorage,
    this.onNavigateToEnrollment,
    this.onNavigateToTrade,
    this.pollInterval = const Duration(seconds: 30),
  });

  @override
  State<FloatingNotification> createState() => _FloatingNotificationState();
}

class _FloatingNotificationState extends State<FloatingNotification>
    with SingleTickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _scaleAnimation;
  late Animation<double> _fadeAnimation;
  late Animation<Offset> _slideAnimation;

  List<NotificationItem> _notifications = [];
  int _unreadCount = 0;
  bool _isOpen = false;
  bool _isLoading = false;
  Timer? _pollTimer;
  OverlayEntry? _overlayEntry;

  @override
  void initState() {
    super.initState();
    _initAnimations();
    _fetchNotifications();
    _startPolling();
  }

  void _initAnimations() {
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 200),
      vsync: this,
    );

    _scaleAnimation = Tween<double>(begin: 0.9, end: 1.0).animate(
      CurvedAnimation(parent: _animationController, curve: Curves.easeOutBack),
    );

    _fadeAnimation = Tween<double>(begin: 0.0, end: 1.0).animate(
      CurvedAnimation(parent: _animationController, curve: Curves.easeOut),
    );

    _slideAnimation = Tween<Offset>(
      begin: const Offset(0, 0.1),
      end: Offset.zero,
    ).animate(CurvedAnimation(parent: _animationController, curve: Curves.easeOut));
  }

  @override
  void dispose() {
    _animationController.dispose();
    _pollTimer?.cancel();
    _removeOverlay();
    super.dispose();
  }

  void _startPolling() {
    _pollTimer = Timer.periodic(widget.pollInterval, (_) {
      if (mounted) _fetchNotifications();
    });
  }

  Future<void> _fetchNotifications() async {
    if (_isLoading) return;

    final isLoggedIn = await widget.authService.isLoggedIn();
    if (!isLoggedIn) return;

    final userId = await widget.authService.getCurrentUserID();
    if (userId == null) return;

    setState(() => _isLoading = true);

    try {
      final response = await widget.apiService.get(
        '/api/v1/users/$userId/notifications',
        queryParameters: {'limit': 10},
      );

      if (response.statusCode == 200) {
        final data = response.data as Map<String, dynamic>;
        if (data['status'] == 'success' && data['data'] != null) {
          final List<dynamic> notificationsJson = data['data'] as List<dynamic>;
          final notifications = notificationsJson
              .map((json) => NotificationItem.fromJson(json as Map<String, dynamic>))
              .toList();

          if (mounted) {
            setState(() {
              _notifications = notifications;
              _unreadCount = notifications.where((n) => n.isUnread).length;
            });
          }
        }
      }
    } catch (e) {
      debugPrint('Failed to fetch notifications: $e');
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _markAsRead(String notificationId) async {
    try {
      final response = await widget.apiService.post(
        '/api/v1/notifications/$notificationId/read',
      );

      if (response.statusCode == 200) {
        setState(() {
          final index = _notifications.indexWhere((n) => n.id == notificationId);
          if (index != -1) {
            _notifications[index] = NotificationItem(
              id: _notifications[index].id,
              userId: _notifications[index].userId,
              enrollmentId: _notifications[index].enrollmentId,
              notificationType: _notifications[index].notificationType,
              channel: _notifications[index].channel,
              title: _notifications[index].title,
              message: _notifications[index].message,
              payload: _notifications[index].payload,
              sentAt: _notifications[index].sentAt,
              readAt: DateTime.now(),
              clickedAt: _notifications[index].clickedAt,
            );
            _unreadCount = _notifications.where((n) => n.isUnread).length;
          }
        });
      }
    } catch (e) {
      debugPrint('Failed to mark notification as read: $e');
    }
  }

  void _handleNotificationTap(NotificationItem notification) {
    _markAsRead(notification.id);
    _toggleDropdown();

    // Navigate based on payload
    if (notification.enrollmentIdFromPayload.isNotEmpty) {
      widget.onNavigateToEnrollment?.call();
    } else if (notification.payload['trade_id'] != null) {
      widget.onNavigateToTrade?.call();
    }
  }

  void _markAllAsRead() {
    for (final notification in _notifications) {
      if (notification.isUnread) {
        _markAsRead(notification.id);
      }
    }
  }

  void _toggleDropdown() {
    if (_isOpen) {
      _removeOverlay();
    } else {
      _showOverlay();
    }
    setState(() => _isOpen = !_isOpen);
    _animationController.forward(from: 0);
  }

  void _showOverlay() {
    _overlayEntry = OverlayEntry(
      builder: (context) => _buildDropdownOverlay(),
    );
    Overlay.of(context).insert(_overlayEntry!);
  }

  void _removeOverlay() {
    _overlayEntry?.remove();
    _overlayEntry = null;
  }

  Widget _buildDropdownOverlay() {
    return Positioned(
      bottom: 80,
      right: 16,
      left: 16,
      child: Material(
        color: Colors.transparent,
        child: ScaleTransition(
          scale: _scaleAnimation,
          child: FadeTransition(
            opacity: _fadeAnimation,
            child: SlideTransition(
              position: _slideAnimation,
              child: _buildDropdown(),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildDropdown() {
    return Container(
      constraints: const BoxConstraints(maxHeight: 500, maxWidth: 400),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.95),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: Colors.grey.shade200),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 20,
            offset: const Offset(0, 8),
          ),
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(16),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            _buildHeader(),
            Flexible(child: _buildList()),
            _buildFooter(),
          ],
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: Colors.grey.shade50,
        border: Border(bottom: BorderSide(color: Colors.grey.shade200)),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          const Text(
            'Notifications',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: Color(0xFF0F172A),
            ),
          ),
          if (_unreadCount > 0)
            TextButton(
              onPressed: _markAllAsRead,
              style: TextButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
              ),
              child: const Text(
                'Mark all read',
                style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w500,
                  color: Color(0xFF1E88E5),
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildList() {
    if (_isLoading) {
      return const Center(
        child: Padding(
          padding: EdgeInsets.all(40),
          child: CircularProgressIndicator(strokeWidth: 2),
        ),
      );
    }

    if (_notifications.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(40),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text('🔔', style: TextStyle(fontSize: 40, color: Colors.grey.shade400)),
              const SizedBox(height: 12),
              Text(
                'No notifications yet',
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w500,
                  color: Colors.grey.shade600,
                ),
              ),
              const SizedBox(height: 4),
              Text(
                'You\'ll see enrollment reminders here',
                style: TextStyle(fontSize: 13, color: Colors.grey.shade500),
              ),
            ],
          ),
        ),
      );
    }

    return ListView.separated(
      shrinkWrap: true,
      padding: const EdgeInsets.all(8),
      itemCount: _notifications.length,
      separatorBuilder: (_, __) => const SizedBox(height: 4),
      itemBuilder: (context, index) {
        final notification = _notifications[index];
        return _buildNotificationItem(notification);
      },
    );
  }

  Widget _buildNotificationItem(NotificationItem notification) {
    final isUnread = notification.isUnread;
    final color = _getNotificationColor(notification.notificationType);
    final icon = _getNotificationIcon(notification.notificationType);

    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: () => _handleNotificationTap(notification),
        borderRadius: BorderRadius.circular(12),
        child: Container(
          padding: const EdgeInsets.all(14),
          decoration: BoxDecoration(
            color: isUnread ? const Color(0xFF1E88E5).withOpacity(0.05) : Colors.transparent,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: isUnread ? const Color(0xFF1E88E5).withOpacity(0.15) : Colors.transparent,
            ),
          ),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: color,
                  borderRadius: BorderRadius.circular(10),
                  boxShadow: [
                    BoxShadow(
                      color: color.withOpacity(0.3),
                      blurRadius: 8,
                      offset: const Offset(0, 2),
                    ),
                  ],
                ),
                child: Center(
                  child: Text(icon, style: const TextStyle(fontSize: 18)),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Expanded(
                          child: Text(
                            notification.title,
                            style: TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w600,
                              color: isUnread ? const Color(0xFF0F172A) : Colors.grey.shade800,
                              height: 1.4,
                            ),
                          ),
                        ),
                        const SizedBox(width: 8),
                        Text(
                          _formatTimeAgo(notification.sentAt),
                          style: TextStyle(
                            fontSize: 11,
                            fontWeight: FontWeight.w500,
                            color: Colors.grey.shade500,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Text(
                      notification.message,
                      style: TextStyle(
                        fontSize: 13,
                        color: Colors.grey.shade700,
                        height: 1.5,
                      ),
                    ),
                    if (notification.daysUntilExpiry > 0) ...[
                      const SizedBox(height: 8),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                        decoration: BoxDecoration(
                          color: notification.daysUntilExpiry <= 1 
                              ? const Color(0xFFEF4444) 
                              : const Color(0xFFF59E0B),
                          borderRadius: BorderRadius.circular(6),
                        ),
                        child: Text(
                          'Expires in ${notification.daysUntilExpiry} day${notification.daysUntilExpiry != 1 ? 's' : ''}',
                          style: const TextStyle(
                            fontSize: 11,
                            fontWeight: FontWeight.w600,
                            color: Colors.white,
                          ),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
              if (isUnread)
                Container(
                  margin: const EdgeInsets.only(top: 4),
                  width: 8,
                  height: 8,
                  decoration: BoxDecoration(
                    color: const Color(0xFF1E88E5),
                    shape: BoxShape.circle,
                    boxShadow: [
                      BoxShadow(
                        color: const Color(0xFF1E88E5).withOpacity(0.3),
                        blurRadius: 4,
                        spreadRadius: 2,
                      ),
                    ],
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildFooter() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: Colors.grey.shade50,
        border: Border(top: BorderSide(color: Colors.grey.shade200)),
      ),
      child: Center(
        child: TextButton(
          onPressed: () {
            _toggleDropdown();
            // Navigate to full notifications screen
          },
          style: TextButton.styleFrom(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
          ),
          child: const Text(
            'View all notifications',
            style: TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w500,
              color: Color(0xFF1E88E5),
            ),
          ),
        ),
      ),
    );
  }

  Color _getNotificationColor(String type) {
    switch (type) {
      case 'expiry_final':
        return const Color(0xFFEF4444);
      case 'expiry_warning':
        return const Color(0xFFF59E0B);
      case 'enrollment_reminder':
        return const Color(0xFF1E88E5);
      case 'course_update':
        return const Color(0xFF8B5CF6);
      default:
        return const Color(0xFF64748B);
    }
  }

  String _getNotificationIcon(String type) {
    switch (type) {
      case 'expiry_final':
        return '⚠️';
      case 'expiry_warning':
        return '⏰';
      case 'enrollment_reminder':
        return '📚';
      case 'course_update':
        return '🔄';
      default:
        return '🔔';
    }
  }

  String _formatTimeAgo(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);

    if (diff.inMinutes < 1) return 'Just now';
    if (diff.inMinutes < 60) return '${diff.inMinutes}m ago';
    if (diff.inHours < 24) return '${diff.inHours}h ago';
    if (diff.inDays < 7) return '${diff.inDays}d ago';
    return '${date.day}/${date.month}/${date.year}';
  }

  @override
  Widget build(BuildContext context) {
    final isLoggedInFuture = widget.authService.isLoggedIn();

    return FutureBuilder<bool>(
      future: isLoggedInFuture,
      builder: (context, snapshot) {
        if (!snapshot.hasData || !snapshot.data!) {
          return const SizedBox.shrink();
        }

        return Stack(
          children: [
            // Floating Action Button
            Positioned(
              bottom: 24,
              right: 24,
              child: AnimatedScale(
                scale: _isOpen ? 1.05 : 1.0,
                duration: const Duration(milliseconds: 150),
                curve: Curves.easeOut,
                child: Material(
                  color: Colors.transparent,
                  child: InkWell(
                    onTap: _toggleDropdown,
                    borderRadius: BorderRadius.circular(28),
                    child: Ink(
                      width: 56,
                      height: 56,
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.9),
                        borderRadius: BorderRadius.circular(28),
                        boxShadow: [
                          BoxShadow(
                            color: Colors.black.withOpacity(0.1),
                            blurRadius: 20,
                            offset: const Offset(0, 4),
                          ),
                          BoxShadow(
                            color: Colors.black.withOpacity(0.05),
                            blurRadius: 6,
                            offset: const Offset(0, 1),
                          ),
                        ],
                        border: Border.all(color: Colors.grey.shade100),
                      ),
                      child: Stack(
                        alignment: Alignment.center,
                        children: [
                          const Icon(
                            Icons.notifications_none_rounded,
                            size: 26,
                            color: Color(0xFF1E88E5),
                          ),
                          if (_unreadCount > 0)
                            Positioned(
                              top: 4,
                              right: 4,
                              child: AnimatedContainer(
                                duration: const Duration(milliseconds: 200),
                                constraints: const BoxConstraints(minWidth: 20, minHeight: 20),
                                padding: const EdgeInsets.symmetric(horizontal: 6),
                                decoration: BoxDecoration(
                                  color: const Color(0xFFEF4444),
                                  borderRadius: BorderRadius.circular(10),
                                  boxShadow: [
                                    BoxShadow(
                                      color: const Color(0xFFEF4444).withOpacity(0.4),
                                      blurRadius: 8,
                                      offset: const Offset(0, 2),
                                    ),
                                  ],
                                ),
                                alignment: Alignment.center,
                                child: Text(
                                  _unreadCount > 9 ? '9+' : '$_unreadCount',
                                  style: const TextStyle(
                                    fontSize: 11,
                                    fontWeight: FontWeight.w600,
                                    color: Colors.white,
                                  ),
                                ),
                              ),
                            ),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            ),

            // Backdrop for closing on outside tap
            if (_isOpen)
              Positioned.fill(
                child: GestureDetector(
                  onTap: _toggleDropdown,
                  behavior: HitTestBehavior.translucent,
                  child: Container(color: Colors.transparent),
                ),
              ),
          ],
        );
      },
    );
  }
}

/// Floating Notification Button - Simpler version for direct use in Scaffold
class FloatingNotificationButton extends StatelessWidget {
  final int unreadCount;
  final VoidCallback onPressed;
  final bool isOpen;
  final VoidCallback onToggle;

  const FloatingNotificationButton({
    super.key,
    required this.unreadCount,
    required this.onPressed,
    required this.isOpen,
    required this.onToggle,
  });

  @override
  Widget build(BuildContext context) {
    return FloatingActionButton(
      onPressed: onToggle,
      backgroundColor: Colors.white,
      elevation: 4,
      shape: const CircleBorder(),
      child: Stack(
        alignment: Alignment.center,
        children: [
          const Icon(Icons.notifications_none_rounded, size: 26, color: Color(0xFF1E88E5)),
          if (unreadCount > 0)
            Positioned(
              top: 4,
              right: 4,
              child: Container(
                constraints: const BoxConstraints(minWidth: 20, minHeight: 20),
                padding: const EdgeInsets.symmetric(horizontal: 6),
                decoration: BoxDecoration(
                  color: const Color(0xFFEF4444),
                  borderRadius: BorderRadius.circular(10),
                ),
                alignment: Alignment.center,
                child: Text(
                  unreadCount > 9 ? '9+' : '$unreadCount',
                  style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w600, color: Colors.white),
                ),
              ),
            ),
        ],
      ),
    );
  }
}