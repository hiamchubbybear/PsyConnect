import 'dart:convert';
import 'dart:io';

import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:http/http.dart' as http;

class CloudinaryApiService {
  final String cloudName = dotenv.env['CLOUDINARY_CLOUD_NAME'] ?? "";
  final String uploadPreset = dotenv.env['CLOUDINARY_UPLOAD_PRESET'] ?? "";

  Future<String?> uploadImage({required File imageFile,required String username}) async {
    try {
      if (imageFile.readAsBytesSync().length > 9000000000) {
        throw Exception("Your image have to be less than 10 mb");
      }
      var url =
          Uri.parse("https://api.cloudinary.com/v1_1/$cloudName/image/upload");

      var request = http.MultipartRequest("POST", url)
        ..fields['upload_preset'] = uploadPreset
        ..files.add(await http.MultipartFile.fromPath("file", imageFile.path))
        ..fields['public_id'] = username.toLowerCase()
        ..fields['filename_override'] = username.toLowerCase();

      var response = await request.send();
      var responseData = await response.stream.bytesToString();
      var jsonResponse = jsonDecode(responseData);

      if (response.statusCode == 200) {
        return jsonResponse["secure_url"];
      } else {
        print("Upload thất bại: ${jsonResponse["error"]["message"]}");
        return null;
      }
    } catch (e) {
      print("Lỗi khi upload ảnh: $e");
      return null;
    }
  }
}
