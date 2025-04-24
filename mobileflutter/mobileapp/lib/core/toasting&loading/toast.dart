import 'package:flutter/material.dart';

class ToastService {
  static void showToast({
    required BuildContext context,
    required String message,
    required String title,
    required ToastType type,
    int duration = 5,
  }) {
    late OverlayEntry overlayEntry;
    overlayEntry = OverlayEntry(
      builder: (context) => Positioned(
        bottom: 50,
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

    Overlay.of(context).insert(overlayEntry);
    Future.delayed(Duration(seconds: duration), () {
      overlayEntry.remove();
    });
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
    Color iconColor;

    switch (type) {
      case ToastType.success:
        bgColor = Colors.green.shade100;
        icon = Icons.check_circle;
        iconColor = Colors.green;
        break;
      case ToastType.info:
        bgColor = Colors.blue.shade100;
        icon = Icons.info;
        iconColor = Colors.blue;
        break;
      case ToastType.warning:
        bgColor = Colors.amber.shade100;
        icon = Icons.warning_amber_rounded;
        iconColor = Colors.amber;
        break;
      case ToastType.error:
        bgColor = Colors.red.shade100;
        icon = Icons.error;
        iconColor = Colors.red;
        break;
    }

    return Material(
      color: Colors.transparent,
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 10),
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: bgColor,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: iconColor, width: 1.5),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.1),
              blurRadius: 6,
              offset: const Offset(0, 2),
            ),
          ],
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Icon(icon, color: iconColor, size: 30),
            const SizedBox(width: 10),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    title,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: iconColor,
                    ),
                  ),
                  Text(
                    message,
                    style: const TextStyle(fontSize: 14, color: Colors.black87),
                  ),
                ],
              ),
            ),
            IconButton(
              icon: const Icon(Icons.close, size: 20, color: Colors.black54),
              onPressed: onClose,
            ),
          ],
        ),
      ),
    );
  }
}
