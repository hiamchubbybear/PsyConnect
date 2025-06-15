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
  LoginService loginService = LoginService();

  @override
  Widget build(BuildContext context) {
    final themeProvider = Provider.of<ThemeProvider>(context);

    return Scaffold(
      backgroundColor: themeProvider.isDarkMode ? Colors.black : Colors.white,
      body: SafeArea(
        child: Column(
          children: [
            _buildAppBar(context),
            Expanded(
              child: _buildLoginForm(context),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAppBar(BuildContext context) {
    final themeProvider = Provider.of<ThemeProvider>(context);
    return Padding(
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
              color: themeProvider.isDarkMode ? Colors.white : Colors.black,
              size: 24,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildLoginForm(BuildContext context) {
    final themeProvider = Provider.of<ThemeProvider>(context);
    UserProfileProvider userProfileProvider =
        Provider.of<UserProfileProvider>(context, listen: false);
    AuthTokenProvider tokenProvider =
        Provider.of<AuthTokenProvider>(context, listen: false);

    final textColor = themeProvider.isDarkMode ? Colors.white : Colors.black;
    final subtitleColor = themeProvider.isDarkMode ? Colors.grey[400] : Colors.grey[600];

    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 32.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 60),

            // Welcome Text
            Text(
              'Welcome',
              style: TextStyle(
                fontSize: 32,
                fontWeight: FontWeight.w300,
                color: textColor,
                letterSpacing: -0.5,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Sign in to continue',
              style: TextStyle(
                fontSize: 16,
                color: subtitleColor,
                fontWeight: FontWeight.w400,
              ),
            ),

            const SizedBox(height: 80),

            // Username Field
            _buildTextField(
              controller: nameController,
              label: 'Username',
              maxLength: 20,
              isDark: themeProvider.isDarkMode,
            ),

            const SizedBox(height: 32),

            // Password Field
            _buildTextField(
              controller: passwordController,
              label: 'Password',
              isPassword: true,
              isDark: themeProvider.isDarkMode,
            ),

            const SizedBox(height: 48),

            // Login Button
            _buildLoginButton(context, tokenProvider, userProfileProvider, themeProvider.isDarkMode),

            const SizedBox(height: 40),

            // Divider
            _buildDivider(themeProvider.isDarkMode),

            const SizedBox(height: 40),

            // Social Login
            _buildSocialLoginSection(themeProvider.isDarkMode),

            const SizedBox(height: 40),

            // Forgot Password
            _buildForgotPasswordSection(context, themeProvider.isDarkMode),

            const SizedBox(height: 40),
          ],
        ),
      ),
    );
  }

  Widget _buildTextField({
    required TextEditingController controller,
    required String label,
    required bool isDark,
    bool isPassword = false,
    int? maxLength,
  }) {
    final borderColor = isDark ? Colors.grey[700] : Colors.grey[300];
    final textColor = isDark ? Colors.white : Colors.black;
    final labelColor = isDark ? Colors.grey[400] : Colors.grey[600];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: labelColor,
          ),
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            border: Border(
              bottom: BorderSide(
                color: borderColor!,
                width: 1,
              ),
            ),
          ),
          child: TextField(
            controller: controller,
            maxLength: maxLength,
            obscureText: isPassword ? !isPasswordVisible : false,
            style: TextStyle(
              fontSize: 16,
              color: textColor,
              fontWeight: FontWeight.w400,
            ),
            decoration: InputDecoration(
              border: InputBorder.none,
              contentPadding: const EdgeInsets.symmetric(vertical: 12),
              counterText: '',
              suffixIcon: isPassword
                  ? IconButton(
                      icon: Icon(
                        isPasswordVisible ? Icons.visibility_outlined : Icons.visibility_off_outlined,
                        color: labelColor,
                        size: 20,
                      ),
                      onPressed: () {
                        setState(() {
                          isPasswordVisible = !isPasswordVisible;
                        });
                      },
                    )
                  : null,
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildLoginButton(BuildContext context, AuthTokenProvider tokenProvider,
      UserProfileProvider userProfileProvider, bool isDark) {
    return SizedBox(
      width: double.infinity,
      height: 48,
      child: ElevatedButton(
        onPressed: () => _handleOnLoginButton(
          username: nameController.text,
          password: passwordController.text,
          loginType: "NORMAL",
          context: context,
          tokenProvider: tokenProvider,
          userProfileProvider: userProfileProvider,
        ),
        style: ElevatedButton.styleFrom(
          backgroundColor: isDark ? Colors.white : Colors.black,
          foregroundColor: isDark ? Colors.black : Colors.white,
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(4),
          ),
        ),
        child: const Text(
          'Sign In',
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w500,
            letterSpacing: 0.5,
          ),
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

  Widget _buildSocialLoginSection(bool isDark) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _buildSocialButton(
          icon: FontAwesomeIcons.google,
          onPressed: _handleOnGoogleLogin,
          isDark: isDark,
        ),
        const SizedBox(width: 24),
        _buildSocialButton(
          icon: FontAwesomeIcons.apple,
          onPressed: _handleOnAppleIdLogin,
          isDark: isDark,
        ),
        const SizedBox(width: 24),
        _buildSocialButton(
          icon: FontAwesomeIcons.facebook,
          onPressed: _handleOnFacebookLogin,
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
        borderRadius: BorderRadius.circular(4),
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
          style: TextStyle(
            color: textColor,
            fontSize: 14,
            fontWeight: FontWeight.w400,
            decoration: TextDecoration.underline,
            decorationColor: textColor,
          ),
        ),
      ),
    );
  }

  _handleOnLoginButton({
    required String username,
    required String password,
    required String loginType,
    required BuildContext context,
    required AuthTokenProvider tokenProvider,
    required UserProfileProvider userProfileProvider,
  }) {
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
