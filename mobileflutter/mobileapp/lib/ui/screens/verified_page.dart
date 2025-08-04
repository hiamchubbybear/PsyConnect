import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/services/account_service/forgot.dart';
import 'package:PsyConnect/ui/screens/login_page.dart';
import 'package:flutter/material.dart';
import 'package:pin_code_fields/pin_code_fields.dart';
import 'package:provider/provider.dart';

class VerifiedPage extends StatefulWidget {
  const VerifiedPage({super.key});

  @override
  State<VerifiedPage> createState() => _VerifiedPageState();
}

class _VerifiedPageState extends State<VerifiedPage> {
  final TextEditingController verifiedCodeController = TextEditingController();
  final ForgotService forgotService = ForgotService();

  @override
  void dispose() {
    verifiedCodeController.dispose();
    super.dispose();
  }

  Future<void> _onVerifiedCode({
    required String token,
    required String email,
  }) async {
    final isSuccess = await forgotService.registerHandle(
      email: email,
      token: token,
      context: context,
    );

    if (isSuccess && mounted) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (_) => const LoginPage()),
      );
    } else {
      verifiedCodeController.clear();
    }
  }

  @override
  Widget build(BuildContext context) {
    final userProvider = Provider.of<UserProvider>(context);
    final username = userProvider.user?["username"] as String?;
    final email = userProvider.user?["email"] as String?;

    final theme = Theme.of(context);

    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 32.0),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  "Verify your account",
                  style: subHeadingStyle,
                ),
                const SizedBox(height: 8),
                if (email != null)
                  Text("We sent a code to $email",
                      style: kSecondarirySubHeadingVerifiedPage),
                const SizedBox(height: 32),
                PinCodeTextField(
                  controller: verifiedCodeController,
                  onChanged: (_) {},
                  pinTheme: PinTheme(
                    shape: PinCodeFieldShape.underline,
                    fieldHeight: 50,
                    fieldWidth: 40,
                    activeFillColor: Colors.transparent,
                    inactiveColor: Colors.grey.shade400,
                    selectedColor: Colors.black,
                  ),
                  keyboardType: TextInputType.number,
                  appContext: context,
                  length: 5,
                ),
                const SizedBox(height: 32),
                SizedBox(
                  width: double.infinity,
                  child: FilledButton(
                    onPressed: (email != null &&
                            verifiedCodeController.text.isNotEmpty)
                        ? () {
                            _onVerifiedCode(
                              token: verifiedCodeController.text,
                              email: email,
                            );

                          }
                        : null,
                    style: FilledButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 16.0),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                    ),
                    child: Text(
                      'Confirm',
                      style: kSecondarirySubHeadingRegisterPage,
                    ),
                  ),
                ),
                const SizedBox(height: 16),
                Center(
                  child: TextButton(
                    onPressed: () {},
                    child: Text(
                      "Resend code",
                      style: kSecondarirySubHeadingRegisterPage,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
