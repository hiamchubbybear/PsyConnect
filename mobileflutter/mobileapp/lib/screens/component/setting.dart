import 'package:PsyConnect/page/login_page.dart';
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
      children: <Widget>[
        IconButton(
            onPressed: () => _handleSignOutPress(), icon: Icon(Icons.logout)),
        IconButton(
            onPressed: () {
              throw new Error();
            },
            icon: Icon(Icons.no_accounts)),
        IconButton(
            onPressed: () {
              throw new Error();
            },
            icon: Icon(Icons.no_accounts)),
        IconButton(
            onPressed: () {
              throw new Error();
            },
            icon: Icon(Icons.no_accounts)),
      ],
    ));
  }

  Future<void> _handleSignOutPress() async {
    Navigator.pushReplacement(
        context, MaterialPageRoute(builder: (context) => LoginPage()));
  }
}
