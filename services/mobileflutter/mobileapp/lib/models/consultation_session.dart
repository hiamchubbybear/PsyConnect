class ConsultationSession {
  final String sessionId;
  final String therapistId;
  final String clientId;
  final String mode;
  final String startTime;
  final String endTime;
  final String? scheduledDate;
  final String? timeZone;
  final String status;
  final double price;
  final String? paymentStatus;
  final String? callSessionId;
  final String createdAt;

  ConsultationSession({
    required this.sessionId,
    required this.therapistId,
    required this.clientId,
    required this.mode,
    required this.startTime,
    required this.endTime,
    this.scheduledDate,
    this.timeZone,
    required this.status,
    required this.price,
    this.paymentStatus,
    this.callSessionId,
    required this.createdAt,
  });

  factory ConsultationSession.fromJson(Map<String, dynamic> json) {
    return ConsultationSession(
      sessionId: json['session_id'] ?? '',
      therapistId: json['therapist_id'] ?? '',
      clientId: json['client_id'] ?? '',
      mode: json['mode'] ?? 'online',
      startTime: json['start_time'] ?? '',
      endTime: json['end_time'] ?? '',
      scheduledDate: json['scheduled_date'],
      timeZone: json['time_zone'],
      status: json['status'] ?? 'pending_payment',
      price: (json['price'] as num?)?.toDouble() ?? 0.0,
      paymentStatus: json['payment_status'],
      callSessionId: json['call_session_id'],
      createdAt: json['created_at'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'session_id': sessionId,
      'therapist_id': therapistId,
      'client_id': clientId,
      'mode': mode,
      'start_time': startTime,
      'end_time': endTime,
      'scheduled_date': scheduledDate,
      'time_zone': timeZone,
      'status': status,
      'price': price,
      'payment_status': paymentStatus,
      'call_session_id': callSessionId,
      'created_at': createdAt,
    };
  }
}
