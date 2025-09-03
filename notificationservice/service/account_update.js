const loadTemplate = require("./template_loader");
const sendEmail = require("./email");

const sendAccountUpdateEmail = async ({ username, email }) => {
  const template = loadTemplate("account-update.html");
  if (!template) throw new Error("Template not found");

  const html = template.replace("{USERNAME}", username);

  await sendEmail({
    to: email,
    subject: `🔔 ${username}, your account information was updated`,
    html,
  });
};
const sendResetPasswordEmail = async ({ username, email, code, fullname }) => {
  const template = loadTemplate("reset-password.html");
  if (!template) throw new Error("Template not found");

  const html = template
    .replace("{USERNAME}", username)
    .replace("{TOKEN}", code)
    .replace("{EMAIL}", email);

  await sendEmail({
    to: email,
    subject: `Hey ${username}!! Reset Your Password (15 min valid)`,
    html,
  });
};

module.exports = { sendAccountUpdateEmail, sendResetPasswordEmail };
