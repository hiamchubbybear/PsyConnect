import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/mood.dart';
import 'package:PsyConnect/services/profile_service/mood.dart';
import 'package:flutter/material.dart';

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
  showGeneralDialog(
    context: context,
    barrierDismissible: true,
    barrierLabel: "Post Dialog",
    transitionDuration: const Duration(milliseconds: 100),
    pageBuilder: (context, animation, secondaryAnimation) {
      return Center(
        child: Material(
          borderRadius: BorderRadius.circular(20),
          color: Colors.white,
          child: StatefulBuilder(
            builder: (context, setState) {
              return Container(
                width: MediaQuery.of(context).size.width * 1.3,
                height: 270,
                padding:
                    const EdgeInsets.symmetric(vertical: 20, horizontal: 10),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text("Share your mind", style: quickSand12Font),
                        Container(
                          child: DropdownButton<String>(
                            value: selectedItem,
                            items: items
                                .map((item) => DropdownMenuItem<String>(
                                    value: item,
                                    child: Text(item, style: quickSand12Font)))
                                .toList(),
                            onChanged: (String? value) {
                              setState(() {
                                selectedItem = value!;
                                visibility = value!;
                              });
                            },
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 10),
                    TextField(
                      maxLines: 3,
                      maxLength: 200,
                      decoration: InputDecoration(
                        hintText:
                            "What are you thinking about? (200 characters)",
                        border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(10)),
                      ),
                      onChanged: (value) => {
                        mood = value,
                      },
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        TextButton(
                          child: const Text("Cancel"),
                          onPressed: () => Navigator.of(context).pop(),
                        ),
                        const SizedBox(width: 8),
                        ElevatedButton(
                          child: const Text("Post"),
                          onPressed: () {
                            if (mood.isEmpty) {
                              ToastService.showToast(
                                  context: context,
                                  message: "Please type what you are thinking",
                                  title: "Warning",
                                  type: ToastType.warning);
                            } else if (visibility.isEmpty) {
                              ToastService.showToast(
                                  context: context,
                                  message: "Please choose post privacy",
                                  title: "Warning",
                                  type: ToastType.warning);
                            } else {
                              MoodModel moodModel = MoodModel(
                                  mood: mood,
                                  moodDescription: moodDescription,
                                  visibility: visibility);
                              moodService.createMoodPost(
                                  mood: moodModel, context: context);
                            }
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

class MoodNoteBubbleWithSmoke extends StatelessWidget {
  final String text;

  const MoodNoteBubbleWithSmoke({super.key, required this.text});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        // Smoke effect
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircleAvatar(radius: 2, backgroundColor: Colors.white),
            const SizedBox(width: 4),
            CircleAvatar(radius: 4, backgroundColor: Colors.white),
            const SizedBox(width: 4),
            CircleAvatar(radius: 6, backgroundColor: Colors.white),
          ],
        ),
        const SizedBox(height: 4),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(10),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.1),
                blurRadius: 4,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: Text(
            text,
            style: const TextStyle(fontSize: 12, color: Colors.black),
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }
}
