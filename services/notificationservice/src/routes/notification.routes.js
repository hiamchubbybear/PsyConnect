import express from "express";
import { NotificationController } from "../controllers/notification.controller.js";

const router = express.Router();

router.post("/mock", NotificationController.sendMock);
router.get("/me", NotificationController.getMyNotifications);
router.patch("/:id/read", NotificationController.markAsRead);

export default router;
