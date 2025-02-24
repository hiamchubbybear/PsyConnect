import 'package:PsyConnect/page/forgot_page.dart';
import 'package:PsyConnect/page/register_page.dart';
import 'package:PsyConnect/provider/auth_token_provider.dart';
import 'package:PsyConnect/provider/user_profile_provider.dart';
import 'package:PsyConnect/service/account_service/login.dart';
import 'package:PsyConnect/variable/variable.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  TextEditingController nameController = TextEditingController();
  TextEditingController passwordController = TextEditingController();
  bool isPasswordVisible = false;

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
                bottom: TabBar(
                  tabs: [
                    Tab(child: Text('Login', style: textStyle)),
                    Tab(child: Text('Register', style: textStyle))
                  ],
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
    UserProfileProvider userProfileProvider =
        Provider.of<UserProfileProvider>(context, listen: false);
    AuthTokenProvider tokenProvider =
        Provider.of<AuthTokenProvider>(context, listen: false);
    return Center(
      child: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(19.0),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: <Widget>[
              TextField(
                maxLength: 20,
                style: textStyle,
                controller: nameController,
                decoration: InputDecoration(
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(20)),
                  labelText: 'Username',
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                style: textStyle,
                obscureText: !isPasswordVisible,
                controller: passwordController,
                decoration: InputDecoration(
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(20)),
                  labelText: 'Password',
                  suffixIcon: IconButton(
                    icon: Icon(
                      isPasswordVisible
                          ? Icons.visibility
                          : Icons.visibility_off,
                    ),
                    onPressed: () {
                      setState(() {
                        isPasswordVisible = !isPasswordVisible;
                      });
                    },
                  ),
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
                      tokenProvider,
                      userProfileProvider);
                },
                child: Text(
                  'Login',
                  style: textStyle,
                ),
              ),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: <Widget>[
                  Text('Forgot your password?', style: textStyle),
                  TextButton(
                    onPressed: () => _handleOnResetPassword(context: context),
                    child: Text(
                      'Reset now',
                      style: textStyle,
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
  _handleOnResetPassword({
    required BuildContext context,
  }) {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => const ForgotPage()),
    );
  }
}
