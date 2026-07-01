package com.example.notification.service.impl;

import com.example.notification.config.Config;
import com.example.notification.model.NotificationRequest;
import com.example.notification.service.NotificationService;
import com.example.notification.util.EmailUtils;
import com.example.notification.util.HttpUtils;

import java.util.Map;

public class NotificationServiceImpl implements NotificationService {
    private final Config config;

    public NotificationServiceImpl(Config config) {
        this.config = config;
    }

    @Override
    public void send(NotificationRequest request) throws Exception {
        if (request.getAgent() == null || request.getAgent().getName() == null) {
            throw new IllegalArgumentException("agent name required");
        }
        if (request.getTitle() == null) {
            throw new IllegalArgumentException("title required");
        }

        if (request.getEmailRecipients() != null && !request.getEmailRecipients().isEmpty()) {
            EmailUtils.sendMail(config, request.getEmailRecipients(), request.getTitle(), request.getMessage(), request.getAgent());
        }

        if (request.getSlackWebhookUrls() != null) {
            for (String url : request.getSlackWebhookUrls()) {
                HttpUtils.sendWebhook(url, Map.of("text", request.getMessage()));
            }
        }

        if (request.getTeamsWebhookUrls() != null) {
            for (String url : request.getTeamsWebhookUrls()) {
                HttpUtils.sendWebhook(url, Map.of("text", request.getMessage()));
            }
        }
    }
}
