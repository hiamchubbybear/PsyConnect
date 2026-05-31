import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:PsyConnect/models/group.dart';
import 'package:PsyConnect/provider/group_provider.dart';
import 'package:PsyConnect/ui/widgets/common/custom_button.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';

class GroupsListScreen extends StatefulWidget {
  const GroupsListScreen({super.key});

  @override
  State<GroupsListScreen> createState() => _GroupsListScreenState();
}

class _GroupsListScreenState extends State<GroupsListScreen> {
  final TextEditingController _searchController = TextEditingController();
  final List<String> _categories = ["", "CBT", "Self-Care", "Anxiety", "Stress", "Depression", "Mindfulness"];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<GroupProvider>(context, listen: false).fetchGroups();
    });
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _refreshGroups() async {
    await Provider.of<GroupProvider>(context, listen: false).fetchGroups();
  }

  void _showCreateGroupDialog(BuildContext context, GroupProvider provider) {
    final formKey = GlobalKey<FormState>();
    String name = "";
    String description = "";
    String selectedCategory = _categories[1]; // Default to first actual category

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) {
        final theme = Theme.of(context);
        final isDark = theme.brightness == Brightness.dark;

        return StatefulBuilder(
          builder: (context, setModalState) {
            return Container(
              decoration: BoxDecoration(
                color: theme.scaffoldBackgroundColor,
                borderRadius: const BorderRadius.vertical(top: Radius.circular(24)),
              ),
              padding: EdgeInsets.only(
                left: 20,
                right: 20,
                top: 24,
                bottom: MediaQuery.of(context).viewInsets.bottom + 24,
              ),
              child: Form(
                key: formKey,
                child: SingleChildScrollView(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Center(
                        child: Container(
                          width: 40,
                          height: 4,
                          decoration: BoxDecoration(
                            color: isDark ? Colors.grey[700] : Colors.grey[300],
                            borderRadius: BorderRadius.circular(2),
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),
                      Text(
                        "Create Community Group",
                        style: GoogleFonts.quicksand(
                          fontSize: 22,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),
                      TextFormField(
                        decoration: InputDecoration(
                          labelText: "Group Name",
                          labelStyle: GoogleFonts.quicksand(),
                          border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                        ),
                        validator: (value) {
                          if (value == null || value.trim().isEmpty) {
                            return "Group name cannot be empty";
                          }
                          return null;
                        },
                        onSaved: (value) => name = value?.trim() ?? "",
                      ),
                      const SizedBox(height: 16),
                      TextFormField(
                        maxLines: 3,
                        decoration: InputDecoration(
                          labelText: "Description",
                          labelStyle: GoogleFonts.quicksand(),
                          border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                        ),
                        validator: (value) {
                          if (value == null || value.trim().isEmpty) {
                            return "Description cannot be empty";
                          }
                          return null;
                        },
                        onSaved: (value) => description = value?.trim() ?? "",
                      ),
                      const SizedBox(height: 16),
                      Text(
                        "Category",
                        style: GoogleFonts.quicksand(
                          fontSize: 14,
                          fontWeight: FontWeight.bold,
                          color: Colors.grey[600],
                        ),
                      ),
                      const SizedBox(height: 8),
                      DropdownButtonFormField<String>(
                        value: selectedCategory,
                        items: _categories
                            .where((cat) => cat.isNotEmpty)
                            .map((cat) => DropdownMenuItem(
                                  value: cat,
                                  child: Text(cat, style: GoogleFonts.quicksand()),
                                ))
                            .toList(),
                        onChanged: (val) {
                          if (val != null) {
                            setModalState(() {
                              selectedCategory = val;
                            });
                          }
                        },
                        decoration: InputDecoration(
                          border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                          contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                        ),
                      ),
                      const SizedBox(height: 24),
                      CustomButton(
                        onPressed: () async {
                          if (formKey.currentState?.validate() ?? false) {
                            formKey.currentState?.save();
                            Navigator.pop(context); // Close bottom sheet
                            
                            final success = await provider.createGroup(
                              name: name,
                              description: description,
                              category: selectedCategory,
                            );

                            if (success && mounted) {
                              ToastService.showToast(
                                context: context,
                                message: "Community Group '$name' created successfully!",
                                title: "Group Created",
                                type: ToastType.success,
                              );
                            } else if (mounted) {
                              ToastService.showToast(
                                context: context,
                                message: provider.errorMessage ?? "Failed to create group",
                                title: "Error",
                                type: ToastType.error,
                              );
                            }
                          }
                        },
                        text: "Create Group",
                        height: 50,
                      ),
                    ],
                  ),
                ),
              ),
            );
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;
    final primaryColor = isDark ? Colors.blue[300]! : Colors.blue;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        title: Text(
          "Community Forums",
          style: GoogleFonts.quicksand(fontWeight: FontWeight.bold),
        ),
        centerTitle: true,
        actions: [
          Consumer<GroupProvider>(
            builder: (context, provider, _) => IconButton(
              icon: const Icon(CupertinoIcons.add_circled_solid),
              color: primaryColor,
              iconSize: 28,
              onPressed: () => _showCreateGroupDialog(context, provider),
            ),
          )
        ],
      ),
      body: Consumer<GroupProvider>(
        builder: (context, provider, child) {
          return Column(
            children: [
              _buildSearchBar(theme, isDark, provider),
              _buildCategoriesRow(theme, isDark, primaryColor, provider),
              Expanded(
                child: RefreshIndicator(
                  onRefresh: _refreshGroups,
                  child: provider.isLoading && provider.groups.isEmpty
                      ? const Center(child: CircularProgressIndicator())
                      : provider.errorMessage != null && provider.groups.isEmpty
                          ? Center(
                              child: Padding(
                                padding: const EdgeInsets.all(24),
                                child: Text(
                                  "Error loading groups: ${provider.errorMessage}",
                                  textAlign: TextAlign.center,
                                  style: GoogleFonts.quicksand(color: Colors.redAccent),
                                ),
                              ),
                            )
                          : provider.groups.isEmpty
                              ? _buildEmptyState()
                              : ListView.builder(
                                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                                  itemCount: provider.groups.length,
                                  itemBuilder: (context, index) {
                                    final group = provider.groups[index];
                                    return _buildGroupCard(group, provider, theme, isDark, primaryColor);
                                  },
                                ),
                ),
              )
            ],
          );
        },
      ),
    );
  }

  Widget _buildSearchBar(ThemeData theme, bool isDark, GroupProvider provider) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: TextField(
        controller: _searchController,
        style: GoogleFonts.quicksand(fontSize: 16),
        onSubmitted: (value) {
          provider.setSearchQuery(value.trim());
        },
        decoration: InputDecoration(
          hintText: "Search forums or topics...",
          hintStyle: GoogleFonts.quicksand(color: Colors.grey[500]),
          prefixIcon: Icon(CupertinoIcons.search, color: Colors.grey[500]),
          suffixIcon: _searchController.text.isNotEmpty
              ? IconButton(
                  icon: const Icon(CupertinoIcons.clear_circled),
                  onPressed: () {
                    _searchController.clear();
                    provider.setSearchQuery("");
                  },
                )
              : null,
          filled: true,
          fillColor: isDark ? Colors.grey[850] : Colors.grey[100],
          contentPadding: const EdgeInsets.symmetric(vertical: 0, horizontal: 16),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(30),
            borderSide: BorderSide.none,
          ),
        ),
      ),
    );
  }

  Widget _buildCategoriesRow(
    ThemeData theme,
    bool isDark,
    Color primaryColor,
    GroupProvider provider,
  ) {
    return Container(
      height: 48,
      margin: const EdgeInsets.only(bottom: 8),
      child: ListView.builder(
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: 12),
        itemCount: _categories.length,
        itemBuilder: (context, index) {
          final cat = _categories[index];
          final label = cat.isEmpty ? "All" : cat;
          final isSelected = provider.selectedCategory == cat;

          return Padding(
            padding: const EdgeInsets.symmetric(horizontal: 4),
            child: ChoiceChip(
              label: Text(label),
              selected: isSelected,
              onSelected: (selected) {
                if (selected) {
                  provider.setCategory(cat);
                }
              },
              selectedColor: primaryColor.withOpacity(0.2),
              backgroundColor: isDark ? Colors.grey[800] : Colors.grey[200],
              labelStyle: GoogleFonts.quicksand(
                color: isSelected
                    ? primaryColor
                    : (isDark ? Colors.white70 : Colors.black87),
                fontWeight: isSelected ? FontWeight.bold : FontWeight.w600,
              ),
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
            ),
          );
        },
      ),
    );
  }

  Widget _buildGroupCard(
    Group group,
    GroupProvider provider,
    ThemeData theme,
    bool isDark,
    Color primaryColor,
  ) {
    final isJoined = provider.joinedGroupIds.contains(group.id);

    // Get a beautiful category color/gradient indicator
    Color categoryBg;
    switch (group.category.toLowerCase()) {
      case 'cbt':
        categoryBg = Colors.purpleAccent;
        break;
      case 'self-care':
        categoryBg = Colors.greenAccent[700]!;
        break;
      case 'anxiety':
        categoryBg = Colors.orangeAccent;
        break;
      case 'stress':
        categoryBg = Colors.redAccent;
        break;
      case 'depression':
        categoryBg = Colors.indigoAccent;
        break;
      case 'mindfulness':
        categoryBg = Colors.teal;
        break;
      default:
        categoryBg = primaryColor;
    }

    return Card(
      elevation: 2,
      margin: const EdgeInsets.only(bottom: 14),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: categoryBg.withOpacity(0.12),
                    borderRadius: BorderRadius.circular(20),
                    border: Border.all(color: categoryBg.withOpacity(0.4), width: 1),
                  ),
                  child: Text(
                    group.category.toUpperCase(),
                    style: GoogleFonts.quicksand(
                      fontSize: 11,
                      fontWeight: FontWeight.bold,
                      color: categoryBg,
                    ),
                  ),
                ),
                Row(
                  children: [
                    Icon(CupertinoIcons.person_2_fill, size: 16, color: Colors.grey[500]),
                    const SizedBox(width: 4),
                    Text(
                      "${group.memberCount} members",
                      style: GoogleFonts.quicksand(
                        fontSize: 12,
                        color: Colors.grey[500],
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 12),
            Text(
              group.name,
              style: GoogleFonts.quicksand(
                fontSize: 18,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 6),
            Text(
              group.description,
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
              style: GoogleFonts.quicksand(
                fontSize: 14,
                color: isDark ? Colors.grey[400] : Colors.grey[600],
              ),
            ),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                if (isJoined) ...[
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
                    decoration: BoxDecoration(
                      color: Colors.green.withOpacity(0.12),
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: Colors.green.withOpacity(0.3)),
                    ),
                    child: Row(
                      children: [
                        const Icon(CupertinoIcons.checkmark_alt_circle_fill, color: Colors.green, size: 18),
                        const SizedBox(width: 6),
                        Text(
                          "Joined Room",
                          style: GoogleFonts.quicksand(
                            color: Colors.green,
                            fontWeight: FontWeight.bold,
                            fontSize: 13,
                          ),
                        ),
                      ],
                    ),
                  ),
                ] else ...[
                  ElevatedButton(
                    onPressed: () async {
                      final success = await provider.joinGroup(group.id);
                      if (success && mounted) {
                        ToastService.showToast(
                          context: context,
                          message: "You have successfully joined '${group.name}' room!",
                          title: "Joined Group",
                          type: ToastType.success,
                        );
                      } else if (mounted) {
                        ToastService.showToast(
                          context: context,
                          message: provider.errorMessage ?? "Failed to join community forum",
                          title: "Error",
                          type: ToastType.error,
                        );
                      }
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: primaryColor,
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    ),
                    child: Text(
                      "Join Group",
                      style: GoogleFonts.quicksand(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ],
            )
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(CupertinoIcons.chat_bubble_2, size: 64, color: Colors.grey[400]),
            const SizedBox(height: 16),
            Text(
              "No Forums Found",
              style: GoogleFonts.quicksand(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              "There are no community rooms matching your criteria.\nCreate a new one to start the conversation!",
              textAlign: TextAlign.center,
              style: GoogleFonts.quicksand(
                fontSize: 14,
                color: Colors.grey[500],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
