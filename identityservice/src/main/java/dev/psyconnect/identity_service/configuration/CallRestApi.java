package dev.psyconnect.identity_service.configuration;

import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpHeaders;

import java.util.Collections;
import org.springframework.http.*;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestTemplate;
@Component
public class CallRestApi {

    public String callGooglePeopleApi(String accessToken) {
        RestTemplate restTemplate = new RestTemplate();

        String url = "https://people.googleapis.com/v1/people/me?personFields=names,genders,birthdays,emailAddresses,addresses";


        HttpHeaders headers = new HttpHeaders();
        headers.setBearerAuth(accessToken);
        headers.setAccept(Collections.singletonList(MediaType.APPLICATION_JSON));

        HttpEntity<String> entity = new HttpEntity<>(headers);

        ResponseEntity<String> response = restTemplate.exchange(url, HttpMethod.GET, entity, String.class);

        return response.getBody(); // JSON string từ Google API
    }
}
