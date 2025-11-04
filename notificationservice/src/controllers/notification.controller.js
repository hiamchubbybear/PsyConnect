import { NotificationService } from "../services/notification.service.js";

export const NotificationController = {
  async saveToken(req, res) {
    try {
      const { userId, token } = req.body;
      if (!userId || !token)
        return res.status(400).json({ message: "Missing userId or token" });

      await NotificationService.saveToken(userId, token);
      res.json({ message: "Token saved successfully" });
    } catch (err) {
      res.status(500).json({ message: err.message });
    }
  },

  async send(req, res) {
    try {
      const { userId, title, body } = req.body;
      await NotificationService.sendToUser(userId, title, body);
      res.json({ message: "Notification sent successfully" });
    } catch (err) {
      res.status(500).json({ message: err.message });
    }
  },

  async sendMock(req, res) {
    try {
      console.log("Mock notification hit");
      const { title, body } = req.body;
      const token =
        "etQIazr8266kObDWWJ3n2w:APA91bEvnr65F2PayA2EptT-Kh3YXF4XXllJw-YgiUX415AWnxV5qGtLN8uYga5Gr3y88smwNrr00ivWZ2la_uyFL6pHxrmRl1jsj6GQ2fXMJkiHTE7wpzM";

      await NotificationService.send(
        token,
        title || "Mock Notification",
        body || "This is a test push notification"
      );
      res.json({ message: "Mock notification sent successfully" });
    } catch (err) {
      console.error("Error sending mock notification:", err);
      res.status(500).json({ message: "Failed to send mock notification" });
    }
  },
};
