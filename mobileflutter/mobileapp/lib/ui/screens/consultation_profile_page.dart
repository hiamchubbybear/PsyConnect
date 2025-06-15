import 'package:flutter/material.dart';
import 'package:multi_select_flutter/multi_select_flutter.dart';

class ConsultationProfilePage extends StatefulWidget {
  const ConsultationProfilePage({super.key});

  @override
  _ConsultationProfilePageState createState() => _ConsultationProfilePageState();
}

class _ConsultationProfilePageState extends State<ConsultationProfilePage> {
  final _formKey = GlobalKey<FormState>();
  final TextEditingController addressController = TextEditingController();
  final TextEditingController experienceController = TextEditingController();
  final TextEditingController ratingController = TextEditingController();
  final TextEditingController priceController = TextEditingController();

  List<String> selectedLanguages = [];
  List<String> selectedSpecializations = [];
  List<String> selectedConsultationModes = [];
  List<String> selectedDays = [];
  List<String> selectedTimeSlots = [];

  final List<String> allLanguages = ['English', 'Spanish', 'French'];
  final List<String> allSpecializations = ['CBT', 'Anxiety', 'Depression'];
  final List<String> allModes = ['Online', 'In-Person'];
  final List<String> allDays = [
    'Monday',
    'Tuesday',
    'Wednesday',
    'Thursday',
    'Friday',
    'Saturday',
    'Sunday'
  ];
  final List<String> allTimeSlots = [
    '9:00 AM - 11:00 AM',
    '10:00 AM - 12:00 PM',
    '2:00 PM - 4:00 PM',
    '4:00 PM - 6:00 PM',
  ];

  void _saveProfile() {
    if (_formKey.currentState!.validate()) {
      final profile = {
        "address": addressController.text,
        "languages": selectedLanguages,
        "specialization": selectedSpecializations,
        "consultation_modes": selectedConsultationModes,
        "experience": int.parse(experienceController.text),
        "rating": double.parse(ratingController.text),
        "price_per_session": double.parse(priceController.text),
        "availability": {
          "days": selectedDays,
          "time_slots": selectedTimeSlots,
        },
      };

      print("Saved Profile: $profile");

      showDialog(
        context: context,
        builder: (_) => AlertDialog(
          title: const Text("Saved"),
          content: const Text("Consultation Profile saved."),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text("OK"),
            ),
          ],
        ),
      );
    }
  }

  Widget _buildSectionTitle(String title) {
    return Padding(
      padding: const EdgeInsets.only(top: 24, bottom: 12),
      child: Text(
        title,
        style: const TextStyle(
          fontSize: 18,
          fontWeight: FontWeight.w500,
          color: Color(0xFF333333),
        ),
      ),
    );
  }

  Widget _buildTextField({
    required TextEditingController controller,
    required String label,
    required String? Function(String?) validator,
    TextInputType keyboardType = TextInputType.text,
  }) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: TextFormField(
        controller: controller,
        keyboardType: keyboardType,
        decoration: InputDecoration(
          labelText: label,
          labelStyle: const TextStyle(color: Color(0xFF666666)),
          border: const UnderlineInputBorder(
            borderSide: BorderSide(color: Color(0xFFCCCCCC)),
          ),
          enabledBorder: const UnderlineInputBorder(
            borderSide: BorderSide(color: Color(0xFFCCCCCC)),
          ),
          focusedBorder: const UnderlineInputBorder(
            borderSide: BorderSide(color: Color(0xFF00695C), width: 1.5),
          ),
        ),
        style: const TextStyle(color: Color(0xFF333333)),
        validator: validator,
      ),
    );
  }

  Widget _buildMultiSelectField({
    required List<String> items,
    required String title,
    required String buttonText,
    required Function(List<dynamic>) onConfirm,
  }) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: MultiSelectDialogField(
        items: items.map((e) => MultiSelectItem(e, e)).toList(),
        title: Text(title, style: const TextStyle(color: Color(0xFF333333))),
        buttonText: Text(buttonText, style: const TextStyle(color: Color(0xFF666666))),
        chipDisplay:  MultiSelectChipDisplay(
          textStyle: TextStyle(color: Color(0xFF333333)),
        ),
        decoration: const BoxDecoration(
          border: Border(bottom: BorderSide(color: Color(0xFFCCCCCC))),
        ),
        onConfirm: onConfirm,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          "Consultation Profile",
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w500,
            color: Color(0xFF333333),
          ),
        ),
        centerTitle: true,
        backgroundColor: const Color(0xFFF7F7F7),
        elevation: 0,
      ),
      backgroundColor: const Color(0xFFF7F7F7),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Basic Info Section
              _buildSectionTitle("Basic Information"),
              _buildTextField(
                controller: addressController,
                label: "Address",
                validator: (value) => value!.isEmpty ? "Required" : null,
              ),
              _buildTextField(
                controller: experienceController,
                label: "Years of Experience",
                keyboardType: TextInputType.number,
                validator: (value) => value!.isEmpty ? "Required" : null,
              ),
              _buildTextField(
                controller: ratingController,
                label: "Rating (e.g. 4.5)",
                keyboardType: TextInputType.number,
                validator: (value) => value!.isEmpty ? "Required" : null,
              ),
              _buildTextField(
                controller: priceController,
                label: "Price per Session",
                keyboardType: TextInputType.number,
                validator: (value) => value!.isEmpty ? "Required" : null,
              ),
              // Consultation Preferences Section
              _buildSectionTitle("Consultation Preferences"),
              _buildMultiSelectField(
                items: allLanguages,
                title: "Languages",
                buttonText: "Select Languages",
                onConfirm: (values) =>
                    setState(() => selectedLanguages = List<String>.from(values)),
              ),
              _buildMultiSelectField(
                items: allSpecializations,
                title: "Specializations",
                buttonText: "Select Specializations",
                onConfirm: (values) => setState(() =>
                    selectedSpecializations = List<String>.from(values)),
              ),
              _buildMultiSelectField(
                items: allModes,
                title: "Consultation Modes",
                buttonText: "Select Modes",
                onConfirm: (values) => setState(() =>
                    selectedConsultationModes = List<String>.from(values)),
              ),
              // Availability Section
              _buildSectionTitle("Availability"),
              _buildMultiSelectField(
                items: allDays,
                title: "Available Days",
                buttonText: "Select Days",
                onConfirm: (values) =>
                    setState(() => selectedDays = List<String>.from(values)),
              ),
              _buildMultiSelectField(
                items: allTimeSlots,
                title: "Time Slots",
                buttonText: "Select Time Slots",
                onConfirm: (values) =>
                    setState(() => selectedTimeSlots = List<String>.from(values)),
              ),
              const SizedBox(height: 32),
              // Save Button
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color.fromARGB(255, 115, 170, 163),
                    foregroundColor: const Color(0xFFF7F7F7),
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    textStyle: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w500,
                    ),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                    elevation: 0,
                  ),
                  onPressed: _saveProfile,
                  child: const Text("Save Profile"),
                ),
              ),
              const SizedBox(height: 16),
            ],
          ),
        ),
      ),
    );
  }
}
