import 'package:flutter/material.dart';

class StatusBadge extends StatelessWidget {
  final String status;
  final double size;

  const StatusBadge({
    super.key,
    required this.status,
    this.size = 24.0,
  });

  Color get _color {
    return switch (status.toLowerCase()) {
      'hired' => Colors.green,
      'shortlisted' => Colors.orange,
      'reviewed' => Colors.blue,
      'applied' => Colors.grey,
      'rejected' => Colors.red,
      _ => Colors.grey,
    };
  }

  IconData get _icon {
    return switch (status.toLowerCase()) {
      'hired' => Icons.check_circle,
      'shortlisted' => Icons.star,
      'reviewed' => Icons.visibility,
      'applied' => Icons.send,
      'rejected' => Icons.cancel,
      _ => Icons.help_outline,
    };
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    return Container(
      width: size * 1.8,
      height: size,
      decoration: BoxDecoration(
        color: _color.withOpacity(0.2),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: _color, width: 1.5),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(_icon, size: size * 0.6, color: _color),
          const SizedBox(width: 2),
          Flexible(
            child: Text(
              status.toUpperCase(),
              style: TextStyle(
                color: _color,
                fontWeight: FontWeight.bold,
                fontSize: size * 0.45,
              ),
              overflow: TextOverflow.ellipsis,
            ),
          ),
        ],
      ),
    );
  }
}

