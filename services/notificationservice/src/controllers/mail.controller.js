import { sendAccountUpdateEmail } from "../services/account_update.service.js";
import { sendActivateEmail } from "../services/activate_email.service.js";

export const MailController = {
  async sendActivationEmail(req, res) {
    try {
      const { username, code, email, fullname } = req.body;
      if (!username || !code || !email || !fullname) {
        return res.status(400).json({ message: "Missing fields" });
      }
      await sendActivateEmail({ username, code, email, fullname });
      res.json({ message: "Activation email sent" });
    } catch (err) {
      console.error("Failed to send activation email:", err);
      res.status(500).json({ message: "Failed to send email" });
    }
  },

  async sendAccountUpdateEmail(req, res) {
    try {
      const { username, email } = req.body;
      if (!username || !email) {
        return res.status(400).json({ message: "Missing fields" });
      }
      await sendAccountUpdateEmail({ username, email });
      res.json({ message: "Account update email sent" });
    } catch (err) {
      console.error("Failed to send account update email:", err);
      res.status(500).json({ message: "Failed to send email" });
    }
  },
};
