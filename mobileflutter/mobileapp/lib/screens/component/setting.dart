import 'package:PsyConnect/page/login_page.dart';
import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/widgets/text/custom_text.dart';
import 'package:flutter/material.dart';

class Setting extends StatefulWidget {
  const Setting({super.key});

  @override
  State<Setting> createState() => _SettingState();
}

class _SettingState extends State<Setting> {
  @override
  Widget build(BuildContext context) {
    return Center(
        child: ListView(
      padding: const EdgeInsets.all(20.0),
      children: [
        Row(mainAxisAlignment: MainAxisAlignment.end, children: [
          CustomTextButton(
            text: UserProvider().getData?["username"] ?? "Unknow user",
            fontSize: 30,
            fontWeight: FontWeight.bold,
          ),
          SizedBox(
            height: 30,
          ),
          CircleAvatar(
            radius: 40,
            backgroundImage: NetworkImage(
              UserProvider().getData?["avatarUri"] ??
                  "https://example.com/default_avatar.png",
            ),
            onBackgroundImageError: (error, stackTrace) {},
            backgroundColor: Colors.grey[200],
            child: UserProvider().user?["avatarUri"] == null
                ? Icon(Icons.person, size: 40, color: Colors.grey)
                : null,
          )
        ]),
        ListTile(
          trailing: Icon(Icons.person),
          title: CustomTextButton(text: "Persional Infomation"),
          onTap: () {},
        ),
        ListTile(
          trailing: Icon(Icons.lock),
          title: CustomTextButton(text: "Security & Login"),
          onTap: () {},
        ),
        ListTile(
          trailing: Icon(Icons.notifications),
          title: CustomTextButton(text: "Notifications"),
          onTap: () {},
        ),
        ListTile(
          trailing: Icon(Icons.palette),
          title: CustomTextButton(text: "Appearance & Theme "),
          onTap: () {},
        ),
        ListTile(
          trailing: Icon(Icons.help),
          title: CustomTextButton(text: "Help & Feedback"),
          onTap: () {},
        ),
        ListTile(
            trailing: Icon(Icons.logout, color: Colors.red),
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
