import 'package:PsyConnect/page/register_page.dart';
import 'package:PsyConnect/service/account_service/login.dart';
import 'package:flutter/material.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  TextEditingController nameController = TextEditingController();
  TextEditingController passwordController = TextEditingController();
  LoginService loginService = LoginService();

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: DefaultTabController(
          length: 2,
          child: Container(
            
            width: MediaQuery.of(context).size.width,
            height: MediaQuery.of(context).size.height,
            decoration: const BoxDecoration(
              image: DecorationImage(
                image: AssetImage(
                    "assets/images/light_mode_backround_login_page.jpg"),
                fit: BoxFit.cover,
              ),
            ),
            child: Scaffold(

              appBar: AppBar(
                bottom: const TabBar(
                  tabs: [Tab(text: 'Login'), Tab(text: 'Register')],
                ),
              ),
              body: TabBarView(
                children: [
                  _LoginForm(context),
                  const RegisterPage(),
                ],
              ),
            ),
          )),
    );
  }

  Widget _LoginForm(BuildContext context) {
    return Center(
      child: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(19.0),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: <Widget>[
              TextField(
                controller: nameController,
                decoration: InputDecoration(
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(30)),
                  labelText: 'Username',
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                obscureText: true,
                controller: passwordController,
                decoration: InputDecoration(
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(30)),
                  labelText: 'Password',
                ),
              ),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  IconButton(
                    onPressed: () => _handleOnAppleIdLogin(),
                    icon: Image.asset(
                      "assets/images/login_appleid_icon.png",
                      width: 30,
                    ),
                  ),
                  IconButton(
                    onPressed: () => _handleOnGoogleLogin(),
                    icon: Image.asset(
                      "assets/images/login_google_icon.jpeg",
                      width: 30,
                    ),
                  ),
                  IconButton(
                    onPressed: () => _handleOnFacebookLogin(),
                    icon: Image.asset(
                      "assets/images/login_facebook_icon.jpeg",
                      width: 30,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () {
                  loginService.loginHandle(
                    nameController.text,
                    passwordController.text,
                    "NORMAL",
                    context,
                    "android",
                  );
                },
                child: const Text('Login'),
              ),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: <Widget>[
                  const Text('Forgot your password?'),
                  TextButton(
                    onPressed: () {},
                    child: const Text(
                      'Reset now',
                      style: TextStyle(fontSize: 15),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  _handleOnGoogleLogin() {}

  _handleOnFacebookLogin() {}

  _handleOnAppleIdLogin() {}

  Widget _RegisterForm(BuildContext context) {
    return SingleChildScrollView(
      child: Center(),
    );
  }

  Widget _ForgotPasswordForm(BuildContext context) {
    return Center(
      child: Text("Forgot Password Form - Work in Progress"),
    );
  }
}
