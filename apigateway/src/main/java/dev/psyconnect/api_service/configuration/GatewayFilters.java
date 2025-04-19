package dev.psyconnect.api_service.configuration;

import org.springframework.context.annotation.Bean;
import org.springframework.web.reactive.function.server.RouterFunction;
import org.springframework.web.reactive.function.server.ServerResponse;

import static io.micrometer.core.instrument.binder.jetty.JettyClientKeyValues.uri;
import static org.springframework.boot.webservices.client.WebServiceMessageSenderFactory.http;
import static org.springframework.web.reactive.function.server.RouterFunctions.route;

public class GatewayFilters {
}
