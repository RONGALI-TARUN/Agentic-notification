package com.example.notification.util;

import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.InputStream;

public final class JsonUtils {
    private static final ObjectMapper MAPPER = new ObjectMapper();

    private JsonUtils() {}

    public static <T> T fromBody(InputStream body, Class<T> clazz) {
        try {
            return MAPPER.readValue(body, clazz);
        } catch (Exception e) {
            return null;
        }
    }
}
