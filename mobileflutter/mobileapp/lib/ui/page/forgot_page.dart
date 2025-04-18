import 'package:PsyConnect/service/account_service/forgot.dart';
import 'package:flutter/material.dart';
import 'package:flutter_otp_text_field/flutter_otp_text_field.dart';

class ForgotPage extends StatefulWidget {
  const ForgotPage({super.key});

  @override
  State<ForgotPage> createState() => _ForgotPageState();
}

class _ForgotPageState extends State<ForgotPage> {
  final TextEditingController verifiedCodeController = TextEditingController();
  ForgotService forgotService = ForgotService();

  bool isSent = false;
  String savedEmail = "";

  void _changeSentStatus(String email) {
    setState(() {
      isSent = true;
      savedEmail = email;
    });
  }

  void _onVerifiedHandle({required String email}) {
    print("Processing email: $email");
    _changeSentStatus(email);
  }

  void _onVerifiedCode({required String token}) {
    if (savedEmail.isNotEmpty) {
      forgotService.registerHandle(
          email: savedEmail, token: token, context: context);
    } else {
      print("Error: Email is missing");
    }
  }

  @override
  Widget build(BuildContext context) {
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
                if (!isSent) ...[
                  ElevatedButton(
                    onPressed: () {
                      String email = "user@example.com";
                      print("Sending request for email: $email");
                      _onVerifiedHandle(email: email);
                    },
                    child: const Text('Send Verification Code'),
                  )
                ] else ...[
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
                    onPressed: () {
                      _onVerifiedCode(token: verifiedCodeController.text);
                    },
                    child: const Text('Confirm Code'),
                  ),
                ]
              ],
            ),
          ),
        ),
      ),
    );
  }
}
