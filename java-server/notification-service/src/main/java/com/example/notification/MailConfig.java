package com.example.notification;

public class MailConfig {
    private final String host;
    private final int port;
    private final String username;
    private final String password;
    private final String sender;

    public MailConfig(String host, int port, String username, String password, String sender) {
        this.host = host;
        this.port = port;
        this.username = username;
        this.password = password;
        this.sender = sender;
    }

    public String getHost() { return host; }
    public int getPort() { return port; }
    public String getUsername() { return username; }
    public String getPassword() { return password; }
    public String getSender() { return sender; }
}
