import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/service/account_service/forgot.dart';
import 'package:flutter/material.dart';
import 'package:flutter_otp_text_field/flutter_otp_text_field.dart';
import 'package:provider/provider.dart';

class VerifiedPage extends StatefulWidget {
  const VerifiedPage({super.key});

  @override
  State<VerifiedPage> createState() => _VerifiedPageState();
}

class _VerifiedPageState extends State<VerifiedPage> {
  final TextEditingController verifiedCodeController = TextEditingController();
  ForgotService forgotService = ForgotService();

  void _onVerifiedCode({required String token, required String email}) {
    forgotService.registerHandle(email: email, token: token);
  }

  @override
  Widget build(BuildContext context) {
    final userProvider = Provider.of<UserProvider>(context);
    String? email = userProvider.user?["email"];

    return Scaffold(
      body: Align(
        alignment: Alignment.center,
        child: Center(
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Image.asset("assets/gifs/forgot_password_animation.gif"),
                const SizedBox(height: 10),
                if (email != null)
                  Text("Verification for: $email",
                      style:
                          TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                const SizedBox(height: 10),
                OtpTextField(
                  numberOfFields: 5,
                  showFieldAsBox: false,
                  borderWidth: 4.0,
                  onSubmit: (String value) {
                    setState(() {
                      verifiedCodeController.text = value;
                    });
                  },
                ),
                const SizedBox(height: 10),
                ElevatedButton(
                  onPressed: email != null
                      ? () {
                          _onVerifiedCode(
                              token: verifiedCodeController.text, email: email);
                        }
                      : null,
                  child: const Text('Confirm Code'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
