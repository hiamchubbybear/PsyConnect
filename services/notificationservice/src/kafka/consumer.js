import { Kafka } from "kafkajs";
import {
  sendAccountUpdateEmail,
  sendResetPasswordEmail,
} from "../services/account_update.service.js";
import { sendActivateEmail } from "../services/activate_email.service.js";
import { NotificationService } from "../services/notification.service.js";

const kafka = new Kafka({
  clientId: "notification-service",
  brokers: [process.env.KAFKA_BROKER || "localhost:9092"],
});

const consumer = kafka.consumer({ groupId: "notification-group" });

const recentNotifications = new Map();
const SPAM_WINDOW_MS = 60 * 1000;

function canSend(userId, type) {
  const key = `${userId}:${type}`;
  const now = Date.now();
  const lastSent = recentNotifications.get(key) || 0;
  if (now - lastSent < SPAM_WINDOW_MS) return false;
  recentNotifications.set(key, now);
  return true;
}

const emailHandlers = {
  "notification.user-create": sendActivateEmail,
  "notification.user-activate": sendActivateEmail,
  "notification.account-change": sendAccountUpdateEmail,
  "notification.user-reset": sendResetPasswordEmail,

  "notification.therapist-approve": async (data) => {
    const { email, fullname } = data;
    await NotificationService.sendToUser(
      data.userId,
      "Your profile was approved!",
      "Congratulations! Your therapist profile has been approved."
    );
    console.log(`[Kafka] Therapist approved email sent to ${email}`);
  },

  "notification.therapist-reject": async (data) => {
    const { email } = data;
    await NotificationService.sendToUser(
      data.userId,
      "Profile rejected",
      "Unfortunately, your profile was rejected. Please update and resubmit."
    );
    console.log(`[Kafka] Therapist rejected email sent to ${email}`);
  },
};

const pushHandlers = {
  "notification.push.new-message": async (data) => {
    const { userId, from } = data;
    if (!canSend(userId, "new-message")) return;
    await NotificationService.sendToUser(
      userId,
      "New message",
      `You’ve received a new message from ${from || "someone"}.`
    );
  },

  "notification.push.consultation-created": async (data) => {
    const { userId } = data;
    await NotificationService.sendToUser(
      userId,
      "New consultation request",
      "A client has just booked a consultation with you."
    );
  },

  "notification.push.consultation-updated": async (data) => {
    const { userId } = data;
    await NotificationService.sendToUser(
      userId,
      "Consultation updated",
      "Your consultation schedule has been updated."
    );
  },

  "notification.push.consultation-reminder": async (data) => {
    const { userId } = data;
    if (!canSend(userId, "consultation-reminder")) return;
    await NotificationService.sendToUser(
      userId,
      "Upcoming consultation",
      "Reminder: You have a consultation starting soon."
    );
  },

  "notification.push.consultation-completed": async (data) => {
    const { userId } = data;
    await NotificationService.sendToUser(
      userId,
      "Consultation completed",
      "Your consultation has been completed successfully."
    );
  },

  "notification.push.profile-update": async (data) => {
    const { userId } = data;
    if (!canSend(userId, "profile-update")) return;
    await NotificationService.sendToUser(
      userId,
      "Profile update required",
      "Your therapist profile needs verification or update."
    );
  },

  "notification.push.profile-approved": async (data) => {
    const { userId } = data;
    await NotificationService.sendToUser(
      userId,
      "Profile approved",
      "Your therapist profile has been approved!"
    );
  },

  "notification.push.profile-rejected": async (data) => {
    const { userId } = data;
    await NotificationService.sendToUser(
      userId,
      "Profile rejected",
      "Your therapist profile was rejected. Please review and resubmit."
    );
  },

  "notification.push.new-review": async (data) => {
    const { userId, reviewer } = data;
    if (!canSend(userId, "new-review")) return;
    await NotificationService.sendToUser(
      userId,
      "New review received",
      `${reviewer || "A client"} just left a review for you.`
    );
  },

  "notification.push.new-client": async (data) => {
    const { userId, clientName } = data;
    await NotificationService.sendToUser(
      userId,
      "New client connected",
      `${clientName || "A new user"} has booked a consultation with you!`
    );
  },

  "notification.push.system": async (data) => {
    const { userId, title, message } = data;
    await NotificationService.sendToUser(userId, title, message);
  },

  // Social Engagement Handlers
  "notification.social.post-upvote": async (data) => {
    const { userId, voterName, postId } = data;
    await NotificationService.sendToUser(
      userId,
      "New Upvote! ⬆️",
      `${voterName || "Someone"} upvoted your post.`,
      "upvote",
      { postId }
    );
  },

  "notification.social.post-bookmark": async (data) => {
    const { userId, bookmarkerName, postId } = data;
    await NotificationService.sendToUser(
      userId,
      "Post Bookmarked! 🔖",
      `${bookmarkerName || "Someone"} saved your post.`,
      "bookmark",
      { postId }
    );
  },

  "notification.social.post-share": async (data) => {
    const { userId, sharerName, postId } = data;
    await NotificationService.sendToUser(
      userId,
      "Post Shared! 🔗",
      `${sharerName || "Someone"} shared your post.`,
      "share",
      { postId }
    );
  },

  "notification.social.post-comment": async (data) => {
    const { userId, commenterName, postId, commentId } = data;
    await NotificationService.sendToUser(
      userId,
      "New Comment! 💬",
      `${commenterName || "Someone"} commented on your post.`,
      "comment",
      { postId, commentId }
    );
  },

  "notification.social.user-follow": async (data) => {
    const { userId, followerName, followerId } = data;
    await NotificationService.sendToUser(
      userId,
      "New Follower! 👤",
      `${followerName || "Someone"} started following you.`,
      "follow",
      { followerId }
    );
  },
};

const topicHandlers = { ...emailHandlers, ...pushHandlers };

const startConsumer = async () => {
  try {
    await consumer.connect();
    console.log("[Kafka] Consumer connected");

    for (const topic of Object.keys(topicHandlers)) {
      await consumer.subscribe({ topic, fromBeginning: true });
      console.log(`[Kafka] Subscribed to: ${topic}`);
    }

    await consumer.run({
      eachMessage: async ({ topic, message }) => {
        try {
          const data = JSON.parse(message.value.toString());
          const handler = topicHandlers[topic];
          if (handler) {
            await handler(data);
            console.log(`[Kafka] Processed ${topic}`);
          } else {
            console.warn(`[Kafka] No handler for topic: ${topic}`);
          }
        } catch (err) {
          console.error(`[Kafka] Error processing ${topic}:`, err);
        }
      },
    });
  } catch (err) {
    console.error("[Kafka] Consumer failed to start:", err);
  }
};

export default startConsumer;
