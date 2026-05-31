import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/models/setting.dart';
import 'package:PsyConnect/services/profile_service/setting.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/ui/widgets/common/custom_button.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  Setting? setting;
  bool isLoading = true;
  bool isSaving = false;

  bool showLastSeen = false;
  bool showProfilePicture = false;
  bool showMood = false;
  bool notificationsEnabled = false;
  bool emailNotifications = false;
  bool pushNotifications = false;
  bool smsNotifications = false;
  bool twoFactorAuth = false;
  bool allowLoginAlerts = false;
  bool autoDeleteOldMoods = false;

  final List<String> languages = ['en', 'vi', 'jp'];
  final List<String> themes = ['light', 'dark'];

  @override
  void initState() {
    super.initState();
    _loadSettings();
  }

  Future<void> _resetDefault() async {
    setState(() {
      showLastSeen = false;
      showProfilePicture = false;
      showMood = false;
      notificationsEnabled = false;
      emailNotifications = false;
      pushNotifications = false;
      smsNotifications = false;
      twoFactorAuth = false;
      allowLoginAlerts = false;
      autoDeleteOldMoods = false;
    });
  }

  Future<void> _loadSettings() async {
    try {
      final settingService = SettingService();
      final fetchedSetting = await settingService.getSetting();

      if (mounted) {
        setState(() {
          setting = fetchedSetting;
          showLastSeen = fetchedSetting.showLastSeen;
          showProfilePicture = fetchedSetting.showProfilePicture;
          showMood = fetchedSetting.showMood;
          notificationsEnabled = fetchedSetting.notificationsEnabled;
          emailNotifications = fetchedSetting.emailNotifications;
          pushNotifications = fetchedSetting.pushNotifications;
          smsNotifications = fetchedSetting.smsNotifications;
          twoFactorAuth = fetchedSetting.twoFactorAuth;
          allowLoginAlerts = fetchedSetting.allowLoginAlerts;
          autoDeleteOldMoods = fetchedSetting.autoDeleteOldMoods;
          isLoading = false;
        });
      }
    } catch (e) {
      print("Error loading settings: $e");
      if (mounted) {
        setState(() => isLoading = false);
        ToastService.showToast(
            context: context,
            message: "Failed to load settings",
            title: "Error",
            type: ToastType.error);
      }
    }
  }

  Future<void> _saveSettings() async {
    if (setting == null || isSaving) return;

    setState(() => isSaving = true);

    try {
      final newSetting = setting!.copyWith(
        showLastSeen: showLastSeen,
        showProfilePicture: showProfilePicture,
        showMood: showMood,
        notificationsEnabled: notificationsEnabled,
        emailNotifications: emailNotifications,
        pushNotifications: pushNotifications,
        smsNotifications: smsNotifications,
        twoFactorAuth: twoFactorAuth,
        allowLoginAlerts: allowLoginAlerts,
        autoDeleteOldMoods: autoDeleteOldMoods,
      );
      final prefs = SharedPreferencesProvider();
      await prefs.setSetting(newSetting);
      final SettingService settingService = SettingService();
      final response = await settingService.updateSetting();

      if (mounted) {
        setState(() {
          isSaving = false;
          if (response) setting = newSetting;
        });

        ToastService.showToast(
            context: context,
            message: response
                ? "Settings saved successfully"
                : "Failed to save settings",
            title: response ? "Success" : "Error",
            type: response ? ToastType.success : ToastType.error);
      }
    } catch (e) {
      print("Error saving settings: $e");
      if (mounted) {
        setState(() => isSaving = false);
        ToastService.showToast(
            context: context,
            message: "Error saving settings",
            title: "Error",
            type: ToastType.error);
      }
    }
  }

  Widget _buildToggleRow(BuildContext context, String title, bool value, Function(bool) onChanged, {IconData? icon}) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
      child: Row(
        children: [
          if (icon != null) ...[
            Icon(
              icon,
              size: 20,
              color: isDark ? Colors.white70 : Colors.grey[600],
            ),
            const SizedBox(width: 16),
          ],
          Expanded(
            child: Text(
              title,
              style: theme.textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.w500),
            ),
          ),
          GestureDetector(
            onTap: () => onChanged(!value),
            child: AnimatedContainer(
              duration: const Duration(milliseconds: 200),
              width: 44,
              height: 24,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(12),
                color: value
                    ? Colors.green[300]
                    : (isDark ? Colors.grey[750] : Colors.grey[300]),
              ),
              child: AnimatedAlign(
                duration: const Duration(milliseconds: 200),
                alignment: value ? Alignment.centerRight : Alignment.centerLeft,
                child: Container(
                  width: 20,
                  height: 20,
                  margin: const EdgeInsets.all(2),
                  decoration: const BoxDecoration(
                    shape: BoxShape.circle,
                    color: Colors.white,
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSection(BuildContext context, String title, List<Widget> children, {IconData? icon}) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: isDark ? Colors.grey[900] : Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.04),
            blurRadius: 6,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.fromLTRB(20, 20, 20, 8),
            child: Row(
              children: [
                if (icon != null) ...[
                  Icon(
                    icon,
                    size: 20,
                    color: isDark ? Colors.blue[300] : Colors.blue,
                  ),
                  const SizedBox(width: 12),
                ],
                Text(
                  title,
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: isDark ? Colors.white : Colors.black87,
                  ),
                ),
              ],
            ),
          ),
          ...children,
        ],
      ),
    );
  }

  Widget _buildDropdownRow(
    BuildContext context,
    String title,
    String value,
    List<String> items,
    Function(String) onChanged, {
    IconData? icon,
  }) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;
    final screenWidth = MediaQuery.of(context).size.width;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
      child: Row(
        children: [
          if (icon != null) ...[
            Icon(
              icon,
              size: 20,
              color: isDark ? Colors.white70 : Colors.grey[600],
            ),
            const SizedBox(width: 16),
          ],
          Expanded(
            child: Text(
              title,
              style: theme.textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.w500),
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
            decoration: BoxDecoration(
              color: isDark ? Colors.grey[800] : Colors.grey[100],
              borderRadius: BorderRadius.circular(8),
            ),
            child: SizedBox(
              width: screenWidth * 0.20,
              child: DropdownButton<String>(
                isDense: true,
                isExpanded: true,
                value: value,
                underline: const SizedBox(),
                icon: Icon(
                  Icons.keyboard_arrow_down,
                  size: 20,
                  color: isDark ? Colors.white70 : Colors.grey[600],
                ),
                style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600),
                items: items.map((item) {
                  return DropdownMenuItem(
                    value: item,
                    child: Text(
                      item.toUpperCase(),
                      overflow: TextOverflow.ellipsis,
                    ),
                  );
                }).toList(),
                onChanged: (val) {
                  if (val != null) onChanged(val);
                },
              ),
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final themeProvider = Provider.of<ThemeProvider>(context);

    if (isLoading) {
      return Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const CircularProgressIndicator(
                strokeWidth: 2,
              ),
              const SizedBox(height: 16),
              Text(
                "Loading settings...",
                style: theme.textTheme.bodyLarge,
              ),
            ],
          ),
        ),
      );
    }

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        title: const Text("Settings"),
        centerTitle: true,
      ),
      body: ListView(
        padding: const EdgeInsets.symmetric(vertical: 8),
        children: [
          _buildSection(
            context,
            "Privacy",
            [
              _buildToggleRow(context, "Show Last Seen", showLastSeen, (val) {
                setState(() => showLastSeen = val);
              }, icon: Icons.visibility_outlined),
              _buildToggleRow(context, "Show Profile Picture", showProfilePicture, (val) {
                setState(() => showProfilePicture = val);
              }, icon: Icons.account_circle_outlined),
              _buildToggleRow(context, "Show Mood", showMood, (val) {
                setState(() => showMood = val);
              }, icon: Icons.mood_outlined),
            ],
            icon: Icons.security_outlined,
          ),
          _buildSection(
            context,
            "Notifications",
            [
              _buildToggleRow(context, "Enable Notifications", notificationsEnabled, (val) {
                setState(() => notificationsEnabled = val);
              }, icon: Icons.notifications_outlined),
              _buildToggleRow(context, "Email Notifications", emailNotifications, (val) {
                setState(() => emailNotifications = val);
              }, icon: Icons.email_outlined),
              _buildToggleRow(context, "Push Notifications", pushNotifications, (val) {
                setState(() => pushNotifications = val);
              }, icon: Icons.push_pin_outlined),
              _buildToggleRow(context, "SMS Notifications", smsNotifications, (val) {
                setState(() => smsNotifications = val);
              }, icon: Icons.sms_outlined),
            ],
            icon: Icons.notifications_active_outlined,
          ),
          _buildSection(
            context,
            "Security",
            [
              _buildToggleRow(context, "Two-Factor Authentication", twoFactorAuth, (val) {
                setState(() => twoFactorAuth = val);
              }, icon: Icons.security_outlined),
              _buildToggleRow(context, "Allow Login Alerts", allowLoginAlerts, (val) {
                setState(() => allowLoginAlerts = val);
              }, icon: Icons.login_outlined),
            ],
            icon: Icons.shield_outlined,
          ),
          _buildSection(
            context,
            "Preferences",
            [
              _buildDropdownRow(context, "Language", setting!.language, languages, (val) {
                setState(() {
                  setting = setting!.copyWith(language: val);
                });
              }, icon: Icons.language_outlined),
              _buildDropdownRow(context, "Theme", setting!.theme, themes, (val) {
                setState(() {
                  setting = setting!.copyWith(theme: val);
                  themeProvider.toggleTheme(val == 'dark');
                });
              }, icon: Icons.palette_outlined),
              _buildToggleRow(context, "Auto Delete Old Moods", autoDeleteOldMoods, (val) {
                setState(() => autoDeleteOldMoods = val);
              }, icon: Icons.auto_delete_outlined),
            ],
            icon: Icons.settings_outlined,
          ),
          const SizedBox(height: 32),
          Container(
            margin: const EdgeInsets.symmetric(horizontal: 16),
            child: Row(
              children: [
                Expanded(
                  child: CustomButton(
                    onPressed: _resetDefault,
                    text: "Reset Default",
                    isOutlined: true,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: CustomButton(
                    onPressed: _saveSettings,
                    text: "Save Changes",
                    isLoading: isSaving,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),
        ],
      ),
    );
  }
}
