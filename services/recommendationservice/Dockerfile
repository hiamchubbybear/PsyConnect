FROM python:3.9-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Environment setup
ENV PORT=8086
ENV FLASK_ENV=production

EXPOSE 8086

# Run using gunicorn for production
CMD ["gunicorn", "--bind", "0.0.0.0:8086", "api.main:app"]
