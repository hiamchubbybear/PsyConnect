import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

export const loadTemplate = (filename) => {
  try {
    return fs.readFileSync(
      path.join(__dirname, "../templates", filename),
      "utf-8"
    );
  } catch (err) {
    console.error(`Error loading template ${filename}:`, err);
    return null;
  }
};
