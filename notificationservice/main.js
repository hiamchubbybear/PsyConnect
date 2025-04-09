const express = require("express");
const cors = require("cors");
const startConsumer = require("./kafka/consumer");
const PORT = process.env.PORT || 8082;

require("dotenv").config();

const app = express();

app.use(express.json());

app.use(cors({
    credentials: true,
}));

app.use("/noti/", require("./controllers/mail_controller"));

startConsumer().catch(err => {
    console.error("Failed to start Kafka consumer:", err);
  });

app.listen(PORT, () => {
    console.log(`Server started on port ${PORT}`);
});
