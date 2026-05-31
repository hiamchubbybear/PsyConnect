import 'package:PsyConnect/core/colors/color.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/provider/post_provider.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';

class CreatePostScreen extends StatefulWidget {
  const CreatePostScreen({super.key});

  @override
  State<CreatePostScreen> createState() => _CreatePostScreenState();
}

class _CreatePostScreenState extends State<CreatePostScreen> {
  final _formKey = GlobalKey<FormState>();
  final TextEditingController _titleController = TextEditingController();
  final TextEditingController _contentController = TextEditingController();
  final TextEditingController _tagController = TextEditingController();

  final List<String> _selectedCategories = [];
  final List<String> _tags = [];

  final List<String> _availableCategories = [
    'Anxiety',
    'Depression',
    'CBT',
    'Stress Management',
    'PTSD',
    'Relationship Issues',
    'Self-Care',
  ];

  bool _isSubmitting = false;

  @override
  void dispose() {
    _titleController.dispose();
    _contentController.dispose();
    _tagController.dispose();
    super.dispose();
  }

  void _addTag(String rawTag) {
    final clean = rawTag.trim().replaceAll('#', '').toLowerCase();
    if (clean.isNotEmpty && !_tags.contains(clean)) {
      setState(() {
        _tags.add(clean);
      });
      _tagController.clear();
    }
  }

  void _removeTag(String tag) {
    setState(() {
      _tags.remove(tag);
    });
  }

  void _toggleCategory(String category) {
    setState(() {
      if (_selectedCategories.contains(category)) {
        _selectedCategories.remove(category);
      } else {
        _selectedCategories.add(category);
      }
    });
  }

  Future<void> _publish(PostProvider postProvider) async {
    if (!_formKey.currentState!.validate()) return;
    if (_contentController.text.trim().isEmpty) return;

    setState(() {
      _isSubmitting = true;
    });

    final success = await postProvider.publishPost(
      title: _titleController.text.trim(),
      content: _contentController.text.trim(),
      tags: _tags,
      categories: _selectedCategories,
    );

    if (success && mounted) {
      ToastService.showToast(
        context: context,
        message: "Your new article has been published successfully.",
        title: "Article Published",
        type: ToastType.success,
      );
      Navigator.pop(context);
    } else if (mounted) {
      ToastService.showToast(
        context: context,
        message: postProvider.errorMessage ?? "Failed to publish post.",
        title: "Error",
        type: ToastType.error,
      );
    }

    setState(() {
      _isSubmitting = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final postProvider = Provider.of<PostProvider>(context);
    final isDark = Provider.of<ThemeProvider>(context).isDarkMode;
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: isDark ? Colors.black : const Color(0xFFF7F9FC),
      appBar: AppBar(
        title: Text(
          "Compose Article",
          style: GoogleFonts.quicksand(
            fontWeight: FontWeight.bold,
            fontSize: 18,
            color: theme.textTheme.titleLarge?.color,
          ),
        ),
        elevation: 0.5,
        backgroundColor: isDark ? theme.cardColor : Colors.white,
        leading: IconButton(
          icon: Icon(Icons.close, color: theme.iconTheme.color),
          onPressed: () => Navigator.pop(context),
        ),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 12, top: 10, bottom: 10),
            child: _isSubmitting
                ? const SizedBox(
                    width: 28,
                    height: 28,
                    child: CircularProgressIndicator(strokeWidth: 2.5),
                  )
                : ElevatedButton(
                    onPressed: () => _publish(postProvider),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: isDark ? Colors.blueGrey[800] : primaryColor,
                      foregroundColor: Colors.white,
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      elevation: 0,
                    ),
                    child: Text(
                      "Publish",
                      style: GoogleFonts.quicksand(fontWeight: FontWeight.bold, fontSize: 13),
                    ),
                  ),
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 1. Categories Multi-Select Chips
              Text(
                "Choose Categories",
                style: GoogleFonts.quicksand(
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white70 : Colors.black87,
                ),
              ),
              const SizedBox(height: 10),
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: _availableCategories.map((category) {
                  final isSelected = _selectedCategories.contains(category);
                  return GestureDetector(
                    onTap: () => _toggleCategory(category),
                    child: AnimatedContainer(
                      duration: const Duration(milliseconds: 150),
                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      decoration: BoxDecoration(
                        color: isSelected
                            ? (isDark ? Colors.blueGrey[900] : primaryColor.withOpacity(0.08))
                            : (isDark ? theme.cardColor : Colors.white),
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(
                          color: isSelected
                              ? (isDark ? Colors.blue[300]! : primaryColor)
                              : (isDark ? Colors.grey[800]! : Colors.grey[300]!),
                          width: 1,
                        ),
                      ),
                      child: Text(
                        category,
                        style: GoogleFonts.quicksand(
                          fontSize: 12,
                          fontWeight: isSelected ? FontWeight.bold : FontWeight.w500,
                          color: isSelected
                              ? (isDark ? Colors.blue[300] : primaryColor)
                              : (isDark ? Colors.grey[300] : Colors.grey[700]),
                        ),
                      ),
                    ),
                  );
                }).toList(),
              ),

              const SizedBox(height: 24),

              // 2. Title Field
              Text(
                "Article Title",
                style: GoogleFonts.quicksand(
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white70 : Colors.black87,
                ),
              ),
              const SizedBox(height: 8),
              TextFormField(
                controller: _titleController,
                style: GoogleFonts.quicksand(color: isDark ? Colors.white : Colors.black87, fontSize: 15),
                decoration: InputDecoration(
                  hintText: "Enter a compelling title...",
                  hintStyle: GoogleFonts.quicksand(color: Colors.grey, fontSize: 13),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                    borderSide: BorderSide(color: isDark ? Colors.grey[800]! : Colors.grey[300]!),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                    borderSide: BorderSide(color: isDark ? Colors.grey[800]! : Colors.grey[300]!),
                  ),
                  fillColor: isDark ? theme.cardColor : Colors.white,
                  filled: true,
                  contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                ),
                validator: (value) => value == null || value.trim().isEmpty ? "Title is required" : null,
              ),

              const SizedBox(height: 24),

              // 3. Post Content Field
              Text(
                "Article Content",
                style: GoogleFonts.quicksand(
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white70 : Colors.black87,
                ),
              ),
              const SizedBox(height: 8),
              TextFormField(
                controller: _contentController,
                maxLines: 8,
                style: GoogleFonts.quicksand(color: isDark ? Colors.white : Colors.black87, fontSize: 14, height: 1.4),
                decoration: InputDecoration(
                  hintText: "Write down your thoughts, resources, or advice...",
                  hintStyle: GoogleFonts.quicksand(color: Colors.grey, fontSize: 13),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                    borderSide: BorderSide(color: isDark ? Colors.grey[800]! : Colors.grey[300]!),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                    borderSide: BorderSide(color: isDark ? Colors.grey[800]! : Colors.grey[300]!),
                  ),
                  fillColor: isDark ? theme.cardColor : Colors.white,
                  filled: true,
                  contentPadding: const EdgeInsets.all(16),
                ),
                validator: (value) => value == null || value.trim().isEmpty ? "Content cannot be empty" : null,
              ),

              const SizedBox(height: 24),

              // 4. Custom Tags Inputs
              Text(
                "Tags",
                style: GoogleFonts.quicksand(
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                  color: isDark ? Colors.white70 : Colors.black87,
                ),
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _tagController,
                      style: GoogleFonts.quicksand(color: isDark ? Colors.white : Colors.black87, fontSize: 13),
                      decoration: InputDecoration(
                        hintText: "Add tags (e.g. selflove, mindfulness)...",
                        hintStyle: GoogleFonts.quicksand(color: Colors.grey, fontSize: 13),
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                          borderSide: BorderSide(color: isDark ? Colors.grey[800]! : Colors.grey[300]!),
                        ),
                        enabledBorder: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8),
                          borderSide: BorderSide(color: isDark ? Colors.grey[800]! : Colors.grey[300]!),
                        ),
                        fillColor: isDark ? theme.cardColor : Colors.white,
                        filled: true,
                        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                      ),
                      onSubmitted: _addTag,
                    ),
                  ),
                  const SizedBox(width: 8),
                  Container(
                    decoration: BoxDecoration(
                      color: isDark ? Colors.blueGrey[800] : primaryColor.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: IconButton(
                      icon: Icon(Icons.add, color: isDark ? Colors.blue[300] : primaryColor),
                      onPressed: () => _addTag(_tagController.text),
                    ),
                  ),
                ],
              ),

              if (_tags.isNotEmpty) ...[
                const SizedBox(height: 12),
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: _tags.map((tag) {
                    return InputChip(
                      label: Text("#$tag", style: GoogleFonts.quicksand(fontSize: 11, fontWeight: FontWeight.bold)),
                      onDeleted: () => _removeTag(tag),
                      deleteIconColor: Colors.red[300],
                      backgroundColor: isDark ? Colors.grey[900] : Colors.grey[200],
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
}
