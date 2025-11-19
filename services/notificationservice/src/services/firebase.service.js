import { fcm } from "../config/firebase.js";

export const FirebaseService = {
  async sendPushNotification(token, title, body) {
    try {
      const message = { token, notification: { title, body } };
      const response = await fcm.send(message);
      console.log("Sent push:", response);
      return response;
    } catch (err) {
      console.error("Push error:", err);
      throw err;
    }
  },
};
