#!/bin/bash

# Đặt profile mặc định là dev
SPRING_PROFILE="dev"

# Mảng chứa các tiến trình cần dừng sau khi chạy xong
PORTS_TO_KILL=(8080 8081 8082 8888)

# Duyệt qua tất cả thư mục trong dự án để chạy các service Spring Boot
for dir in /Users/chessy/CodeSpace/PsyConnect/*/ ; do
    if [[ -f "$dir/pom.xml" ]]; then
        echo "▶️  Đang chạy service: $(basename "$dir")"
        (cd "$dir" && mvn spring-boot:run -Dspring-boot.run.profiles=$SPRING_PROFILE &)
    fi
done

# Chạy Notification Service (Node.js)
NOTIFICATION_SERVICE_DIR="/Users/chessy/CodeSpace/PsyConnect/notificationservice"
if [[ -f "$NOTIFICATION_SERVICE_DIR/main.js" ]]; then
    echo "▶️  Đang chạy Notification Service (Node.js)"
    (cd "$NOTIFICATION_SERVICE_DIR" && npm install && node main.js &)
else
    echo "❌ Không tìm thấy main.js trong notificationservice!"
fi

# Đợi tất cả các tiến trình chạy xong
wait

echo "🛑 Đang dừng các service trên cổng: ${PORTS_TO_KILL[*]}"

# Dừng tiến trình dựa trên danh sách cổng
for port in "${PORTS_TO_KILL[@]}"; do
    pid=$(lsof -ti :$port)  # Tìm PID của tiến trình chạy trên cổng
    if [[ -n "$pid" ]]; then
        echo "⚡ Killing process on port $port (PID: $pid)"
        kill -9 "$pid"
    else
        echo "✅ Không có tiến trình nào chạy trên cổng $port"
    fi
done

echo "🎯 Tất cả tiến trình đã được dừng!"