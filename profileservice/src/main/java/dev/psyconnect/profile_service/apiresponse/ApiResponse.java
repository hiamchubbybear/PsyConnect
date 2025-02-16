package dev.psyconnect.profile_service.apiresponse;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonInclude;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import org.springframework.http.HttpStatus;

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
