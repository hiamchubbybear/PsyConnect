import 'package:PsyConnect/screens/home/my_home_page.dart';
import 'package:PsyConnect/service/identityService.dart';
import 'package:PsyConnect/themes/colors/color.dart';
import 'package:flutter/material.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  TextEditingController nameController = TextEditingController();
  TextEditingController passwordController = TextEditingController();

  Future<void> _loginHandle() async {
    if (nameController.text.isEmpty || passwordController.text.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
            content: Text('Please enter both username and password before proceeding',
                style: TextStyle(fontSize: 20))),
      );
    } else {
      var response = await loginHandleIdentityService(
          nameController.text, passwordController.text, "NORMAL");
      if (response != null && response.body.isNotEmpty) {
        if (response.statusCode == 200) {
          print(response);
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
                showCloseIcon: true,
                content: Text('Login sucess!',
                    style: TextStyle(fontSize: 20))),
          );
          Navigator.pushReplacement(
            context,
            MaterialPageRoute(
                builder: (context) => const MyHomePage(title: 'Home Page')),
          );
        } else if (response.statusCode == 500) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
                showCloseIcon: true,
                content: Text('Incorrect username or password!',
                    style: TextStyle(fontSize: 20))),
          );
        } else {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
                showCloseIcon: true,
                content: Text(
                  'Sai mật khẩu!!',
                  style: TextStyle(fontSize: 20),
                )),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(10),
        child: ListView(
          children: <Widget>[
            Container(
                alignment: Alignment.center,
                padding: const EdgeInsets.all(10),
                child: Text(
                  'PsyConnect',
                  style: TextStyle(
                      color: primaryColor,
                      fontWeight: FontWeight.w500,
                      fontSize: 30),
                )),
            Container(
                alignment: Alignment.center,
                padding: const EdgeInsets.all(10),
                child: const Text(
                  'Sign in',
                  style: TextStyle(fontSize: 20),
                )),
            Container(
              padding: const EdgeInsets.all(10),
              child: TextField(
                controller: nameController,
                decoration: const InputDecoration(
                  border: OutlineInputBorder(),
                  labelText: 'Username',
                ),
              ),
            ),
            Container(
              padding: const EdgeInsets.fromLTRB(10, 10, 10, 0),
              child: TextField(
                obscureText: true,
                controller: passwordController,
                decoration: const InputDecoration(
                  border: OutlineInputBorder(),
                  labelText: 'Password',
                ),
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                IconButton(
                    onPressed: () => _handleOnAppleIdLogin(),
                    icon: Image.asset(
                      "assets/images/login_appleid_icon.png",
                      width: 30,
                    )),
                IconButton(
                    onPressed: () => _handleOnGoogleLogin(),
                    icon: Image.asset(
                      "assets/images/login_google_icon.jpeg",
                      width: 30,
                    )),
                IconButton(
                    onPressed: () => _handleOnFacebookLogin(),
                    icon: Image.asset(
                      "assets/images/login_facebook_icon.jpeg",
                      width: 30,
                    )),
              ],
            ),
            Container(
                height: 50,
                padding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
                child: ElevatedButton(
                  child: const Text('Login'),
                  onPressed: _loginHandle,
                )),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                const Text('Forgot your password?'),
                TextButton(
                  child: const Text(
                    'Reset now',
                    style: TextStyle(fontSize: 15),
                  ),
                  onPressed: () {},
                )
              ],
            ),
          ],
        ),
      ),
    );
  }

  _handleOnGoogleLogin() {}

  _handleOnFacebookLogin() {}

  _handleOnAppleIdLogin() {}
}
