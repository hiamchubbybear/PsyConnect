import 'dart:core';
import 'dart:io';

import 'package:PsyConnect/page/forgot_page.dart';
import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/service/account_service/image.dart';
import 'package:PsyConnect/service/account_service/register.dart';
import 'package:PsyConnect/service/api/cloudinary_api_service.dart';
import 'package:PsyConnect/service/logic.dart';
import 'package:PsyConnect/toasting&loading/toast.dart';
import 'package:flutter/material.dart';
import 'package:flutter_datetime_picker_plus/flutter_datetime_picker_plus.dart';
import 'package:image_picker/image_picker.dart';
import 'package:provider/provider.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({Key? key}) : super(key: key);

  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

final _accountFormKey = GlobalKey<FormState>();

class _RegisterPageState extends State<RegisterPage> {
  bool _isLoading = false;
  RegisterService registerService = RegisterService();
  File? _image;
  final ImagePicker _picker = ImagePicker();
  bool showProfileFields = false;
  _setStartLoading() {
    setState(() {
      _isLoading = true;
    });
    Future.delayed(const Duration(seconds: 3), () {
      setState(() {
        _isLoading = false;
      });
    });
  }

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

  final CloudinaryApiService cloudinaryApiService = CloudinaryApiService();
  final ImageService imageService = ImageService();
  BusinessLogic businessLogic = BusinessLogic();
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
    final userProvider = context.read<UserProvider>();
    return Scaffold(
      body: Form(
        key: _accountFormKey,
        child: SingleChildScrollView(
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
                            title: const Text('Choose from gallery'),
                            onTap: () {
                              _pickImage(ImageSource.gallery);
                              Navigator.pop(context);
                            },
                          ),
                          ListTile(
                            contentPadding: EdgeInsets.all(10),
                            leading: const Icon(Icons.camera_alt),
                            title: const Text('Take a photo'),
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
                    radius: 60,
                    backgroundImage: _image != null
                        ? FileImage(_image!)
                        : const NetworkImage(
                            'https://static.vecteezy.com/system/resources/previews/014/194/215/original/avatar-icon-human-a-person-s-badge-social-media-profile-symbol-the-symbol-of-a-person-vector.jpg',
                          ) as ImageProvider,
                    backgroundColor: Colors.transparent,
                  )),
              const SizedBox(height: 16),
              ElevatedButton(
                style: ButtonStyle(
                  elevation: MaterialStateProperty.all(
                      showProfileFields ? 10 : 2), // Increase shadow
                  shape: MaterialStateProperty.all(
                    RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(20), // Border Radius
                    ),
                  ),
                  padding: MaterialStateProperty.all(
                    EdgeInsets.symmetric(
                        vertical: 10, horizontal: 15), // Distances in button
                  ),
                ),
                onPressed: () {
                  setState(() {
                    showProfileFields = !showProfileFields;
                  });
                },
                child: Text(showProfileFields ? 'Profile' : 'Profile',
                    style: (showProfileFields)
                        ? const TextStyle(
                            fontSize: 16,
                          )
                        : const TextStyle(
                            fontSize: 16,
                          )),
              ),
              if (showProfileFields)
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                            child: TextFormField(
                          textCapitalization: TextCapitalization.sentences,
                          validator: (value) => value == null
                              ? "Hey you forgot field first name"
                              : null,
                          controller: firstNameController,
                          decoration: const InputDecoration(
                              labelText: 'My first name is'),
                          onChanged: (value) => setState(() {
                            firstNameController.text = value;
                          }),
                        )),
                        const SizedBox(width: 8),
                        Expanded(
                          child: TextFormField(
                            textCapitalization: TextCapitalization.sentences,
                            validator: (value) => value == null
                                ? "Hey you forgot field last name"
                                : null,
                            controller: lastNameController,
                            decoration: const InputDecoration(
                                labelText: 'My last name is'),
                            onChanged: (value) => setState(() {
                              lastNameController.text = value;
                            }),
                          ),
                        )
                      ],
                    ),
                    const SizedBox(height: 8),
                    InkWell(
                        onTap: () {
                          DatePicker.showDatePicker(
                            context,
                            showTitleActions: true,
                            minTime: DateTime(1900, 1, 1),
                            maxTime: DateTime.now(),
                            onConfirm: (date) {
                              setState(() {
                                dobController.text =
                                    "${date.year}-${date.month}-${date.day}";
                                print("${dobController.text}");
                              });
                            },
                            currentTime: DateTime.now(),
                            locale: LocaleType.vi,
                          );
                        },
                        child: InputDecorator(
                          decoration: const InputDecoration(
                            labelText: 'My birthday in',
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
                            validator: (value) =>
                                value == null ? "What do you wanna be" : null,
                            decoration:
                                const InputDecoration(labelText: 'Wanna be a'),
                            value: roleController.text.isEmpty
                                ? null
                                : roleController.text,
                            items: const [
                              DropdownMenuItem(
                                  value: "Client",
                                  child: Text("Ordinary user")),
                              DropdownMenuItem(
                                  value: "Therapist", child: Text("Thrapist ")),
                            ],
                            onChanged: (value) =>
                                setState(() => roleController.text = value!),
                          ),
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: DropdownButtonFormField<String>(
                            decoration:
                                const InputDecoration(labelText: 'I am a'),
                            value: genderController.text.isEmpty
                                ? null
                                : genderController.text,
                            items: const [
                              DropdownMenuItem(
                                  value: "Male", child: Text("Male")),
                              DropdownMenuItem(
                                  value: "Female", child: Text("Female")),
                              DropdownMenuItem(
                                  value: "None", child: Text("Keep it secret")),
                            ],
                            onChanged: (value) =>
                                setState(() => genderController.text = value!),
                            validator: (value) => value == null
                                ? "Please choose your gender"
                                : null,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    TextFormField(
                        validator: (value) {
                          (value == null || value.isEmpty)
                              ? 'Address is required'
                              : null;
                        },
                        controller: addressController,
                        textInputAction: TextInputAction.route,
                        decoration: const InputDecoration(labelText: 'Address'),
                        onChanged: (value) => setState(() {
                              addressController.text = value;
                            })),
                  ],
                ),
              const SizedBox(height: 16),
              TextFormField(
                validator: (username) {
                  (username == null || username!.length <= 3)
                      ? "Username must be greater than 3 characters"
                      : null;
                },
                controller: usernameController,
                decoration: const InputDecoration(labelText: 'Username'),
              ),
              const SizedBox(height: 8),
              TextFormField(
                controller: passwordController,
                obscureText: true,
                decoration: const InputDecoration(labelText: 'Password'),
                validator: (password) {
                  RegExp passReg = RegExp(
                      r"^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$");
                  if (password == null ||
                      !passReg.hasMatch(password!) ||
                      password.length > 3) {
                    return "Password must than 3 characters and include character , number  ";
                  }
                  return null;
                },
              ),
              const SizedBox(height: 8),
              TextFormField(
                onChanged: (value) => {retypePasswordController.text = value},
                controller: retypePasswordController,
                obscureText: true,
                decoration:
                    const InputDecoration(labelText: 'Confirm password'),
                validator: (retypePass) {
                  if (retypePass != passwordController.text ||
                      retypePass == null) {
                    return "Password and confirm password doesn't match";
                  }
                  return null;
                },
              ),
              const SizedBox(height: 8),
              TextFormField(
                  controller: emailController,
                  decoration: const InputDecoration(labelText: 'Email'),
                  keyboardType: TextInputType.emailAddress,
                  validator: (email) {
                    RegExp mailRegex = RegExp(
                        r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$");
                    if (email == null || email.isEmpty) {
                      return "Email cannot be empty";
                    }
                    if (!mailRegex.hasMatch(email)) {
                      return "Invalid email format";
                    }
                    return null;
                  }),
              const SizedBox(height: 8),
              SizedBox(height: 16),
              ElevatedButton(
                onPressed: () {
                  if (!_accountFormKey.currentState!.validate()) {
                    _setStartLoading();
                    _handleOnRegisterButton(userProvider: userProvider , context: context);

                  }
                },
                child: (!_isLoading)
                    ? const Text('Register')
                    : const CircularProgressIndicator(),
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
                  const Text("Forgot your password?"),
                  TextButton(
                    onPressed: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                            builder: (context) => const ForgotPage()),
                      );
                    },
                    child: const Text('Reset it'),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  _handleOnAppleIdRegister() {}
  _handleOnGoogleRegister() {}

  _handleOnFacebookRegister() {}

  Future<void> _handleOnRegisterButton({required dynamic userProvider , required BuildContext context}) async {
    if (_image != null) {
      try {
        String? cloudinaryUrlImage = await cloudinaryApiService.uploadImage(
            imageFile: _image as File, username: usernameController.text);
        print(cloudinaryUrlImage);
        if (cloudinaryUrlImage == null) {
          ToastService.showErrorToast(message: "Image doesn't exceptable");
        }
        Map<String, String> jsonData = {
          "username": usernameController.text.toLowerCase(),
          "password": passwordController.text,
          "firstName": firstNameController.text,
          "lastName": lastNameController.text,
          "address": addressController.text,
          "gender": genderController.text,
          "email": emailController.text,
          "role": roleController.text,
          "avatarUri": cloudinaryUrlImage.toString() ?? ""
        };
        await registerService.registerHandle(
            requestBody: jsonData, userProvider: userProvider , context : context);
      } catch (e) {
        ToastService.showErrorToast(message: "Some error ${e}");
      }
    }
  }
}
