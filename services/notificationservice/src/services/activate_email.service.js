import { sendEmail } from "./email.service.js";
import { loadTemplate } from "./template_loader.js";

export const sendActivateEmail = async ({ username, code, email, fullname }) => {
  console.log("[DEBUG] sendActivateEmail called with:", {
    username,
    code,
    email,
    fullname,
  });

  const template = loadTemplate("verified.html");
  if (!template) throw new Error("Template not found");

  const html = template
    .replace("{USERNAME}", username)
    .replace("{CODE}", code);

  console.log("[DEBUG] Sending activation email...");
  await sendEmail({
    to: email,
    subject: `Hey ${fullname}!! Your Verification Code`,
    html,
  });
  console.log("[DEBUG] Activation email sent successfully!");
};
