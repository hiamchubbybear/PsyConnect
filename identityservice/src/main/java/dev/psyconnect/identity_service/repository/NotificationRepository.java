package dev.psyconnect.identity_service.repository;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import com.google.gson.Gson;

import dev.psyconnect.identity_service.dto.request.CreateAccountNotificationRequest;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;

@Component
public class NotificationRepository {
    private final Gson gson;

    @Autowired
    public NotificationRepository(Gson gson) {
        this.gson = gson;
    }

    public HttpResponse<String> notificationSend(CreateAccountNotificationRequest request) {
        HttpClient httpClient = HttpClient.newHttpClient();
        String tojsonRequest = gson.toJson(request);
        HttpRequest httpRequest = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:8082/noti/internal/user"))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(tojsonRequest))
                .build();
        try {
            HttpResponse<String> response = httpClient.send(httpRequest, HttpResponse.BodyHandlers.ofString());
            return response;
        } catch (Exception e) {
            throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
        }
    }
}
