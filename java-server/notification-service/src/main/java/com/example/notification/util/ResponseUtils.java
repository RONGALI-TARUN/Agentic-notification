package com.example.notification.util;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.sun.net.httpserver.HttpExchange;

import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.util.Map;

public final class ResponseUtils {
    private static final ObjectMapper MAPPER = new ObjectMapper();

    private ResponseUtils() {}

    public static void writeJson(HttpExchange exchange, int status, Object payload) {
        try {
            byte[] bytes = MAPPER.writeValueAsString(payload).getBytes(StandardCharsets.UTF_8);
            exchange.getResponseHeaders().set("Content-Type", "application/json");
            exchange.sendResponseHeaders(status, bytes.length);
            try (OutputStream os = exchange.getResponseBody()) {
                os.write(bytes);
            }
        } catch (Exception e) {
            try {
                byte[] bytes = MAPPER.writeValueAsString(error("internal server error")).getBytes(StandardCharsets.UTF_8);
                exchange.getResponseHeaders().set("Content-Type", "application/json");
                exchange.sendResponseHeaders(500, bytes.length);
                try (OutputStream os = exchange.getResponseBody()) {
                    os.write(bytes);
                }
            } catch (Exception ignored) {
            }
        }
    }

    public static Map<String, String> error(String message) {
        return Map.of("error", message);
    }

    public static Map<String, String> message(String message) {
        return Map.of("message", message);
    }
}
