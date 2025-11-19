import 'package:PsyConnect/core/variable/variable.dart';
import 'package:flutter/material.dart';

class PostOptionsMenu extends StatelessWidget {
  final void Function(String)? onSelected;

  const PostOptionsMenu({Key? key, this.onSelected}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<String>(
      icon: Icon(
        Icons.more_vert,
        color: themeProvider.isDarkMode ? Colors.white70 : Colors.grey.shade600,
        size: 20,
      ),
      onSelected: onSelected,
      color: themeProvider.isDarkMode ? Colors.grey.shade800 : Colors.white,
      elevation: 8,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(
          color: themeProvider.isDarkMode
              ? Colors.grey.shade700
              : Colors.grey.shade200,
          width: 1,
        ),
      ),
      offset: const Offset(0, 8),
      itemBuilder: (context) => [
        PopupMenuItem<String>(
          value: 'edit',
          height: 48,
          child: Row(
            children: [
              Icon(
                Icons.edit_outlined,
                size: 18,
                color: themeProvider.isDarkMode ? Colors.white70 : Colors.grey.shade700,
              ),
              const SizedBox(width: 12),
              Text(
                'Edit',
                style: kSubHeadingStyle.copyWith(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ],
          ),
        ),
        PopupMenuItem<String>(
          value: 'delete',
          height: 48,
          child: Row(
            children: [
              Icon(
                Icons.delete_outline,
                size: 18,
                color: themeProvider.isDarkMode ? Colors.red.shade300 : Colors.red.shade600,
              ),
              const SizedBox(width: 12),
              Text(
                'Delete',
                style: kSubHeadingStyle.copyWith(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: themeProvider.isDarkMode ? Colors.red.shade300 : Colors.red.shade600,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}
