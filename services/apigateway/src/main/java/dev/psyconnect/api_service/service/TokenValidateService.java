package dev.psyconnect.api_service.service;


import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.client.RestTemplate;

@Service
public class TokenValidateService {

    private final RestTemplate restTemplate;
    private final String identityServiceUrl;

    @Autowired
    public TokenValidateService(RestTemplate restTemplate,
            @Value("${IDENTITY_SERVICE_URL:http://identityservice:8080}") String identityServiceUrl) {
        this.restTemplate = restTemplate;
        this.identityServiceUrl = identityServiceUrl;
    }

    public boolean isTokenValid(String token) {
        String url = identityServiceUrl + "/auth/internal/valid";
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<String> entity = new HttpEntity<>(token, headers);
        Boolean isValid = restTemplate.postForObject(url, entity, Boolean.class);
        return Boolean.TRUE.equals(isValid);
    }
}
