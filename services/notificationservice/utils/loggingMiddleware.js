// Logging middleware for Express (ES modules)
export default function loggingMiddleware(kafkaLogger) {
    return (req, res, next) => {
        const start = Date.now();

        // Capture response
        const originalSend = res.send;
        res.send = function (data) {
            res.send = originalSend;

            const duration = Date.now() - start;
            const fields = {
                ip: req.ip || req.connection.remoteAddress,
                userAgent: req.get('User-Agent'),
            };

            // Get userId from request if available
            if (req.userId) {
                fields.userId = req.userId;
            }

            kafkaLogger.logHttpRequest(
                req.method,
                req.path,
                res.statusCode,
                duration,
                fields
            );

            return originalSend.call(this, data);
        };

        next();
    };
}
