import 'package:PsyConnect/core/colors/color.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/provider/session_provider.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';

class BookingScreen extends StatefulWidget {
  final String therapistId;
  final String therapistName;
  final double price;

  const BookingScreen({
    super.key,
    required this.therapistId,
    required this.therapistName,
    required this.price,
  });

  @override
  State<BookingScreen> createState() => _BookingScreenState();
}

class _BookingScreenState extends State<BookingScreen> {
  DateTime _selectedDate = DateTime.now().add(const Duration(days: 1)); // Default tomorrow
  String _selectedSlot = 'morning'; // 'morning', 'afternoon', 'evening'
  String _selectedMode = 'online'; // 'online', 'in_person'

  bool _isSubmitting = false;

  final List<Map<String, String>> _slots = [
    {'value': 'morning', 'label': 'Morning', 'hours': '09:00 - 11:00'},
    {'value': 'afternoon', 'label': 'Afternoon', 'hours': '14:00 - 16:00'},
    {'value': 'evening', 'label': 'Evening', 'hours': '19:00 - 21:00'},
  ];

  final List<Map<String, String>> _modes = [
    {'value': 'online', 'label': 'Online Video Call', 'desc': 'Telehealth virtual consulting room'},
    {'value': 'in_person', 'label': 'In Person Meeting', 'desc': 'Face to face office visit'},
  ];

  Future<void> _selectDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime.now().add(const Duration(days: 1)),
      lastDate: DateTime.now().add(const Duration(days: 30)),
      builder: (context, child) {
        return Theme(
          data: Theme.of(context).copyWith(
            colorScheme: ColorScheme.light(
              primary: primaryColor,
              onPrimary: Colors.white,
              onSurface: Colors.black87,
            ),
          ),
          child: child!,
        );
      },
    );
    if (picked != null && picked != _selectedDate) {
      setState(() {
        _selectedDate = picked;
      });
    }
  }

  Future<void> _confirmBooking(SessionProvider provider) async {
    setState(() {
      _isSubmitting = true;
    });

    // Format start and end times mapping back to slot definitions
    String startHour = "09:00:00";
    String endHour = "11:00:00";
    if (_selectedSlot == 'afternoon') {
      startHour = "14:00:00";
      endHour = "16:00:00";
    } else if (_selectedSlot == 'evening') {
      startHour = "19:00:00";
      endHour = "21:00:00";
    }

    final dateStr = DateFormat('yyyy-MM-dd').format(_selectedDate);
    final startTimeISO = "${dateStr}T$startHour.000Z";
    final endTimeISO = "${dateStr}T$endHour.000Z";

    final success = await provider.bookSession(
      therapistId: widget.therapistId,
      startTime: startTimeISO,
      endTime: endTimeISO,
      scheduledDate: dateStr,
      mode: _selectedMode,
      price: widget.price,
    );

    if (success && mounted) {
      ToastService.showToast(
        context: context,
        message: "Your consulting session was booked successfully. Please complete the invoice payment.",
        title: "Session Booked",
        type: ToastType.success,
      );
      
      // Pop back twice to get back to Schedule listings
      Navigator.pop(context); // Pop BookingScreen
      Navigator.pop(context); // Pop TherapistDetailScreen
    } else if (mounted) {
      ToastService.showToast(
        context: context,
        message: provider.errorMessage ?? "Failed to book session.",
        title: "Booking Error",
        type: ToastType.error,
      );
    }

    setState(() {
      _isSubmitting = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final sessionProvider = Provider.of<SessionProvider>(context);
    final isDark = Provider.of<ThemeProvider>(context).isDarkMode;
    final theme = Theme.of(context);

    final formattedPrice = widget.price > 0
        ? NumberFormat.currency(locale: 'vi_VN', symbol: 'đ').format(widget.price)
        : "Free";

    return Scaffold(
      backgroundColor: isDark ? Colors.black : const Color(0xFFF7F9FC),
      appBar: AppBar(
        title: Text(
          "Reserve Appointment",
          style: GoogleFonts.quicksand(
            fontWeight: FontWeight.bold,
            fontSize: 18,
            color: theme.textTheme.titleLarge?.color,
          ),
        ),
        elevation: 0.5,
        backgroundColor: isDark ? theme.cardColor : Colors.white,
        leading: IconButton(
          icon: Icon(Icons.arrow_back_ios_new_rounded, size: 20, color: theme.iconTheme.color),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  // Summary Banner
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: isDark ? theme.cardColor : Colors.white,
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(color: isDark ? Colors.grey[850]! : Colors.grey[200]!, width: 1),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          widget.therapistName,
                          style: GoogleFonts.quicksand(fontSize: 16, fontWeight: FontWeight.bold, color: isDark ? Colors.white : Colors.black87),
                        ),
                        const SizedBox(height: 6),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              "Session Price:",
                              style: GoogleFonts.quicksand(fontSize: 12, color: Colors.grey),
                            ),
                            Text(
                              formattedPrice,
                              style: GoogleFonts.quicksand(fontSize: 14, fontWeight: FontWeight.bold, color: Colors.green),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 24),

                  // 1. Date Selector Box
                  _buildSectionHeader("1. Choose Date", isDark),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                    decoration: BoxDecoration(
                      color: isDark ? theme.cardColor : Colors.white,
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(color: isDark ? Colors.grey[850]! : Colors.grey[200]!, width: 1),
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          children: [
                            Icon(Icons.calendar_today_rounded, color: isDark ? Colors.blue[300] : primaryColor),
                            const SizedBox(width: 12),
                            Text(
                              DateFormat('EEEE, dd/MM/yyyy').format(_selectedDate),
                              style: GoogleFonts.quicksand(fontSize: 14, fontWeight: FontWeight.bold, color: isDark ? Colors.white : Colors.black87),
                            ),
                          ],
                        ),
                        TextButton(
                          onPressed: () => _selectDate(context),
                          child: Text(
                            "Change",
                            style: GoogleFonts.quicksand(fontWeight: FontWeight.bold, color: isDark ? Colors.blue[300] : primaryColor),
                          ),
                        ),
                      ],
                    ),
                  ),

                  const SizedBox(height: 24),

                  // 2. Time Slot Selection
                  _buildSectionHeader("2. Choose Time Slot", isDark),
                  Column(
                    children: _slots.map((slot) {
                      final isSelected = _selectedSlot == slot['value'];
                      return GestureDetector(
                        onTap: () => setState(() => _selectedSlot = slot['value']!),
                        child: Container(
                          margin: const EdgeInsets.only(bottom: 10),
                          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                          decoration: BoxDecoration(
                            color: isSelected
                                ? (isDark ? Colors.blueGrey[900] : primaryColor.withOpacity(0.08))
                                : (isDark ? theme.cardColor : Colors.white),
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(
                              color: isSelected
                                  ? (isDark ? Colors.blue[300]! : primaryColor)
                                  : (isDark ? Colors.grey[850]! : Colors.grey[200]!),
                              width: 1.5,
                            ),
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    slot['label']!,
                                    style: GoogleFonts.quicksand(fontSize: 14, fontWeight: FontWeight.bold, color: isDark ? Colors.white : Colors.black87),
                                  ),
                                  const SizedBox(height: 2),
                                  Text(
                                    slot['hours']!,
                                    style: GoogleFonts.quicksand(fontSize: 12, color: Colors.grey),
                                  ),
                                ],
                              ),
                              if (isSelected)
                                Icon(Icons.check_circle, color: isDark ? Colors.blue[300] : primaryColor)
                              else
                                const Icon(Icons.circle_outlined, color: Colors.grey, size: 20),
                            ],
                          ),
                        ),
                      );
                    }).toList(),
                  ),

                  const SizedBox(height: 24),

                  // 3. Consultation Mode Selection
                  _buildSectionHeader("3. Choose Consultation Mode", isDark),
                  Column(
                    children: _modes.map((mode) {
                      final isSelected = _selectedMode == mode['value'];
                      return GestureDetector(
                        onTap: () => setState(() => _selectedMode = mode['value']!),
                        child: Container(
                          margin: const EdgeInsets.only(bottom: 10),
                          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                          decoration: BoxDecoration(
                            color: isSelected
                                ? (isDark ? Colors.blueGrey[900] : primaryColor.withOpacity(0.08))
                                : (isDark ? theme.cardColor : Colors.white),
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(
                              color: isSelected
                                  ? (isDark ? Colors.blue[300]! : primaryColor)
                                  : (isDark ? Colors.grey[850]! : Colors.grey[200]!),
                              width: 1.5,
                            ),
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      mode['label']!,
                                      style: GoogleFonts.quicksand(fontSize: 14, fontWeight: FontWeight.bold, color: isDark ? Colors.white : Colors.black87),
                                    ),
                                    const SizedBox(height: 2),
                                    Text(
                                      mode['desc']!,
                                      style: GoogleFonts.quicksand(fontSize: 11, color: Colors.grey),
                                    ),
                                  ],
                                ),
                              ),
                              if (isSelected)
                                Icon(Icons.check_circle, color: isDark ? Colors.blue[300] : primaryColor)
                              else
                                const Icon(Icons.circle_outlined, color: Colors.grey, size: 20),
                            ],
                          ),
                        ),
                      );
                    }).toList(),
                  ),
                ],
              ),
            ),
          ),

          // 4. Reserve Submission Button
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: isDark ? theme.cardColor : Colors.white,
              border: Border(top: BorderSide(color: isDark ? Colors.grey[850]! : Colors.grey[200]!, width: 1)),
            ),
            child: _isSubmitting || sessionProvider.isLoading
                ? const Center(child: CircularProgressIndicator())
                : ElevatedButton(
                    onPressed: () => _confirmBooking(sessionProvider),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: isDark ? Colors.blueGrey[800] : primaryColor,
                      foregroundColor: Colors.white,
                      padding: const EdgeInsets.symmetric(vertical: 16),
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                      elevation: 0,
                    ),
                    child: Text(
                      "Submit Booking Request",
                      style: GoogleFonts.quicksand(fontWeight: FontWeight.bold, fontSize: 15),
                    ),
                  ),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(String title, bool isDark) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Text(
        title,
        style: GoogleFonts.quicksand(
          fontSize: 14,
          fontWeight: FontWeight.bold,
          color: isDark ? Colors.white70 : Colors.black87,
        ),
      ),
    );
  }
}
