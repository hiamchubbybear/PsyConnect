const loadTemplate = require("./template_loader");
const sendEmail = require("./email");

const sendActivateEmail = async ({ username, code, email, fullname }) => {
    const template = loadTemplate("verified.html");
    if (!template) throw new Error("Template not found");

    const html = template
        .replace("{USERNAME}", username)
        .replace("{CODE}", code);

    await sendEmail({
        to: email,
        subject: `Hey ${fullname}!! Your Verification Code`,
        html
    });
};

module.exports = sendActivateEmail;
