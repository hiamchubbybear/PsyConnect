import 'package:PsyConnect/provider/auth_token_provider.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/provider/user_profile_provider.dart';
import 'package:PsyConnect/services/account_service/login.dart';
import 'package:PsyConnect/ui/screens/forgot_page.dart';
import 'package:PsyConnect/ui/screens/my_home_page.dart';
import 'package:PsyConnect/ui/screens/register_page.dart';
import 'package:PsyConnect/ui/widgets/common/custom_text_field.dart';
import 'package:PsyConnect/ui/widgets/common/custom_button.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final TextEditingController nameController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  bool isPasswordVisible = false;
  final LoginService loginService = LoginService();
  final GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      loginService.initDeepLinkListener(context, () {
        _navigateToHome();
      });
    });
  }

  @override
  void dispose() {
    nameController.dispose();
    passwordController.dispose();
    loginService.dispose();
    super.dispose();
  }

  void _navigateToHome() {
    final navigator = navigatorKey.currentState ?? Navigator.of(context);
    navigator.pushReplacement(
      MaterialPageRoute(
        builder: (_) => const MyHomePage(title: 'Home Page'),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;
    final themeProvider = Provider.of<ThemeProvider>(context);

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  IconButton(
                    onPressed: () {
                      themeProvider.toggleTheme(!themeProvider.isDarkMode);
                    },
                    icon: Icon(
                      themeProvider.isDarkMode ? Icons.light_mode : Icons.dark_mode,
                      color: isDark ? Colors.white : Colors.black87,
                      size: 24,
                    ),
                  ),
                ],
              ),
            ),
            Expanded(
              child: _buildLoginForm(context),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildLoginForm(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;
    final userProfileProvider = Provider.of<UserProfileProvider>(context, listen: false);
    final tokenProvider = Provider.of<AuthTokenProvider>(context, listen: false);

    final textColor = isDark ? Colors.white : Colors.black87;
    final subtitleColor = isDark ? Colors.grey[400] : Colors.grey[600];

    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 32.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 60),
            Text(
              'Welcome',
              style: GoogleFonts.quicksand(
                fontSize: 32,
                fontWeight: FontWeight.w400,
                color: textColor,
                letterSpacing: -0.5,
              ),
            ),
            const SizedBox(height: 8),
            RichText(
              text: TextSpan(
                style: TextStyle(
                  fontSize: 16,
                  color: subtitleColor,
                  fontWeight: FontWeight.w400,
                ),
                children: [
                  const TextSpan(text: 'Sign in to continue '),
                  TextSpan(
                    text: 'or sign up',
                    style: TextStyle(
                      decoration: TextDecoration.underline,
                      color: subtitleColor,
                    ),
                    recognizer: TapGestureRecognizer()
                      ..onTap = () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => const MultiStepRegisterPage(),
                          ),
                        );
                      },
                  ),
                ],
              ),
            ),
            const SizedBox(height: 8),
            RichText(
              text: TextSpan(
                style: TextStyle(
                  fontSize: 16,
                  color: subtitleColor,
                  fontWeight: FontWeight.w400,
                ),
                children: [
                  const TextSpan(text: 'Forgot password '),
                  TextSpan(
                    text: 'reset now',
                    style: TextStyle(
                      decoration: TextDecoration.underline,
                      color: subtitleColor,
                    ),
                    recognizer: TapGestureRecognizer()
                      ..onTap = () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => const ForgotPage(),
                          ),
                        );
                      },
                  ),
                ],
              ),
            ),
            const SizedBox(height: 60),
            CustomTextField(
              controller: nameController,
              labelText: 'Username',
              prefixIcon: Icon(Icons.person_outline, color: isDark ? Colors.grey[400] : Colors.grey[600]),
              maxLength: 20,
            ),
            const SizedBox(height: 24),
            CustomTextField(
              controller: passwordController,
              labelText: 'Password',
              obscureText: !isPasswordVisible,
              prefixIcon: Icon(Icons.lock_outline, color: isDark ? Colors.grey[400] : Colors.grey[600]),
              suffixIcon: IconButton(
                icon: Icon(
                  isPasswordVisible ? Icons.visibility_outlined : Icons.visibility_off_outlined,
                  color: isDark ? Colors.grey[400] : Colors.grey[600],
                  size: 20,
                ),
                onPressed: () {
                  setState(() {
                    isPasswordVisible = !isPasswordVisible;
                  });
                },
              ),
            ),
            const SizedBox(height: 40),
            CustomButton(
              onPressed: () => _handleOnLoginButton(
                username: nameController.text,
                password: passwordController.text,
                provider: "NORMAL",
                context: context,
                tokenProvider: tokenProvider,
                userProfileProvider: userProfileProvider,
              ),
              text: 'Sign In',
              color: isDark ? Colors.white : Colors.black,
              textColor: isDark ? Colors.black : Colors.white,
            ),
            const SizedBox(height: 40),
            _buildDivider(isDark),
            const SizedBox(height: 40),
            _buildSocialLoginSection(isDark, context),
            const SizedBox(height: 40),
            _buildForgotPasswordSection(context, isDark),
            const SizedBox(height: 40),
          ],
        ),
      ),
    );
  }

  Widget _buildDivider(bool isDark) {
    final dividerColor = isDark ? Colors.grey[800] : Colors.grey[200];
    final textColor = isDark ? Colors.grey[500] : Colors.grey[500];

    return Row(
      children: [
        Expanded(
          child: Container(
            height: 1,
            color: dividerColor,
          ),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Text(
            'or',
            style: TextStyle(
              color: textColor,
              fontSize: 14,
              fontWeight: FontWeight.w400,
            ),
          ),
        ),
        Expanded(
          child: Container(
            height: 1,
            color: dividerColor,
          ),
        ),
      ],
    );
  }

  Widget _buildSocialLoginSection(bool isDark, BuildContext ctx) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _buildSocialButton(
          icon: FontAwesomeIcons.google,
          onPressed: () => _handleOnGoogleLogin(context: ctx),
          isDark: isDark,
        ),
        const SizedBox(width: 24),
        _buildSocialButton(
          icon: FontAwesomeIcons.apple,
          onPressed: () => _handleOnAppleLogin(context: ctx),
          isDark: isDark,
        ),
        const SizedBox(width: 24),
        _buildSocialButton(
          icon: FontAwesomeIcons.facebook,
          onPressed: () => _handleOnFacebookLogin(context: ctx),
          isDark: isDark,
        ),
      ],
    );
  }

  Widget _buildSocialButton({
    required IconData icon,
    required VoidCallback onPressed,
    required bool isDark,
  }) {
    final borderColor = isDark ? Colors.grey[700] : Colors.grey[300];
    final iconColor = isDark ? Colors.white : Colors.black;

    return Container(
      width: 48,
      height: 48,
      decoration: BoxDecoration(
        border: Border.all(
          color: borderColor!,
          width: 1,
        ),
        borderRadius: BorderRadius.circular(12),
      ),
      child: IconButton(
        onPressed: onPressed,
        icon: Icon(
          icon,
          color: iconColor,
          size: 18,
        ),
      ),
    );
  }

  Widget _buildForgotPasswordSection(BuildContext context, bool isDark) {
    final textColor = isDark ? Colors.grey[400] : Colors.grey[600];

    return Center(
      child: TextButton(
        onPressed: () => _handleOnResetPassword(context: context),
        style: TextButton.styleFrom(
          padding: EdgeInsets.zero,
          minimumSize: Size.zero,
          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
        ),
        child: Text(
          'Forgot password?',
          style: GoogleFonts.quicksand(
            color: textColor,
            fontSize: 15,
            fontWeight: FontWeight.w400,
            decoration: TextDecoration.underline,
            decorationColor: textColor,
          ),
        ),
      ),
    );
  }

  void _handleOnLoginButton({
    required String username,
    required String password,
    required String provider,
    required BuildContext context,
    required AuthTokenProvider tokenProvider,
    required UserProfileProvider userProfileProvider,
  }) {
    const String platform = "MOBILE";
    loginService.loginHandle(username, password, provider, context, platform,
        tokenProvider, userProfileProvider);
  }
}

void _handleOnGoogleLogin({required BuildContext context}) {
  const String provider = "google";
  final LoginService loginService = LoginService();
  loginService.oauth2LoginHandle(provider);
}

void _handleOnAppleLogin({required BuildContext context}) {
  const String provider = "apple";
  final LoginService loginService = LoginService();
  loginService.oauth2LoginHandle(provider);
}

void _handleOnFacebookLogin({required BuildContext context}) {
  const String provider = "facebook";
  final LoginService loginService = LoginService();
  loginService.oauth2LoginHandle(provider);
}

void _handleOnResetPassword({required BuildContext context}) {
  Navigator.push(
    context,
    MaterialPageRoute(builder: (context) => const ForgotPage()),
  );
}
