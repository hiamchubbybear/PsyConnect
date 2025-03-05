package dev.psyconnect.kafka.entity;

import lombok.*;
import org.springframework.messaging.handler.annotation.SendTo;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class Greeting {
    String message;
    String name;
}
