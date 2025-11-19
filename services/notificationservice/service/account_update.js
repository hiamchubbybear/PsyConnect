import { sendEmail } from "./email.service.js";
import { loadTemplate } from "./template_loader.js";

export const sendAccountUpdateEmail = async ({ username, email }) => {
  console.log("[DEBUG] sendAccountUpdateEmail called with:", {
    username,
    email,
  });

  const template = loadTemplate("account-update.html");
  if (!template) throw new Error("Template not found");

  const html = template.replace("{USERNAME}", username);

  console.log("[DEBUG] Sending account update email...");
  await sendEmail({
    to: email,
    subject: `${username}, your account information was updated`,
    html,
  });
  console.log("[DEBUG] Account update email sent successfully!");
};

export const sendResetPasswordEmail = async ({ username, email, code }) => {
  console.log("[DEBUG] sendResetPasswordEmail called with:", {
    username,
    email,
    code,
  });

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
  console.log("[DEBUG] Reset password email sent successfully!");
};
