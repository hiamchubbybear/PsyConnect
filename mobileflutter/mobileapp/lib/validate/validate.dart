import 'package:http/http.dart' as http;

Future<bool> checkImageExists(String url) async {
  try {
    final response = await http.head(Uri.parse(url));
    return response.statusCode == 200;
  } catch (_) {
    return false;
  }
}
