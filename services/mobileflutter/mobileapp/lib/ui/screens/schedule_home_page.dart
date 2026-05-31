import 'package:PsyConnect/core/colors/color.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/consultation_session.dart';
import 'package:PsyConnect/provider/session_provider.dart';
import 'package:PsyConnect/provider/chat_provider.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/ui/screens/payment_webview.dart';
import 'package:PsyConnect/ui/screens/call_screen.dart';
import 'package:PsyConnect/ui/screens/smart_match_screen.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:intl/intl.dart';

class ScheduleHomePage extends StatefulWidget {
  const ScheduleHomePage({
    super.key,
    required this.size,
  });

  final Size size;

  @override
  State<ScheduleHomePage> createState() => _ScheduleHomePageState();
}

class _ScheduleHomePageState extends State<ScheduleHomePage> {
  String _activeTab = 'upcoming'; // 'upcoming', 'completed', 'cancelled'

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<SessionProvider>(context, listen: false).fetchSessions();
      Provider.of<ChatProvider>(context, listen: false).initMyProfile();
    });
  }

  void _changeTab(String tabName) {
    setState(() {
      _activeTab = tabName;
    });
  }

  List<ConsultationSession> _getFilteredSessions(List<ConsultationSession> allSessions) {
    final now = DateTime.now();
    switch (_activeTab) {
      case 'upcoming':
        return allSessions.where((s) {
          if (s.status == 'cancelled' || s.status == 'completed') return false;
          // Return pending, active or future sessions
          if (s.status == 'pending_payment' || s.status == 'active' || s.status == 'pending') return true;
          try {
            final start = DateTime.parse(s.startTime);
            return start.isAfter(now);
          } catch (_) {
            return true;
          }
        }).toList();
      case 'completed':
        return allSessions.where((s) => s.status == 'completed').toList();
      case 'cancelled':
        return allSessions.where((s) => s.status == 'cancelled').toList();
      default:
        return allSessions;
    }
  }

  String _formatSessionDate(String dateStr) {
    if (dateStr.isEmpty) return '—';
    try {
      final dt = DateTime.parse(dateStr).toLocal();
      return DateFormat('dd/MM/yyyy').format(dt);
    } catch (_) {
      return dateStr;
    }
  }

  String _formatSessionTime(String startStr, String endStr) {
    try {
      final start = DateTime.parse(startStr).toLocal();
      final end = DateTime.parse(endStr).toLocal();
      final startFormatted = DateFormat('HH:mm').format(start);
      final endFormatted = DateFormat('HH:mm').format(end);
      return "$startFormatted - $endFormatted";
    } catch (_) {
      return "TBD";
    }
  }

  String _getWeekDay(String dateStr) {
    if (dateStr.isEmpty) return '';
    try {
      final dt = DateTime.parse(dateStr).toLocal();
      final weekdays = ['Chủ Nhật', 'Thứ Hai', 'Thứ Ba', 'Thứ Tư', 'Thứ Năm', 'Thứ Sáu', 'Thứ Bảy'];
      return weekdays[dt.weekday % 7];
    } catch (_) {
      return '';
    }
  }

  String _getDayNum(String dateStr) {
    if (dateStr.isEmpty) return '—';
    try {
      final dt = DateTime.parse(dateStr).toLocal();
      return dt.day.toString().padLeft(2, '0');
    } catch (_) {
      return '';
    }
  }

  @override
  Widget build(BuildContext context) {
    final sessionProvider = Provider.of<SessionProvider>(context);
    final chatProvider = Provider.of<ChatProvider>(context);
    final isDark = Provider.of<ThemeProvider>(context).isDarkMode;
    final theme = Theme.of(context);

    final filteredSessions = _getFilteredSessions(sessionProvider.sessions);

    final myProfileId = chatProvider.myProfileId ?? '';

    return Scaffold(
      backgroundColor: isDark ? Colors.black : const Color(0xFFF7F9FC),
      body: Column(
        children: [
          // Dynamic Custom AppBar Header
          Container(
            padding: EdgeInsets.only(
              left: kDefault,
              right: kDefault,
              top: MediaQuery.of(context).padding.top * 1.2,
              bottom: kDefault / 1.5,
            ),
            decoration: BoxDecoration(
              color: isDark ? theme.cardColor : Colors.white,
              borderRadius: const BorderRadius.vertical(
                bottom: Radius.circular(kDefault * 1.5),
              ),
              boxShadow: [
                BoxShadow(
                  offset: const Offset(0, 1),
                  color: isDark ? Colors.black38 : Colors.grey.withOpacity(.12),
                  blurRadius: 10.0,
                )
              ],
            ),
            child: Column(
              children: [
                Padding(
                  padding: const EdgeInsets.only(bottom: kDefault / 2),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        "Consultation Schedules",
                        style: GoogleFonts.quicksand(
                          fontWeight: FontWeight.bold,
                          fontSize: 20,
                          color: isDark ? Colors.white : Colors.black87,
                        ),
                      ),
                      Row(
                        children: [
                          IconButton(
                            tooltip: "Smart Therapist Match",
                            icon: const Icon(
                              Icons.favorite_rounded,
                              color: Colors.redAccent,
                            ),
                            onPressed: () {
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                  builder: (context) => const SmartMatchScreen(),
                                ),
                              );
                            },
                          ),
                          IconButton(
                            icon: Icon(
                              Icons.refresh_rounded,
                              color: isDark ? secondaryColor : primaryColor,
                            ),
                            onPressed: () => sessionProvider.fetchSessions(),
                          ),
                        ],
                      )
                    ],
                  ),
                ),

                // Inline Navigation tabs matching web app layout
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    _buildTabButton("upcoming", "Upcoming", isDark),
                    _buildTabButton("completed", "Completed", isDark),
                    _buildTabButton("cancelled", "Cancelled", isDark),
                  ],
                ),
              ],
            ),
          ),

          Expanded(
            child: sessionProvider.isLoading
                ? Center(
                    child: CircularProgressIndicator(
                      color: isDark ? secondaryColor : primaryColor,
                    ),
                  )
                : sessionProvider.errorMessage != null
                    ? Center(
                        child: Padding(
                          padding: const EdgeInsets.all(kDefault * 2),
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              const Icon(Icons.error_outline_rounded, color: Colors.redAccent, size: 48),
                              const SizedBox(height: 12),
                              Text(
                                "Failed to load schedule: ${sessionProvider.errorMessage}",
                                textAlign: TextAlign.center,
                                style: GoogleFonts.quicksand(
                                  color: Colors.red[300],
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                              const SizedBox(height: 16),
                              ElevatedButton(
                                onPressed: () => sessionProvider.fetchSessions(),
                                style: ElevatedButton.styleFrom(
                                  backgroundColor: isDark ? secondaryColor : primaryColor,
                                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
                                ),
                                child: Text("Retry", style: GoogleFonts.quicksand(fontWeight: FontWeight.bold)),
                              )
                            ],
                          ),
                        ),
                      )
                    : filteredSessions.isEmpty
                        ? _buildEmptyState(isDark)
                        : RefreshIndicator(
                            onRefresh: () => sessionProvider.fetchSessions(),
                            child: ListView.builder(
                              itemCount: filteredSessions.length,
                              padding: const EdgeInsets.symmetric(vertical: kDefault),
                              itemBuilder: (context, index) {
                                final session = filteredSessions[index];
                                return _buildSessionCard(session, myProfileId, chatProvider, sessionProvider, isDark);
                              },
                            ),
                          ),
          ),
        ],
      ),
    );
  }

  Widget _buildTabButton(String tabKey, String label, bool isDark) {
    final isActive = _activeTab == tabKey;
    final activeBg = isDark ? Colors.blueGrey[900] : primaryColor.withOpacity(0.1);
    final activeText = isDark ? Colors.blue[300] : primaryColor;

    return GestureDetector(
      onTap: () => _changeTab(tabKey),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        decoration: BoxDecoration(
          color: isActive ? activeBg : Colors.transparent,
          borderRadius: BorderRadius.circular(20),
        ),
        child: Text(
          label,
          style: GoogleFonts.quicksand(
            fontWeight: isActive ? FontWeight.bold : FontWeight.w500,
            fontSize: 14,
            color: isActive 
                ? activeText 
                : (isDark ? Colors.grey[400] : Colors.grey[600]),
          ),
        ),
      ),
    );
  }

  Widget _buildEmptyState(bool isDark) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(kDefault * 2),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.calendar_month_outlined,
              size: 72,
              color: isDark ? Colors.grey[800] : Colors.grey[300],
            ),
            const SizedBox(height: 16),
            Text(
              "No consultation sessions here.",
              style: GoogleFonts.quicksand(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: isDark ? Colors.grey[400] : Colors.grey[700],
              ),
            ),
            const SizedBox(height: 6),
            Text(
              "Any booked meetings or invoices will appear in this tab.",
              textAlign: TextAlign.center,
              style: GoogleFonts.quicksand(
                color: Colors.grey,
                fontSize: 13,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSessionCard(
    ConsultationSession session,
    String myProfileId,
    ChatProvider chatProvider,
    SessionProvider sessionProvider,
    bool isDark,
  ) {
    final theme = Theme.of(context);

    // Deduce participant details
    final isClient = myProfileId == session.clientId;
    final companionId = isClient ? session.therapistId : session.clientId;

    // Load companion profile silently in background
    final companionProfile = chatProvider.cachedProfiles[companionId];
    if (companionProfile == null) {
      chatProvider.fetchUserProfileSilently(companionId);
    }

    final companionName = companionProfile != null
        ? "${companionProfile.firstName ?? ''} ${companionProfile.lastName ?? ''}".trim()
        : (isClient ? "Consulting Specialist" : "Patient Client");

    final avatarUrl = companionProfile?.avatarUri ?? 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';

    // Pricing Format
    final currencyFormat = NumberFormat.currency(locale: 'vi_VN', symbol: 'đ');
    final formattedPrice = session.price > 0 ? currencyFormat.format(session.price) : "Free";

    // Status attributes
    Color statusColor;
    String statusLabel = session.status.toUpperCase();
    if (session.status == 'completed') {
      statusColor = Colors.green;
      statusLabel = "Completed";
    } else if (session.status == 'cancelled') {
      statusColor = Colors.redAccent;
      statusLabel = "Cancelled";
    } else if (session.paymentStatus == 'paid' || session.status == 'active') {
      statusColor = Colors.blue;
      statusLabel = "Paid / Active";
    } else {
      statusColor = Colors.orangeAccent;
      statusLabel = "Unpaid";
    }

    return Container(
      margin: const EdgeInsets.symmetric(vertical: 8, horizontal: 16),
      decoration: BoxDecoration(
        color: isDark ? theme.cardColor : Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            offset: const Offset(0, 4),
            color: isDark ? Colors.black45 : Colors.grey.withOpacity(.08),
            blurRadius: 12.0,
          )
        ],
        border: Border.all(
          color: isDark ? Colors.grey[850]! : Colors.grey[200]!,
          width: 1,
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Card Upper Area: Date & Status
          Padding(
            padding: const EdgeInsets.only(left: 16, right: 16, top: 16, bottom: 8),
            child: Row(
              children: [
                // Highlighted visual Date block
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: isDark ? Colors.blueGrey[900] : primaryColor.withOpacity(0.08),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Column(
                    children: [
                      Text(
                        _getWeekDay(session.startTime),
                        style: GoogleFonts.quicksand(
                          fontSize: 10,
                          fontWeight: FontWeight.bold,
                          color: isDark ? Colors.blue[300] : primaryColor,
                        ),
                      ),
                      Text(
                        _getDayNum(session.startTime),
                        style: GoogleFonts.quicksand(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: isDark ? Colors.blue[300] : primaryColor,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 16),

                // Session Metadata & Times
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        _formatSessionDate(session.startTime),
                        style: GoogleFonts.quicksand(
                          fontSize: 14,
                          fontWeight: FontWeight.bold,
                          color: isDark ? Colors.white : Colors.black87,
                        ),
                      ),
                      const SizedBox(height: 2),
                      Row(
                        children: [
                          Icon(Icons.access_time, size: 13, color: isDark ? Colors.grey : Colors.grey[600]),
                          const SizedBox(width: 4),
                          Text(
                            _formatSessionTime(session.startTime, session.endTime),
                            style: GoogleFonts.quicksand(
                              fontSize: 12,
                              fontWeight: FontWeight.w500,
                              color: isDark ? Colors.grey : Colors.grey[600],
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),

                // Status tag
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: statusColor.withOpacity(0.12),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    statusLabel,
                    style: GoogleFonts.quicksand(
                      fontSize: 11,
                      fontWeight: FontWeight.bold,
                      color: statusColor,
                    ),
                  ),
                ),
              ],
            ),
          ),

          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 16),
            child: Divider(height: 8),
          ),

          // Card Middle Area: Avatar & Info
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Row(
              children: [
                CircleAvatar(
                  radius: 24,
                  backgroundImage: NetworkImage(avatarUrl),
                  backgroundColor: isDark ? Colors.grey[800] : Colors.grey[350],
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        companionName,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: GoogleFonts.quicksand(
                          fontWeight: FontWeight.bold,
                          fontSize: 15,
                          color: isDark ? Colors.white : Colors.black87,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Row(
                        children: [
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                            decoration: BoxDecoration(
                              color: isDark ? Colors.grey[800] : Colors.grey[200],
                              borderRadius: BorderRadius.circular(4),
                            ),
                            child: Text(
                              session.mode.toUpperCase(),
                              style: GoogleFonts.quicksand(
                                fontSize: 9,
                                fontWeight: FontWeight.bold,
                                color: isDark ? Colors.white70 : Colors.grey[800],
                              ),
                            ),
                          ),
                          const SizedBox(width: 8),
                          Text(
                            formattedPrice,
                            style: GoogleFonts.quicksand(
                              fontSize: 13,
                              fontWeight: FontWeight.bold,
                              color: isDark ? Colors.green[300] : Colors.green[700],
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),

          // Card Actions Footer
          if (session.status != 'cancelled')
            Padding(
              padding: const EdgeInsets.only(left: 16, right: 16, bottom: 16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  // Payment Button (unpaid & active status)
                  if (session.price > 0 && session.paymentStatus != 'paid')
                    ElevatedButton.icon(
                      icon: const Icon(Icons.payment, size: 16),
                      label: Text(
                        "Pay Session Fee",
                        style: GoogleFonts.quicksand(fontWeight: FontWeight.bold, fontSize: 13),
                      ),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.orangeAccent,
                        foregroundColor: Colors.white,
                        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
                        elevation: 0,
                      ),
                      onPressed: () async {
                        final paymentUrl = await sessionProvider.getPaymentUrl(session.sessionId);
                        if (paymentUrl != null && mounted) {
                          final success = await Navigator.push<bool>(
                            context,
                            MaterialPageRoute(
                              builder: (context) => PaymentWebView(
                                paymentUrl: paymentUrl,
                                sessionId: session.sessionId,
                              ),
                            ),
                          );
                          if (success == true) {
                            sessionProvider.fetchSessions();
                          }
                        }
                      },
                    ),

                  // Call/Join consulting room button
                  if (session.mode == 'online' && (session.paymentStatus == 'paid' || session.status == 'active' || session.status == 'pending')) ...[
                    if (session.price > 0 && session.paymentStatus != 'paid') const SizedBox(width: 8),
                    ElevatedButton.icon(
                      icon: const Icon(Icons.video_call, size: 18),
                      label: Text(
                        "Join consultation Room",
                        style: GoogleFonts.quicksand(fontWeight: FontWeight.bold, fontSize: 13),
                      ),
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.green,
                        foregroundColor: Colors.white,
                        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
                        elevation: 0,
                      ),
                      onPressed: () {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => CallScreen(
                              conversationId: session.sessionId,
                              receiverId: companionId,
                              receiverName: companionName,
                              isCaller: true,
                            ),
                          ),
                        );
                      },
                    ),
                  ],
                ],
              ),
            ),
        ],
      ),
    );
  }
}
