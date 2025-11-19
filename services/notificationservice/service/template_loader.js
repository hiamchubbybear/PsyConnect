import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export const loadTemplate = (templateName) => {
  try {
    const templatePath = path.join(__dirname, "../templates", templateName);
    return fs.readFileSync(templatePath, "utf-8");
  } catch (err) {
    console.error(`Error loading template ${templateName}:`, err);
    return null;
  }
};
