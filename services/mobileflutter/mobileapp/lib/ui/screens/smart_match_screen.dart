import 'dart:math';
import 'package:PsyConnect/core/colors/color.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/swipe_card.dart';
import 'package:PsyConnect/provider/swipe_provider.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/ui/screens/consultation_profile_page.dart';
import 'package:PsyConnect/ui/screens/therapist_detail_screen.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';

class SmartMatchScreen extends StatefulWidget {
  const SmartMatchScreen({super.key});

  @override
  State<SmartMatchScreen> createState() => _SmartMatchScreenState();
}

class _SmartMatchScreenState extends State<SmartMatchScreen> with SingleTickerProviderStateMixin {
  Offset _dragOffset = Offset.zero;
  double _angle = 0.0;
  bool _isDragging = false;

  late AnimationController _animController;
  late Animation<Offset> _offsetAnimation;
  late Animation<double> _angleAnimation;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<SwipeProvider>(context, listen: false).fetchRecommendations();
    });

    _animController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 300),
    );

    _animController.addListener(() {
      setState(() {
        _dragOffset = _offsetAnimation.value;
        _angle = _angleAnimation.value;
      });
    });
  }

  @override
  void dispose() {
    _animController.dispose();
    super.dispose();
  }

  void _onPanStart(DragStartDetails details) {
    if (_animController.isAnimating) return;
    setState(() {
      _isDragging = true;
    });
  }

  void _onPanUpdate(DragUpdateDetails details) {
    if (_animController.isAnimating) return;
    setState(() {
      _dragOffset += details.delta;
      // Rotation angle proportional to horizontal drag distance
      _angle = (_dragOffset.dx / 400.0) * (pi / 12.0); 
    });
  }

  void _onPanEnd(DragEndDetails details, SwipeProvider swipeProvider, String therapistId) {
    if (_animController.isAnimating) return;
    setState(() {
      _isDragging = false;
    });

    final threshold = MediaQuery.of(context).size.width * 0.35;

    if (_dragOffset.dx > threshold) {
      // Swiped Right (Match)
      _animateAndComplete(const Offset(600, 0), pi / 10.0, () {
        swipeProvider.swipeRight(therapistId);
      });
    } else if (_dragOffset.dx < -threshold) {
      // Swiped Left (Pass)
      _animateAndComplete(const Offset(-600, 0), -pi / 10.0, () {
        swipeProvider.swipeLeft(therapistId);
      });
    } else {
      // Spring back to center
      _offsetAnimation = Tween<Offset>(begin: _dragOffset, end: Offset.zero).animate(
        CurvedAnimation(parent: _animController, curve: Curves.elasticOut),
      );
      _angleAnimation = Tween<double>(begin: _angle, end: 0.0).animate(
        CurvedAnimation(parent: _animController, curve: Curves.elasticOut),
      );
      _animController.forward(from: 0.0);
    }
  }

  void _animateAndComplete(Offset endOffset, double endAngle, VoidCallback onComplete) {
    _offsetAnimation = Tween<Offset>(begin: _dragOffset, end: endOffset).animate(
      CurvedAnimation(parent: _animController, curve: Curves.easeOut),
    );
    _angleAnimation = Tween<double>(begin: _angle, end: endAngle).animate(
      CurvedAnimation(parent: _animController, curve: Curves.easeOut),
    );

    _animController.forward(from: 0.0).then((_) {
      onComplete();
      setState(() {
        _dragOffset = Offset.zero;
        _angle = 0.0;
      });
    });
  }

  void _swipeBtnAction(bool right, SwipeProvider swipeProvider, String therapistId) {
    if (_animController.isAnimating) return;
    final endOffset = right ? const Offset(600, 0) : const Offset(-600, 0);
    final endAngle = right ? pi / 10.0 : -pi / 10.0;

    _offsetAnimation = Tween<Offset>(begin: Offset.zero, end: endOffset).animate(
      CurvedAnimation(parent: _animController, curve: Curves.easeInOut),
    );
    _angleAnimation = Tween<double>(begin: 0.0, end: endAngle).animate(
      CurvedAnimation(parent: _animController, curve: Curves.easeInOut),
    );

    _animController.forward(from: 0.0).then((_) {
      if (right) {
        swipeProvider.swipeRight(therapistId);
      } else {
        swipeProvider.swipeLeft(therapistId);
      }
      setState(() {
        _dragOffset = Offset.zero;
        _angle = 0.0;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    final swipeProvider = Provider.of<SwipeProvider>(context);
    final isDark = Provider.of<ThemeProvider>(context).isDarkMode;
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: isDark ? Colors.black : const Color(0xFFF7F9FC),
      appBar: AppBar(
        title: Text(
          "Smart Therapist Match",
          style: GoogleFonts.quicksand(
            fontWeight: FontWeight.bold,
            fontSize: 18,
            color: theme.textTheme.titleLarge?.color,
          ),
        ),
        centerTitle: true,
        elevation: 0.5,
        backgroundColor: isDark ? theme.cardColor : Colors.white,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: swipeProvider.isLoading
            ? Center(
                child: CircularProgressIndicator(
                  color: isDark ? secondaryColor : primaryColor,
                ),
              )
            : swipeProvider.missingProfile
                ? _buildMissingProfileState(isDark)
                : swipeProvider.cards.isEmpty
                    ? _buildEmptyRecommendationsState(isDark, swipeProvider)
                    : Column(
                        children: [
                          Expanded(
                            child: Stack(
                              clipBehavior: Clip.none,
                              children: [
                                // Bottom/Next Card preview
                                if (swipeProvider.cards.length > 1)
                                  Positioned.fill(
                                    child: Transform.scale(
                                      scale: 0.95,
                                      child: _buildTinderCard(swipeProvider.cards[1], false, isDark),
                                    ),
                                  ),

                                // Top interactive card
                                Positioned.fill(
                                  child: GestureDetector(
                                    onPanStart: _onPanStart,
                                    onPanUpdate: _onPanUpdate,
                                    onPanEnd: (details) => _onPanEnd(details, swipeProvider, swipeProvider.cards[0].profileId),
                                    child: Transform.translate(
                                      offset: _dragOffset,
                                      child: Transform.rotate(
                                        angle: _angle,
                                        child: _buildTinderCard(swipeProvider.cards[0], true, isDark),
                                      ),
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(height: 16),
                          
                          // Control Actions buttons
                          _buildControlRow(swipeProvider, swipeProvider.cards[0].profileId),
                        ],
                      ),
      ),
    );
  }

  Widget _buildMissingProfileState(bool isDark) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(kDefault * 1.5),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.assignment_ind_outlined, size: 64, color: Colors.teal[200]),
            const SizedBox(height: 16),
            Text(
              "Missing Consultation Profile",
              style: GoogleFonts.quicksand(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: isDark ? Colors.white : Colors.black87,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              "Smart matchmaking matches you to specialists based on your issues, preferred fees, and modes. Please fill out your preferences first.",
              textAlign: TextAlign.center,
              style: GoogleFonts.quicksand(
                color: Colors.grey,
                fontSize: 13,
              ),
            ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => const ConsultationProfilePage()),
                );
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: const Color.fromARGB(255, 115, 170, 163),
                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
              ),
              child: Text(
                "Create Profile Now",
                style: GoogleFonts.quicksand(fontWeight: FontWeight.bold, color: Colors.white),
              ),
            )
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyRecommendationsState(bool isDark, SwipeProvider swipeProvider) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(kDefault * 1.5),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.done_all_rounded, size: 64, color: Colors.grey[400]),
            const SizedBox(height: 16),
            Text(
              "You're all swiped up!",
              style: GoogleFonts.quicksand(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: isDark ? Colors.white : Colors.black87,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              "We have run out of therapist recommendations for today. We will generate fresh matching lists for you shortly.",
              textAlign: TextAlign.center,
              style: GoogleFonts.quicksand(
                color: Colors.grey,
                fontSize: 13,
              ),
            ),
            const SizedBox(height: 24),
            ElevatedButton.icon(
              icon: const Icon(Icons.refresh),
              label: Text(
                "Check Again",
                style: GoogleFonts.quicksand(fontWeight: FontWeight.bold),
              ),
              onPressed: () => swipeProvider.fetchRecommendations(),
              style: ElevatedButton.styleFrom(
                backgroundColor: isDark ? Colors.blueGrey[800] : primaryColor,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
              ),
            )
          ],
        ),
      ),
    );
  }

  Widget _buildTinderCard(SwipeCard card, bool active, bool isDark) {
    final theme = Theme.of(context);
    final avatarUrl = card.avatarUrl.isNotEmpty ? card.avatarUrl : 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';

    return Container(
      decoration: BoxDecoration(
        color: isDark ? theme.cardColor : Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(
          color: isDark ? Colors.grey[850]! : Colors.grey[200]!,
          width: 1,
        ),
        boxShadow: active
            ? [
                BoxShadow(
                  color: isDark ? Colors.black54 : Colors.grey.withOpacity(0.12),
                  blurRadius: 16,
                  offset: const Offset(0, 8),
                )
              ]
            : [],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(20),
        child: SingleChildScrollView(
          physics: const NeverScrollableScrollPhysics(),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // 1. Therapist Photo with Match points overlay
              Stack(
                children: [
                  Image.network(
                    avatarUrl,
                    height: 280,
                    width: double.infinity,
                    fit: BoxFit.cover,
                  ),
                  
                  // Match Percentage badge
                  if (card.points > 0)
                    Positioned(
                      top: 16,
                      left: 16,
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                        decoration: BoxDecoration(
                          color: const Color(0xFF00695C).withOpacity(0.85),
                          borderRadius: BorderRadius.circular(20),
                          boxShadow: const [
                            BoxShadow(color: Colors.black26, blurRadius: 4, offset: Offset(0, 2))
                          ],
                        ),
                        child: Row(
                          children: [
                            const Icon(Icons.flash_on, color: Colors.amber, size: 14),
                            const SizedBox(width: 4),
                            Text(
                              "${(card.points * 100).toInt()}% Match",
                              style: GoogleFonts.quicksand(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 12,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  
                  // Bottom gradient for card photo
                  Positioned.fill(
                    child: Container(
                      decoration: const BoxDecoration(
                        gradient: LinearGradient(
                          colors: [Colors.transparent, Colors.black54],
                          begin: Alignment.topCenter,
                          end: Alignment.bottomCenter,
                          stops: [0.6, 1.0],
                        ),
                      ),
                    ),
                  ),
                  
                  // Details overlay
                  Positioned(
                    bottom: 12,
                    left: 16,
                    right: 16,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          card.name,
                          style: GoogleFonts.quicksand(
                            fontSize: 22,
                            fontWeight: FontWeight.bold,
                            color: Colors.white,
                          ),
                        ),
                        Text(
                          card.title,
                          style: GoogleFonts.quicksand(
                            fontSize: 13,
                            fontWeight: FontWeight.w500,
                            color: Colors.white70,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),

              // 2. Info area
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Experience / Rating / Price indicators row
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        _buildStatIcon(Icons.work_outline_rounded, "${card.experience} Yrs Exp", isDark),
                        _buildStatIcon(Icons.star_rounded, "${card.rating}", isDark, color: Colors.amber),
                        _buildStatIcon(Icons.payments_outlined, "${NumberFormat.compactSimpleCurrency(locale: 'vi_VN').format(card.price)}/Session", isDark, color: Colors.green),
                      ],
                    ),
                    const SizedBox(height: 12),
                    const Divider(height: 8),
                    const SizedBox(height: 8),

                    // Specializations tags
                    if (card.specialization.isNotEmpty) ...[
                      Text(
                        "Specializations",
                        style: GoogleFonts.quicksand(fontSize: 12, fontWeight: FontWeight.bold, color: Colors.grey),
                      ),
                      const SizedBox(height: 6),
                      Wrap(
                        spacing: 6,
                        runSpacing: 4,
                        children: card.specialization.take(4).map((s) => _buildBadge(s, Colors.teal, isDark)).toList(),
                      ),
                      const SizedBox(height: 12),
                    ],

                    // Swipe reasons points
                    if (card.reasons.isNotEmpty) ...[
                      Text(
                        "Why you match",
                        style: GoogleFonts.quicksand(fontSize: 12, fontWeight: FontWeight.bold, color: Colors.grey),
                      ),
                      const SizedBox(height: 6),
                      Column(
                        children: card.reasons.take(3).map((reason) {
                          return Padding(
                            padding: const EdgeInsets.symmetric(vertical: 2),
                            child: Row(
                              children: [
                                const Icon(Icons.check_circle_outline, color: Colors.teal, size: 13),
                                const SizedBox(width: 6),
                                Expanded(
                                  child: Text(
                                    reason,
                                    maxLines: 1,
                                    overflow: TextOverflow.ellipsis,
                                    style: GoogleFonts.quicksand(
                                      fontSize: 11,
                                      fontWeight: FontWeight.w500,
                                      color: isDark ? Colors.white70 : Colors.black87,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          );
                        }).toList(),
                      ),
                    ],
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildStatIcon(IconData icon, String text, bool isDark, {Color? color}) {
    return Row(
      children: [
        Icon(icon, size: 16, color: color ?? Colors.grey),
        const SizedBox(width: 4),
        Text(
          text,
          style: GoogleFonts.quicksand(
            fontSize: 11,
            fontWeight: FontWeight.bold,
            color: isDark ? Colors.white70 : Colors.black87,
          ),
        ),
      ],
    );
  }

  Widget _buildBadge(String text, Color color, bool isDark) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(
        color: color.withOpacity(0.08),
        borderRadius: BorderRadius.circular(6),
        border: Border.all(color: color.withOpacity(0.2), width: 0.5),
      ),
      child: Text(
        text,
        style: GoogleFonts.quicksand(
          fontSize: 9,
          fontWeight: FontWeight.bold,
          color: isDark ? color.withOpacity(0.8) : color,
        ),
      ),
    );
  }

  Widget _buildControlRow(SwipeProvider swipeProvider, String therapistId) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        // 1. Pass button
        _buildActionBtn(
          icon: Icons.close_rounded,
          color: Colors.red,
          onTap: () => _swipeBtnAction(false, swipeProvider, therapistId),
        ),
        
        // 2. View Info details button
        _buildActionBtn(
          icon: Icons.info_outline_rounded,
          color: Colors.blueAccent,
          isSmall: true,
          onTap: () {
            Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => TherapistDetailScreen(therapistId: therapistId),
              ),
            );
          },
        ),

        // 3. Match button
        _buildActionBtn(
          icon: Icons.favorite_rounded,
          color: Colors.green,
          onTap: () => _swipeBtnAction(true, swipeProvider, therapistId),
        ),
      ],
    );
  }

  Widget _buildActionBtn({
    required IconData icon,
    required Color color,
    required VoidCallback onTap,
    bool isSmall = false,
  }) {
    final size = isSmall ? 52.0 : 64.0;
    return GestureDetector(
      onTap: onTap,
      child: Container(
        height: size,
        width: size,
        decoration: BoxDecoration(
          color: Colors.white,
          shape: BoxShape.circle,
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.1),
              blurRadius: 10,
              offset: const Offset(0, 4),
            )
          ],
          border: Border.all(
            color: color.withOpacity(0.4),
            width: 1.5,
          ),
        ),
        child: Center(
          child: Icon(
            icon,
            color: color,
            size: isSmall ? 22 : 30,
          ),
        ),
      ),
    );
  }
}
