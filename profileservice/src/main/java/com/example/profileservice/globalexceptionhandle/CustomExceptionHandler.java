package com.example.profileservice.globalexceptionhandle;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Setter
@Getter
@NoArgsConstructor
public class CustomExceptionHandler extends IllegalStateException {
    public CustomExceptionHandler(ErrorCode errorCode) {
        super(errorCode.getStatus());
        this.errorCode = errorCode;
    }

    private ErrorCode errorCode;
}
