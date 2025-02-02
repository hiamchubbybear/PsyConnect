import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/service/api/api_service.dart';
import 'package:PsyConnect/toasting&loading/toast.dart';

class RegisterService {
  Future<void> registerHandle({
    required Map<String, String> requestBody,
    required UserProvider userProvider,
  }) async {
    try {
      final responseData = await ApiService.post(
        endpoint: "identity/create",
        body: requestBody,
      );
      userProvider.setUser(responseData);
      ToastService.showSuccessToast(message: "Register successful!");
    } catch (error) {
      ToastService.showErrorToast(
        message: "There was an error while creating the account: ${error.toString()}",
      );
    }
  }
}
