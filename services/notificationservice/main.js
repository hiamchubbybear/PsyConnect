import cors from "cors";
import dotenv from "dotenv";
import express from "express";
import mailController from "./controllers/mail_controller.js";
import startConsumer from "./kafka/consumer.js";

dotenv.config();

const app = express();
app.use(cors());
app.use(express.json());

app.use("/noti", mailController);

startConsumer().catch(console.error);

const PORT = process.env.PORT || 8082;
app.listen(PORT, () =>
  console.log(`Notification service running on port ${PORT}`)
);
