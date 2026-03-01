package dev.psyconnect.identity_service.controller.v1;

import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.dto.request.CloudinarySignRequest;
import dev.psyconnect.identity_service.dto.response.CloudinarySignResponse;
import dev.psyconnect.identity_service.service.CloudinaryService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequestMapping("/v1/identity/cloudinary")
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
@RequiredArgsConstructor
public class CloudinaryController {
    
    CloudinaryService cloudinaryService;

    @PostMapping("/sign")
    public ApiResponse<CloudinarySignResponse> sign(@RequestBody CloudinarySignRequest request) {
        return new ApiResponse<>(cloudinaryService.generateSignature(request));
    }
}
