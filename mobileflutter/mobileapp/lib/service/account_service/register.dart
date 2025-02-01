import 'dart:convert';
import 'dart:io';

import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/toasting&loading/toast.dart';
import 'package:PsyConnect/variable/variable.dart';
import 'package:http/http.dart' as http;

class RegisterService {
  Future<void> registerHandle({
    required Map<String, dynamic> requestBody,
    required UserProvider userProvider,
    required File file,
  }) async {
    var uri = Uri.parse((Platform.isAndroid)
        ? androidBaseUrlProfileService
        : iosBaseUrlProfileService);

    var request = http.MultipartRequest("POST", uri);
    request.files.add(
      await http.MultipartFile.fromPath(
        "file",
        file.path,
        filename: requestBody['username'],
      ),
    );
    request.fields['json'] = jsonEncode(requestBody);

    request.headers['Content-Type'] = "multipart/form-data";
    print("${request.headers}");
    print("${request.method}");
    print("${request.url}");
    try {
      var streamedResponse = await request.send();
      var response = await http.Response.fromStream(streamedResponse);

      if (response.statusCode == 200) {
        var responseData = jsonDecode(response.body);
        userProvider.setUser(responseData);
        ToastService.showSuccessToast(message: "Register successful!");
      } else {
        ToastService.showErrorToast(
          message: "Register failed: ${response.body}",
        );
      }
    } catch (error) {
      ToastService.showErrorToast(
        message: "There was an error while creating the account: $error",
      );
    }
  }
}
