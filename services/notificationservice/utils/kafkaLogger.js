import { Kafka } from 'kafkajs';
import os from 'os';

class KafkaLogger {
    constructor(config) {
        this.serviceName = config.serviceName || 'notification-service';
        this.environment = config.environment || process.env.ENVIRONMENT || 'development';
        this.version = config.version || '1.0.0';
        this.hostname = os.hostname();
        this.topic = config.topic || 'logging-service';

        // Initialize Kafka
        const kafka = new Kafka({
            clientId: this.serviceName,
            brokers: config.brokers || ['localhost:9092'],
        });

        this.producer = kafka.producer();
        this.connected = false;

        // Connect to Kafka
        this.connect();
    }

    async connect() {
        try {
            await this.producer.connect();
            this.connected = true;
            console.log('✅ Kafka logger connected');
        } catch (error) {
            console.error('❌ Failed to connect Kafka logger:', error.message);
        }
    }

    async log(level, message, fields = {}) {
        if (!this.connected) {
            console.warn('Kafka logger not connected, skipping log');
            return;
        }

        try {
            const event = {
                timestamp: new Date().toISOString(),
                level,
                service: this.serviceName,
                message,
                environment: this.environment,
                version: this.version,
                hostname: this.hostname,
            };

            // Extract known fields
            if (fields.traceId) event.traceId = fields.traceId;
            if (fields.userId) event.userId = fields.userId;
            if (fields.action) event.action = fields.action;
            if (fields.error) event.error = fields.error;
            if (fields.stackTrace) event.stackTrace = fields.stackTrace;
            if (fields.method) event.method = fields.method;
            if (fields.path) event.path = fields.path;
            if (fields.statusCode) event.statusCode = fields.statusCode;
            if (fields.duration) event.duration = fields.duration;
            if (fields.ip) event.ip = fields.ip;
            if (fields.userAgent) event.userAgent = fields.userAgent;

            // Remaining fields go to metadata
            const metadata = { ...fields };
            delete metadata.traceId;
            delete metadata.userId;
            delete metadata.action;
            delete metadata.error;
            delete metadata.stackTrace;
            delete metadata.method;
            delete metadata.path;
            delete metadata.statusCode;
            delete metadata.duration;
            delete metadata.ip;
            delete metadata.userAgent;

            if (Object.keys(metadata).length > 0) {
                event.metadata = metadata;
            }

            // Send to Kafka
            await this.producer.send({
                topic: this.topic,
                messages: [
                    {
                        value: JSON.stringify(event),
                    },
                ],
            });
        } catch (error) {
            console.error('Failed to send log to Kafka:', error.message);
        }
    }

    debug(message, fields) {
        return this.log('DEBUG', message, fields);
    }

    info(message, fields) {
        return this.log('INFO', message, fields);
    }

    warn(message, fields) {
        return this.log('WARN', message, fields);
    }

    error(message, fields) {
        return this.log('ERROR', message, fields);
    }

    fatal(message, fields) {
        return this.log('FATAL', message, fields);
    }

    audit(message, fields) {
        return this.log('AUDIT', message, fields);
    }

    logHttpRequest(method, path, statusCode, duration, fields = {}) {
        fields.method = method;
        fields.path = path;
        fields.statusCode = statusCode;
        fields.duration = duration;

        let level = 'INFO';
        if (statusCode >= 500) {
            level = 'ERROR';
        } else if (statusCode >= 400) {
            level = 'WARN';
        }

        const message = `${method} ${path} - ${statusCode} (${duration}ms)`;
        return this.log(level, message, fields);
    }

    logAction(action, userId, message, fields = {}) {
        fields.action = action;
        fields.userId = userId;
        return this.log('AUDIT', message, fields);
    }

    async disconnect() {
        if (this.connected) {
            await this.producer.disconnect();
            this.connected = false;
            console.log('Kafka logger disconnected');
        }
    }
}

export default KafkaLogger;
