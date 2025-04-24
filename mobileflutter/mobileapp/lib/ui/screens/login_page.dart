import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/provider/auth_token_provider.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/provider/user_profile_provider.dart';
import 'package:PsyConnect/services/account_service/login.dart';
import 'package:PsyConnect/ui/screens/forgot_page.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
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
  bool light = true;
  LoginService loginService = LoginService();

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: DefaultTabController(
          length: 1,
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
              appBar: _appBar(context),
              body: TabBarView(
                children: [
                  _LoginForm(context),
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
                style: quickSand15Font,
                controller: nameController,
                decoration: InputDecoration(
                  border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(20)),
                  labelText: 'Username',
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                style: quickSand15Font,
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
                      icon: const Icon(
                        FontAwesomeIcons.apple,
                        size: 32,
                      )),
                  IconButton(
                    onPressed: () => _handleOnGoogleLogin(),
                    icon: const Icon(FontAwesomeIcons.facebook),
                  ),
                  IconButton(
                      onPressed: () => _handleOnFacebookLogin(),
                      icon: const Icon(FontAwesomeIcons.google)),
                ],
              ),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => _handleOnLoginButton(
                  username: nameController.text,
                  password: passwordController.text,
                  loginType: "NORMAL",
                  context: context,
                  tokenProvider: tokenProvider,
                  userProfileProvider: userProfileProvider,
                ),
                child: Text(
                  'Login',
                  style: quickSand15Font,
                ),
              ),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: <Widget>[
                  Text('Forgot your password?', style: quickSand15Font),
                  TextButton(
                    onPressed: () => _handleOnResetPassword(context: context),
                    child: Text(
                      'Reset now',
                      style: quickSand15Font,
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

  _handleOnLoginButton(
      {required String username,
      required String password,
      required String loginType,
      required BuildContext context,
      required AuthTokenProvider tokenProvider,
      required UserProfileProvider userProfileProvider}) {
    loginService.loginHandle(username, password, loginType, context, loginType,
        tokenProvider, userProfileProvider);
  }
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

AppBar _appBar(BuildContext context) {
  final themeProvider = Provider.of<ThemeProvider>(context);
  return AppBar(
    automaticallyImplyLeading: false,
    actions: [
      PopupMenuButton<String>(
        icon: const Icon(Icons.more_horiz_outlined),
        itemBuilder: (context) => [
          PopupMenuItem(
            enabled: true,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Brightness'),
                Transform.scale(
                  scale: 0.8,
                  child: Switch.adaptive(
                    activeTrackColor: Colors.black,
                    value: themeProvider.isDarkMode,
                    onChanged: (value) {
                      Navigator.pop(context);
                      themeProvider.toggleTheme(value);
                    },
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    ],
  );
}
