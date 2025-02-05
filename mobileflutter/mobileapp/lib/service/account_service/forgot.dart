import 'package:PsyConnect/service/api/api_service.dart';

class ForgotService {
  registerHandle({required String email}) async {
    final response = await ApiService.postWithParameters(
        endpoint: "emai;", param: "email", paramValue: email);
    if (response.statusCode == 200) {
      // Navigator.push();
    }
  }
}
