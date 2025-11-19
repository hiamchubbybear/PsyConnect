import express from "express";
import { MailController } from "../controllers/mail.controller.js";

const router = express.Router();

router.post("/activate", MailController.sendActivationEmail);
router.post("/account-update", MailController.sendAccountUpdateEmail);

export default router;
