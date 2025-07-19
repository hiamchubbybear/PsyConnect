import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/setting.dart';
import 'package:PsyConnect/services/profile_service/setting.dart';
import 'package:flutter/material.dart';

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
      // themeProvider.toggleTheme(!themeProvider.isDarkMode);
      final prefs = SharedPreferencesProvider();
      await prefs.setSetting(newSetting);
      SettingService settingService = SettingService();
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

  Widget _buildToggleRow(String title, bool value, Function(bool) onChanged,
      {IconData? icon}) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
      child: Row(
        children: [
          if (icon != null) ...[
            Icon(
              icon,
              size: 20,
              color: themeProvider.isDarkMode
                  ? Colors.white70
                  : Colors.grey.shade600,
            ),
            const SizedBox(width: 16),
          ],
          Expanded(
            child: Text(
              title,
              style: kSubHeadingStyle.copyWith(
                fontWeight: FontWeight.w500,
                fontSize: 15,
              ),
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
                    ? acceptColor
                    : (themeProvider.isDarkMode
                        ? Colors.grey.shade700
                        : Colors.grey.shade300),
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

  Widget _buildSection(String title, List<Widget> children, {IconData? icon}) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      decoration: BoxDecoration(
        color: themeProvider.isDarkMode ? Colors.grey.shade800 : Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.all(20),
            child: Row(
              children: [
                if (icon != null) ...[
                  Icon(
                    icon,
                    size: 20,
                    color: themeProvider.isDarkMode
                        ? Colors.white
                        : Colors.black87,
                  ),
                  const SizedBox(width: 12),
                ],
                Text(
                  title,
                  style: kSubHeadingStyle.copyWith(
                    fontWeight: FontWeight.w600,
                    fontSize: 16,
                    color: themeProvider.isDarkMode
                        ? Colors.white
                        : Colors.black87,
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

  Widget _buildDropdownRow(String title, String value, List<String> items,
      Function(String) onChanged,
      {IconData? icon}) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
      child: Row(
        children: [
          if (icon != null) ...[
            Icon(
              icon,
              size: 20,
              color: themeProvider.isDarkMode
                  ? Colors.white70
                  : Colors.grey.shade600,
            ),
            const SizedBox(width: 16),
          ],
          Expanded(
            child: Text(
              title,
              style: kSubHeadingStyle.copyWith(
                fontWeight: FontWeight.w500,
                fontSize: 15,
              ),
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            decoration: BoxDecoration(
              color: themeProvider.isDarkMode
                  ? Colors.grey.shade700
                  : Colors.grey.shade100,
              borderRadius: BorderRadius.circular(8),
            ),
            child: DropdownButton<String>(
              value: value,
              underline: const SizedBox(),
              style: kSubHeadingStyle.copyWith(fontSize: 14),
              items: items.map((item) {
                return DropdownMenuItem(
                  value: item,
                  child: Text(
                    item.toUpperCase(),
                    style: kSubHeadingStyle.copyWith(fontSize: 14),
                  ),
                );
              }).toList(),
              onChanged: (val) {
                if (val != null) onChanged(val);
              },
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (isLoading) {
      return Scaffold(
        backgroundColor: themeProvider.isDarkMode
            ? Colors.grey.shade900
            : Colors.grey.shade50,
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              CircularProgressIndicator(
                color: acceptColor,
                strokeWidth: 2,
              ),
              const SizedBox(height: 16),
              Text(
                "Loading settings...",
                style: kSubHeadingStyle,
              ),
            ],
          ),
        ),
      );
    }

    return Scaffold(
      backgroundColor:
          themeProvider.isDarkMode ? Colors.grey.shade900 : Colors.grey.shade50,
      appBar: AppBar(
        title: Text(
          "Settings",
          style: kHeadingStyle.copyWith(fontSize: 20),
        ),
        centerTitle: true,
        elevation: 0,
        backgroundColor: themeProvider.isDarkMode
            ? Colors.grey.shade900
            : Colors.grey.shade50,
        foregroundColor:
            themeProvider.isDarkMode ? Colors.white : Colors.black87,
      ),
      body: ListView(
        padding: const EdgeInsets.symmetric(vertical: 8),
        children: [
          _buildSection(
            "Privacy",
            [
              _buildToggleRow("Show Last Seen", showLastSeen, (val) {
                setState(() => showLastSeen = val);
              }, icon: Icons.visibility_outlined),
              _buildToggleRow("Show Profile Picture", showProfilePicture,
                  (val) {
                setState(() => showProfilePicture = val);
              }, icon: Icons.account_circle_outlined),
              _buildToggleRow("Show Mood", showMood, (val) {
                setState(() => showMood = val);
              }, icon: Icons.mood_outlined),
            ],
            icon: Icons.security_outlined,
          ),
          _buildSection(
            "Notifications",
            [
              _buildToggleRow("Enable Notifications", notificationsEnabled,
                  (val) {
                setState(() => notificationsEnabled = val);
              }, icon: Icons.notifications_outlined),
              _buildToggleRow("Email Notifications", emailNotifications, (val) {
                setState(() => emailNotifications = val);
              }, icon: Icons.email_outlined),
              _buildToggleRow("Push Notifications", pushNotifications, (val) {
                setState(() => pushNotifications = val);
              }, icon: Icons.push_pin_outlined),
              _buildToggleRow("SMS Notifications", smsNotifications, (val) {
                setState(() => smsNotifications = val);
              }, icon: Icons.sms_outlined),
            ],
            icon: Icons.notifications_active_outlined,
          ),
          _buildSection(
            "Security",
            [
              _buildToggleRow("Two-Factor Authentication", twoFactorAuth,
                  (val) {
                setState(() => twoFactorAuth = val);
              }, icon: Icons.security_outlined),
              _buildToggleRow("Allow Login Alerts", allowLoginAlerts, (val) {
                setState(() => allowLoginAlerts = val);
              }, icon: Icons.login_outlined),
            ],
            icon: Icons.shield_outlined,
          ),
          _buildSection(
            "Preferences",
            [
              _buildDropdownRow("Language", setting!.language, languages,
                  (val) {
                setState(() {
                  setting = setting!.copyWith(language: val);
                });
              }, icon: Icons.language_outlined),
              _buildDropdownRow("Theme", setting!.theme, themes, (val) {
                setState(() {
                  setting = setting!.copyWith(theme: val);
                  themeProvider.toggleTheme(val == 'dark');
                });
              }, icon: Icons.palette_outlined),
              _buildToggleRow("Auto Delete Old Moods", autoDeleteOldMoods,
                  (val) {
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
                  child: TextButton(
                    onPressed: _resetDefault,
                    style: TextButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                        side: BorderSide(
                          color: themeProvider.isDarkMode
                              ? Colors.grey.shade600
                              : Colors.grey.shade300,
                        ),
                      ),
                    ),
                    child: Text(
                      "Reset Default",
                      style: kSubHeadingStyle.copyWith(
                        fontWeight: FontWeight.w300,
                        color: themeProvider.isDarkMode
                            ? Colors.white70
                            : Colors.grey.shade700,
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton(
                    onPressed: isSaving ? null : _saveSettings,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: acceptColor,
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: isSaving
                        ? const SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              valueColor:
                                  AlwaysStoppedAnimation<Color>(Colors.white),
                            ),
                          )
                        : Text(
                            "Save Changes",
                            style: kSubHeadingStyle.copyWith(
                              fontWeight: FontWeight.w600,
                              color: Colors.white,
                            ),
                          ),
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
