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
  ForgotService forgotService = ForgotService();

  void _onVerifiedCode({required String token, required String email}) async {
    var resource = await forgotService.registerHandle(
        email: email, token: token, context: context);
    print("Response from server $resource");
    if (resource) {
      Navigator.pushReplacement((context),
          MaterialPageRoute(builder: (context) => const LoginPage()));
    } else {
      verifiedCodeController.clear();
    }
  }

  @override
  Widget build(BuildContext context) {
    final userProvider = Provider.of<UserProvider>(context);
    String? username = userProvider.user?["username"];
    String? email = userProvider.user?["email"];

    return Scaffold(
      body: Align(
        alignment: Alignment.center,
        child: Center(
          child: Padding(
            padding: const EdgeInsets.all(60.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Image.asset("assets/gifs/forgot_password_animation.gif"),
                const SizedBox(height: 10),
                if (email != null)
                  Text("Verification for: $username",
                      style:
                          TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                const SizedBox(height: 10),
                PinCodeTextField(
                  onChanged: (String value) {
                    setState(() {
                      verifiedCodeController.text = value;
                    });
                  },
                  pinTheme: PinTheme(
                    shape: PinCodeFieldShape.underline,
                    borderRadius: BorderRadius.circular(5),
                    fieldHeight: 40,
                    fieldWidth: 30,
                    inactiveColor: blackColor,
                    activeFillColor: whiteColor,
                  ),
                  keyboardType: TextInputType.number,
                  appContext: context,
                  length: 5,
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
