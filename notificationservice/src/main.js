import dotenv from "dotenv";
import app from "./app.js";
import { sequelize } from "./config/db.js";
import startConsumer from "./kafka/consumer.js";

dotenv.config();

const PORT = process.env.PORT || 8082;

const startServer = async () => {
  try {
    await sequelize.authenticate();
    console.log("Connected to MySQL");

    app.listen(PORT, async () => {
      console.log(`Notification service running on port ${PORT}`);
      await startConsumer();
    });
  } catch (err) {
    console.error("Failed to start server:", err);
  }
};

startServer();
