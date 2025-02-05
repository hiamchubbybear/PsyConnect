import 'dart:convert';

import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/service/api/api_service.dart';
import 'package:PsyConnect/toasting&loading/toast.dart';

class RegisterService {
  Future<void> registerHandle({
    required Map<String, String> requestBody,
    required UserProvider userProvider,
  }) async {
    final response = await ApiService.post(
      endpoint: "identity/create",
      body: requestBody,
    );
    final responseDecoded = jsonDecode(response.body);
    if (response.statusCode == 200) {
      userProvider.setUser(responseDecoded);
      ToastService.showSuccessToast(message: "Register successful!");
      return jsonDecode(response.body);
    } else if (response.statusCode == 409) {
      ToastService.showErrorToast(message: "username or email already used");
      throw Exception("username or email already used");
    } else if (response.statusCode == 500) {
      throw Exception("Server error : ${response.statusCode}");
    } else {
      ToastService.showErrorToast(message: "There are an orcur error");
      throw Exception("There are an orcur error");
    }
  }
}
