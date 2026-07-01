package com.example.notification.util;

import com.example.notification.config.Config;
import com.example.notification.model.AgentInfo;

import javax.mail.Authenticator;
import javax.mail.Message;
import javax.mail.PasswordAuthentication;
import javax.mail.Session;
import javax.mail.Transport;
import javax.mail.internet.InternetAddress;
import javax.mail.internet.MimeMessage;
import java.util.List;
import java.util.Properties;
import java.util.stream.Collectors;

public final class EmailUtils {
    private EmailUtils() {}

    public static void sendMail(Config config, List<String> recipients, String title, String body, AgentInfo agent) throws Exception {
        Properties props = new Properties();
        props.put("mail.smtp.host", config.getSmtpHost());
        props.put("mail.smtp.port", String.valueOf(config.getSmtpPort()));
        props.put("mail.smtp.auth", "true");
        props.put("mail.smtp.starttls.enable", "true");

        Session session = Session.getInstance(props, new Authenticator() {
            @Override
            protected PasswordAuthentication getPasswordAuthentication() {
                return new PasswordAuthentication(config.getSmtpUsername(), config.getSmtpPassword());
            }
        });

        Message message = new MimeMessage(session);
        message.setFrom(new InternetAddress(config.getSmtpSender()));
        message.setRecipients(Message.RecipientType.TO, InternetAddress.parse(recipients.stream().collect(Collectors.joining(","))));
        message.setSubject("[n8n agent] " + title);
        message.setText(body + "\n\nAgent: " + agent.getName() + "\nOwner: " + agent.getOwner() + "\nGroup: " + agent.getGroup());
        Transport.send(message);
    }
}
