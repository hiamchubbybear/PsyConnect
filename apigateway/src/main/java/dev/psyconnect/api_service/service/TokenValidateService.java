package dev.psyconnect.api_service.service;


import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
public class TokenValidateService {

    private final RestTemplate restTemplate;
    private final String identityServiceUrl;
    private final String identityServicePort;
    @Autowired
    public TokenValidateService(RestTemplate restTemplate, @Value("${base.url}") String identityServiceUrl,
                                @Value("${identity.port}") String identityServicePort) {
        this.restTemplate = restTemplate;
        this.identityServiceUrl = identityServiceUrl;
        this.identityServicePort = identityServicePort;
    }

    public boolean isTokenValid(String token) {
        String url = identityServiceUrl+ ":"+identityServicePort  + "/auth/internal/valid/" + token;
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<String> entity = new HttpEntity<>(token, headers);
        Boolean isValid = restTemplate.postForObject(url, entity, Boolean.class);
        return isValid;
    }
}
