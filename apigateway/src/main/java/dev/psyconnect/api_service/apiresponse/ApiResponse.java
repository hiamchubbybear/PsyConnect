package dev.psyconnect.api_service.apiresponse;

import com.fasterxml.jackson.annotation.JsonInclude;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import org.springframework.http.HttpStatus;

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
