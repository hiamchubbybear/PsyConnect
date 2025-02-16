package dev.psyconnect.identity_service.apiresponse;

import com.fasterxml.jackson.annotation.JsonInclude;
import org.springframework.http.HttpStatus;

import lombok.*;

@Getter
@AllArgsConstructor
@Builder
@NoArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
public class ApiResponse<T> {
    private int code = HttpStatus.OK.value();
    private String message;
    private T data;

    public ApiResponse(T data) {
        this.code = HttpStatus.OK.value();
        this.message = "Success";
        this.data = data;
    }
}
