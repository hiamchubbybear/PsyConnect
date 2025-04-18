const fs = require("fs");
const path = require("path");

const loadTemplate = (filename) => {
    try {
        return fs.readFileSync(path.join(__dirname, "../templates", filename), "utf-8");
    } catch (err) {
        console.error(`Error loading template ${filename}:`, err);
        return null;
    }
};

module.exports = loadTemplate;
