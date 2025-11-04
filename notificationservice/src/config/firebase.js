import admin from "firebase-admin";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const serviceAccount = path.join(
  __dirname,
  "../../psyconnect-041125-firebase-adminsdk-fbsvc-531ba83b3b.json"
);

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
});

export const fcm = admin.messaging();
