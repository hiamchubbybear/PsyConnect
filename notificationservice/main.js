const express = require("express");
const cors = require("cors");
const dotenv = require("dotenv");
const startConsumer = require("./kafka/consumer");

dotenv.config();

const app = express();
app.use(cors());
app.use(express.json());

app.use("/noti", require("./controllers/mail_controller"));

startConsumer().catch(console.error);

const PORT = process.env.PORT || 8082;
app.listen(PORT, () => console.log(`Notification service running on port ${PORT}`));
