package com.example.notification.server;

import com.example.notification.config.Config;
import com.example.notification.controller.NotificationController;
import com.example.notification.service.NotificationService;
import com.example.notification.service.impl.NotificationServiceImpl;
import com.sun.net.httpserver.HttpServer;

import java.net.InetSocketAddress;

public class Application {
    public static void main(String[] args) throws Exception {
        Config config = Config.load();
        NotificationService notificationService = new NotificationServiceImpl(config);
        NotificationController controller = new NotificationController(notificationService);

        HttpServer server = HttpServer.create(new InetSocketAddress(config.getServerPort()), 0);
        server.createContext("/notify", controller);
        server.createContext("/health", exchange -> {
            controller.handleHealth(exchange);
        });
        server.createContext("/metrics", exchange -> {
            controller.handleMetrics(exchange);
        });
        server.setExecutor(null);
        server.start();
        System.out.println("Notification service started on port " + config.getServerPort());
    }
}
