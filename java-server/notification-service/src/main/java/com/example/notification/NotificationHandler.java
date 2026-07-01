package com.example.notification;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class NotificationHandler implements HttpHandler {
    private final MailConfig mailConfig;
    private final ObjectMapper objectMapper;

    public NotificationHandler(MailConfig mailConfig, ObjectMapper objectMapper) {
        this.mailConfig = mailConfig;
        this.objectMapper = objectMapper;
    }

    @Override
    public void handle(HttpExchange exchange) throws IOException {
        if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            exchange.sendResponseHeaders(405, 0);
            return;
        }

        NotificationRequest request = objectMapper.readValue(exchange.getRequestBody(), NotificationRequest.class);
        String subject = "[n8n agent] " + request.getTitle();
        String body = request.getMessage() + "\n\nAgent: " + request.getAgent().getName()
                + "\nOwner: " + request.getAgent().getOwner()
                + "\nGroup: " + request.getAgent().getGroup();

        StringBuilder result = new StringBuilder();
        boolean failed = false;

        if (request.getEmailRecipients() != null && !request.getEmailRecipients().isEmpty()) {
            try {
                JavaMailUtils.sendMail(mailConfig, request.getEmailRecipients(), subject, body);
                result.append("email sent\n");
            } catch (Exception e) {
                failed = true;
                result.append("email failure: ").append(e.getMessage()).append("\n");
            }
        }

        if (request.getSlackWebhookUrls() != null) {
            for (String url : request.getSlackWebhookUrls()) {
                try {
                    JavaHttpUtils.sendWebhook(url, objectMapper.writeValueAsString(Map.of("text", body)));
                    result.append("slack sent\n");
                } catch (Exception e) {
                    failed = true;
                    result.append("slack failure: ").append(e.getMessage()).append("\n");
                }
            }
        }

        if (request.getTeamsWebhookUrls() != null) {
            for (String url : request.getTeamsWebhookUrls()) {
                try {
                    JavaHttpUtils.sendWebhook(url, objectMapper.writeValueAsString(Map.of("text", body)));
                    result.append("teams sent\n");
                } catch (Exception e) {
                    failed = true;
                    result.append("teams failure: ").append(e.getMessage()).append("\n");
                }
            }
        }

        byte[] response = result.toString().getBytes(StandardCharsets.UTF_8);
        exchange.sendResponseHeaders(failed ? 500 : 200, response.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(response);
        }
    }
}
