package com.example.notification.controller;

import com.example.notification.model.NotificationRequest;
import com.example.notification.service.NotificationService;
import com.example.notification.util.JsonUtils;
import com.example.notification.util.ResponseUtils;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;

import java.io.IOException;

public class NotificationController implements HttpHandler {
    private final NotificationService notificationService;

    public NotificationController(NotificationService notificationService) {
        this.notificationService = notificationService;
    }

    @Override
    public void handle(HttpExchange exchange) throws IOException {
        if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            ResponseUtils.writeJson(exchange, 405, ResponseUtils.error("method not allowed"));
            return;
        }

        NotificationRequest request = JsonUtils.fromBody(exchange.getRequestBody(), NotificationRequest.class);
        if (request == null) {
            ResponseUtils.writeJson(exchange, 400, ResponseUtils.error("invalid request body"));
            return;
        }

        try {
            notificationService.send(request);
            ResponseUtils.writeJson(exchange, 202, ResponseUtils.message("accepted"));
        } catch (Exception e) {
            ResponseUtils.writeJson(exchange, 500, ResponseUtils.error(e.getMessage()));
        }
    }

    public void handleHealth(HttpExchange exchange) throws IOException {
        ResponseUtils.writeJson(exchange, 200, ResponseUtils.message("ok"));
    }

    public void handleMetrics(HttpExchange exchange) throws IOException {
        ResponseUtils.writeJson(exchange, 200, ResponseUtils.message("metrics not implemented"));
    }
}
