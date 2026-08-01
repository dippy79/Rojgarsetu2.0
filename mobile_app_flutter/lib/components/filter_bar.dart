import 'package:flutter/material.dart';

class FilterOption {
  final String label;
  final String value;

  const FilterOption(this.label, this.value);
}

class FilterBar extends StatefulWidget {
  final List<FilterOption> options;
  final Function(String? selectedValue) onSelected;
  final String? initialValue;

  const FilterBar({
    super.key,
    required this.options,
    required this.onSelected,
    this.initialValue,
  });

  @override
  State<FilterBar> createState() => _FilterBarState();
}

class _FilterBarState extends State<FilterBar> {
  String? _selectedValue;

  @override
  void initState() {
    super.initState();
    _selectedValue = widget.initialValue;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      height: 40,
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        itemCount: widget.options.length + 1,
        separatorBuilder: (_, __) => const SizedBox(width: 8),
        itemBuilder: (context, index) {
          if (index == 0) {
            return GestureDetector(
              onTap: () {
                setState(() => _selectedValue = null);
                widget.onSelected(null);
              },
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                decoration: BoxDecoration(
                  color: _selectedValue == null 
                    ? theme.colorScheme. primary.withValues(alpha:0.1)
                    : Colors.grey[100],
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(
                    color: _selectedValue == null
                      ? theme.colorScheme.primary
                      : Colors.grey,
                  ),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.filter_list, size: 16, color: Colors.grey[600]),
                    const SizedBox(width: 4),
                    Text('All', style: TextStyle(color: Colors.grey[600])),
                  ],
                ),
              ),
            );
          }

          final option = widget.options[index - 1];
          final isSelected = _selectedValue == option.value;
          
          return GestureDetector(
            onTap: () {
              setState(() => _selectedValue = option.value);
              widget.onSelected(option.value);
            },
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
              decoration: BoxDecoration(
                color: isSelected 
                  ? theme.colorScheme.primary.withValues(alpha: 0.1)
                  : Colors.grey[100],
                borderRadius: BorderRadius.circular(20),
                border: Border.all(
                  color: isSelected 
                    ? theme.colorScheme.primary 
                    : Colors.grey,
                ),
              ),
              child: Text(
                option.label,
                style: TextStyle(
                  fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                  color: isSelected ? theme.colorScheme.primary : Colors.grey[700],
                  fontSize: 14,
                ),
              ),
            ),
          );
        },
      ),
    );
  }
}

