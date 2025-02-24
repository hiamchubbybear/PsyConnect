import 'package:PsyConnect/page/login_page.dart';
import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/widgets/text/custom_text.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class Setting extends StatefulWidget {
  const Setting({super.key});

  @override
  State<Setting> createState() => _SettingState();
}

class _SettingState extends State<Setting> {
  @override
  Widget build(BuildContext context) {
    final userProvider = Provider.of<UserProvider>(context);

    print(userProvider.user?["imageUri"]);
    return Center(
        child: ListView(
      padding: const EdgeInsets.all(20.0),
      children: [
        Row(mainAxisAlignment: MainAxisAlignment.end, children: [
          CircleAvatar(
            radius: 30,
            backgroundImage: NetworkImage(
              userProvider.user?["imageUri"] ??
                  "https://example.com/default_avatar.png",
            ),
            onBackgroundImageError: (error, stackTrace) {},
            backgroundColor: Colors.grey[200],
            child: userProvider.user?["avatarUri"] == null
                ? Icon(Icons.person, size: 50, color: Colors.grey)
                : null,
          ),
          CustomTextButton(
            text: userProvider.user?["username"] ?? "Unknown user",
            fontSize: 30,
            fontWeight: FontWeight.bold,
          ),
        ]),
        ListTile(
          leading: Icon(Icons.person),
          title: CustomTextButton(text: "Persional Infomation"),
          onTap: () {},
        ),
        ListTile(
          leading: Icon(Icons.lock),
          title: CustomTextButton(text: "Security & Login"),
          onTap: () {},
        ),
        ListTile(
          leading: Icon(Icons.notifications),
          title: CustomTextButton(text: "Notifications"),
          onTap: () {},
        ),
        ListTile(
          leading: Icon(Icons.palette),
          title: CustomTextButton(text: "Appearance & Theme "),
          onTap: () {},
        ),
        ListTile(
          leading: Icon(Icons.help),
          title: CustomTextButton(text: "Help & Feedback"),
          onTap: () {},
        ),
        ListTile(
            leading: Icon(Icons.logout, color: Colors.red),
            title: CustomTextButton(text: "Logout"),
            onTap: () => _handleSignOutPress()),
      ],
    ));
  }

  Future<void> _handleSignOutPress() async {
    UserProvider().clearUser();
    Navigator.pushReplacement(
        context, MaterialPageRoute(builder: (context) => LoginPage()));
  }
}
