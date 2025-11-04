import { Kafka } from "kafkajs";
import {
    sendAccountUpdateEmail,
    sendResetPasswordEmail,
} from "../services/account_update.service.js";
import { sendActivateEmail } from "../services/activate_email.service.js";

const kafka = new Kafka({
  clientId: "notification-service",
  brokers: ["localhost:9092"],
});

const consumer = kafka.consumer({ groupId: "notification-group" });
const topicHandlers = {
  "notification.user-create": async (data) => {
    await sendActivateEmail(data);
  },
  "notification.user-activate": async (data) => {
    await sendActivateEmail(data);
  },
  "notification.account-change": async (data) => {
    await sendAccountUpdateEmail(data);
  },
  "notification.user-reset": async (data) => {
    await sendResetPasswordEmail(data);
  },
};

const startConsumer = async () => {
  try {
    await consumer.connect();
    console.log("[Kafka] Consumer connected");
    for (const topic of Object.keys(topicHandlers)) {
      await consumer.subscribe({ topic, fromBeginning: true });
      console.log(`[Kafka] Subscribed to topic: ${topic}`);
    }
    await consumer.run({
      eachMessage: async ({ topic, message }) => {
        try {
          const data = JSON.parse(message.value.toString());
          if (topicHandlers[topic]) {
            await topicHandlers[topic](data);
          } else {
            console.warn("[Kafka] No handler for topic:", topic);
          }
        } catch (err) {
          console.error("[Kafka] Message processing error:", err);
        }
      },
    });
  } catch (err) {
    console.error("[Kafka] Consumer failed to start:", err);
    throw err;
  }
};

export default startConsumer;
