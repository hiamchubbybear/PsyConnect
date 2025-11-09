import { FcmToken } from "../models/fcmToken.model.js";
import { FirebaseService } from "./firebase.service.js";

export const NotificationService = {
  async saveToken(userId, token) {
    const existing = await FcmToken.findOne({ where: { userId } });
    if (existing) {
      existing.token = token;
      await existing.save();
    } else {
      await FcmToken.create({ userId, token });
    }
  },

  async sendToUser(userId, title, body) {
    const record = await FcmToken.findOne({ where: { userId } });
    if (!record) throw new Error("User has no token");
    return FirebaseService.sendPushNotification(record.token, title, body);
  },
};
