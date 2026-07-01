package com.example.notification.config;

public class Config {
    private final int serverPort;
    private final String smtpHost;
    private final int smtpPort;
    private final String smtpUsername;
    private final String smtpPassword;
    private final String smtpSender;

    public Config(int serverPort, String smtpHost, int smtpPort, String smtpUsername, String smtpPassword, String smtpSender) {
        this.serverPort = serverPort;
        this.smtpHost = smtpHost;
        this.smtpPort = smtpPort;
        this.smtpUsername = smtpUsername;
        this.smtpPassword = smtpPassword;
        this.smtpSender = smtpSender;
    }

    public static Config load() {
        int serverPort = Integer.parseInt(System.getenv().getOrDefault("SERVER_PORT", "8080"));
        int smtpPort = Integer.parseInt(System.getenv().getOrDefault("SMTP_PORT", "587"));
        return new Config(
                serverPort,
                System.getenv().getOrDefault("SMTP_HOST", "smtp.example.com"),
                smtpPort,
                System.getenv().getOrDefault("SMTP_USERNAME", "user@example.com"),
                System.getenv().getOrDefault("SMTP_PASSWORD", "password"),
                System.getenv().getOrDefault("SMTP_SENDER", "noreply@example.com")
        );
    }

    public int getServerPort() { return serverPort; }
    public String getSmtpHost() { return smtpHost; }
    public int getSmtpPort() { return smtpPort; }
    public String getSmtpUsername() { return smtpUsername; }
    public String getSmtpPassword() { return smtpPassword; }
    public String getSmtpSender() { return smtpSender; }
}
