import dotenv from "dotenv";
import KafkaLogger from "../utils/kafkaLogger.js";
import loggingMiddleware from "../utils/loggingMiddleware.js";
import app from "./app.js";
import { sequelize } from "./config/db.js";
import startConsumer from "./kafka/consumer.js";

dotenv.config();

const PORT = process.env.PORT || 8082;

// Initialize Kafka Logger
const kafkaBrokers = (process.env.KAFKA_BROKERS || 'localhost:9092').split(',');
const kafkaLogger = new KafkaLogger({
  brokers: kafkaBrokers,
  serviceName: 'notification-service',
  environment: process.env.ENVIRONMENT || 'development',
  version: '1.0.0',
  topic: 'logging-service',
});

// Add logging middleware to app
app.use(loggingMiddleware(kafkaLogger));

const startServer = async () => {
  try {
    await sequelize.authenticate();
    console.log("Connected to MySQL");
    kafkaLogger.info("Notification service connected to MySQL", null);

    app.listen(PORT, async () => {
      console.log(`Notification service running on port ${PORT}`);
      kafkaLogger.info("Notification service started", { port: PORT });

      await startConsumer();
      kafkaLogger.info("Kafka consumer started", null);
    });
  } catch (err) {
    console.error("Failed to start server:", err);
    kafkaLogger.fatal("Failed to start notification service", {
      error: err.message,
      stackTrace: err.stack,
    });
  }
};

// Graceful shutdown
process.on('SIGTERM', async () => {
  kafkaLogger.info("Notification service shutting down", null);
  await kafkaLogger.disconnect();
  process.exit(0);
});

startServer();

// Export logger for use in other modules
export { kafkaLogger };
