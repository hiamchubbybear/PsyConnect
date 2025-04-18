const loadTemplate = require("./template_loader");
const sendEmail = require("./email");

const sendAccountUpdateEmail = async ({ username, email }) => {
    const template = loadTemplate("account-update.html");
    if (!template) throw new Error("Template not found");

    const html = template.replace("{USERNAME}", username);

    await sendEmail({
        to: email,
        subject: `🔔 ${username}, your account information was updated`,
        html
    });
};

module.exports = sendAccountUpdateEmail;
