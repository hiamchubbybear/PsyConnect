import 'package:PsyConnect/service/api/api_service.dart';
import 'package:PsyConnect/toasting&loading/toast.dart';
import 'package:flutter/material.dart';

class ForgotService {
  DateTime now = DateTime.now();
  Future<bool> registerHandle(
      {required String email,
      required String token,
      required BuildContext context}) async {
    print("Token $token");
    Map<String, dynamic> jsonBody = {
      "token": "$token",
      "email": "$email",
      "verifedTime": now.microsecondsSinceEpoch
    };
    print("${jsonBody}");
    final response =
        await ApiService.post(endpoint: "identity/activate", body: jsonBody);
    if (response.statusCode == 200) {
      ToastService.showToast(
          context: context,
          message: "Acitvate your account successfull",
          title: "Success",
          type: ToastType.success);
      return true;
    } else {
      ToastService.showToast(
          context: context,
          message: "Wrong activate token!!",
          title: "Acitvate failed",
          type: ToastType.warning);
      return false;
    }
  }
}
