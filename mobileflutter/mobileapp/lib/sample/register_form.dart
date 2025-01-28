import 'dart:io';

import 'package:PsyConnect/page/forgot_page.dart';
import 'package:flutter/material.dart';
import 'package:flutter_datetime_picker_plus/flutter_datetime_picker_plus.dart';
import 'package:image_picker/image_picker.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({Key? key}) : super(key: key);

  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {
  File? _image;
  final ImagePicker _picker = ImagePicker();

  Future<void> _pickImage(ImageSource source) async {
    final pickedFile = await _picker.pickImage(source: source);
    if (pickedFile != null) {
      setState(() {
        _image = File(pickedFile.path);
      });
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('No image selected')),
      );
    }
  }

  final TextEditingController usernameController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  final TextEditingController retypePasswordController =
      TextEditingController();
  final TextEditingController emailController = TextEditingController();
  final TextEditingController firstNameController = TextEditingController();
  final TextEditingController lastNameController = TextEditingController();
  final TextEditingController dobController = TextEditingController();
  final TextEditingController addressController = TextEditingController();
  final TextEditingController genderController = TextEditingController();
  final TextEditingController roleController = TextEditingController();
  final TextEditingController imageController = TextEditingController();
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: <Widget>[
            GestureDetector(
              onTap: () {
                showModalBottomSheet(
                  context: context,
                  builder: (context) => Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      ListTile(
                        leading: const Icon(Icons.photo_library),
                        title: const Text('Chọn từ thư viện'),
                        onTap: () {
                          _pickImage(ImageSource.gallery);
                          Navigator.pop(context);
                        },
                      ),
                      ListTile(
                        leading: const Icon(Icons.camera_alt),
                        title: const Text('Chụp ảnh'),
                        onTap: () {
                          _pickImage(ImageSource.camera);
                          Navigator.pop(context);
                        },
                      ),
                    ],
                  ),
                );
              },
              child: CircleAvatar(
                radius: 50,
                backgroundImage: _image != null
                    ? FileImage(_image!)
                    : const NetworkImage(
                        'https://static.vecteezy.com/system/resources/previews/014/194/215/original/avatar-icon-human-a-person-s-badge-social-media-profile-symbol-the-symbol-of-a-person-vector.jpg',
                      ) as ImageProvider,
              ),
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: usernameController,
              decoration: const InputDecoration(labelText: 'User name'),
            ),
            const SizedBox(height: 8),
            TextFormField(
              controller: passwordController,
              obscureText: true,
              decoration: const InputDecoration(labelText: 'Password'),
            ),
            const SizedBox(height: 8),
            TextFormField(
              controller: retypePasswordController,
              obscureText: true,
              decoration: const InputDecoration(labelText: 'Confirm password'),
            ),
            const SizedBox(height: 8),
            TextFormField(
              controller: emailController,
              decoration: const InputDecoration(labelText: 'Email'),
              keyboardType: TextInputType.emailAddress,
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    controller: firstNameController,
                    decoration: const InputDecoration(labelText: 'First name'),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: TextFormField(
                    controller: lastNameController,
                    decoration: const InputDecoration(labelText: 'Last name'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            InkWell(
                onTap: () {
                  DatePicker.showDatePicker(
                    context,
                    showTitleActions: true,
                    minTime: DateTime(1900, 1, 1),
                    maxTime: DateTime.now().add(const Duration(days: 3650)),
                    onConfirm: (date) {
                      setState(() {
                        dobController.text = date.toString().substring(0, 10);
                      });
                    },
                    currentTime: DateTime.now(),
                    locale: LocaleType.vi,
                  );
                },
                child: InputDecorator(
                  decoration: const InputDecoration(
                    labelText: 'Date of birth',
                    hintText: 'Choose date',
                  ),
                  child: Text(dobController.text.isEmpty
                      ? "Choose date"
                      : dobController.text),
                )),
            const SizedBox(height: 8),
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<String>(
                    decoration: const InputDecoration(labelText: 'Wanna be ?'),
                    value: roleController.text.isEmpty
                        ? null
                        : roleController.text,
                    items: const [
                      DropdownMenuItem(value: "Client", child: Text("Client")),
                      DropdownMenuItem(
                          value: "Therapist", child: Text("Thrapist")),
                    ],
                    onChanged: (value) =>
                        setState(() => roleController.text = value!),
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: DropdownButtonFormField<String>(
                    decoration: const InputDecoration(labelText: 'Giới tính'),
                    value: genderController.text.isEmpty
                        ? null
                        : genderController.text,
                    items: const [
                      DropdownMenuItem(value: "Male", child: Text("Nam")),
                      DropdownMenuItem(value: "Female", child: Text("Nữ")),
                      DropdownMenuItem(
                          value: "None", child: Text("Không muốn tiết lộ")),
                    ],
                    onChanged: (value) =>
                        setState(() => genderController.text = value!),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            TextFormField(
              controller: addressController,
              decoration: const InputDecoration(labelText: 'Địa chỉ'),
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () {},
              child: const Text('Đăng ký'),
            ),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                IconButton(
                  onPressed: () => _handleOnAppleIdRegister(),
                  icon: Image.asset("assets/images/login_appleid_icon.png",
                      width: 30),
                ),
                IconButton(
                  onPressed: () => _handleOnGoogleRegister(),
                  icon: Image.asset("assets/images/login_google_icon.jpeg",
                      width: 30),
                ),
                IconButton(
                  onPressed: () => _handleOnFacebookRegister(),
                  icon: Image.asset("assets/images/login_facebook_icon.jpeg",
                      width: 30),
                ),
              ],
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Text("Quên mật khẩu?"),
                TextButton(
                  onPressed: () {
                    Navigator.push(
                      context,
                      MaterialPageRoute(builder: (context) => ForgotPassword()),
                    );
                  },
                  child: const Text('Đặt lại ngay'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  _handleOnAppleIdRegister() {}
  _handleOnGoogleRegister() {}

  _handleOnFacebookRegister() {}
}
