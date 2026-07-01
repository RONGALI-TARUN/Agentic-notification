package com.example.notification.model;

import java.util.List;

public class NotificationRequest {
    private AgentInfo agent;
    private String title;
    private String message;
    private List<String> emailRecipients;
    private List<String> slackWebhookUrls;
    private List<String> teamsWebhookUrls;

    public AgentInfo getAgent() { return agent; }
    public void setAgent(AgentInfo agent) { this.agent = agent; }
    public String getTitle() { return title; }
    public void setTitle(String title) { this.title = title; }
    public String getMessage() { return message; }
    public void setMessage(String message) { this.message = message; }
    public List<String> getEmailRecipients() { return emailRecipients; }
    public void setEmailRecipients(List<String> emailRecipients) { this.emailRecipients = emailRecipients; }
    public List<String> getSlackWebhookUrls() { return slackWebhookUrls; }
    public void setSlackWebhookUrls(List<String> slackWebhookUrls) { this.slackWebhookUrls = slackWebhookUrls; }
    public List<String> getTeamsWebhookUrls() { return teamsWebhookUrls; }
    public void setTeamsWebhookUrls(List<String> teamsWebhookUrls) { this.teamsWebhookUrls = teamsWebhookUrls; }
}
