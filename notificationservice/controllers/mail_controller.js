const express = require("express");
const router = express.Router();
const { sendEmail } = require("../service/email.js");

router.post("/internal/user", async (req, res) => {
    try {
        const { username, code, email, fullname } = req.body;
        if (!username || !code || !email  || !fullname) {
            return res.status(400).json({ message: "Missing required fields" });
        }
        await sendEmail({ username, code, email , fullname });
        res.status(200).json({ message: "Email sent successfully" });
    } catch (error) {
        console.error(error);
        res.status(500).json({ message: "Failed to send email" });
    }
});

module.exports = router;
