const { Kafka } = require("kafkajs");

const kafka = new Kafka({
  clientId: "notification-service",
  brokers: ["localhost:9092"]
});

const consumer = kafka.consumer({
  groupId: "notification.user-create",
  minBytes: 5,
  maxBytes: 10000,
});

const startConsumer = async () => {
  await consumer.connect();
  await consumer.subscribe({
    topic: "notification.user-create",
    fromBeginning: true
  });
  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
        res = JSON.parse(message.value.toString());
        console.log(res.username, res.code , res.email , res.fullname);
        sendActivateEmailKafka(res.username, res.code , res.email , res.fullname);
    }
  });
};
const listenOnActivate = async () => {
    await consumer.connect();
    await consumer.subscribe({
        topic : "notification.user-activate",
        fromBeginning : true
    })
    await consumer.run(
        {
        eachMessage: async ({ topic, partition, message }) => {
        res = JSON.parse(message.value.toString());
        console.log(res.username, res.code , res.email , res.fullname);
        sendActivateEmailKafka(res.username, res.code , res.email , res.fullname);
            }
        }
    )
}

module.exports = startConsumer;
