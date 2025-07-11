import 'dart:core';
import 'dart:io';

import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/services/account_service/cloudinary_service.dart';
import 'package:PsyConnect/services/account_service/image.dart';
import 'package:PsyConnect/services/account_service/register.dart';
import 'package:PsyConnect/services/api/address_auto_complete.dart';
import 'package:PsyConnect/services/logic.dart';
import 'package:flutter/material.dart';
import 'package:flutter_datetime_picker_plus/flutter_datetime_picker_plus.dart';
import 'package:flutter_typeahead/flutter_typeahead.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:image_picker/image_picker.dart';
import 'package:intl/intl.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:provider/provider.dart';

class MultiStepRegisterPage extends StatefulWidget {
  const MultiStepRegisterPage({super.key});

  @override
  State<MultiStepRegisterPage> createState() => _MultiStepRegisterPageState();
}

class _MultiStepRegisterPageState extends State<MultiStepRegisterPage> {
  ThemeProvider themeProvider = ThemeProvider();
  final PageController _pageController = PageController();
  int _currentStep = 0;
  final int _totalSteps = 7;
  bool _isLoading = false;
  RegisterService registerService = RegisterService();
  final CloudinaryApiService cloudinaryApiService = CloudinaryApiService();
  final ImageService imageService = ImageService();
  BusinessLogic businessLogic = BusinessLogic();

  File? _image;
  final ImagePicker _picker = ImagePicker();

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
  List<String> suggestions = [];
  bool isLoading = false;

  void onAddressChanged(String query) async {
    if (query.isEmpty) {
      setState(() => suggestions = []);
      return;
    }

    setState(() => isLoading = true);

    final result = await AddressAutoComplete().getAddressSuggestions(query);
    setState(() {
      suggestions = result;
      isLoading = false;
    });
  }

  final List<GlobalKey<FormState>> _formKeys =
      List.generate(7, (index) => GlobalKey<FormState>());

  void _nextStep() {
    if (_currentStep < _totalSteps - 1) {
      if (_formKeys[_currentStep].currentState?.validate() ?? true) {
        setState(() {
          _currentStep++;
        });
        _pageController.nextPage(
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeInOut,
        );
      }
    }
  }

  void _previousStep() {
    if (_currentStep > 0) {
      setState(() {
        _currentStep--;
      });
      _pageController.previousPage(
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeInOut,
      );
    }
  }

  Future<void> _pickImage(ImageSource source) async {
    try {
      PermissionStatus status;
      if (source == ImageSource.camera) {
        status = await Permission.camera.request();
      } else {
        status = await Permission.photos.request();
      }
      if (status.isDenied) {
        ToastService.showToast(
          message:
              "Permission denied. Please allow access to ${source == ImageSource.camera ? 'camera' : 'photos'}.",
          context: context,
          title: 'Permission Error',
          type: ToastType.error,
        );
        return;
      }

      if (status.isPermanentlyDenied) {
        ToastService.showToast(
          message:
              "Permission permanently denied. Please enable it in Settings.",
          context: context,
          title: 'Permission Error',
          type: ToastType.error,
        );
        await openAppSettings();
        return;
      }

      final pickedFile = await _picker.pickImage(
        source: source,
        maxHeight: 800,
        maxWidth: 800,
        imageQuality: 85,
      );

      if (pickedFile != null) {
        setState(() {
          _image = File(pickedFile.path);
        });
      } else {
        ToastService.showToast(
          message: "No image selected!",
          context: context,
          title: 'Image Error',
          type: ToastType.error,
        );
      }
    } catch (e) {
      ToastService.showToast(
        message: "Error picking image: $e",
        context: context,
        title: 'Image Error',
        type: ToastType.error,
      );
    }
  }

  void _setStartLoading() {
    setState(() {
      _isLoading = true;
    });
    Future.delayed(const Duration(seconds: 3), () {
      setState(() {
        _isLoading = false;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    final userProvider = context.read<UserProvider>();
    final isDarkMode = themeProvider.isDarkMode;

    return Scaffold(
      backgroundColor: isDarkMode ? Colors.grey[900] : Colors.grey[50],
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: _currentStep > 0
            ? IconButton(
                icon: Icon(Icons.arrow_back_ios,
                    color: isDarkMode ? Colors.white : Colors.black87),
                onPressed: _previousStep,
              )
            : null,
        title: Text(
          'Create Account',
          style: GoogleFonts.quicksand(
            fontSize: 20,
            fontWeight: FontWeight.w600,
            color: isDarkMode ? Colors.white : Colors.black87,
          ),
        ),
        centerTitle: true,
      ),
      body: Column(
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
            child: LinearProgressIndicator(
              value: (_currentStep + 1) / _totalSteps,
              backgroundColor: isDarkMode ? Colors.grey[700] : Colors.grey[300],
              valueColor: AlwaysStoppedAnimation<Color>(
                  isDarkMode ? Colors.blue[300]! : Colors.blue),
              minHeight: 4,
            ),
          ),
          Text(
            'Step ${_currentStep + 1} of $_totalSteps',
            style: GoogleFonts.quicksand(
              fontSize: 14,
              color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
            ),
          ),
          const SizedBox(height: 20),
          Expanded(
            child: PageView(
              controller: _pageController,
              physics: const NeverScrollableScrollPhysics(),
              children: [
                _buildAvatarStep(),
                _buildNameStep(),
                _buildBirthGenderStep(),
                _buildAddressStep(),
                _buildEmailStep(),
                _buildRoleStep(),
                _buildCredentialsStep(),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.all(20),
            child: Row(
              children: [
                if (_currentStep > 0)
                  Expanded(
                    child: OutlinedButton(
                      onPressed: _previousStep,
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 15),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(10),
                        ),
                        side: BorderSide(
                            color:
                                isDarkMode ? Colors.grey[600]! : Colors.grey),
                      ),
                      child: Text(
                        'Back',
                        style: GoogleFonts.quicksand(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: isDarkMode ? Colors.white : Colors.black87,
                        ),
                      ),
                    ),
                  ),
                if (_currentStep > 0) const SizedBox(width: 15),
                Expanded(
                  flex: _currentStep == 0 ? 1 : 1,
                  child: ElevatedButton(
                    onPressed: _currentStep == _totalSteps - 1
                        ? () => _handleFinalSubmit(userProvider, context)
                        : _nextStep,
                    style: ElevatedButton.styleFrom(
                      backgroundColor:
                          isDarkMode ? Colors.blue[700] : Colors.blue,
                      padding: const EdgeInsets.symmetric(vertical: 15),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(10),
                      ),
                    ),
                    child: _isLoading && _currentStep == _totalSteps - 1
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              valueColor:
                                  AlwaysStoppedAnimation<Color>(Colors.white),
                            ),
                          )
                        : Text(
                            _currentStep == _totalSteps - 1
                                ? 'Create Account'
                                : 'Continue',
                            style: GoogleFonts.quicksand(
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                              color: Colors.white,
                            ),
                          ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAvatarStep() {
    final isDarkMode = context.watch<ThemeProvider>().isDarkMode;
    return Form(
      key: _formKeys[0],
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Choose Your Profile Picture',
              style: kSubHeadingRegisterPage.copyWith(
                color: isDarkMode ? Colors.white : Colors.black87,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 10),
            Text(
              'This helps others recognize you',
              style: kSecondarirySubHeadingRegisterPage.copyWith(
                color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 40),
            GestureDetector(
              onTap: () {
                showModalBottomSheet(
                  context: context,
                  shape: const RoundedRectangleBorder(
                    borderRadius:
                        BorderRadius.vertical(top: Radius.circular(20)),
                  ),
                  backgroundColor: isDarkMode ? Colors.grey[900] : Colors.white,
                  builder: (context) => Container(
                    padding: const EdgeInsets.all(20),
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          'Select Photo',
                          style: GoogleFonts.quicksand(
                            fontSize: 18,
                            fontWeight: FontWeight.w600,
                            color: isDarkMode ? Colors.white : Colors.black87,
                          ),
                        ),
                        const SizedBox(height: 20),
                        ListTile(
                          leading: Icon(Icons.photo_library,
                              color:
                                  isDarkMode ? Colors.blue[300] : Colors.blue),
                          title: Text(
                            'Choose from gallery',
                            style: GoogleFonts.quicksand(
                              fontSize: 16,
                              color: isDarkMode ? Colors.white : Colors.black87,
                            ),
                          ),
                          onTap: () {
                            _pickImage(ImageSource.gallery);
                            Navigator.pop(context);
                          },
                        ),
                        ListTile(
                          leading: Icon(Icons.camera_alt,
                              color:
                                  isDarkMode ? Colors.blue[300] : Colors.blue),
                          title: Text(
                            'Take a photo',
                            style: GoogleFonts.quicksand(
                              fontSize: 16,
                              color: isDarkMode ? Colors.white : Colors.black87,
                            ),
                          ),
                          onTap: () {
                            _pickImage(ImageSource.camera);
                            Navigator.pop(context);
                          },
                        ),
                      ],
                    ),
                  ),
                );
              },
              child: Container(
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  border: Border.all(
                    color: _image != null
                        ? (isDarkMode ? Colors.blue[300]! : Colors.blue)
                        : (isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                    width: 3,
                  ),
                ),
                child: CircleAvatar(
                  radius: 80,
                  backgroundImage: _image != null
                      ? FileImage(_image!)
                      : const NetworkImage(
                          'assets/images/avatar.jpg',
                        ) as ImageProvider,
                  backgroundColor:
                      isDarkMode ? Colors.grey[800] : Colors.grey[100],
                ),
              ),
            ),
            const SizedBox(height: 20),
            if (_image == null)
              Text(
                'Tap to add photo',
                style: GoogleFonts.quicksand(
                  fontSize: 14,
                  color: isDarkMode ? Colors.grey[400] : Colors.grey[500],
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildNameStep() {
    final isDarkMode = themeProvider.isDarkMode;
    return Form(
      key: _formKeys[1],
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'What\'s Your Name?',
              style: GoogleFonts.quicksand(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: isDarkMode ? Colors.white : Colors.black87,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 10),
            Text(
              'Let us know how to address you',
              style: GoogleFonts.quicksand(
                fontSize: 16,
                color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 40),
            TextFormField(
              controller: firstNameController,
              textCapitalization: TextCapitalization.words,
              style: GoogleFonts.quicksand(
                  fontSize: 16,
                  color: isDarkMode ? Colors.white : Colors.black87),
              decoration: InputDecoration(
                labelText: 'First Name',
                labelStyle: GoogleFonts.quicksand(
                    color: isDarkMode ? Colors.grey[400] : Colors.grey[600]),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color:
                          isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color: isDarkMode ? Colors.blue[300]! : Colors.blue,
                      width: 2),
                ),
                fillColor: isDarkMode ? Colors.grey[800] : Colors.white,
                filled: true,
              ),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter your first name';
                }
                return null;
              },
            ),
            const SizedBox(height: 20),
            TextFormField(
              controller: lastNameController,
              textCapitalization: TextCapitalization.words,
              style: GoogleFonts.quicksand(
                  fontSize: 16,
                  color: isDarkMode ? Colors.white : Colors.black87),
              decoration: InputDecoration(
                labelText: 'Last Name',
                labelStyle: GoogleFonts.quicksand(
                    color: isDarkMode ? Colors.grey[400] : Colors.grey[600]),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color:
                          isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color: isDarkMode ? Colors.blue[300]! : Colors.blue,
                      width: 2),
                ),
                fillColor: isDarkMode ? Colors.grey[800] : Colors.white,
                filled: true,
              ),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter your last name';
                }
                return null;
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildBirthGenderStep() {
    final isDarkMode = themeProvider.isDarkMode;
    return Form(
      key: _formKeys[2],
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Personal Information',
              style: GoogleFonts.quicksand(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: isDarkMode ? Colors.white : Colors.black87,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 10),
            Text(
              'Tell us your date of birth and gender',
              style: GoogleFonts.quicksand(
                fontSize: 16,
                color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 40),
            InkWell(
              onTap: () {
                DatePicker.showDatePicker(
                  context,
                  showTitleActions: true,
                  minTime: DateTime(1950, 1, 1),
                  maxTime: DateTime.now(),
                  onConfirm: (date) {
                    setState(() {
                      DateTime currentDate = DateTime.now();
                      String selectedDate =
                          DateFormat('yyyy-MM-dd').format(date);
                      Duration difference = currentDate.difference(date);
                      if (difference.inDays < 30 * 12 * 18) {
                        ToastService.showToast(
                            context: context,
                            message: "You must be 18",
                            title: "Warning",
                            type: ToastType.warning);
                        Navigator.pop(context);
                      } else {
                        dobController.text = selectedDate;
                      }
                    });
                  },
                  currentTime: DateTime.now(),
                  locale: LocaleType.vi,
                );
              },
              child: Container(
                padding: const EdgeInsets.all(15),
                decoration: BoxDecoration(
                  border: Border.all(
                      color:
                          isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                  borderRadius: BorderRadius.circular(10),
                  color: isDarkMode ? Colors.grey[800] : Colors.white,
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      dobController.text.isEmpty
                          ? 'Select Date of Birth'
                          : dobController.text,
                      style: GoogleFonts.quicksand(
                        fontSize: 16,
                        color: dobController.text.isEmpty
                            ? (isDarkMode ? Colors.grey[400] : Colors.grey[600])
                            : (isDarkMode ? Colors.white : Colors.black87),
                      ),
                    ),
                    Icon(Icons.calendar_today,
                        color: isDarkMode ? Colors.blue[300] : Colors.blue),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 20),
            DropdownButtonFormField<String>(
              value:
                  genderController.text.isEmpty ? null : genderController.text,
              style: GoogleFonts.quicksand(
                  fontSize: 16,
                  color: isDarkMode ? Colors.white : Colors.black87),
              decoration: InputDecoration(
                labelText: 'Gender',
                labelStyle: GoogleFonts.quicksand(
                    color: isDarkMode ? Colors.grey[400] : Colors.grey[600]),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color:
                          isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color: isDarkMode ? Colors.blue[300]! : Colors.blue,
                      width: 2),
                ),
                fillColor: isDarkMode ? Colors.grey[800] : Colors.white,
                filled: true,
              ),
              items: const [
                DropdownMenuItem(value: "Male", child: Text("Male")),
                DropdownMenuItem(value: "Female", child: Text("Female")),
                DropdownMenuItem(
                    value: "None", child: Text("Prefer not to say")),
              ],
              onChanged: (value) =>
                  setState(() => genderController.text = value!),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please select your gender';
                }
                return null;
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAddressStep() {
    final isDarkMode = themeProvider.isDarkMode;

    return Form(
      key: _formKeys[3],
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            Text(
              'Where Are You Located?',
              style: GoogleFonts.quicksand(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: isDarkMode ? Colors.white : Colors.black87,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 10),
            Text(
              'This helps us provide better services',
              style: GoogleFonts.quicksand(
                fontSize: 16,
                color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 40),
            TypeAheadField<String>(
              textFieldConfiguration: TextFieldConfiguration(
                controller: addressController,
                style: textStyle,
                decoration: InputDecoration(
                  labelText: 'Address',
                  labelStyle: GoogleFonts.quicksand(
                      color: isDarkMode ? Colors.grey[400] : Colors.grey[600]),
                  hintText: 'Enter your full address',
                  hintStyle: GoogleFonts.quicksand(
                      color: isDarkMode ? Colors.grey[500] : Colors.grey[500]),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                    borderSide: BorderSide(
                      color: isDarkMode ? Colors.grey[700]! : Colors.grey[300]!,
                    ),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                    borderSide: BorderSide(
                      color: isDarkMode ? Colors.blue[300]! : Colors.blue,
                      width: 2,
                    ),
                  ),
                  fillColor: isDarkMode ? Colors.grey[800] : Colors.white,
                  filled: true,
                ),
              ),
              suggestionsCallback: (pattern) async {
                final suggestions =
                    await AddressAutoComplete().getAddressSuggestions(pattern);
                return suggestions;
              },
              itemBuilder: (context, suggestion) {
                return ListTile(
                  title: Text(
                    suggestion,
                    style: GoogleFonts.quicksand(
                      fontSize: 14,
                      color: isDarkMode ? Colors.white : Colors.black,
                    ),
                  ),
                );
              },
              onSuggestionSelected: (suggestion) {
                addressController.text = suggestion;
              },
              suggestionsBoxDecoration: SuggestionsBoxDecoration(
                borderRadius: BorderRadius.circular(10),
                color: isDarkMode ? Colors.grey[800] : Colors.white,
                elevation: 4,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmailStep() {
    final isDarkMode = themeProvider.isDarkMode;
    return Form(
      key: _formKeys[4],
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Email Address',
              style: GoogleFonts.quicksand(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: isDarkMode ? Colors.white : Colors.black87,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 10),
            Text(
              'We\'ll use this to send you important updates',
              style: GoogleFonts.quicksand(
                fontSize: 16,
                color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 40),
            TextFormField(
              controller: emailController,
              keyboardType: TextInputType.emailAddress,
              style: GoogleFonts.quicksand(
                  fontSize: 16,
                  color: isDarkMode ? Colors.white : Colors.black87),
              decoration: InputDecoration(
                labelText: 'Email',
                labelStyle: GoogleFonts.quicksand(
                    color: isDarkMode ? Colors.grey[400] : Colors.grey[600]),
                hintText: 'your.email@example.com',
                hintStyle: GoogleFonts.quicksand(
                    color: isDarkMode ? Colors.grey[500] : Colors.grey[500]),
                prefixIcon: Icon(Icons.email_outlined,
                    color: isDarkMode ? Colors.blue[300] : Colors.blue),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color:
                          isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color: isDarkMode ? Colors.blue[300]! : Colors.blue,
                      width: 2),
                ),
                fillColor: isDarkMode ? Colors.grey[800] : Colors.white,
                filled: true,
              ),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter your email';
                }
                RegExp emailRegex =
                    RegExp(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$");
                if (!emailRegex.hasMatch(value)) {
                  return 'Please enter a valid email address';
                }
                return null;
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRoleStep() {
    final isDarkMode = themeProvider.isDarkMode;
    return Form(
      key: _formKeys[5],
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Choose Your Role',
              style: GoogleFonts.quicksand(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: isDarkMode ? Colors.white : Colors.black87,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 10),
            Text(
              'How would you like to use PsyConnect?',
              style: GoogleFonts.quicksand(
                fontSize: 16,
                color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 40),
            Container(
              decoration: BoxDecoration(
                border: Border.all(
                    color: isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                borderRadius: BorderRadius.circular(10),
                color: isDarkMode ? Colors.grey[800] : Colors.white,
              ),
              child: Column(
                children: [
                  RadioListTile<String>(
                    title: Text(
                      'Client',
                      style: GoogleFonts.quicksand(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                        color: isDarkMode ? Colors.white : Colors.black87,
                      ),
                    ),
                    subtitle: Text(
                      'Looking for mental health support',
                      style: GoogleFonts.quicksand(
                        fontSize: 14,
                        color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
                      ),
                    ),
                    value: 'Client',
                    groupValue: roleController.text.isEmpty
                        ? null
                        : roleController.text,
                    onChanged: (value) =>
                        setState(() => roleController.text = value!),
                    activeColor: isDarkMode ? Colors.blue[300] : Colors.blue,
                  ),
                  Divider(
                      height: 1,
                      color: isDarkMode ? Colors.grey[700] : Colors.grey[300]),
                  RadioListTile<String>(
                    title: Text(
                      'Therapist',
                      style: GoogleFonts.quicksand(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                        color: isDarkMode ? Colors.white : Colors.black87,
                      ),
                    ),
                    subtitle: Text(
                      'Providing mental health services',
                      style: GoogleFonts.quicksand(
                        fontSize: 14,
                        color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
                      ),
                    ),
                    value: 'Therapist',
                    groupValue: roleController.text.isEmpty
                        ? null
                        : roleController.text,
                    onChanged: (value) =>
                        setState(() => roleController.text = value!),
                    activeColor: isDarkMode ? Colors.blue[300] : Colors.blue,
                  ),
                ],
              ),
            ),
            if (roleController.text.isEmpty)
              Padding(
                padding: const EdgeInsets.only(top: 10),
                child: Text(
                  'Please select a role',
                  style: GoogleFonts.quicksand(
                    fontSize: 12,
                    color: Colors.red,
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildCredentialsStep() {
    final isDarkMode = themeProvider.isDarkMode;
    return Form(
      key: _formKeys[6],
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Create Your Account',
              style: GoogleFonts.quicksand(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: isDarkMode ? Colors.white : Colors.black87,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 10),
            Text(
              'Choose a username and secure password',
              style: GoogleFonts.quicksand(
                fontSize: 16,
                color: isDarkMode ? Colors.grey[400] : Colors.grey[600],
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 40),
            TextFormField(
              controller: usernameController,
              style: GoogleFonts.quicksand(
                  fontSize: 16,
                  color: isDarkMode ? Colors.white : Colors.black87),
              decoration: InputDecoration(
                labelText: 'Username',
                labelStyle: GoogleFonts.quicksand(
                    color: isDarkMode ? Colors.grey[400] : Colors.grey[600]),
                prefixIcon: Icon(Icons.person_outline,
                    color: isDarkMode ? Colors.blue[300] : Colors.blue),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color:
                          isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color: isDarkMode ? Colors.blue[300]! : Colors.blue,
                      width: 2),
                ),
                fillColor: isDarkMode ? Colors.grey[800] : Colors.white,
                filled: true,
              ),
              validator: (value) {
                if (value == null || value.length <= 3) {
                  return 'Username must be greater than 3 characters';
                }
                return null;
              },
            ),
            const SizedBox(height: 20),
            TextFormField(
              controller: passwordController,
              obscureText: true,
              style: GoogleFonts.quicksand(
                  fontSize: 16,
                  color: isDarkMode ? Colors.white : Colors.black87),
              decoration: InputDecoration(
                labelText: 'Password',
                labelStyle: GoogleFonts.quicksand(
                    color: isDarkMode ? Colors.grey[400] : Colors.grey[600]),
                prefixIcon: Icon(Icons.lock_outline,
                    color: isDarkMode ? Colors.blue[300] : Colors.blue),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color:
                          isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color: isDarkMode ? Colors.blue[300]! : Colors.blue,
                      width: 2),
                ),
                fillColor: isDarkMode ? Colors.grey[800] : Colors.white,
                filled: true,
              ),
              validator: (value) {
                if (value == null || value.length < 8) {
                  return 'Password must be at least 8 characters';
                }
                RegExp passReg =
                    RegExp(r"^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$");
                if (!passReg.hasMatch(value)) {
                  return 'Password must include characters, number and special character';
                }
                return null;
              },
            ),
            const SizedBox(height: 20),
            TextFormField(
              controller: retypePasswordController,
              obscureText: true,
              style: GoogleFonts.quicksand(
                  fontSize: 16,
                  color: isDarkMode ? Colors.white : Colors.black87),
              decoration: InputDecoration(
                labelText: 'Confirm Password',
                labelStyle: GoogleFonts.quicksand(
                    color: isDarkMode ? Colors.grey[400] : Colors.grey[600]),
                prefixIcon: Icon(Icons.lock_outline,
                    color: isDarkMode ? Colors.blue[300] : Colors.blue),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color:
                          isDarkMode ? Colors.grey[700]! : Colors.grey[300]!),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: BorderSide(
                      color: isDarkMode ? Colors.blue[300]! : Colors.blue,
                      width: 2),
                ),
                fillColor: isDarkMode ? Colors.grey[800] : Colors.white,
                filled: true,
              ),
              validator: (value) {
                if (value != passwordController.text || value == null) {
                  return 'Passwords do not match';
                }
                return null;
              },
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _handleFinalSubmit(
      dynamic userProvider, BuildContext context) async {
    if (!_formKeys[6].currentState!.validate()) return;
    if (roleController.text.isEmpty) {
      ToastService.showToast(
        message: "Please select your role!",
        context: context,
        title: 'Role Required!',
        type: ToastType.error,
      );
      return;
    }

    _setStartLoading();

    try {
      String? cloudinaryUrlImage;

      if (_image != null) {
        cloudinaryUrlImage = await cloudinaryApiService.uploadImage(
          imageFile: _image as File,
          username: usernameController.text,
        );

        if (cloudinaryUrlImage == null) {
          ToastService.showToast(
            message: "Image upload failed!",
            context: context,
            title: 'Image Error!',
            type: ToastType.error,
          );
          return;
        }
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
        "avatarUri": cloudinaryUrlImage ?? "",
        "dob": dobController.text,
      };

      await registerService.registerHandle(
        requestBody: jsonData,
        userProvider: userProvider,
        context: context,
      );
    } catch (e) {
      ToastService.showToast(
        message: "Registration failed. Please try again.",
        context: context,
        title: 'Error!',
        type: ToastType.error,
      );
    }
  }
}
