import 'package:PsyConnect/core/colors/color.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/swipe_card.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/services/profile_service/therapist_service.dart';
import 'package:PsyConnect/ui/screens/booking_screen.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';

class TherapistDetailScreen extends StatefulWidget {
  final String therapistId;

  const TherapistDetailScreen({
    super.key,
    required this.therapistId,
  });

  @override
  State<TherapistDetailScreen> createState() => _TherapistDetailScreenState();
}

class _TherapistDetailScreenState extends State<TherapistDetailScreen> {
  final TherapistService _service = TherapistService();
  
  SwipeCard? _therapist;
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadDetails();
  }

  Future<void> _loadDetails() async {
    try {
      final data = await _service.getTherapistById(widget.therapistId);
      setState(() {
        _therapist = data;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString().replaceAll("Exception: ", "");
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Provider.of<ThemeProvider>(context).isDarkMode;
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: isDark ? Colors.black : const Color(0xFFF7F9FC),
      appBar: AppBar(
        title: Text(
          "Specialist Details",
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
      body: _isLoading
          ? Center(
              child: CircularProgressIndicator(
                color: isDark ? secondaryColor : primaryColor,
              ),
            )
          : _error != null
              ? Center(
                  child: Padding(
                    padding: const EdgeInsets.all(kDefault * 2),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        const Icon(Icons.error_outline_rounded, color: Colors.redAccent, size: 48),
                        const SizedBox(height: 12),
                        Text(
                          "Failed to load specialist details: $_error",
                          textAlign: TextAlign.center,
                          style: GoogleFonts.quicksand(color: Colors.red[300], fontWeight: FontWeight.w500),
                        ),
                      ],
                    ),
                  ),
                )
              : _buildContent(isDark, theme),
    );
  }

  Widget _buildContent(bool isDark, ThemeData theme) {
    final card = _therapist!;
    final avatarUrl = card.avatarUrl.isNotEmpty ? card.avatarUrl : 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';
    final priceStr = card.price > 0 
        ? NumberFormat.currency(locale: 'vi_VN', symbol: 'đ').format(card.price)
        : "Free";

    return Column(
      children: [
        Expanded(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // 1. Big Avatar, Name and Title Card
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: isDark ? theme.cardColor : Colors.white,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: isDark ? Colors.grey[850]! : Colors.grey[200]!, width: 1),
                  ),
                  child: Row(
                    children: [
                      CircleAvatar(
                        radius: 36,
                        backgroundImage: NetworkImage(avatarUrl),
                        backgroundColor: isDark ? Colors.grey[800] : Colors.grey[300],
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              card.name,
                              style: GoogleFonts.quicksand(fontSize: 18, fontWeight: FontWeight.bold, color: isDark ? Colors.white : Colors.black87),
                            ),
                            const SizedBox(height: 2),
                            Text(
                              card.title,
                              style: GoogleFonts.quicksand(fontSize: 13, fontWeight: FontWeight.w500, color: isDark ? Colors.blue[300] : primaryColor),
                            ),
                            const SizedBox(height: 6),
                            Row(
                              children: [
                                const Icon(Icons.star_rounded, color: Colors.amber, size: 16),
                                const SizedBox(width: 4),
                                Text(
                                  "${card.rating} / 5.0 Rating",
                                  style: GoogleFonts.quicksand(fontSize: 12, fontWeight: FontWeight.bold),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),

                const SizedBox(height: 16),

                // 2. Profile Details Section
                _buildSectionHeader("Consultation Details", isDark),
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: isDark ? theme.cardColor : Colors.white,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: isDark ? Colors.grey[850]! : Colors.grey[200]!, width: 1),
                  ),
                  child: Column(
                    children: [
                      _buildDetailRow(Icons.work_history_outlined, "Experience", "${card.experience} Years of Practice", isDark),
                      const Divider(height: 20),
                      _buildDetailRow(Icons.payments_outlined, "Consultation Fee", "$priceStr per session", isDark),
                      const Divider(height: 20),
                      _buildDetailRow(Icons.location_on_outlined, "Office Address", card.address.isNotEmpty ? card.address : "Online Gateway Support", isDark),
                      const Divider(height: 20),
                      _buildDetailRow(Icons.language_rounded, "Languages Spoken", card.languages.isNotEmpty ? card.languages.join(", ") : "English, Vietnamese", isDark),
                    ],
                  ),
                ),

                const SizedBox(height: 16),

                // 3. Specializations Badges
                if (card.specialization.isNotEmpty) ...[
                  _buildSectionHeader("Clinical Focus & Specializations", isDark),
                  Wrap(
                    spacing: 8,
                    runSpacing: 8,
                    children: card.specialization.map((s) => _buildBadge(s, Colors.teal, isDark)).toList(),
                  ),
                  const SizedBox(height: 16),
                ],

                // 4. Consultation Modes Badges
                if (card.consultationModes.isNotEmpty) ...[
                  _buildSectionHeader("Consultation Support Modes", isDark),
                  Wrap(
                    spacing: 8,
                    runSpacing: 8,
                    children: card.consultationModes.map((m) => _buildBadge(m.toUpperCase(), Colors.indigo, isDark)).toList(),
                  ),
                  const SizedBox(height: 16),
                ],

                // 5. Reasons for Matching points list
                if (card.reasons.isNotEmpty) ...[
                  _buildSectionHeader("Why you matched with this specialist", isDark),
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: isDark ? theme.cardColor : Colors.white,
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(color: isDark ? Colors.grey[850]! : Colors.grey[200]!, width: 1),
                    ),
                    child: Column(
                      children: card.reasons.map((reason) {
                        return Padding(
                          padding: const EdgeInsets.symmetric(vertical: 4),
                          child: Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const Padding(
                                padding: EdgeInsets.only(top: 2),
                                child: Icon(Icons.verified_outlined, color: Colors.teal, size: 14),
                              ),
                              const SizedBox(width: 8),
                              Expanded(
                                child: Text(
                                  reason,
                                  style: GoogleFonts.quicksand(fontSize: 12, height: 1.4, color: isDark ? Colors.white70 : Colors.black87),
                                ),
                              ),
                            ],
                          ),
                        );
                      }).toList(),
                    ),
                  ),
                  const SizedBox(height: 16),
                ],
              ],
            ),
          ),
        ),

        // 6. Book Consultation Button Block
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: isDark ? theme.cardColor : Colors.white,
            border: Border(top: BorderSide(color: isDark ? Colors.grey[850]! : Colors.grey[200]!, width: 1)),
          ),
          child: ElevatedButton(
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => BookingScreen(
                    therapistId: card.profileId,
                    therapistName: card.name,
                    price: card.price,
                  ),
                ),
              );
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: isDark ? Colors.blueGrey[800] : primaryColor,
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(vertical: 16),
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
              elevation: 0,
            ),
            child: Text(
              "Book a Consultation Session",
              style: GoogleFonts.quicksand(fontWeight: FontWeight.bold, fontSize: 15),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildSectionHeader(String title, bool isDark) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 10, top: 12),
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

  Widget _buildDetailRow(IconData icon, String label, String value, bool isDark) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Icon(icon, color: Colors.grey, size: 20),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                label,
                style: GoogleFonts.quicksand(fontSize: 11, color: Colors.grey[500], fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 2),
              Text(
                value,
                style: GoogleFonts.quicksand(fontSize: 13, color: isDark ? Colors.white70 : Colors.black87, fontWeight: FontWeight.w500),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildBadge(String text, Color color, bool isDark) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: color.withOpacity(0.08),
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: color.withOpacity(0.2), width: 0.5),
      ),
      child: Text(
        text,
        style: GoogleFonts.quicksand(
          fontSize: 11,
          fontWeight: FontWeight.bold,
          color: isDark ? color.withOpacity(0.8) : color,
        ),
      ),
    );
  }
}
