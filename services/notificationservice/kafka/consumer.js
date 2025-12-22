import { Kafka } from "kafkajs";
import {
  sendAccountUpdateEmail,
  sendResetPasswordEmail,
} from "../service/account_update.js";
import { sendActivateEmail } from "../service/activate_email.js";

const kafka = new Kafka({
  clientId: "notification-service",
  brokers: [process.env.KAFKA_BROKER || "localhost:9092"],
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
    console.log(data);
    await sendResetPasswordEmail(data);
  },
};

const startConsumer = async () => {
  let connected = false;
  while (!connected) {
    try {
      await consumer.connect();
      connected = true;
      console.log("Connected to Kafka");
    } catch (err) {
      console.error(
        "Failed to connect to Kafka, retrying in 5s...",
        err.message
      );
      await new Promise((resolve) => setTimeout(resolve, 5000));
    }
  }

  for (const topic of Object.keys(topicHandlers)) {
    await consumer.subscribe({ topic, fromBeginning: true });
  }

  await consumer.run({
    eachMessage: async ({ topic, message }) => {
      try {
        const data = JSON.parse(message.value.toString());
        if (topicHandlers[topic]) {
          await topicHandlers[topic](data);
        } else {
          console.warn("No handler for topic:", topic);
        }
      } catch (err) {
        console.error("Kafka message error:", err);
      }
    },
  });
};

export default startConsumer;
