import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/setting.dart';
import 'package:PsyConnect/services/profile_service/setting.dart';
import 'package:animated_toggle_switch/animated_toggle_switch.dart';
import 'package:flutter/material.dart';

class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});
  final String usersetting = 'user_settings';
  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  Setting? setting;

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
        });
      }
    } catch (e) {
      print("Error loading settings: $e");
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error loading settings: $e')),
        );
      }
    }
  }

  Future<void> _saveSettings() async {
    if (setting == null) return;

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
      SettingService settingService = SettingService();
      final response = await settingService.updateSetting();
      if (response != null && response) {
        if (mounted) {
          setState(() {
            setting = newSetting;
          });

          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Settings saved')),
          );
        }
      } else if (!response) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Failed to save setting')),
        );
      }
    } catch (e) {
      print("Error saving settings: $e");
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error saving settings: $e')),
        );
      }
    }
  }

  Widget buildCompactToggle(
      String title, bool value, Function(bool) onChanged) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Expanded(
            child: Text(
              title,
              style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w500),
            ),
          ),
          AnimatedToggleSwitch<bool>.dual(
            current: value,
            first: false,
            second: true,
            spacing: 3,
            height: 24,
            indicatorSize: const Size.square(18),
            animationDuration: const Duration(milliseconds: 250),
            animationCurve: Curves.easeInOut,
            borderWidth: 0,
            customStyleBuilder: (context, local, global) {
              return ToggleStyle(
                backgroundColor:
                    value ? Colors.green.shade300 : Colors.grey.shade300,
                borderRadius: BorderRadius.circular(12),
                indicatorColor: Colors.white,
              );
            },
            iconBuilder: (value) => Icon(
              value ? Icons.check : Icons.close,
              size: 14,
              color: value ? Colors.green.shade800 : Colors.grey.shade600,
            ),
            onChanged: onChanged,
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (setting == null) {
      return const Scaffold(
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              CircularProgressIndicator(),
              SizedBox(height: 16),
              Text("Loading settings..."),
            ],
          ),
        ),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text("Settings"),
        elevation: 1,
        backgroundColor: Colors.white,
        foregroundColor: Colors.black87,
      ),
      backgroundColor: Colors.grey.shade100,
      body: ListView(
        children: [
          buildSection("Privacy", [
            buildCompactToggle("Show Last Seen", showLastSeen, (val) {
              setState(() => showLastSeen = val);
            }),
            buildCompactToggle("Show Profile Picture", showProfilePicture,
                (val) {
              setState(() => showProfilePicture = val);
            }),
            buildCompactToggle("Show Mood", showMood, (val) {
              setState(() => showMood = val);
            }),
          ]),
          buildSection("Notifications", [
            buildCompactToggle("Enable Notifications", notificationsEnabled,
                (val) {
              setState(() => notificationsEnabled = val);
            }),
            buildCompactToggle("Email Notifications", emailNotifications,
                (val) {
              setState(() => emailNotifications = val);
            }),
            buildCompactToggle("Push Notifications", pushNotifications, (val) {
              setState(() => pushNotifications = val);
            }),
            buildCompactToggle("SMS Notifications", smsNotifications, (val) {
              setState(() => smsNotifications = val);
            }),
          ]),
          buildSection("Security", [
            buildCompactToggle("Two-Factor Authentication", twoFactorAuth,
                (val) {
              setState(() => twoFactorAuth = val);
            }),
            buildCompactToggle("Allow Login Alerts", allowLoginAlerts, (val) {
              setState(() => allowLoginAlerts = val);
            }),
          ]),
          buildSection("Preferences", [
            ListTile(
              title: const Text("Language", style: TextStyle(fontSize: 14)),
              trailing: DropdownButton<String>(
                value: setting!.language,
                items: languages.map((lang) {
                  return DropdownMenuItem(
                    value: lang,
                    child: Text(lang.toUpperCase()),
                  );
                }).toList(),
                onChanged: (val) {
                  if (val != null) {
                    setState(() {
                      setting = setting!.copyWith(language: val);
                    });
                  }
                },
              ),
            ),
            ListTile(
              title: const Text("Theme", style: TextStyle(fontSize: 14)),
              trailing: DropdownButton<String>(
                value: setting!.theme,
                items: themes.map((t) {
                  return DropdownMenuItem(
                    value: t,
                    child: Text(t[0].toUpperCase() + t.substring(1)),
                  );
                }).toList(),
                onChanged: (val) {
                  if (val != null) {
                    setState(() {
                      setting = setting!.copyWith(theme: val);
                    });
                  }
                },
              ),
            ),
            buildCompactToggle("Auto Delete Old Moods", autoDeleteOldMoods,
                (val) {
              setState(() => autoDeleteOldMoods = val);
            }),
          ]),
          const SizedBox(height: 20),
          Container(
            child: Center(
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  TextButton(
                    onPressed: () => _resetDefault(),
                    child: Text(
                      "Reset Default",
                      style: quickSand15Font,
                    ),
                  ),
                  TextButton(
                    onPressed: () => _saveSettings(),
                    child: Text("Save", style: quickSand15Font),
                  ),
                ],
              ),
            ),
          )
        ],
      ),
    );
  }

  Widget buildSection(String title, List<Widget> children) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 8, vertical: 6),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(12),
            child: Text(
              title,
              style: const TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 15,
                color: Colors.black87,
              ),
            ),
          ),
          ...children,
          const SizedBox(height: 8),
        ],
      ),
    );
  }
}
