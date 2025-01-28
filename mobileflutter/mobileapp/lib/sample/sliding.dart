import 'package:flutter/material.dart';

void main() => runApp(MyApp());

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: SlidingButtonWithDialog(),
    );
  }
}

class SlidingButtonWithDialog extends StatefulWidget {
  @override
  _SlidingButtonWithDialogState createState() =>
      _SlidingButtonWithDialogState();
}

class _SlidingButtonWithDialogState extends State<SlidingButtonWithDialog> {
  double _dragPosition = 0.0;
  final double _dragThreshold = 200.0;

  void _onDragUpdate(DragUpdateDetails details) {
    setState(() {
      _dragPosition += details.delta.dx;
      if (_dragPosition > _dragThreshold) {
        _dragPosition = _dragThreshold;
      } else if (_dragPosition < 0) {
        _dragPosition = 0;
      }
    });
  }

  void _onDragEnd(DragEndDetails details) {
    if (_dragPosition >= _dragThreshold) {
      _showDialog();
    }
    setState(() {
      _dragPosition = 0.0; // Reset vị trí trượt
    });
  }

  void _showDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text('Thông báo'),
          content: Text('Bạn đã trượt thành công!'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: Text('Đóng'),
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Sliding Button Dialog')),
      body: Center(
        child: Stack(
          children: [
            Container(
              width: 300,
              height: 60,
              decoration: BoxDecoration(
                color: Colors.grey[300],
                borderRadius: BorderRadius.circular(30),
              ),
              alignment: Alignment.centerLeft,
              padding: EdgeInsets.only(left: 20),
              child: Text(
                'Lướt sang để hiển thị',
                style: TextStyle(color: Colors.black54),
              ),
            ),
            Positioned(
              left: _dragPosition,
              child: GestureDetector(
                onHorizontalDragUpdate: _onDragUpdate,
                onHorizontalDragEnd: _onDragEnd,
                child: Container(
                  width: 60,
                  height: 60,
                  decoration: BoxDecoration(
                    color: Colors.blue,
                    borderRadius: BorderRadius.circular(30),
                  ),
                  child: Icon(Icons.arrow_forward, color: Colors.white),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
