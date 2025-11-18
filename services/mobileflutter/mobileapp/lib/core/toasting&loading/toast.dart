import 'package:flutter/material.dart';

class ToastService {
  static void showToast({
    required BuildContext context,
    required String message,
    required String title,
    required ToastType type,
    int duration = 4,
  }) {
    final mediaQueryData = MediaQuery.of(context);

    late OverlayEntry overlayEntry;
    overlayEntry = OverlayEntry(
      builder: (context) => Positioned(
        bottom: mediaQueryData.size.height*0.155,
        left: 20,
        right: 20,
        child: _ToastWidget(
          message: message,
          title: title,
          type: type,
          onClose: () => overlayEntry.remove(),
        ),
      ),
    );
    final overlay = Overlay.maybeOf(context);
    if (overlay != null) {
      overlay.insert(overlayEntry);
      Future.delayed(Duration(seconds: duration), () {
        if (overlayEntry.mounted) overlayEntry.remove();
      });
    }
  }
}

enum ToastType { success, info, warning, error }

class _ToastWidget extends StatelessWidget {
  final String title;
  final String message;
  final ToastType type;
  final VoidCallback onClose;

  const _ToastWidget({
    Key? key,
    required this.title,
    required this.message,
    required this.type,
    required this.onClose,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    Color bgColor;
    IconData icon;
    Color accentColor;

    switch (type) {
      case ToastType.success:
        bgColor = Colors.green.shade50;
        icon = Icons.check;
        accentColor = Colors.green.shade600;
        break;
      case ToastType.info:
        bgColor = Colors.blue.shade50;
        icon = Icons.info_outline;
        accentColor = Colors.blue.shade600;
        break;
      case ToastType.warning:
        bgColor = Colors.amber.shade50;
        icon = Icons.warning_amber_outlined;
        accentColor = Colors.amber.shade700;
        break;
      case ToastType.error:
        bgColor = Colors.red.shade50;
        icon = Icons.error_outline;
        accentColor = Colors.red.shade600;
        break;
    }

    return Material(
      color: Colors.transparent,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        decoration: BoxDecoration(
          color: bgColor,
          borderRadius: BorderRadius.circular(10),
          border: Border.all(color: accentColor.withOpacity(0.4)),
        ),
        child: Row(
          children: [
            Icon(icon, color: accentColor, size: 20),
            const SizedBox(width: 10),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    title,
                    style: TextStyle(
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                      color: accentColor,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    message,
                    style: const TextStyle(
                      fontSize: 13,
                      color: Colors.black87,
                    ),
                  ),
                ],
              ),
            ),
            GestureDetector(
              onTap: onClose,
              child: const Icon(Icons.close, size: 18, color: Colors.black45),
            ),
          ],
        ),
      ),
    );
  }
}
