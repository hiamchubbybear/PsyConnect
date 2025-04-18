const { Kafka } = require("kafkajs");
const sendActivateEmail = require("../service/activate_email");
const sendAccountUpdateEmail = require("../service/account_update");

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
};

const startConsumer = async () => {
    await consumer.connect();
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

module.exports = startConsumer;
