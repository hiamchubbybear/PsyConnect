package dev.psyconnect.profile_service.apiresponse;

import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class ApiResponse<T> {
    int code;
    String message;
    T data;

    public ApiResponse(T data) {
        this.code = 200;
        this.message = "Success";
        this.data = data;
    }
}
