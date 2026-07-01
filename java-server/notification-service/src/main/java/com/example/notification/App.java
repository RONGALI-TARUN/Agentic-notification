package com.example.notification;

import com.sun.net.httpserver.HttpServer;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpExchange;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import com.fasterxml.jackson.databind.ObjectMapper;

public class App {
    public static void main(String[] args) throws Exception {
        MailConfig mailConfig = new MailConfig(
                System.getenv().getOrDefault("SMTP_HOST", "smtp.example.com"),
                Integer.parseInt(System.getenv().getOrDefault("SMTP_PORT", "587")),
                System.getenv().getOrDefault("SMTP_USERNAME", "user@example.com"),
                System.getenv().getOrDefault("SMTP_PASSWORD", "password"),
                System.getenv().getOrDefault("SMTP_SENDER", "noreply@example.com")
        );

        HttpServer server = HttpServer.create(new InetSocketAddress(8080), 0);
        ObjectMapper mapper = new ObjectMapper();
        server.createContext("/notify", new NotificationHandler(mailConfig, mapper));
        server.createContext("/health", new HttpHandler() {
            @Override
            public void handle(HttpExchange exchange) throws IOException {
                byte[] resp = "ok".getBytes(StandardCharsets.UTF_8);
                exchange.sendResponseHeaders(200, resp.length);
                try (OutputStream os = exchange.getResponseBody()) { os.write(resp); }
            }
        });

        server.start();
        System.out.println("Java production server listening on http://localhost:8080");
    }
}
