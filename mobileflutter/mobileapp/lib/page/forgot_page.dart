import 'package:flutter/material.dart';

class VerifiedPage extends StatelessWidget {
  const VerifiedPage({super.key});

  @override
  Widget build(BuildContext context) {
    final emailConfirmController = TextEditingController();
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
                TextFormField(
                    decoration: const InputDecoration(
                      border: UnderlineInputBorder(),
                      labelText: 'Verified email',
                    ),
                    onChanged: (value) {
                      setState() {
                        emailConfirmController.text = value!;
                      }
                    }),
                const SizedBox(height: 10),
                ElevatedButton(
                  onPressed: () =>
                      _onVerifiedHandle(email: emailConfirmController.text),
                  child: const Text('Send'),
                ),
              ],
            ),
          ),
        ),
      ),
      
    );
  }

  void _onVerifiedHandle({required String email}) {

  }
}
