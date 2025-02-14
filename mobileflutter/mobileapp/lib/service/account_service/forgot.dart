import 'package:PsyConnect/service/api/api_service.dart';

class ForgotService {
    DateTime now = DateTime.now();
  registerHandle({required String email, required String token}) async {
    Map<String, dynamic> jsonBody = {
      "token": "$token",
      "email": "$email",
      "verifedTime": now.microsecondsSinceEpoch
    };
    print("${jsonBody}");
    final response =
        await ApiService.post(endpoint: "identity/activate", body: jsonBody);
    if (response.statusCode == 200) {
      print("Activate success");
    } else {

      print("Activate failed");
      print("${response.body}");
    }
  }
}
