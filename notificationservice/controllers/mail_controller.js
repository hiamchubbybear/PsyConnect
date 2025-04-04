const express = require("express");
const router = express.Router();
const {sendEmail} = require("../service/email.js");
const {request} = require("express");

// Middleware để log thông tin request
router.use((req, res, next) => {
    console.log("----- Request Log -----");
    console.log(`Method: ${req.method}`);
    console.log(`URL: ${req.originalUrl}`);
    console.log("Headers:", req.headers);
    console.log("Query Parameters:", req.query);
    console.log("Body:", req.body);
    console.log("------------------------");
    next();
});

sendActivateEmail = async (req, res) => {
    try {
        const {username, code, email, fullname} = req.body;
        console.log("username: " , username , " code : " ,code," email : ",email," fullname : ",fullname )
        if (!username || !code || !email || !fullname) {
            return res.status(400).json({message: "Missing required fields"});
        }
        await sendEmail({username, code, email, fullname});
        res.status(200).json({message: "Email sent successfully"});
    } catch (error) {
        console.error(error);
        res.status(500).json({message: "Failed to send email"});
    }
};
sendActivateEmailKafka = async (username, code, email, fullname) =>  {
    if (!username || !code || !email || !fullname) {
        console.error("Missing field");
    }
    await sendEmail({username, code, email, fullname}).catch((err) => {
        console.error(error);
    } );
    return true;

}
router.post("/internal/user",sendActivateEmail);
router.post("/internal/account/req/activate",sendActivateEmail);
    module.exports = router;
