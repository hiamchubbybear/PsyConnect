importScripts(
  "https://www.gstatic.com/firebasejs/10.13.0/firebase-app-compat.js"
);
importScripts(
  "https://www.gstatic.com/firebasejs/10.13.0/firebase-messaging-compat.js"
);

const firebaseConfig = {
  apiKey: "AIzaSyBjyhkn71RnXenG4dzhdiRIlIc95YM2bAk",
  authDomain: "psyconnect-041125.firebaseapp.com",
  projectId: "psyconnect-041125",
  storageBucket: "psyconnect-041125.firebasestorage.app",
  messagingSenderId: "130716148033",
  appId: "1:130716148033:web:bafc7ed2934cf9768241c0",
  measurementId: "G-EREPPG37PL",
};
firebase.initializeApp(firebaseConfig);

const messaging = firebase.messaging();
messaging.onBackgroundMessage((payload) => {
  console.log(
    "[firebase-messaging-sw.js] Received background message ",
    payload
  );
  const { title, body } = payload.notification;
  self.registration.showNotification(title || "Notification", {
    body: body || "",
  });
});
