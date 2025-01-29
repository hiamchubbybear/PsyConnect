import 'dart:io';

import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/toasting&loading/toast.dart';
import 'package:http/http.dart' as http;
import 'package:http_parser/http_parser.dart';
class RegisterService {

Future<void> registerHandle({
  required Map<String, String> requestBody,
  required UserProvider userProvider,
  required File? image
}) async {
  try {
    var request = http.MultipartRequest(
      'POST',
      Uri.parse('http://localhost:8080/identity/create'),
    );
    requestBody.forEach((key, value) {
      request.fields[key] = value;
    });
    if (image != null) {
      request.files.add(
        await http.MultipartFile.fromPath(
          'avatar', image!.path,
          contentType: MediaType('image', 'jpg'),
        ),
      );
    }
    print("Request ${requestBody} and ${image.toString()}");

    var response = await request.send();
    var responseBody = await response;

    if (response.statusCode == 200) {
      userProvider.setUser(responseBody as Map<String, dynamic>);
      ToastService.showSuccessToast(message: "Register successful!");
    } else {
      ToastService.showErrorToast(
        message: "Register failed: ${response.reasonPhrase}",
      );
    }
  } catch (error) {
    ToastService.showErrorToast(
      message: "There was an error while creating the account: $error",
    );
  }
}
}
