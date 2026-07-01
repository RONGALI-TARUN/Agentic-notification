package com.example.notification.util;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.Map;

public final class HttpUtils {
    private static final HttpClient CLIENT = HttpClient.newHttpClient();

    private HttpUtils() {}

    public static void sendWebhook(String url, Map<String, String> payload) throws Exception {
        String body = String.format("{\"text\":\"%s\"}", payload.get("text"));
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(body))
                .build();
        HttpResponse<String> response = CLIENT.send(request, HttpResponse.BodyHandlers.ofString());
        if (response.statusCode() < 200 || response.statusCode() >= 300) {
            throw new RuntimeException("webhook failure: " + response.body());
        }
    }
}
