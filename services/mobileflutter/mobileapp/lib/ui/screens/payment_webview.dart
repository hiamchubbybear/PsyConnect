import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';

class PaymentWebView extends StatefulWidget {
  final String paymentUrl;
  final String sessionId;

  const PaymentWebView({
    super.key,
    required this.paymentUrl,
    required this.sessionId,
  });

  @override
  State<PaymentWebView> createState() => _PaymentWebViewState();
}

class _PaymentWebViewState extends State<PaymentWebView> {
  late final WebViewController _controller;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _controller = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setNavigationDelegate(
        NavigationDelegate(
          onPageStarted: (String url) {
            setState(() {
              _isLoading = true;
            });
          },
          onPageFinished: (String url) {
            setState(() {
              _isLoading = false;
            });
          },
          onNavigationRequest: (NavigationRequest request) {
            final url = request.url;
            print("WebView navigating to: $url");
            
            // Detect transaction callback redirects
            if (url.contains("vnp_ResponseCode") || url.contains("/payment-success") || url.contains("/payment/success") || url.contains("status=success")) {
              bool success = false;
              if (url.contains("vnp_ResponseCode=00") || url.contains("status=success") || url.contains("success")) {
                success = true;
              }
              
              _handlePaymentResult(success);
              return NavigationDecision.prevent;
            }
            return NavigationDecision.navigate;
          },
        ),
      )
      ..loadRequest(Uri.parse(widget.paymentUrl));
  }

  void _handlePaymentResult(bool success) {
    if (!mounted) return;
    if (success) {
      ToastService.showToast(
        context: context,
        message: "Your consulting session invoice paid successfully.",
        title: "Payment Success",
        type: ToastType.success,
      );
      Navigator.pop(context, true);
    } else {
      ToastService.showToast(
        context: context,
        message: "Invoice payment cancelled or encountered errors.",
        title: "Payment Failed",
        type: ToastType.error,
      );
      Navigator.pop(context, false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        title: Text(
          'Invoice Checkout',
          style: GoogleFonts.quicksand(
            fontWeight: FontWeight.bold, 
            fontSize: 18,
            color: theme.textTheme.titleLarge?.color,
          ),
        ),
        leading: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () => _handlePaymentResult(false),
        ),
        elevation: 0.5,
        backgroundColor: theme.cardColor,
      ),
      body: Stack(
        children: [
          WebViewWidget(controller: _controller),
          if (_isLoading)
            Center(
              child: CircularProgressIndicator(
                color: theme.primaryColor,
              ),
            ),
        ],
      ),
    );
  }
}
