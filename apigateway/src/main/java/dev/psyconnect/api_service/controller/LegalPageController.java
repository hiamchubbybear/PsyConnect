package dev.psyconnect.api_service.controller;

import java.io.IOException;

import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class LegalPageController {

    @GetMapping("/privacy-policy")
    public ResponseEntity<Resource> privacyPolicy() throws IOException {
        Resource html = new ClassPathResource("static/privacy-policy.html");
        return ResponseEntity.ok()
                .contentType(MediaType.TEXT_HTML)
                .body(html);
    }

    @GetMapping("/delete-data-page")
    public ResponseEntity<Resource> deleteData() throws IOException {
        Resource html = new ClassPathResource("static/delete-data.html");
        return ResponseEntity.ok()
                .contentType(MediaType.TEXT_HTML)
                .body(html);
    }

    @GetMapping("/")
    public ResponseEntity<Resource> indexPage() throws IOException {
        Resource html = new ClassPathResource("static/index.html");
        return ResponseEntity.ok()
                .contentType(MediaType.TEXT_HTML)
                .body(html);
    }
}
