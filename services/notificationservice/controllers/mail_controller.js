import express from "express";
import { sendAccountUpdateEmail } from "../service/account_update.js";
import { sendActivateEmail } from "../service/activate_email.js";
const router = express.Router();

router.post("/activate", async (req, res) => {
  try {
    const { username, code, email, fullname } = req.body;
    if (!username || !code || !email || !fullname) {
      return res.status(400).json({ message: "Missing fields" });
    }
    await sendActivateEmail({ username, code, email, fullname });
    res.json({ message: "Activation email sent" });
  } catch (err) {
    console.error(err);
    res.status(500).json({ message: "Failed to send email" });
  }
});

router.post("/account-update", async (req, res) => {
  try {
    const { username, email } = req.body;
    if (!username || !email) {
      return res.status(400).json({ message: "Missing fields" });
    }
    await sendAccountUpdateEmail({ username, email });
    res.json({ message: "Account update email sent" });
  } catch (err) {
    console.error(err);
    res.status(500).json({ message: "Failed to send email" });
  }
});

export default router;
