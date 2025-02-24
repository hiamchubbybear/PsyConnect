package dev.psyconnect.profile_service.service;

import lombok.Getter;
import org.springframework.context.ApplicationEvent;

@Getter
public class OnProfileCreatedEvent extends ApplicationEvent {
    private final String profileId;

    public OnProfileCreatedEvent(Object source, String profileId) {
        super(source);
        this.profileId = profileId;
    }
}
