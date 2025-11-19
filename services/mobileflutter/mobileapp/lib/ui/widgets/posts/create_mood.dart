import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/mood.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/services/profile_service/mood.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class MoodWidget extends StatefulWidget {
  @override
  State<MoodWidget> createState() => _MoodWidgetState();
}

String mood = "";
String moodDescription = "";
String visibility = "";
String selectedItem = items[0];
List<String> items = [
  "Private",
  "Public",
  "Friends only",
];
MoodService moodService = MoodService();

class _MoodWidgetState extends State<MoodWidget> {
  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () => showPostDialog(context),
      child: Container(
        width: 20,
        height: 20,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          color: Colors.green[300],
          border: Border.all(color: Colors.white, width: 2),
        ),
        child: const Icon(
          Icons.add,
          size: 14,
          color: Colors.white,
        ),
      ),
    );
  }
}

void showPostDialog(BuildContext context) {
  final themeProvider = Provider.of<ThemeProvider>(context, listen: false);
  final TextEditingController moodController = TextEditingController();
  String visibility = selectedItem;

  showGeneralDialog(
    context: context,
    barrierDismissible: true,
    barrierLabel: "Post Dialog",
    transitionDuration: const Duration(milliseconds: 200),
    pageBuilder: (context, animation, secondaryAnimation) {
      return Center(
        child: Material(
          borderRadius: BorderRadius.circular(20),
          color: themeProvider.isDarkMode ? Colors.grey.shade900 : Colors.white,
          elevation: 8,
          child: StatefulBuilder(
            builder: (context, setState) {
              return Container(
                width: MediaQuery.of(context).size.width * 0.9,
                height: 320,
                padding: const EdgeInsets.all(24),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    /// Title & Dropdown
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          "Share your mind",
                          style: kSubHeadingStyle.copyWith(
                            fontWeight: FontWeight.w600,
                            fontSize: 18,
                          ),
                        ),
                        Container(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 12, vertical: 4),
                          decoration: BoxDecoration(
                            color: themeProvider.isDarkMode
                                ? Colors.grey.shade700
                                : Colors.grey.shade100,
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: DropdownButton<String>(
                            value: visibility,
                            underline: const SizedBox(),
                            icon: Icon(
                              Icons.keyboard_arrow_down,
                              color: themeProvider.isDarkMode
                                  ? Colors.white70
                                  : Colors.grey.shade600,
                            ),
                            items: items
                                .map((item) => DropdownMenuItem<String>(
                                      value: item,
                                      child: Text(
                                        item,
                                        style: kSubHeadingStyle.copyWith(
                                          fontSize: 13,
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                    ))
                                .toList(),
                            onChanged: (String? value) {
                              setState(() {
                                visibility = value!;
                              });
                            },
                          ),
                        ),
                      ],
                    ),

                    const SizedBox(height: 20),

                    /// TextField
                    Expanded(
                      child: TextField(
                        controller: moodController,
                        maxLines: null,
                        expands: true,
                        maxLength: 200,
                        style: kSubHeadingStyle.copyWith(
                          fontSize: 15,
                          height: 1.4,
                        ),
                        decoration: InputDecoration(
                          hintText: "What are you thinking about?",
                          hintStyle: kSubHintStyle,
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(12),
                            borderSide: BorderSide(
                              color: themeProvider.isDarkMode
                                  ? Colors.grey.shade600
                                  : Colors.grey.shade300,
                            ),
                          ),
                          enabledBorder: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(12),
                            borderSide: BorderSide(
                              color: themeProvider.isDarkMode
                                  ? Colors.grey.shade600
                                  : Colors.grey.shade300,
                            ),
                          ),
                          focusedBorder: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(12),
                            borderSide: BorderSide(
                              color: acceptColor,
                              width: 2,
                            ),
                          ),
                          fillColor: themeProvider.isDarkMode
                              ? Colors.grey.shade800
                              : Colors.grey.shade50,
                          filled: true,
                          contentPadding: const EdgeInsets.all(16),
                        ),
                      ),
                    ),

                    const SizedBox(height: 20),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        TextButton(
                          style: TextButton.styleFrom(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 20, vertical: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(10),
                            ),
                          ),
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
                        const SizedBox(width: 12),
                        ElevatedButton(
                          style: ElevatedButton.styleFrom(
                            backgroundColor: acceptColor,
                            padding: const EdgeInsets.symmetric(
                                horizontal: 24, vertical: 12),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(10),
                            ),
                            elevation: 2,
                          ),
                          child: Text(
                            "Post",
                            style: kSubHeadingStyle.copyWith(
                              color: Colors.white,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                          onPressed: () {
                            final text = moodController.text.trim();
                            if (text.isEmpty) {
                              ToastService.showToast(
                                  context: context,
                                  message: "Please type what you are thinking",
                                  title: "Warning",
                                  type: ToastType.warning);
                              return;
                            }

                            if (visibility.isEmpty) {
                              ToastService.showToast(
                                  context: context,
                                  message: "Please choose post privacy",
                                  title: "Warning",
                                  type: ToastType.warning);
                              return;
                            }

                            MoodModel moodModel = MoodModel(
                              mood: text,
                              moodDescription: text,
                              visibility: visibility,
                            );

                            moodService.createMoodPost(
                                mood: moodModel, context: context);
                          },
                        ),
                      ],
                    )
                  ],
                ),
              );
            },
          ),
        ),
      );
    },
    transitionBuilder: (context, anim1, anim2, child) {
      return FadeTransition(
        opacity: anim1,
        child: ScaleTransition(scale: anim1, child: child),
      );
    },
  );
}
