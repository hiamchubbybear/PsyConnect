import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/profile_mood.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/services/profile_service/mood.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class MoodNoteBubbleWithSmoke extends StatelessWidget {
  final String text;

  const MoodNoteBubbleWithSmoke({super.key, required this.text});

  @override
  Widget build(BuildContext context) {
    final themeProvider = Provider.of<ThemeProvider>(context);
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircleAvatar(
              radius: 2,
              backgroundColor: themeProvider.isDarkMode
                  ? Colors.grey.shade500
                  : Colors.grey.shade300,
            ),
            const SizedBox(width: 4),
            CircleAvatar(
              radius: 4,
              backgroundColor: themeProvider.isDarkMode
                  ? Colors.grey.shade400
                  : Colors.grey.shade400,
            ),
            const SizedBox(width: 4),
            CircleAvatar(
              radius: 6,
              backgroundColor: themeProvider.isDarkMode
                  ? Colors.grey.shade300
                  : Colors.grey.shade500,
            ),
          ],
        ),
        const SizedBox(height: 6),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            color:
                themeProvider.isDarkMode ? Colors.grey.shade800 : Colors.white,
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: themeProvider.isDarkMode
                  ? Colors.grey.shade600
                  : Colors.grey.shade200,
              width: 1,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black
                    .withOpacity(themeProvider.isDarkMode ? 0.3 : 0.1),
                blurRadius: 6,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: Text(
            text,
            style: quickSand12FontMoodCreate,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }
}

void _showDeleteConfirmation(
    BuildContext context, ProfileMoodModel mood, VoidCallback? onMoodCreated) {
  showDialog(
    context: context,
    builder: (BuildContext context) {
      return AlertDialog(
        backgroundColor:
            themeProvider.isDarkMode ? Colors.grey.shade900 : Colors.white,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(15),
        ),
        title: Text(
          "Delete Mood",
          style: kSubHeadingStyle.copyWith(
            fontWeight: FontWeight.w600,
            fontSize: 18,
          ),
        ),
        content: Text(
          "Are you sure you want to delete this mood? This action cannot be undone.",
          style: kSubHeadingStyle.copyWith(
            fontSize: 14,
            color: themeProvider.isDarkMode
                ? Colors.white70
                : Colors.grey.shade600,
          ),
        ),
        actions: [
          TextButton(
            child: Text(
              "Cancel",
              style: kSubHeadingStyle.copyWith(
                color: themeProvider.isDarkMode
                    ? Colors.white70
                    : Colors.grey.shade600,
                fontWeight: FontWeight.w500,
              ),
            ),
            onPressed: () => Navigator.of(context).pop(),
          ),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.grey[300],
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
            child: Text(
              "Delete",
              style: kSubHeadingStyle.copyWith(
                color: Colors.white,
                fontWeight: FontWeight.w600,
              ),
            ),
            onPressed: () async {
              Navigator.of(context).pop();
              try {
                await MoodService().deleteMoodPost(
                  moodId: mood.moodId,
                  context: context,
                );

                if (onMoodCreated != null) {
                  onMoodCreated();
                }

                ToastService.showToast(
                    context: context,
                    message: "Mood deleted successfully",
                    title: "Success",
                    type: ToastType.success);
              } catch (e) {
                ToastService.showToast(
                    context: context,
                    message: "Failed to delete mood",
                    title: "Error",
                    type: ToastType.error);
              }
            },
          ),
        ],
      );
    },
  );
}
