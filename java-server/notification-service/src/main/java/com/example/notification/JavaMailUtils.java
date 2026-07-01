package com.example.notification;

import javax.mail.Authenticator;
import javax.mail.Message;
import javax.mail.PasswordAuthentication;
import javax.mail.Session;
import javax.mail.Transport;
import javax.mail.internet.InternetAddress;
import javax.mail.internet.MimeMessage;
import java.util.List;
import java.util.Properties;

public final class JavaMailUtils {
    private JavaMailUtils() {}

    public static void sendMail(MailConfig cfg, List<String> recipients, String subject, String body) throws Exception {
        Properties props = new Properties();
        props.put("mail.smtp.host", cfg.getHost());
        props.put("mail.smtp.port", String.valueOf(cfg.getPort()));
        props.put("mail.smtp.auth", "true");
        props.put("mail.smtp.starttls.enable", "true");

        Session session = Session.getInstance(props, new Authenticator() {
            protected PasswordAuthentication getPasswordAuthentication() {
                return new PasswordAuthentication(cfg.getUsername(), cfg.getPassword());
            }
        });

        Message msg = new MimeMessage(session);
        msg.setFrom(new InternetAddress(cfg.getSender()));
        InternetAddress[] addrs = recipients.stream().map(r -> {
            try { return new InternetAddress(r); } catch (Exception e) { throw new RuntimeException(e); }
        }).toArray(InternetAddress[]::new);
        msg.setRecipients(Message.RecipientType.TO, addrs);
        msg.setSubject(subject);
        msg.setText(body);
        Transport.send(msg);
    }
}
