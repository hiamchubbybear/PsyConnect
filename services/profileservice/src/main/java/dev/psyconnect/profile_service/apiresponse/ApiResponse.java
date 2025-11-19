package dev.psyconnect.profile_service.apiresponse;

import org.springframework.http.HttpStatus;

import com.fasterxml.jackson.annotation.JsonInclude;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;

@JsonInclude(JsonInclude.Include.NON_NULL)
@Getter
@AllArgsConstructor
@Builder
@NoArgsConstructor
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
