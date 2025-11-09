import nodemailer from "nodemailer";

const createTransporter = () =>
  nodemailer.createTransport({
    service: process.env.SMPT_SERVICE,
    auth: {
      user: process.env.SMPT_MAIL,
      pass: process.env.SMPT_APP_PASS,
    },
  });

export const sendEmail = async ({ to, subject, html }) => {
  const transporter = createTransporter();

  const mailOptions = {
    from: process.env.SMPT_MAIL,
    to,
    subject,
    html,
  };

  await transporter.sendMail(mailOptions);
};
