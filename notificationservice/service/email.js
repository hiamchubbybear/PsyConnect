const fs = require("fs");
const path = require("path");
const nodemailer = require("nodemailer");

const sendEmail = async (options) => {
    let emailTemplate;
    try {
        emailTemplate = fs.readFileSync(path.join(__dirname, "verified.html"), "utf-8");
    } catch (err) {
        console.error("Error reading email template:", err);
        throw new Error("Email template not found");
    }

    emailTemplate = emailTemplate.replace("{USERNAME}", options.username)
                                 .replace("{CODE}", options.code);
    const transporter = nodemailer.createTransport({
        service: process.env.SMPT_SERVICE,
        auth: {
            user: process.env.SMPT_MAIL,
            pass: process.env.SMPT_APP_PASS,
        },
    });
    const mailOptions = {
        from: process.env.SMPT_MAIL,
        to: options.email,
        subject: `Hey ${options.fullname}!! Your Verification Code`,
        html: emailTemplate,
    };

    await transporter.sendMail(mailOptions);
};

module.exports = { sendEmail };
