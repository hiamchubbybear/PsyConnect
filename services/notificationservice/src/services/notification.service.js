import { FcmToken } from "../models/fcmToken.model.js";
import { Notification } from "../models/notification.model.js";
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

  async sendToUser(userId, title, body, type = "system", metadata = {}) {
    // Persist notification first
    await this.createNotification(userId, { title, body, type, metadata });

    const record = await FcmToken.findOne({ where: { userId } });
    if (record) {
      return FirebaseService.sendPushNotification(record.token, title, body);
    }
    console.log(
      `[NotificationService] No push token for user ${userId}, notification only persisted.`
    );
  },

  async createNotification(userId, data) {
    return Notification.create({
      userId,
      title: data.title,
      body: data.body,
      type: data.type || "system",
      metadata: data.metadata || {},
      isRead: false,
    });
  },

  async getNotifications(userId, limit = 20, skip = 0) {
    return Notification.findAll({
      where: { userId },
      limit: parseInt(limit),
      offset: parseInt(skip),
      order: [["createdAt", "DESC"]],
    });
  },

  async markAsRead(notificationId, userId) {
    return Notification.update(
      { isRead: true },
      { where: { id: notificationId, userId } }
    );
  },
};
