import 'dart:convert';

import 'package:http/http.dart' as http;

class AddressAutoComplete {
  Future<List<String>> getAddressSuggestions(String query) async {
    final url = Uri.parse(
        'https://nominatim.openstreetmap.org/search?q=$query&format=json&addressdetails=1&limit=10');
    final response = await http.get(url,
        headers: {'User-Agent': 'PsyConnect/1.0 (codemail.noreply@gmail.com)'});

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return (data as List)
          .map((item) => item['display_name'] as String)
          .toList();
    }
    return [];
  }
}
