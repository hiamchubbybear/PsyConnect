import express from "express";
import { NotificationController } from "../controllers/notification.controller.js";

const router = express.Router();

router.post("/mock", NotificationController.sendMock);

export default router;
