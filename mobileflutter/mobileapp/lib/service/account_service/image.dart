import 'dart:io';
import 'dart:typed_data';

class ImageService {
  Future<Uint8List?> readFileByteFromImage(File imageFile) async {
    try {
      Uint8List bytes = await imageFile.readAsBytes();
      print('Reading of bytes is completed');
      return bytes;
    } catch (e) {
      print('Exception Error while reading file from path: $e');
      return null;
    }
  }
}
