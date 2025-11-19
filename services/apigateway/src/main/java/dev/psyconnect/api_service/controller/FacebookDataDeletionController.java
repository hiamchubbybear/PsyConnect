package dev.psyconnect.api_service.controller;

import java.util.HashMap;
import java.util.Map;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class FacebookDataDeletionController {

    @GetMapping("/delete-data")
    public Map<String, String> handleFacebookDataDeletion(@RequestParam("user_id") String userId) {
        Map<String, String> response = new HashMap<>();
        response.put("url", "https://api.chessy.dev/deletion-status?id=" + userId);
        response.put("confirmation_code", userId);
        return response;
    }

    @GetMapping("/deletion-status")
    public Map<String, String> getDeletionStatus(@RequestParam("id") String id) {
        Map<String, String> response = new HashMap<>();
        response.put("status", "User data deletion in progress (or completed)");
        response.put("id", id);
        return response;
    }
}
