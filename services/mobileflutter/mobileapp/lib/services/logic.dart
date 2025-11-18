import 'dart:convert';
import 'dart:io';

class BusinessLogic {
   String? encodeImage(File? image) {
    if (image == null) return null;
    List<int> imageConvert = image.readAsBytesSync();
    return base64Encode(imageConvert);
  }
}
