import cors from "cors";
import express from "express";
import mailRoutes from "./routes/mail.routes.js";
import notificationRoutes from "./routes/notification.routes.js";

const app = express();

app.use(
  cors({
    origin: ["http://localhost:4200", "https://chessy.dev"],
    methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    credentials: true,
  })
);

app.use(express.json());
app.options("*", cors());

app.use("/mail", mailRoutes);
app.use("/notification", notificationRoutes);

export default app;
